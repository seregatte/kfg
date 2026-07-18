# Build Result Management Specification

## Purpose

Specifies how kfg generates and manages the build result YAML file in shell code, ensuring global-scope accessibility, single-file setup with no duplication across command wrappers. The build result temp file is NOT automatically cleaned up — cleanup is the responsibility of the user or explicit step configuration.

## Requirements

### Requirement: Build result file created at global scope

The generated shell code SHALL create the `KFG_BUILD_RESULT_FILE` temp file once in global scope (outside all Cmd wrapper functions) when build result YAML is present. The file SHALL be accessible to all Cmd wrappers without re-creation or duplication.

#### Scenario: Single build result file in workflow with multiple commands
- **WHEN** generating shell code for a workflow with 3+ Cmd wrappers
- **THEN** exactly one `mktemp -t kfg-build-XXXXXX.yaml` call appears in global scope
- **AND** the base64-encoded build result is decoded exactly once into that file
- **AND** `KFG_BUILD_RESULT_FILE` is exported as a global variable

### Requirement: Global helper function for build result access

The generated shell code SHALL define `__kfg_build_result()` as a global function that outputs the contents of `KFG_BUILD_RESULT_FILE`. All Cmd wrappers and Steps SHALL use this helper to access the build result.

#### Scenario: Step accessing build result
- **WHEN** a Step needs to read the build result YAML
- **THEN** the Step can call `__kfg_build_result | kfg assets convert --use <converter>`
- **AND** the command accesses the same shared global file

### Requirement: No automatic cleanup of build result file

The generated shell code SHALL NOT register any EXIT trap or automatic cleanup mechanism for the build result temp file. The file persists after shell exit to allow consecutive runs without regeneration. Cleanup is the responsibility of the user or explicit workflow step configuration.

#### Scenario: Temp file persists after shell exit
- **WHEN** the generated shell is sourced interactively and the user exits the shell
- **THEN** the `$KFG_BUILD_RESULT_FILE` remains on disk
- **AND** consecutive `kfg apply/run` invocations can reuse the same file path

#### Scenario: Explicit cleanup via workflow step
- **WHEN** a workflow includes a cleanup step that removes `$KFG_BUILD_RESULT_FILE`
- **THEN** the file is removed at the end of workflow execution
- **AND** cleanup is controlled by the workflow author, not automatic

### Requirement: No per-Cmd build result duplication

The generated Cmd wrapper functions SHALL NOT contain mktemp, base64 decode, helper function, or cleanup trap code for build result. All build result logic SHALL be at global scope only.

#### Scenario: Cmd wrapper inspection
- **WHEN** examining a generated Cmd wrapper function
- **THEN** no `local __kfg_build_result_file` variable appears inside the function
- **AND** no `mktemp -t kfg-build` call appears inside the function
- **AND** no `trap ... RETURN` for build result appears inside the function

### Requirement: Shell integration for build result
The shell code generation system SHALL emit build result setup in global scope between metadata environment variables and helper functions, ensuring all subsequent Steps and Cmds can access `KFG_BUILD_RESULT_FILE`.

#### Scenario: Code structure ordering
- **WHEN** generating shell code via `GenerateKustomization()`
- **THEN** the global structure is: header → metadata env → **build result setup** → helpers → step functions → cmd wrappers
- **AND** build result setup is emitted only if `Generator.buildResultYAML` is non-empty

#### Scenario: Backward compatibility with no build result
- **WHEN** generating shell code with no build result YAML (e.g., for simple manifest without assets)
- **THEN** no build result setup code is emitted
- **AND** no `KFG_BUILD_RESULT_FILE` variable exists
- **AND** Cmds and Steps that don't use build result remain unaffected
