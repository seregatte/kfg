# Build Result Management Specification

## Purpose

Specifies how kfg generates and manages the build result YAML file in shell code, ensuring global-scope accessibility, single-file setup with no duplication across command wrappers, and proper cleanup via EXIT trap.

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

### Requirement: Cleanup via global EXIT trap

The generated shell code SHALL register a global `EXIT` trap that removes the build result temp file when the shell exits (sourced shell session ends, or non-interactive shell completes).

#### Scenario: Interactive shell session cleanup
- **WHEN** the generated shell is sourced interactively and the user exits the shell
- **THEN** the `EXIT` trap fires and removes `$KFG_BUILD_RESULT_FILE`
- **AND** no dangling temp files remain

#### Scenario: Non-interactive shell execution cleanup
- **WHEN** generated shell code runs in a subshell or non-interactive context
- **THEN** the `EXIT` trap fires when the shell exits
- **AND** cleanup occurs identically to interactive mode

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
