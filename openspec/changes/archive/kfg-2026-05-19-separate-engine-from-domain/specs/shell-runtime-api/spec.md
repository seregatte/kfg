## ADDED Requirements

### Requirement: Generated shell runtime contract
The engine SHALL generate a stable shell runtime API that framework steps can consume without compile-time coupling to specific manifests.

#### Scenario: Runtime helper availability
- **WHEN** the engine generates shell code from a `CmdWorkflow`
- **THEN** the generated shell SHALL expose helper functions and metadata required by framework steps

### Requirement: Artifact tracking API
The generated shell SHALL expose artifact tracking primitives for framework and domain steps.

#### Scenario: Registering artifacts
- **WHEN** a step calls `__kfg_add_artifact()` with a path
- **THEN** the path SHALL be appended to `KFG_ARTIFACTS`
- **AND** the tracked path SHALL remain available for later cleanup during the same invocation

### Requirement: Shell logging API
The generated shell SHALL expose logging helpers that route step logs through the engine logging command.

#### Scenario: Emitting structured shell logs
- **WHEN** a step calls `_kfg.log.info()` or another `_kfg.log.*()` helper
- **THEN** the helper SHALL invoke `kfg sys log` with the requested level and message payload

### Requirement: Build result API
The generated shell SHALL expose build result helpers when build result YAML is available.

#### Scenario: Reading build output
- **WHEN** the engine embeds build result YAML in the generated shell
- **THEN** it SHALL export `KFG_BUILD_RESULT_FILE`
- **AND** it SHALL provide `__kfg_build_result()` to read the current build result

### Requirement: Conditional execution API
The generated shell SHALL expose helper functions for evaluating `when` clauses declared in manifests.

#### Scenario: Evaluating when conditions
- **WHEN** a workflow step has a `when` clause
- **THEN** the generated shell SHALL provide `__kfg_when_*()` helpers needed to evaluate that clause during execution

### Requirement: Runtime metadata API
The generated shell SHALL expose metadata environment variables for the current invocation.

#### Scenario: Accessing runtime metadata
- **WHEN** a generated command runs
- **THEN** it SHALL export `KFG_SESSION_ID`, `KFG_KUSTOMIZATION_NAME`, `KFG_WORKFLOW_NAME`, and `KFG_SHELL`
