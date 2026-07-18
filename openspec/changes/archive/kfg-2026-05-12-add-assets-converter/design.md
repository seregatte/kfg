## Context

kfg currently supports three resource kinds: Step, Cmd, and CmdWorkflow. These are execution-oriented—they generate shell code that can be sourced or executed. The pipeline is: YAML manifests → Load → Kustomize Merge → Validate → Generate Shell.

There is no mechanism for declarative data payloads or format transformation. Users who need to convert structured data between formats (YAML→JSON, TOML→YAML, etc.) or generate text artifacts (markdown with frontmatter, config files) must use external tools.

This change introduces two new resource kinds—Assets and Converter—that operate in a separate "source layer" alongside the existing execution kinds. Assets declare data payloads, Converters declare transformations using yq-go expressions, and the `apply` command gains a conversion mode to process them.

## Goals / Non-Goals

**Goals:**
- Enable declarative data transformation within kfg manifests
- Support all yq-go input/output formats (yaml, json, xml, props, csv, tsv, toml, hcl, lua, ini, shell, base64, uri, kyaml)
- Provide `raw` output format for plain text generation (markdown with frontmatter, etc.)
- Integrate cleanly with existing manifest pipeline (kustomize loading, validation)
- Use yq-go as a Go module dependency (no external CLI required)
- Extend `apply` command with `--convert` and `--use` flags

**Non-Goals:**
- Schema validation of Asset data (no `kind: Schema` in this change)
- Template-based conversion engine (only yq-go expressions)
- Multiple asset/converter pairs in a single invocation (one `--convert` + one `--use`)
- Modifying or replacing existing shell generation functionality
- Bidirectional conversion (one-way: Asset → Converter → output)

## Decisions

### Decision 1: Assets and Converter as source kinds (not execution kinds)

**Choice**: Assets and Converter are explicitly skipped during resolution and shell generation. They exist only for data transformation.

**Rationale**: Mixing data transformation with shell generation would complicate the execution model. Keeping them separate means:
- Resolver ignores Assets/Converter (no index entries)
- Shell generator never sees Assets/Converter
- `apply` command has two distinct modes with mutually exclusive flags

**Alternatives considered**:
- *Unified pipeline*: Process Assets/Converter through the same pipeline as Cmds. Rejected because it conflates data transformation with shell execution.
- *Separate command*: `kfg convert` subcommand. Rejected by user preference—integrated into `apply` with flags.

### Decision 2: yq-go as Go module (not external CLI)

**Choice**: Use `github.com/mikefarah/yq/v4/pkg/yqlib` as a Go dependency, evaluated in-process.

**Rationale**:
- No external binary dependency (self-contained kfg binary)
- Access to all yq features via Go API
- `StringEvaluator` provides simple string-in/string-out interface
- Format encoders/decoders available via `FormatFromString()` + `GetConfiguredEncoder()`

**Alternatives considered**:
- *Shell out to `yq` binary*: Simpler but requires external dependency. Rejected.
- *Custom expression language*: More control but massive scope increase. Rejected.

### Decision 3: `spec.input.format` on both Assets and Converter

**Choice**: Both Assets and Converter use `spec.input.format` to declare their data format.

**Rationale**:
- Consistent naming convention across resource kinds
- Default is `yaml` (can be omitted for YAML data)
- Asset declares what format its data is in
- Converter declares what format it expects as input
- If formats differ, engine converts Asset data to Converter's expected format before evaluation

**Alternatives considered**:
- *`metadata.format`*: More visible but semantically incorrect—format describes the data, not the resource identity.
- *`spec.dataFormat`*: Redundant with `spec.input.format` pattern on Converter.

### Decision 4: `spec.data` is map for YAML, string for other formats

**Choice**: When `spec.input.format` is `yaml` (default), `spec.data` is a YAML map parsed by kustomize. For all other formats, `spec.data` is a string containing data in that format.

**Rationale**:
- YAML is the native format of kustomize—maps are parsed natively
- Other formats (JSON, TOML, XML, etc.) must be represented as strings since kustomize cannot parse them
- This avoids forcing YAML-string wrapping for the common case (YAML data)

**Alternatives considered**:
- *Always string*: Consistent but verbose for YAML (requires `|` block scalar).
- *Always map*: Cannot represent non-YAML formats natively.

### Decision 5: `raw` as output format

**Choice**: Add `raw` as a special output format that returns the yq result as a plain string without YAML/JSON encoding.

**Rationale**:
- Enables markdown with frontmatter generation
- yq expressions can build arrays of strings and join them
- Without `raw`, array results get YAML-encoded (e.g., `- "line1"\n- "line2"`)
- `raw` joins array elements with newlines, outputs string as-is

**Implementation**:
- If `output.format == "raw"` and result is array of strings → join with `\n`
- If `output.format == "raw"` and result is string → output directly
- No YAML/JSON encoding applied

### Decision 6: Mutual exclusivity of shell/conversion modes

**Choice**: Shell generation flags (`-w`, `-c`) and conversion flags (`--convert`, `--use`) cannot be mixed. Validation produces clear error messages.

**Rationale**:
- Two distinct modes with different output semantics
- Mixing would create ambiguous behavior
- Clear error messages guide users to correct usage

**Validation rules**:
- `--convert` and `--use` must be used together (both or neither)
- `-w`/`-c` cannot be used with `--convert`/`--use`
- `--convert`/`--use` cannot be used with `-w`/`-c`

## Risks / Trade-offs

**[Risk] yq-go dependency size** → yq-go pulls in transitive dependencies. Mitigated by Go module vendoring and the fact that kfg already has significant dependencies (kustomize, cobra, viper).

**[Risk] yq expression complexity** → Users must learn yq syntax for transformations. Mitigated by yq's extensive documentation and jq-like syntax familiarity.

**[Risk] Format conversion fidelity** → Converting between formats (e.g., YAML→TOML) may lose type information. Mitigated by yq's native format handling and documented limitations.

**[Risk] No schema validation** → Asset data is not validated against schemas. Mitigated by deferring schema support to a future change—MVP focuses on transformation.

**[Trade-off] One asset per invocation** → Simplifies CLI but limits batch processing. Acceptable for MVP; can be extended later with `--convert-all` or multiple `--convert` flags.

**[Trade-off] No `raw` input format** → Plain text cannot be processed with yq path expressions. Acceptable since yq is designed for structured data; string operations are limited.

## Migration Plan

No migration needed—this is a purely additive change. Existing `apply` behavior is unchanged when `--convert`/`--use` flags are not provided.

## Open Questions

None—all design decisions have been resolved during the planning conversation.
