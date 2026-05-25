## Context

The kfg repository recently refactored its `.manifests/base/` structure across three major commits:
- `e41ab21`: Reorganized AI agents under `extensions/ai-agents/`
- `44f186a`: Consolidated all extensions under unified `extensions/ai/` directory
- `19fa1e8`: Introduced `kfg.materialize` step to replace `kfg.convert` and `kfg.aggregate-mcp`

The lifeos repository's workflow file references the kfg base via relative path (`../../../../kfg/.manifests/base`), making it dependent on kfg's internal structure. The refactoring broke ~80+ references in `lifeos-workflow.yaml`.

## Goals / Non-Goals

**Goals:**
- Update all references in `lifeos-workflow.yaml` to match new kfg naming conventions
- Replace deprecated step kinds with `kfg.materialize` equivalent
- Validate changes using `kfg apply` command
- Maintain backward compatibility for local lifeos assets (commands, MCPs, subagents)

**Non-Goals:**
- Refactoring local lifeos asset structure (commands, MCPs, subagents remain unchanged)
- Adding new capabilities to lifeos workflow
- Modifying kfg base manifests

## Decisions

1. **Use `kfg.materialize` for all conversions**: The new unified step replaces both `kfg.convert` (per-item mode) and `kfg.aggregate-mcp` (aggregate mode). This simplifies the workflow by using a single step kind with different parameters.

2. **Maintain local asset naming**: LifeOS-specific assets (`lifeos.commands.*`, `lifeos.mcp.*`, `lifeos.subagents.*`) keep their current naming since they're defined in lifeos's own base, not kfg's.

3. **WRAP_KEY for MCP aggregation**: The old `kfg.aggregate-mcp` used `TARGET` parameter; the new `kfg.materialize` uses `OUTPUTS` and `WRAP_KEY`. For OpenCode, `WRAP_KEY: "mcp"` is used; for Claude/Gemini, `WRAP_KEY: "mcpServers"`.

4. **Direct file modification**: Rather than creating a migration script, directly modify `lifeos-workflow.yaml` since it's a single file with well-defined replacements.

## Risks / Trade-offs

- **[Risk]**: kfg apply may fail if some referenced steps don't exist in the new structure → **Mitigation**: Validate each reference against kfg's actual manifest files before applying
- **[Risk]**: Parameter semantics change between old and new steps (e.g., `TARGET` → `OUTPUTS`, `ASSET` → `ASSETS`) → **Mitigation**: Carefully map each parameter according to kfg.materialize spec
- **[Trade-off]**: No automated test suite for lifeos manifests → **Mitigation**: Manual validation via `kfg apply` and syntax checking of generated shell code
