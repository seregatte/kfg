## MODIFIED Requirements

### Requirement: After phase final cleanup

The workflow SHALL include a final cleanup step in the `after` phase that removes tracked artifact paths recorded in `KFG_ARTIFACTS` during workflow execution.

#### Scenario: After cleanup runs regardless of agent
- **WHEN** any agent command completes
- **THEN** after-cleanup removes artifact files and directories recorded in `KFG_ARTIFACTS`

#### Scenario: After cleanup with no tracked artifacts
- **WHEN** any agent command completes and `KFG_ARTIFACTS` is empty
- **THEN** after-cleanup succeeds without removing unrelated files

### Requirement: Step reference artifacts

Steps referenced in workflow entries (before/after) MAY specify additional artifact paths via the `artifacts` field in the StepReference.

#### Scenario: Step reference artifacts are tracked
- **GIVEN** a workflow step with `artifacts` specified in the StepReference
- **WHEN** the step is executed as part of a command
- **THEN** each artifact path is added to `KFG_ARTIFACTS` via `__kfg_add_artifact()`
- **AND** the artifact paths are available for cleanup during the same command execution

#### Scenario: Cleanup with step reference artifacts
- **GIVEN** step reference artifacts are tracked during command execution
- **WHEN** `kfg.cleanup` runs (either as a step or in after phase)
- **THEN** all tracked artifacts (including step reference artifacts) are removed
