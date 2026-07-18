## 1. Go: Extend apply.go with --with flag and raw input support

- [x] 1.1 Add `--with` flag variable and registration in `apply.go` init()
- [x] 1.2 Add mutual exclusivity validation: `--with` cannot be used with `--use`
- [x] 1.3 Add validation: `--with` requires `--convert` or `-f -` (stdin)
- [x] 1.4 Modify `runConversion()` to detect when `--convert` value is not an Asset name and treat as raw string input
- [x] 1.5 Modify `runConversion()` to support `--with` inline expression mode (skip Converter lookup, use expression directly)
- [x] 1.6 Add stdin raw mode: when `-f -` and `--with` are used together without `--convert`, pass stdin directly to yq-go engine
- [x] 1.7 Improve error messages: list available Assets before falling back to raw string interpretation

## 2. Go: Converter engine - raw input support

- [x] 2.1 Add `ApplyRaw()` method to Engine that accepts raw string input (not an Asset struct) and a Converter or inline expression
- [x] 2.2 Add format inference: detect if raw input is JSON or YAML based on content
- [x] 2.3 Update `MapManifestAssets()` or add `RawAsset()` constructor that creates an Asset from raw string with format detection
- [x] 2.4 Unit tests: raw JSON input with converter expression
- [x] 2.5 Unit tests: raw YAML input with inline expression
- [x] 2.6 Unit tests: format detection edge cases

## 3. Go: Apply command unit tests

- [x] 3.1 Test: `--convert` with Asset name (existing behavior unchanged)
- [x] 3.2 Test: `--convert` with raw JSON string (fallback when no Asset match)
- [x] 3.3 Test: `--with` inline expression with Asset lookup
- [x] 3.4 Test: `--with` with raw string input (no Asset)
- [x] 3.5 Test: `--with` and `--use` mutual exclusivity (exit code 2)
- [x] 3.6 Test: `-f -` with `--with` stdin pipeline
- [x] 3.7 Test: `-f -` with `--with` and `-o` output file
- [x] 3.8 Test: `--without --convert` and without stdin (exit code 2)

## 4. Manifests: Create aggregate-mcp Step

- [x] 4.1 Create `.manifests/base/steps/aggregate-mcp.yaml` Step definition
- [x] 4.2 Step uses `kfg apply --convert ASSET --use CONVERTER` to convert asset to temp variable
- [x] 4.3 Step uses `kfg apply -f - --with` for multi-document merge into TARGET file
- [x] 4.4 Step handles TARGET file existence (create vs merge)
- [x] 4.5 Step calls `__kfg_add_artifact` for TARGET file
- [x] 4.6 Step includes structured logging for aggregate operations

## 5. Manifests: Update dev workflow Phase 5

- [x] 5.1 Uncomment and rework Phase 5 in `.manifests/overlay/dev/agents-workflow.yaml`
- [x] 5.2 Add `kfg.aggregate-mcp` step calls for claude (context7, chrome-devtools, playwright)
- [x] 5.3 Add `kfg.aggregate-mcp` step calls for gemini (context7, chrome-devtools, playwright)
- [x] 5.4 Add `kfg.aggregate-mcp` step calls for opencode (context7, chrome-devtools, playwright)
- [x] 5.5 Add `when` conditions per agent using `kfg.detect-agent` output
- [x] 5.6 Update cleanup step to remove any temp files if used

## 6. Validation and testing

- [x] 6.1 Run `make test` — all Go unit tests pass
- [x] 6.2 Run `make test-bats` — 98/106 pass (8 failures are pre-existing: bats-assert compat, version format, workspace)
- [x] 6.3 Manual test: `kfg apply -f manifest.yaml --convert '{"key":"value"}' --with '.key'` outputs `value`
- [x] 6.4 Manual test: `echo '{"a":1}---{"b":2}' | kfg apply -f - --with 'select(di == 0) * select(di == 1)'` outputs merged JSON
- [ ] 6.5 Manual test: Run `dev.workflows.dev` workflow with claude agent and verify `.claude/settings.local.json` contains all 3 MCPs
- [x] 6.6 Run `make fmt && make lint && make vet` — no issues
