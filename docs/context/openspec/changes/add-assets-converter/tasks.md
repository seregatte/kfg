## 1. Dependencies and Setup

- [ ] 1.1 Add `github.com/mikefarah/yq/v4` to go.mod
- [ ] 1.2 Run `go mod tidy` to resolve dependencies

## 2. Manifest Types

- [ ] 2.1 Add `Assets` struct to `src/internal/manifest/types.go` with `APIVersion`, `Kind`, `Metadata`, `Spec` fields
- [ ] 2.2 Add `AssetsSpec` struct with `Input` and `Data` fields
- [ ] 2.3 Add `InputSpec` struct with `Format` field
- [ ] 2.4 Add `Converter` struct to `src/internal/manifest/types.go` with `APIVersion`, `Kind`, `Metadata`, `Spec` fields
- [ ] 2.5 Add `ConverterSpec` struct with `Input`, `Engine`, `Output` fields
- [ ] 2.6 Add `EngineSpec` struct with `Expression` field
- [ ] 2.7 Add `OutputSpec` struct with `Format` field
- [ ] 2.8 Add `SupportedInputFormats` and `SupportedOutputFormats` constants
- [ ] 2.9 Add validation methods: `Assets.Validate()`, `Converter.Validate()`
- [ ] 2.10 Update `SupportedKinds` to include `"Assets"` and `"Converter"`
- [ ] 2.11 Update `ParsedResource` struct to include `Assets *Assets` and `Converter *Converter` fields
- [ ] 2.12 Update `ParsedResource` methods: `Kind()`, `Name()`, `Identity()`, `Validate()`

## 3. Manifest Parser

- [ ] 3.1 Add `case "Assets":` to `src/internal/manifest/parser.go` parseNode switch
- [ ] 3.2 Add `case "Converter":` to parseNode switch
- [ ] 3.3 Handle `spec.input.format` defaulting to `yaml` for Assets
- [ ] 3.4 Handle `spec.input.format` and `spec.output.format` defaulting for Converter

## 4. Kustomize Adapter

- [ ] 4.1 Add `case "Assets":` to `src/internal/kustomize/adapter.go` parseYamlNode switch
- [ ] 4.2 Add `case "Converter":` to parseYamlNode switch
- [ ] 4.3 Call validation for Assets and Converter after decode

## 5. Resolver

- [ ] 5.1 Verify Assets and Converter are skipped in `src/internal/resolve/resolve.go` NewIndex
- [ ] 5.2 Add logging when source kinds are skipped (optional, for debug)

## 6. Converter Package

- [ ] 6.1 Create `src/internal/converter/types.go` with public types
- [ ] 6.2 Create `src/internal/converter/engine.go` with yq-go integration
- [ ] 6.3 Implement `Apply(converter, asset) (string, error)` function
- [ ] 6.4 Implement format conversion (Asset format → Converter input format)
- [ ] 6.5 Implement output encoding for all supported formats
- [ ] 6.6 Implement `raw` output format (array join with newlines, string passthrough)
- [ ] 6.7 Implement error handling for invalid expressions, format failures

## 7. Apply Command

- [ ] 7.1 Add `--convert` and `--use` flags to `src/cmd/kfg/apply.go`
- [ ] 7.2 Add validation: `--convert` and `--use` must be used together
- [ ] 7.3 Add validation: `--convert`/`--use` incompatible with `-w`/`-c`
- [ ] 7.4 Implement Asset lookup by `metadata.name`
- [ ] 7.5 Implement Converter lookup by `metadata.name`
- [ ] 7.6 Implement conversion pipeline execution
- [ ] 7.7 Handle output to file (`-o`) or stdout
- [ ] 7.8 Update help text to document both modes

## 8. Unit Tests

- [ ] 8.1 Create `src/internal/converter/engine_test.go`
- [ ] 8.2 Test YAML→JSON conversion
- [ ] 8.3 Test JSON→YAML conversion
- [ ] 8.4 Test YAML→TOML conversion
- [ ] 8.5 Test raw output format (array join)
- [ ] 8.6 Test raw output format (string passthrough)
- [ ] 8.7 Test invalid expression error handling
- [ ] 8.8 Test format mismatch conversion
- [ ] 8.9 Test unsupported format validation
- [ ] 8.10 Test Assets validation (missing name, missing data, invalid format)
- [ ] 8.11 Test Converter validation (missing name, missing expression, invalid formats)

## 9. Integration Tests

- [ ] 9.1 Create `tests/bats/assets.bats`
- [ ] 9.2 Test `kfg apply -f manifest.yaml --convert <asset> --use <converter>` success
- [ ] 9.3 Test `--convert` without `--use` error
- [ ] 9.4 Test `--use` without `--convert` error
- [ ] 9.5 Test `--convert`/`--use` with `-w` error
- [ ] 9.6 Test `--convert`/`--use` with `-c` error
- [ ] 9.7 Test Asset not found error
- [ ] 9.8 Test Converter not found error
- [ ] 9.9 Test output to file with `-o`
- [ ] 9.10 Test raw output format end-to-end

## 10. Documentation

- [ ] 10.1 Update `docs/AGENTS.md` to document Assets and Converter resource kinds
- [ ] 10.2 Update `--help` text for `kfg apply` command
