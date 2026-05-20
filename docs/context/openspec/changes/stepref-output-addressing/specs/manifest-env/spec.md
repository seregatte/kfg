## MODIFIED Requirements

### Requirement: Step Reference Environment Override

Step references MUST support `env` overrides in workflows.

#### Scenario: Reference env overrides step default
- **GIVEN** a Step defines `spec.env.FOO=default`
- **AND** workflow reference defines `env.FOO=override`
- **WHEN** Step executes through that reference
- **THEN** `FOO` equals `override`

#### Scenario: Reference env adds new variables
- **GIVEN** a Step defines `spec.env.FOO=default`
- **AND** workflow reference defines `env.BAR=new`
- **WHEN** Step executes through that reference
- **THEN** both `FOO=default` and `BAR=new` available

#### Scenario: Same step with different env values
- **GIVEN** same Step referenced twice
- **AND** first reference defines `env.VALUE=a`
- **AND** second reference defines `env.VALUE=b`
- **WHEN** each reference executes
- **THEN** first invocation sees `VALUE=a`
- **AND** second invocation sees `VALUE=b`
- **AND** env values do not leak between invocations

#### Scenario: Reference env reads prior step output
- **GIVEN** a workflow step reference named `agents.ctx7.install.claude` produces a Step output
- **AND** a later workflow step reference defines `env.CTX7_CONTEXT="$kfg.output(agents.ctx7.install.claude)"`
- **WHEN** the later workflow step executes in the same command invocation
- **THEN** `CTX7_CONTEXT` SHALL equal the referenced step output value

#### Scenario: Reference env output lookup requires named step reference
- **GIVEN** a workflow step reference env value uses `$kfg.output(<value>)`
- **WHEN** the workflow is validated
- **THEN** `<value>` SHALL match an existing `StepReference.name` in the same CmdWorkflow
- **AND** the referenced workflow step SHALL point to a Step resource with `spec.output`

### Requirement: Env Scope

Environment MUST be scoped appropriately.

#### Scenario: Step env scope (step-scoped)
- **GIVEN** Step defines `spec.env`
- **WHEN** step executes
- **THEN** env exported inside step function
- **AND** not visible to other steps

#### Scenario: Cmd env scope (wrapper-scoped)
- **GIVEN** Cmd defines `spec.env`
- **WHEN** wrapper executes
- **THEN** env visible to all phases (before, run, after)
- **AND** isolated per invocation

#### Scenario: Step reference env (isolated)
- **GIVEN** step reference defines `env`
- **WHEN** reference executes
- **THEN** the referenced env values SHALL be applied only to that step invocation
- **AND** they SHALL NOT leak to later step invocations
- **AND** they SHALL NOT require a subshell that prevents the step from persisting its output
