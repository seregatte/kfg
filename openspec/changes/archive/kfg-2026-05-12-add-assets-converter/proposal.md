## Why

kfg currently only supports execution-oriented resource kinds (Step, Cmd, CmdWorkflow) that generate shell code. There is no way to declaratively define generic data payloads and transform them into different formats. This limits kfg's use cases to shell function generation, preventing it from being used as a general-purpose YAML/data transformation tool. Adding Assets and Converter resource kinds enables kfg to process structured data and output it in any format supported by yq-go (YAML, JSON, TOML, XML, etc.), making it a flexible data transformation engine alongside its existing shell generation capabilities.

## What Changes

- **New `Assets` resource kind**: Declares typed data payloads with `spec.data` (YAML map or string) and `spec.input.format` (data format: yaml, json, toml, xml, etc.)
- **New `Converter` resource kind**: Declares transformations using yq-go expressions (`spec.engine.expression`), with explicit input format (`spec.input.format`) and output format (`spec.output.format`)
- **New yq-go dependency**: `github.com/mikefarah/yq/v4` as Go module for expression evaluation (no external CLI dependency)
- **Extended `apply` command**: New `--convert <asset-name>` and `--use <converter-name>` flags for conversion mode
- **Mutually exclusive modes**: Shell generation flags (`-w`, `-c`) and conversion flags (`--convert`, `--use`) cannot be mixed; validation produces clear error messages
- **All yq-go formats supported**: Input and output support for yaml, json, xml, props, csv, tsv, toml, hcl, lua, ini, shell, base64, uri, kyaml
- **`raw` output format**: Special format that outputs string result without YAML/JSON encoding (for markdown with frontmatter, plain text, etc.)

## Capabilities

### New Capabilities

- `assets-converter-model`: Defines the Assets and Converter resource kinds, their YAML structure, validation rules, and how they integrate with the existing manifest pipeline (parser, adapter, resolver)
- `converter-engine`: Implements the yq-go based transformation engine that evaluates expressions against Assets and produces output in the specified format, including format conversion and `raw` output handling
- `apply-conversion-mode`: Extends the `apply` command with `--convert` and `--use` flags, mutual exclusivity validation with shell generation flags, and conversion pipeline execution

### Modified Capabilities

- `apply-command`: Add conversion mode flags (`--convert`, `--use`) and mutual exclusivity validation with existing flags (`-w`, `-c`)
- `manifest-model`: Add `Assets` and `Converter` to `SupportedKinds`, update `ParsedResource` struct, and extend parser/adapter to handle new kinds

## Impact

- **New packages**: `src/internal/converter/` (types.go, engine.go, engine_test.go)
- **Modified packages**: `src/internal/manifest/` (types.go, parser.go), `src/internal/kustomize/` (adapter.go), `src/internal/resolve/` (resolve.go), `src/cmd/kfg/` (apply.go)
- **New dependency**: `github.com/mikefarah/yq/v4` added to go.mod
- **CLI impact**: `kfg apply` gains `--convert` and `--use` flags; help text updated to document two modes
- **Testing**: New unit tests for converter engine (all formats), integration tests for Assets/Converter loading and conversion
