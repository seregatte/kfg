## Context

The Go codebase contains two categories of terminology issues:

1. **AI-specific terminology in public API**: The `kfg run` command is the primary entry point for one-shot execution. Its help text, usage string, function names, error messages, and listing output all reference "agents", even though the command operates generically on any Cmd resource.

2. **Legacy `nixai` references**: Before the project was renamed to kfg, it was called NixAI. Several references remain in Go code, tests, and versioned fixtures.

These issues are independent of the manifest reorganization (which is covered by the `reorganize-ai-agents-manifests` change). This change focuses exclusively on the Go codebase.

## Goals / Non-Goals

**Goals:**
- Make the `kfg run` command terminology generic and use-case-neutral.
- Remove all legacy `nixai` references from Go code, tests, and fixtures.
- Keep behavioral logic completely unchanged.
- Ensure all Go tests pass after renaming.

**Non-Goals:**
- Change manifest files or their resource names.
- Modify the run command's resolution logic or shell generation.
- Add new features.
- Update Bats tests (those are separate from Go tests).

## Decisions

### Use `cmd` as the canonical term in Go

The `kfg run` command operates on Cmd resources from the manifest index. The term `cmd` is:
- Short and consistent with the manifest model (`kind: Cmd`).
- Already used in other parts of the codebase (`runCmds` flag, `Cmds` in workflows).
- Generic and use-case-neutral.

**Naming changes in `run.go`:**

| Current | New |
|---------|-----|
| `Use: "run [agent] [-- extra-args...]"` | `Use: "run [cmd] [-- extra-args...]"` |
| `Short: "Run an agent ..."` | `Short: "Run a command ..."` |
| `findAgent()` | `findCmd()` |
| `generateForAgent()` | `generateForCmd()` |
| `executeAgent()` | `executeCmd()` |
| `listAvailableAgents()` | `listAvailableCmds()` |
| `agentName` | `cmdName` |
| `"No agents found in manifests"` | `"No commands found in manifests"` |
| `"Available agents:"` | `"Available commands:"` |
| `expectedAgent` (test field) | `expectedCmd` |
| `TestFindAgent` | `TestFindCmd` |
| `TestListAvailableAgents` | `TestListAvailableCmds` |

### Keep behavioral logic unchanged

The resolution logic (matching by `commandName`, finding the containing workflow) remains identical. Only names, messages, and examples change.

### Update Long description and examples

Current examples reference `.nixai/overlay/dev` and mention "claude". Replace with generic examples:

```
kfg run -k .manifests/overlay/dev my-cmd
kfg run -k .manifests/overlay/dev my-cmd -- --flag value
kfg run -k https://github.com/owner/repo//manifests my-cmd
kfg run -k .manifests/overlay/dev -w dev my-cmd
kfg run -f manifest.yaml my-cmd
kfg run -k .manifests/overlay/dev (lists available commands)
KFG_KPATH=./manifests kfg run my-cmd
```

### Remove legacy `nixai` references

**1. `src/internal/generate/generate.go:297`**

Current:
```go
code.WriteString("__kfg_build_result_file=$(mktemp -t nixai-build-XXXXXX.yaml)\n")
```

New:
```go
code.WriteString("__kfg_build_result_file=$(mktemp -t kfg-build-XXXXXX.yaml)\n")
```

**2. `src/internal/manifest/types.go:362`**

Current:
```go
// APIVersion is the expected API version for NixAI manifests.
```

New:
```go
// APIVersion is the expected API version for kfg manifests.
```

**3. `src/internal/image/store.go:60`**

Current:
```go
storeDir = ".nixai/store"
```

New:
```go
storeDir = ".kfg/store"
```

**4. `src/internal/kustomize/openapi.json`**

Replace all `nixai.dev` references with `kfg.dev`:
- `info.title`: `"NixAI Custom Resource Definitions"` -> `"kfg Custom Resource Definitions"`
- `definitions`: `nixai.dev.v1alpha1.*` -> `kfg.dev.v1alpha1.*`
- `x-kubernetes-group-version-kind`: `group: "nixai.dev"` -> `group: "kfg.dev"`
- `enum` values: `["nixai.dev/v1alpha1"]` -> `["kfg.dev/v1alpha1"]`

**5. `src/internal/generate/templates/testdata/golden_basic.bash`**

Replace `__nixai_*` prefixes with `__kfg_*` and update comments.

### Update Go tests

`src/cmd/kfg/run_test.go` must be updated to reflect all terminology changes:
- Variable names (`expectedAgent` -> `expectedCmd`)
- Test names (`TestFindAgent` -> `TestFindCmd`)
- String assertions (`"No agents found"` -> `"No commands found"`)
- Test data (Cmd names in fixtures can stay generic like `test.agent`)

Other test files may need minor updates:
- `src/internal/store/artifacts_test.go` (temp dir prefixes)
- `src/internal/config/config_test.go` (path references)

## Data Contract

### CLI surface changes

```
# Before
$ kfg run --help
run [agent] [-- extra-args...]

Run an agent with one-shot execution.

# After
$ kfg run --help
run [cmd] [-- extra-args...]

Run a command with one-shot execution.
```

### Error message changes

```
# Before
Error: agent 'my-cmd' not found

# After
Error: command 'my-cmd' not found
```

### Listing output changes

```
# Before
No agents found in manifests
Available agents:

# After
No commands found in manifests
Available commands:
```

## Risks / Trade-offs

- [Breaking user scripts that parse CLI output] -> The `kfg run` command is primarily interactive; scripted usage should use `kfg apply` or `kfg build`. The terminology change is a minor UX improvement, not a contract-breaking change.
- [Test fixture churn] -> The golden test file uses `__nixai_*` which is a legacy artifact. Updating it improves consistency but may require regeneration if it's auto-generated.
- [OpenAPI schema changes] -> The `openapi.json` is used for schema validation. Changing `nixai.dev` to `kfg.dev` must align with the actual `apiVersion` used in manifests (`kfg.dev/v1alpha1`), which is already the case.

## Migration Plan

1. Update `run.go` terminology (functions, variables, messages, examples).
2. Update `run_test.go` to match new terminology.
3. Remove `nixai` references from `generate.go`, `types.go`, `store.go`.
4. Update `openapi.json` schema definitions.
5. Update testdata fixtures (`golden_basic.bash`).
6. Run `make test` to confirm all Go tests pass.
