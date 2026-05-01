# Shell Code Generation Specification

## Purpose

Specifies how kfg generates valid bash code from manifests, with emphasis on global-scope build result setup and shared helper functions accessible to all command wrappers and steps.

## Requirements

### Requirement: Shell Code Generation

The CLI MUST generate valid bash code from manifests. The generated shell structure MUST place build result setup at global scope (outside all Cmd wrapper functions) when build result YAML is provided, ensuring all Cmd wrappers and steps can access `KFG_BUILD_RESULT_FILE` without duplication.

#### Scenario: Function structure with global build result
- GIVEN a generated shell with multiple Cmd wrappers and build result YAML
- WHEN the shell code is examined
- THEN exactly one global build result setup appears (mktemp, base64 decode, helper function)
- AND each Cmd wrapper executes `before` steps first (now without internal build result setup)
- AND each Cmd wrapper executes the Cmd's `run` body
- AND each Cmd wrapper executes `after` steps last
- AND the function preserves `"$@"` passed to the function

#### Scenario: Global scope build result access
- GIVEN a Step that uses `__kfg_build_result()`
- WHEN the step executes
- THEN the step accesses the shared global `__kfg_build_result()` helper function
- AND the function outputs the same build result YAML content to all steps and commands

### Requirement: Internal Helpers

Generated shell code MUST use namespaced internal helpers. The build result helper MUST be defined once at global scope.

#### Scenario: Global helper function definition
- GIVEN the generated shell includes build result YAML
- WHEN the shell code is examined
- THEN `__kfg_build_result()` is defined exactly once in global scope
- AND all Cmd wrappers can call `__kfg_build_result` without defining the helper themselves
- AND the helper outputs the contents of `$KFG_BUILD_RESULT_FILE`

### Requirement: Command Execution Flow

Generated functions MUST execute steps and the command in the correct order. Build result initialization MUST occur before any Cmd wrappers are defined, ensuring all steps can access it.

#### Scenario: Build result initialization ordering
- GIVEN the generated shell code structure
- WHEN examined in order
- THEN: header → metadata env → **global build result setup** → helpers → step functions → cmd wrappers
- AND build result setup (mktemp, base64 decode, export, helper function) appears exactly once before any Cmd functions
- AND no per-Cmd build result setup code appears inside function bodies
