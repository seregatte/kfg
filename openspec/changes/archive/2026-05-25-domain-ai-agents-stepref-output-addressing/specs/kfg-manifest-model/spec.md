## MODIFIED Requirements

### Requirement: When Condition Schema

The `when` field MUST support specific condition types.

#### Scenario: Output condition
- GIVEN a `when` clause
- WHEN referencing a Step output
- THEN `step` SHALL specify the `StepReference.name` of a workflow step invocation
- AND `name` SHALL specify the output name declared by the referenced Step resource
- AND exactly one operator SHALL be specified

#### Scenario: Comparison operators
- GIVEN a `when.output` condition
- WHEN an operator is used
- THEN `equals` SHALL compare for exact string equality
- AND `in` SHALL check membership in a list
- AND `contains` SHALL check substring presence
- AND `matches` SHALL evaluate a regular expression

#### Scenario: Boolean composition
- GIVEN a `when` clause
- WHEN composition is needed
- THEN `allOf` SHALL require all conditions to be true
- AND `anyOf` SHALL require at least one condition to be true
- AND `not` SHALL negate the nested condition

#### Scenario: Missing output handling
- GIVEN a `when` condition referencing a non-existent output
- WHEN the condition is evaluated
- THEN the condition SHALL evaluate to false
- AND no error SHALL be raised

### Requirement: CmdWorkflow Resource Schema

The CmdWorkflow resource MUST follow a defined schema.

#### Scenario: Required fields
- GIVEN a CmdWorkflow resource
- WHEN the resource is validated
- THEN `apiVersion` SHALL be `kfg.dev/v1alpha1`
- AND `kind` SHALL be `CmdWorkflow`
- AND `metadata.name` SHALL be present
- AND `spec.shell` SHALL be present
- AND `spec.cmds` SHALL be present

#### Scenario: Shell specification
- GIVEN a CmdWorkflow resource
- WHEN `spec.shell` is specified
- THEN it SHALL be `bash` (in the current version)
- AND unsupported shells SHALL cause a validation error

#### Scenario: Cmds list
- GIVEN a CmdWorkflow resource
- WHEN `spec.cmds` is specified
- THEN it SHALL be a list of Cmd names
- AND each name SHALL reference an existing Cmd resource
- AND the list SHALL contain at least one cmd

#### Scenario: Named step references
- GIVEN a workflow `before` or `after` entry
- WHEN the resource is validated
- THEN the entry SHALL declare `name`
- AND `name` SHALL be unique within that CmdWorkflow
- AND `name` SHALL be a valid namespace identifier

#### Scenario: Output references resolve by step-reference name
- GIVEN a workflow step reference with `when.output.step`
- WHEN the resource is validated
- THEN the referenced value SHALL match an existing `StepReference.name` in the same CmdWorkflow
- AND the referenced workflow step SHALL point to a Step resource with `spec.output`
- AND `when.output.name` SHALL match the referenced Step resource output name

#### Scenario: Runtime output identity
- GIVEN a workflow step reference that executes a Step with `spec.output`
- WHEN the output is stored for later conditions or env expansion
- THEN the output SHALL be keyed by `StepReference.name`
- AND NOT by the underlying Step resource `metadata.name`
