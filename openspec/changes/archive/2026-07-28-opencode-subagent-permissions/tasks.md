## 1. Update Converter Expression

- [x] 1.1 Update `manifests/agents/opencode/converters/subagent.yaml` to emit `tools:` list
      from `.tools` when non-empty
- [x] 1.2 Add `permission:` block using `capture` + `ireduce` + `downcase` to transform
      `.permissions.allow[]` entries into per-tool deny-with-exceptions objects
- [x] 1.3 Handle edge cases: empty `.tools`, empty `.permissions.allow`, no permissions key

## 2. Add Converter Tests

- [x] 2.1 Add `subagentExpression` constant to `src/internal/converter/engine_test.go`
- [x] 2.2 Add `TestSubagentMinimal` — identity fields only, no tools/permissions
- [x] 2.3 Add `TestSubagentPermissions` — single tool with one pattern
- [x] 2.4 Add `TestSubagentMultiPatternPermissions` — multiple tools with patterns
- [x] 2.5 Add `TestSubagentNoPermissions` — empty permissions.allow array

## 3. Update Documentation

- [x] 3.1 Update `manifests/subagents/README.md` with permission materialization behavior
- [ ] 3.2 Update `manifests/agents/opencode/README.md` Limitations section to reflect
      that the subagent converter now emits tools and permissions

## 4. Validation

- [x] 4.1 Run `go test ./internal/converter/ -v -run "TestSubagent" -count=1` — all pass
- [x] 4.2 Run `go test ./internal/converter/ -count=1` — full suite passes
- [x] 4.3 Verify generated output includes `tools:`, `permission:`, and correct deny/allow
      structure for test assets
