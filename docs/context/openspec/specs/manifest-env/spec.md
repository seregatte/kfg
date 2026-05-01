# Manifest Environment Specification

## Purpose

Define environment variable semantics for Step and Cmd resources.

## Requirements

### Requirement: Step Environment Defaults

Steps MUST support default environment variables in `spec.env`.

#### Scenario: Step default env available during execution
- **GIVEN** a Step defines `spec.env` with key-value pairs
- **WHEN** the Step executes
- **THEN** all specified variables are available to `spec.run`
- **AND** variables are exported in step's execution scope

#### Scenario: Step without env works unchanged
- **GIVEN** a Step without `spec.env`
- **WHEN** the Step executes
- **THEN** no error occurs

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

### Requirement: Cmd Environment Defaults

Cmds MUST support default environment variables in `spec.env`.

#### Scenario: Cmd default env available during execution
- **GIVEN** a Cmd defines `spec.env` with key-value pairs
- **WHEN** Cmd wrapper executes
- **THEN** variables available to before steps
- **AND** variables available to `spec.run`
- **AND** variables available to after steps

#### Scenario: Cmd env exported at wrapper start
- **GIVEN** a Cmd defines `spec.env.PROVIDER=anthropic`
- **WHEN** generated wrapper function starts
- **THEN** `PROVIDER` exported before any before steps run

### Requirement: Cmd Environment Direct Assignment

Cmd env variables MUST use direct assignment in generated shell code, NOT `${KEY:-default}` expansion syntax. This prevents env leakage between sequential Cmd invocations in the same shell.

#### Scenario: Cmd env does not leak between sequential invocations
- **GIVEN** a CmdWorkflow with multiple Cmds (e.g., opencode, gemini, pi)
- **AND** each Cmd defines `KFG_AGENT` in `spec.env` with different values
- **WHEN** the generated shell code executes all Cmds sequentially
- **THEN** each Cmd wrapper MUST export its own `KFG_AGENT` value unconditionally
- **AND** a preceding Cmd's `KFG_AGENT` value MUST NOT persist to subsequent Cmds
- **AND** generated code MUST be `export KFG_AGENT="opencode"` (direct assignment)
- **AND** generated code MUST NOT be `export KFG_AGENT="${KFG_AGENT:-opencode}"` (default expansion)

#### Scenario: Cmd env overrides external shell environment
- **GIVEN** a Cmd with `spec.env.FOO=bar`
- **AND** the shell environment already has `FOO=external`
- **WHEN** the Cmd wrapper function is called
- **THEN** `FOO` MUST be set to `"bar"` (Cmd's value wins, not the external value)

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
- **THEN** wrapped in subshell for isolation

### Requirement: Env Merge Behavior

Env values MUST merge correctly.

#### Scenario: Override wins (shallow merge)
- **GIVEN** Step default env and reference env have same key
- **WHEN** merged
- **THEN** reference value wins

#### Scenario: New keys added
- **GIVEN** Step default env has key `A`
- **AND** reference env has key `B`
- **WHEN** merged
- **THEN** both `A` and `B` present