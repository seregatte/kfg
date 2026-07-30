## Why

The OpenCode subagent converter (`ai.opencode.conv.subagent`) currently only outputs
`name`, `mode`, `description`, `model`, and `prompt` fields. The `tools` and
`permissions.allow` fields defined in the portable asset contract are silently dropped
during conversion, making subagents non-functional when they require specific tools
and permission scopes.

## What Changes

- Update the `packages/domains/ai-agents/manifests/agents/opencode/converters/subagent.yaml`
  Converter expression to materialize `tools` and `permissions.allow` into the OpenCode
  YAML output.
- Generate a `permission:` block following the broad-deny → narrow-allow pattern
  (each tool starts with `{"*": "deny"}` then narrows to specific patterns).
- Add converter tests to the engine test suite.
- Update the subagents README with the new permission behavior.

## Capabilities

### Modified Capabilities

- `opencode-converters`: The subagent converter now materializes `tools` and `permissions`
  from the portable asset contract into the OpenCode subagent YAML.

## Impact

- Affected manifest: `manifests/agents/opencode/converters/subagent.yaml` (Converter expression)
- Affected engine tests: `src/internal/converter/engine_test.go` (new subagent test cases)
- Affected documentation: `manifests/subagents/README.md` (document permission materialization)
- Converter output behavior changes: subagent files now include `tools:` and `permission:` blocks
  when the source asset specifies them, matching the OpenCode subagent contract
