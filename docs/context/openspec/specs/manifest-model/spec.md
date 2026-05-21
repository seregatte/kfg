## MODIFIED Requirements

### Requirement: Step Resource Schema

The Step resource MUST follow a defined schema.

#### Scenario: Required fields
- **GIVEN** a Step resource
- **WHEN** the resource is validated
- **THEN** `apiVersion` SHALL be `kfg.dev/v1alpha1`
- **AND** `kind` SHALL be `Step`
- **AND** `metadata.name` SHALL be present
- **AND** `spec.run` SHALL be present

#### Scenario: Optional needs field
- **GIVEN** a Step resource
- **WHEN** `spec.needs` is specified
- **THEN** it SHALL be a list of Step names
- **AND** each name SHALL reference an existing Step resource

#### Scenario: Optional when field
- **GIVEN** a Step resource
- **WHEN** `spec.when` is specified
- **THEN** it SHALL define a condition for step execution
- **AND** the condition SHALL reference only Step outputs

#### Scenario: Optional output field
- **GIVEN** a Step resource
- **WHEN** `spec.output` is specified
- **THEN** `spec.output.name` SHALL be present
- **AND** `spec.output.type` SHALL be `string`
- **AND** only one output SHALL be allowed per Step

#### Scenario: Optional cache field
- **GIVEN** a Step resource
- **WHEN** `spec.cache` is specified
- **THEN** the schema SHALL accept `enabled` and `key`
- **AND** the cache configuration SHALL define the default cache behavior for workflow references to that Step

#### Scenario: Optional failurePolicy field
- **GIVEN** a Step resource
- **WHEN** `spec.failurePolicy` is specified
- **THEN** it SHALL be either `Fail` or `Ignore`
- **AND** the default SHALL be `Fail`

### Requirement: Step Output Semantics

Step outputs MUST follow specific semantics.

#### Scenario: Output capture
- **GIVEN** a Step with `spec.output`
- **WHEN** the step's `run` executes successfully
- **THEN** stdout SHALL be captured as the output value
- **AND** stderr SHALL NOT be captured

#### Scenario: Output absence
- **GIVEN** a Step without `spec.output`
- **WHEN** the step executes
- **THEN** stdout is NOT stored as an output
- **AND** stdout passes through to the user

#### Scenario: Output scope after uncached execution
- **GIVEN** a Step output
- **WHEN** the output is produced by executing the Step in the current invocation
- **THEN** it SHALL be available within the current command invocation
- **AND** the output SHALL be discarded when the command completes unless the Step cache persists it

#### Scenario: Output scope after cache restore
- **GIVEN** a cacheable Step output
- **WHEN** the Step is restored from cache in a later command invocation
- **THEN** the output SHALL be repopulated in the current invocation runtime context
- **AND** later workflow logic SHALL observe it as if the Step had just executed successfully

### Requirement: CmdWorkflow Resource Schema

The CmdWorkflow resource MUST follow a defined schema.

#### Scenario: Required fields
- **GIVEN** a CmdWorkflow resource
- **WHEN** the resource is validated
- **THEN** `apiVersion` SHALL be `kfg.dev/v1alpha1`
- **AND** `kind` SHALL be `CmdWorkflow`
- **AND** `metadata.name` SHALL be present
- **AND** `spec.shell` SHALL be present
- **AND** `spec.cmds` SHALL be present

#### Scenario: Shell specification
- **GIVEN** a CmdWorkflow resource
- **WHEN** `spec.shell` is specified
- **THEN** it SHALL be `bash` (in the current version)
- **AND** unsupported shells SHALL cause a validation error

#### Scenario: Cmds list
- **GIVEN** a CmdWorkflow resource
- **WHEN** `spec.cmds` is specified
- **THEN** it SHALL be a list of Cmd names
- **AND** each name SHALL reference an existing Cmd resource
- **AND** the list SHALL contain at least one cmd

#### Scenario: Named step references in before and after
- **GIVEN** a CmdWorkflow resource
- **WHEN** `spec.before` or `spec.after` is specified
- **THEN** each entry SHALL be a step reference with `name` and `step` fields
- **AND** `name` SHALL be unique within that CmdWorkflow
- **AND** `name` SHALL be a valid namespace identifier
- **AND** `step` SHALL reference an existing Step resource
- **AND** `name` SHALL serve as the runtime output identity for conditions and env expansion

#### Scenario: StepReference cache override
- **GIVEN** a workflow step reference
- **WHEN** it declares `cache`
- **THEN** the schema SHALL accept `enabled` and `key`
- **AND** those values SHALL override the referenced Step cache defaults for that workflow invocation
