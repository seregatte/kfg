## MODIFIED Requirements

### Requirement: ctx7 install step defaults

The step `kfg.extension.ctx7.install` SHALL have empty defaults for `AGENT_HOME` and `OUTPUT_DIR`, requiring explicit values at invocation time. The step MUST validate that `INSTALL_CMD`, `INSTALL_FLAGS`, and `AGENT_HOME` are set before execution.

#### Scenario: Install with explicit agent paths
- **WHEN** `INSTALL_CMD` is `ctx7 setup --cli --project`, `INSTALL_FLAGS` is `--yes`, `AGENT_HOME` is `.opencode`, `OUTPUT_DIR` is `.opencode/skills/`
- **THEN** the step runs `ctx7 setup --cli --project --yes` and copies skills from `.opencode/skills/` to `OUTPUT_DIR`

#### Scenario: Missing required env var
- **WHEN** `AGENT_HOME` is empty or unset
- **THEN** the step logs an error and exits with code 1

#### Scenario: Skills directory does not exist
- **WHEN** `$AGENT_HOME/skills/` does not exist after install
- **THEN** the step logs a warning but exits with success (install may have created no skills)

#### Scenario: Install command fails
- **WHEN** `ctx7 setup --cli --project --yes` returns non-zero exit code
- **THEN** the step propagates the failure
