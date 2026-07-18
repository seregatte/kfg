## MODIFIED Requirements

### Requirement: Run flags

The run command MUST support the same input flags as apply and expose cache refresh control.

#### Scenario: Kustomize path
- **WHEN** user runs `kfg run -k .kfg/overlay/dev claude`
- **THEN** short flag `-k` works same as `--kustomize`

#### Scenario: Manifest file
- **WHEN** user runs `kfg run -f manifest.yaml claude`
- **THEN** short flag `-f` works same as `--file`

#### Scenario: Workflow selection
- **WHEN** user runs `kfg run -w dev claude`
- **THEN** short flag `-w` works same as `--workflow`

#### Scenario: Refresh flag
- **WHEN** user runs `kfg run -k path claude --refresh`
- **THEN** the runtime SHALL set refresh state for the generated shell execution
- **AND** cacheable Steps SHALL bypass existing cache entries during that run

### Requirement: Process execution

The run command MUST execute the agent as a child process.

#### Scenario: Stream inheritance
- **WHEN** agent is executed
- **THEN** stdin, stdout, and stderr are inherited from the parent process

#### Scenario: Exit code propagation
- **WHEN** agent exits with code N
- **THEN** `kfg run` exits with code N

#### Scenario: Cached Step reuse
- **WHEN** a generated agent command reaches a cacheable Step with a valid cache entry
- **THEN** the runtime SHALL restore that Step's artifacts and outputs instead of re-running the Step
