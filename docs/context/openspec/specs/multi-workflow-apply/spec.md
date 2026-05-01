## ADDED Requirements

### Requirement: Multi-workflow resolution
The system SHALL support resolving multiple CmdWorkflows from a single kustomization for shell generation.

#### Scenario: Resolve all workflows
- **WHEN** `ResolveAllWorkflows()` is called with a valid index containing multiple workflows
- **THEN** all workflows in the index SHALL be resolved and returned in an array
- **AND** each resolved workflow SHALL be processed independently

#### Scenario: Resolve named workflows
- **WHEN** `ResolveWorkflowsByName(names)` is called with a list of workflow names
- **THEN** only workflows with matching names SHALL be resolved from the index
- **AND** workflows not found in the index SHALL cause a ResolutionError

#### Scenario: Multi-workflow generation
- **WHEN** multiple resolved workflows are provided to the generator
- **THEN** a single shell output SHALL be produced containing all commands from all workflows
- **AND** shared step functions SHALL be deduplicated (generated once)

### Requirement: Comma-separated workflow flag
The apply command SHALL accept comma-separated workflow names in the workflow flag.

#### Scenario: Parse comma-separated workflows
- **WHEN** user runs `kfg apply -k . -w workflow1,workflow2,workflow3`
- **THEN** the workflow names SHALL be parsed from the single flag value
- **AND** whitespace around commas SHALL be trimmed from each name

#### Scenario: Invalid workflow name
- **WHEN** a specified workflow name does not exist in the manifests
- **THEN** the system SHALL return a ResolutionError with available workflow names

### Requirement: Unified multi-workflow generation
The generator SHALL produce a single shell output from multiple workflows.

#### Scenario: Deduplicate step functions
- **WHEN** multiple workflows reference the same step by name
- **THEN** the step function SHALL be generated only once in the output
- **AND** all command wrappers referencing that step SHALL use the same function

#### Scenario: Combine all commands
- **WHEN** multiple workflows contain commands
- **THEN** all commands from all workflows SHALL be included in the single shell output
- **AND** each command SHALL generate its own wrapper function

#### Scenario: Single metadata context
- **WHEN** generating from multiple workflows
- **THEN** only one header SHALL be generated for the entire output
- **AND** `KFG_WORKFLOW_NAME` variable SHALL NOT be included (meaningless with multiple workflows)
- **AND** `KFG_KUSTOMIZATION_NAME` SHALL identify the kustomization origin