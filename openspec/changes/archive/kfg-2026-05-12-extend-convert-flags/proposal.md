## Why

The current `kfg apply --convert` mode only supports asset name lookup from manifests. This forces MCP aggregation workflows to write to temporary files and merge with external tools like `jq`. The new Assets/Converters paradigm needs `--convert` to accept raw string input and a `--with` flag for inline yq expressions, enabling incremental JSON merge operations without any external dependencies.

## What Changes

- `--convert` flag now accepts both asset names (existing behavior) and raw string input (JSON/YAML literal or piped data)
- New `--with` flag accepts a raw yq expression directly, bypassing Converter lookup by `metadata.name`
- `--with` can read from stdin when used with `-f -`, enabling multi-document merge pipelines
- New `kfg.aggregate-mcp` Step that uses `--with` to perform read-modify-write merge on MCP config files
- MCP Phase 5 in `dev.workflows.dev` uncommented and reworked to use per-MCP conversion + incremental aggregation

## Capabilities

### New Capabilities
- `convert-raw-input`: `--convert` accepts raw string input in addition to asset names; `--with` accepts inline yq expressions; stdin support for multi-document merge

### Modified Capabilities
- `apply-command`: Conversion mode now supports `--with` flag and raw string input for `--convert`

## Impact

- **CLI**: `kfg apply` gains `--with` flag and expanded `--convert` behavior
- **Go code**: `src/cmd/kfg/apply.go` — `runConversion()` extended to handle raw input and inline expressions
- **Manifests**: `.manifests/base/steps/aggregate-mcp.yaml` new Step; `.manifests/overlay/dev/agents-workflow.yaml` Phase 5 uncommented
- **No external dependencies**: Eliminates need for `jq` in MCP aggregation workflows
- **Backward compatible**: Existing `--convert asset --use converter` usage unchanged
