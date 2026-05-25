## Context

kfg currently requires users to provide a kustomization source path explicitly via:
- `build` command: mandatory positional argument
- `apply` command: `-k` flag or positional argument
- `run` command: `-k` flag

There is no way to:
1. Use a GitHub repository URL without manually cloning it first
2. Configure a default source via environment variable

The kustomize library (`sigs.k8s.io/kustomize/api`) already supports GitHub repository URLs natively via its internal git cloner. URLs like `https://github.com/owner/repo//path` are automatically cloned to a temp directory and processed.

## Goals / Non-Goals

**Goals:**
- Enable GitHub URL usage in `build`, `apply`, `run` commands
- Add `KFG_KPATH` environment variable for default source
- Maintain backward compatibility with existing local path workflows
- Keep implementation minimal — leverage kustomize's existing git cloning

**Non-Goals:**
- HTTP URL support (non-GitHub) — out of scope
- GitHub private repositories with authentication — out of scope
- Local caching of cloned repositories — out of scope (kustomize clones fresh each time)
- Customizing clone depth or timeout — out of scope (use kustomize defaults)

## Decisions

### D1: URL Detection Package Location

**Decision:** Create new package `src/internal/urlresolve/`

**Alternatives considered:**
1. Inline detection in each command → rejected: duplicates logic, harder to test
2. Add to existing `kustomize` package → rejected: mixing concerns, kustomize package handles loading not detection
3. New dedicated package → chosen: single responsibility, easy to test, reusable

**Rationale:** URL detection is a distinct concern from kustomization loading. A dedicated package keeps the logic isolated and testable.

### D2: URL Detection Logic

**Decision:** Use simple string matching for GitHub URLs

```go
func IsGitHubURL(arg string) bool {
    return strings.Contains(arg, "github.com") ||
           strings.HasPrefix(arg, "https://github.com") ||
           strings.HasPrefix(arg, "http://github.com")
}
```

**Alternatives considered:**
1. Parse URL with `url.Parse()` → rejected: overkill for this use case, kustomize already handles validation
2. Regex matching → rejected: unnecessary complexity
3. Simple string matching → chosen: sufficient for GitHub URLs, fast, no external deps

**Rationale:** Kustomize will validate the URL during cloning. We only need to detect GitHub URLs to pass them through to the loader. Invalid URLs will fail at clone time with a clear error message.

### D3: Environment Variable Binding

**Decision:** Bind `KFG_KPATH` via Viper in `config.Initialize()`

```go
viper.BindEnv("kpath", "KFG_KPATH")
```

Add getter:
```go
func GetKPath() string {
    return viper.GetString("kpath")
}
```

**Alternatives considered:**
1. Direct `os.Getenv()` in each command → rejected: inconsistent with existing pattern, no viper integration
2. Viper binding with default value → rejected: no sensible default (empty is correct)
3. Viper binding without default → chosen: matches existing env var pattern in codebase

**Rationale:** Follows the established pattern for `KFG_VERBOSE`, `KFG_STORE_DIR`, etc. Viper provides consistent access and allows future override via config file if needed.

### D4: Source Resolution Priority

**Decision:** Priority chain: flag > positional arg > env var > error

| Command | Priority Chain |
|---------|----------------|
| `build` | arg[0] > `KFG_KPATH` > error |
| `apply` | `-k` flag > arg[0] > `KFG_KPATH` > error |
| `run` | `-k` flag > `KFG_KPATH` > error |

**Alternatives considered:**
1. Env var highest priority → rejected: unexpected behavior, flags should override env
2. Error if both flag and env var set → rejected: overly strict, flag should win silently
3. Priority chain (flag > arg > env) → chosen: matches CLI conventions, intuitive

**Rationale:** Users expect explicit CLI flags to override environment defaults. This matches standard CLI behavior (e.g., `git`, `docker`).

### D5: Argument Validation Change

**Decision:** Change `build` command from `cobra.ExactArgs(1)` to `cobra.MaximumNArgs(1)`

**Alternatives considered:**
1. Keep `ExactArgs(1)` → rejected: would require env var to be treated as error case
2. Use `MaximumNArgs(1)` → chosen: allows 0 args when env var is set
3. Use `NoArgs` + special handling → rejected: more complex, less clear

**Rationale:** `MaximumNArgs(1)` allows 0 or 1 argument, matching the new optional behavior. Error handling for "no source" happens in the Run function after checking env var.

### D6: Kustomize Loader Interaction

**Decision:** Pass GitHub URLs directly to `kustomize.NewLoader(nil).Load(url)` without preprocessing

**Alternatives considered:**
1. Clone repo manually before passing path → rejected: duplicates kustomize logic, adds complexity
2. Pre-validate URL format → rejected: kustomize already validates, double validation is redundant
3. Pass URL directly → chosen: simplest, leverages kustomize's existing git cloner

**Rationale:** Kustomize's `krusty.Kustomizer` handles GitHub URLs transparently:
- Parses URL format (owner/repo, // separator, ?ref= parameter)
- Clones to temp directory with `--depth=1`
- Processes kustomization.yaml
- Returns ResMap

No additional code needed beyond passing the URL through.

### D7: Error Messages

**Decision:** Clear error message when no source is available

```
Error: kustomization source required. Provide a path, use -k flag, or set KFG_KPATH.
```

**Alternatives considered:**
1. Generic "invalid argument" → rejected: unhelpful
2. Per-command specific messages → rejected: unnecessary duplication
3. Single clear message → chosen: actionable, consistent

**Rationale:** One message for all commands guides users to the three available options.

## Risks / Trade-offs

**R1: GitHub URL validation**
- **Risk:** Invalid GitHub URLs fail at clone time with potentially confusing git errors
- **Mitigation:** Kustomize provides clear error messages; we log the URL being processed

**R2: Shallow clone performance**
- **Risk:** Large repos with deep history may still take time to clone (--depth=1)
- **Mitigation:** Kustomize handles this; users can use ?ref= to target smaller branches

**R3: Network dependency**
- **Risk:** Commands fail without network access when using GitHub URLs
- **Mitigation:** Clear error from kustomize; fallback to local paths works offline

**R4: URL format confusion**
- **Risk:** Users may confuse GitHub URL formats (missing //, wrong path separators)
- **Mitigation:** Help examples show correct formats; kustomize documentation covers this

**R5: Env var precedence confusion**
- **Risk:** Users may not understand flag > env priority
- **Mitigation:** Help documentation explicitly states priority; behavior matches CLI conventions