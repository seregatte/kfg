# Framework Reusable Steps Specification

## Purpose

Framework steps are reusable manifest-level primitives exported by the framework package. This specification documents the contract and behavior of the core framework steps.

## Requirements

### Requirement: Step exports

The framework MUST export the following reusable steps:

- `kfg.materialize` - Generate shell code from manifests using converters
- `kfg.materialize-scaffold` - Generate scaffolding from templates
- `kfg.cleanup` - Clean up generated artifacts
- `kfg.ensure-gitignore` - Manage repository `.gitignore` entries
- `kfg.copy-context` - Copy context files into artifacts

### Requirement: Materialize step

The materialize step MUST generate shell code from manifests using a specified converter.

#### Scenario: Per-item mode
- **WHEN** MODE is set to `per-item`
- **THEN** the step SHALL process each asset in ASSETS list individually
- **AND** each asset SHALL be converted using the specified CONVERTER
- **AND** results SHALL be written to corresponding paths in OUTPUTS list
- **AND** the counts of ASSETS and OUTPUTS MUST match

#### Scenario: Aggregate mode
- **WHEN** MODE is set to `aggregate`
- **THEN** the step SHALL process all assets in ASSETS list
- **AND** assets SHALL be converted and merged using the specified CONVERTER
- **AND** the merged result SHALL be written to the first path in OUTPUTS list
- **AND** WRAP_KEY MAY be used to wrap the merged result under a specified key

#### Scenario: Materialize usage
- **WHEN** a workflow needs to generate shell code
- **THEN** it SHALL use the materialize step with appropriate MODE, ASSETS, CONVERTER, and OUTPUTS
- **AND** the step SHALL produce artifacts registered via `__kfg_add_artifact()`

### Requirement: Cleanup step

The cleanup step MUST remove generated artifacts registered in the current session.

#### Scenario: Artifact cleanup
- **WHEN** cleanup is invoked
- **THEN** it SHALL iterate through all artifacts in `KFG_ARTIFACTS`
- **AND** each artifact file or directory SHALL be removed
- **AND** if `KFG_ARTIFACTS` is empty, the step SHALL exit gracefully

#### Scenario: Cleanup usage
- **WHEN** a workflow needs to remove generated artifacts
- **THEN** it MAY use the cleanup step as an after-step
- **AND** cleanup MUST run even if earlier steps failed

### Requirement: Ensure-gitignore step

The ensure-gitignore step MUST manage `.gitignore` entries for generated artifacts.

#### Scenario: Adding gitignore entries
- **WHEN** GITIGNORE_ENTRIES is set to a colon-separated list
- **THEN** the step SHALL add each entry to the GITIGNORE_FILE
- **AND** entries MUST NOT be duplicated if already present
- **AND** the GITIGNORE_FILE SHALL be created if it does not exist

#### Scenario: Ensure-gitignore usage
- **WHEN** a workflow generates artifacts that should be ignored
- **THEN** it MAY use the ensure-gitignore step to register gitignore entries
- **AND** entries SHOULD follow the format `/<path>/` for directories or `<file>` for files

### Requirement: Copy-context step

The copy-context step MUST copy context files into generated artifacts.

#### Scenario: Context file copying
- **WHEN** SOURCE_FILES is set to a colon-separated list
- **AND** DEST_PATH is set to a target directory
- **THEN** the step SHALL copy each source file to the destination
- **AND** directory structure SHALL be preserved where possible
- **AND** the step SHALL fail if SOURCE_FILES is invalid

#### Scenario: Copy-context usage
- **WHEN** a workflow generates artifacts that need supporting context files
- **THEN** it MAY use the copy-context step to populate artifacts
- **AND** context files might include templates, configs, or documentation

### Requirement: Materialize-scaffold step

The materialize-scaffold step MUST generate scaffolding structures from templates.

#### Scenario: Scaffold generation
- **WHEN** SCAFFOLD_TEMPLATE is set to a template asset
- **AND** SCAFFOLD_OUTPUT is set to a target path
- **THEN** the step SHALL generate scaffolding from the template
- **AND** template variables MAY be substituted if TEMPLATE_VARS is provided
- **AND** generated scaffolding SHALL be registered as artifacts

#### Scenario: Materialize-scaffold usage
- **WHEN** a workflow needs to create scaffolding structures
- **THEN** it MAY use the materialize-scaffold step
- **AND** this is typically used for new project or feature setup

### Requirement: Step execution contract

All framework steps MUST follow the shell runtime API contract.

#### Scenario: Artifact registration
- **WHEN** a framework step generates artifacts
- **THEN** it SHALL register each artifact via `__kfg_add_artifact <path>`
- **AND** artifacts SHALL be discoverable via the `KFG_ARTIFACTS` variable

#### Scenario: Logging
- **WHEN** a framework step executes
- **THEN** it SHALL use the structured logging API (`__kfg_log_*`)
- **AND** Step-originated log events SHALL rely on runtime-provided `step_name` attribution instead of encoding Step identity inside the component string
- **AND** it MUST NOT use unstructured echo or printf

#### Scenario: Conditional execution
- **WHEN** a framework step has prerequisites or conditions
- **THEN** it SHALL use the `when` condition mechanism
- **AND** it MUST NOT implement its own conditional logic

### Requirement: Backward compatibility

Framework steps MUST maintain API stability.

#### Scenario: Step name stability
- **WHEN** a framework step is exported
- **THEN** its metadata.name MUST remain constant
- **AND** the step SHALL not be renamed without a deprecation period

#### Scenario: Step behavior stability
- **WHEN** a framework step is used
- **THEN** its behavior as documented MUST remain consistent
- **AND** breaking changes require advance notice and migration documentation

#### Scenario: New step addition
- **WHEN** new framework steps are added
- **THEN** they SHALL be spec'd in this document
- **AND** they SHALL be exported through the framework kustomization
- **AND** they SHALL follow the same execution contract as existing steps
