# Internal Shell Generation Specification

## Purpose

This specification defines internal shell code generation including helper functions, artifact tracking, and cleanup.
## Requirements

### Requirement: Internal Helpers

Generated shell code MUST use namespaced internal helpers. The build result helper MUST be defined once at global scope. Generated shell code MUST also define a shared artifact tracking array for explicit cleanup within the same shell session.

#### Scenario: Global helper function definition
- GIVEN the generated shell includes build result YAML
- WHEN the shell code is examined
- THEN `__kfg_build_result()` is defined exactly once in global scope
- AND all Cmd wrappers can call `__kfg_build_result` without defining the helper themselves
- AND the helper outputs the contents of `$KFG_BUILD_RESULT_FILE`

#### Scenario: Artifact tracking helpers and variable
- GIVEN generated shell code with artifact-producing steps
- WHEN the shell code is examined
- THEN it declares a global bash array named `KFG_ARTIFACTS`
- AND it defines `__kfg_add_artifact()` in global scope
- AND each call to `__kfg_add_artifact()` appends the artifact path to `KFG_ARTIFACTS`
- AND cleanup steps and command wrappers can read `KFG_ARTIFACTS` later in the same shell session for explicit cleanup

### Requirement: Artifact Declaration in Workflow

Generated shell code MUST support artifact declaration at multiple levels:
- Step resources define `artifacts` with static paths
- Cmd resources define `artifacts` with static paths
- Workflow steps (before/after) can specify additional `artifacts` via `StepReference`

#### Scenario: Artifact tracking from step resources
- GIVEN a Step resource with `spec.artifacts` array
- WHEN the shell code is generated for that step
- THEN `__kfg_add_artifact()` is called for each artifact in the Step's artifacts array

#### Scenario: Artifact tracking from StepReference
- GIVEN a workflow step (before or after) with `artifacts` specified in the StepReference
- WHEN the shell code is generated
- THEN `__kfg_add_artifact()` is called for each artifact in the StepReference's artifacts array
- AND these artifacts are included in the command wrapper

#### Scenario: Artifact deduplication
- GIVEN multiple steps or workflow entries specify the same artifact path
- WHEN the shell code is generated
- THEN each artifact path appears exactly once in the `__kfg_add_artifact()` calls

#### Scenario: Cleanup with tracked artifacts
- GIVEN `KFG_ARTIFACTS` is populated with artifact paths
- WHEN `kfg.cleanup` step is executed
- THEN it iterates over `KFG_ARTIFACTS` and removes each path that exists
- AND cleanup succeeds without error when `KFG_ARTIFACTS` is empty
