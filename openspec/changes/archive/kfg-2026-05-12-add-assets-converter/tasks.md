## 1. Dependencies and Setup

- [x] 1.1 Add `github.com/mikefarah/yq/v4` to go.mod
- [x] 1.2 Run `go mod tidy` to resolve dependencies

## 2. Manifest Types

- [x] 2.1 Add `Assets` struct to `src/internal/manifest/types.go` with `APIVersion`, `Kind`, `Metadata`, `Spec` fields
- [x] 2.2 Add `AssetsSpec` struct with `Input` and `Data` fields
- [x] 2.3 Add `InputSpec` struct with `Format` field
- [x] 2.4 Add `Converter` struct to `src/internal/manifest/types.go` with `APIVersion`, `Kind`, `Metadata`, `Spec` fields
- [x] 2.5 Add `ConverterSpec` struct with `Input`, `Engine`, `Output` fields
- [x] 2.6 Add `EngineSpec` struct with `Expression` field
- [x] 2.7 Add `OutputSpec` struct with `Format` field
- [x] 2.8 Add `SupportedInputFormats` and `SupportedOutputFormats` constants
- [x] 2.9 Add validation methods: `Assets.Validate()`, `Converter.Validate()`
- [x] 2.10 Update `SupportedKinds` to include `"Assets"` and `"Converter"`
- [x] 2.11 Update `ParsedResource` struct to include `Assets *Assets` and `Converter *Converter` fields
- [x] 2.12 Update `ParsedResource` methods: `Kind()`, `Name()`, `Identity()`, `Validate()`

## 3. Manifest Parser

- [x] 3.1 Add `case "Assets":` to `src/internal/manifest/parser.go` parseNode switch
- [x] 3.2 Add `case "Converter":` to parseNode switch
- [x] 3.3 Handle `spec.input.format` defaulting to `yaml` for Assets
- [x] 3.4 Handle `spec.input.format` and `spec.output.format` defaulting for Converter

## 4. Kustomize Adapter

- [x] 4.1 Add `case "Assets":` to `src/internal/kustomize/adapter.go` parseYamlNode switch
- [x] 4.2 Add `case "Converter":` to parseYamlNode switch
- [x] 4.3 Call validation for Assets and Converter after decode

## 5. Resolver

- [x] 5.1 Verify Assets and Converter are skipped in `src/internal/resolve/resolve.go` NewIndex
- [x] 5.2 Add logging when source kinds are skipped (optional, for debug)

## 6. Converter Package

- [x] 6.1 Create `src/internal/converter/types.go` with public types
- [x] 6.2 Create `src/internal/converter/engine.go` with yq-go integration
- [x] 6.3 Implement `Apply(converter, asset) (string, error)` function
- [x] 6.4 Implement format conversion (Asset format → Converter input format)
- [x] 6.5 Implement output encoding for all supported formats
- [x] 6.6 Implement `raw` output format (array join with newlines, string passthrough)
- [x] 6.7 Implement error handling for invalid expressions, format failures

## 7. Apply Command

- [x] 7.1 Add `--convert` and `--use` flags to `src/cmd/kfg/apply.go`
- [x] 7.2 Add validation: `--convert` and `--use` must be used together
- [x] 7.3 Add validation: `--convert`/`--use` incompatible with `-w`/`-c`
- [x] 7.4 Implement Asset lookup by `metadata.name`
- [x] 7.5 Implement Converter lookup by `metadata.name`
- [x] 7.6 Implement conversion pipeline execution
- [x] 7.7 Handle output to file (`-o`) or stdout
- [x] 7.8 Update help text to document both modes

## 8. Unit Tests

- [x] 8.1 Create `src/internal/converter/engine_test.go`
- [x] 8.2 Test YAML→JSON conversion
- [x] 8.3 Test JSON→YAML conversion
- [x] 8.4 Test YAML→TOML conversion
- [x] 8.5 Test raw output format (array join)
- [x] 8.6 Test raw output format (string passthrough)
- [x] 8.7 Test invalid expression error handling
- [x] 8.8 Test format mismatch conversion
- [x] 8.9 Test unsupported format validation
- [x] 8.10 Test Assets validation (missing name, missing data, invalid format)
- [x] 8.11 Test Converter validation (missing name, missing expression, invalid formats)

## 9. Integration Tests

- [x] 9.1 Create `tests/bats/assets.bats`
- [x] 9.2 Test `kfg apply -f manifest.yaml --convert <asset> --use <converter>` success
- [x] 9.3 Test `--convert` without `--use` error
- [x] 9.4 Test `--use` without `--convert` error
- [x] 9.5 Test `--convert`/`--use` with `-w` error
- [x] 9.6 Test `--convert`/`--use` with `-c` error
- [x] 9.7 Test Asset not found error
- [x] 9.8 Test Converter not found error
- [x] 9.9 Test output to file with `-o`
- [x] 9.10 Test raw output format end-to-end

## 10. Documentation

- [x] 10.1 Update `docs/AGENTS.md` to document Assets and Converter resource kinds
- [x] 10.2 Update `--help` text for `kfg apply` command
