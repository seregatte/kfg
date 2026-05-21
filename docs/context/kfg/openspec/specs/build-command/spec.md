# Build Command Specification

## Purpose

Define the `kfg build` command for building kustomizations and outputting YAML.
## Requirements
### Requirement: Build command syntax

The CLI MUST provide `kfg build [path-or-url]` for building kustomizations with optional argument.

#### Scenario: Basic build
- **WHEN** user runs `kfg build .kfg/overlay/dev`
- **THEN** system loads kustomization.yaml from path
- **AND** processes HTTP resources, overlays
- **AND** outputs resulting YAML to stdout

#### Scenario: Build with output file
- **WHEN** user runs `kfg build .kfg/overlay/dev -o output.yaml`
- **THEN** writes YAML to `output.yaml`

#### Scenario: Build remote kustomization
- **WHEN** user runs `kfg build https://example.com/kustomization.yaml`
- **THEN** loads from URL
- **AND** outputs resulting YAML

#### Scenario: Build GitHub repository
- **WHEN** user runs `kfg build https://github.com/owner/repo//manifests`
- **THEN** clones GitHub repository
- **AND** processes kustomization.yaml
- **AND** outputs resulting YAML

#### Scenario: Build GitHub repository with ref
- **WHEN** user runs `kfg build https://github.com/owner/repo//manifests?ref=v1.0.0`
- **THEN** clones specified tag
- **AND** processes kustomization.yaml
- **AND** outputs resulting YAML

#### Scenario: Build without argument with KFG_KPATH
- **WHEN** user runs `kfg build` (no argument)
- **AND** `KFG_KPATH=./manifests` is set
- **THEN** uses `KFG_KPATH` as source
- **AND** outputs resulting YAML

#### Scenario: Build without argument or KFG_KPATH
- **WHEN** user runs `kfg build` (no argument)
- **AND** `KFG_KPATH` is not set
- **THEN** exit code 1
- **AND** error message: "kustomization source required. Provide a path, use -k flag, or set KFG_KPATH."

### Requirement: Kustomize alias

The CLI MUST provide `kfg kustomize` as alias.

#### Scenario: Kustomize alias
- **WHEN** user runs `kfg kustomize .kfg/overlay/dev`
- **THEN** behaves identically to `kfg build`

### Requirement: Build flags

The CLI MUST support specific flags.

#### Scenario: Output flag
- **WHEN** user runs `kfg build -o output.yaml`
- **THEN** writes to specified file instead of stdout

#### Scenario: Short flag
- **WHEN** user runs `kfg build -o output.yaml`
- **THEN** short flag `-o` works same as `--output`

### Requirement: Kustomization processing

The build MUST process kustomize features.

#### Scenario: Strategic merge patches
- **WHEN** kustomization includes patches
- **THEN** patches applied to base resources

#### Scenario: Resource generators
- **WHEN** kustomization includes generators
- **THEN** generators produce resources

#### Scenario: Overlays
- **WHEN** kustomization references overlays
- **THEN** overlay resources merged

### Requirement: Output format

The output MUST be valid YAML.

#### Scenario: YAML output
- **WHEN** build succeeds
- **THEN** output is valid YAML
- **AND** contains all processed resources

#### Scenario: Multi-document
- **WHEN** multiple resources processed
- **THEN** output uses `---` separators

### Requirement: Exit codes

The CLI MUST use consistent exit codes.

#### Scenario: Success
- **WHEN** build succeeds
- **THEN** exit code 0

#### Scenario: Path not found
- **WHEN** specified path doesn't exist
- **THEN** exit code 1
- **AND** error message indicates path issue

#### Scenario: Invalid YAML
- **WHEN** kustomization contains invalid YAML
- **THEN** exit code 1
- **AND** error message indicates syntax issue

#### Scenario: GitHub clone failure
- **WHEN** GitHub URL fails to clone
- **THEN** exit code 1
- **AND** error message indicates clone issue

#### Scenario: No source provided
- **WHEN** no argument and no `KFG_KPATH`
- **THEN** exit code 1
- **AND** error message indicates source required

