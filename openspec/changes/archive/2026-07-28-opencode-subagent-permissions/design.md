## Context

The subagent converter uses a yq v4 expression evaluated at build time to transform
the portable asset manifest into OpenCode subagent YAML (Markdown frontmatter).
The current expression only outputs identity fields (name, mode, description, model)
and the prompt body, ignoring the structured tools/permissions that OpenCode agents
depend on.

## Goals / Non-Goals

**Goals:**

- Emit `tools:` list when the source asset has a non-empty `.tools` array
- Emit `permission:` object when the source asset has a non-empty `.permissions.allow`
  array, using broad-deny → narrow-allow per-tool semantics
- Preserve backward compatibility for assets without tools/permissions (empty output)
- All converter tests pass

**Non-Goals:**

- Changing the converter engine framework (`NewEngine`, `ApplyRaw`, yq expression format)
- Adding new resource kinds or converter types
- Modifying other converters (skill, command, MCP, config)

## Decisions

### Expression in the Converter manifest

The converter expression lives in `subagent.yaml` as a raw yq expression string evaluated
by `engine.ApplyRaw()` using the Go yqlib (v4.53.2). The Go library supports `ireduce`,
`capture`, `downcase`, `if`/`then`/`else`, and `*` merge — unlike the CLI binary (v4.53.3).

### Permission pattern: broad-deny → narrow-allow

For each permission entry in `.permissions.allow[]` of the form `ToolName(pattern)`:

1. Parse the entry with `capture("^(?P<tool>[A-Za-z][A-Za-z0-9_-]*)\\((?P<pattern>.*)\\)$")`
2. Normalize tool name with `downcase`
3. Build `{<tool>: {"*": "deny"}} * {<tool>: {<pattern>: "allow"}}` per rule
4. `ireduce` across all rules to produce the merged permission object

Example for `Read(/**)` and `Bash(git diff*)`:
```yaml
permission:
  read:
    "/**": "allow"
    "*": "deny"
  bash:
    "git diff*": "allow"
    "*": "deny"
```

### Test strategy

Add a `subagentExpression` constant and dedicated test cases:
- `TestSubagentMinimal` — no tools/permissions, only identity fields
- `TestSubagentPermissions` — single tool with one pattern
- `TestSubagentMultiPatternPermissions` — multiple tools with multiple patterns
- `TestSubagentNoPermissions` — empty permissions.allow array

## Risks / Trade-offs

- The yq expression is complex and relies on Go library features not available in the
  CLI binary. Testing through the Go test suite (`go test ./internal/converter/`) is
  the only reliable verification path.
- The `downcase` function must be used instead of `ascii_downcase` (yq v4.53.2 lexer
  tokenizes `as` inside "ascii" as `AssignAsVariable`, breaking the expression).

## Migration Plan

No migration required. The converter change is forward-only: new output includes
`tools`/`permission` fields; old assets without those fields produce the same output
as before.
