# Manifest Model Specification

## Purpose

kfg uses YAML manifests to define 3 resource kinds organized in one layer. These manifests are loaded from configured directories, merged according to precedence rules, and compiled into shell functions. This specification defines the structure and semantics of the manifest model.

**Execution Layer** (orchestration):
- Step: Reusable units of work
- Cmd: Shell function definitions
- CmdWorkflow: Entry point for shell generation

## Requirements
### Requirement: Resource Types

The manifest model MUST support 3 resource types.

#### Scenario: Step resource
- GIVEN a manifest file
- WHEN a `kind: Step` resource is defined
- THEN the resource defines a reusable unit of work
- AND the resource can produce outputs
- AND the resource can be referenced by CmdWorkflow

#### Scenario: Cmd resource
- GIVEN a manifest file
- WHEN a `kind: Cmd` resource is defined
- THEN the resource defines a shell function to be generated
- AND the resource defines the main execution body

#### Scenario: CmdWorkflow resource
- GIVEN a manifest file
- WHEN a `kind: CmdWorkflow` resource is defined
- THEN the resource defines a named collection of Cmds
- AND the resource is the required entry point for shell generation
- AND the resource specifies before/after steps

### Requirement: Resource Identity

Each resource MUST have a unique identity within its kind.

#### Scenario: Identity fields
- GIVEN any resource
- WHEN identity is determined
- THEN `apiVersion`, `kind`, and `metadata.name` SHALL form the identity
- AND resources with identical identity in the same layer SHALL cause an error

#### Scenario: Name constraints
- GIVEN a resource's `metadata.name`
- WHEN the name is validated
- THEN it SHALL be a valid shell identifier
- AND it SHALL contain only alphanumeric characters, hyphens, and underscores
- AND it SHALL not start with a digit

### Requirement: Step Resource Schema

The Step resource MUST follow a defined schema.

#### Scenario: Required fields
- GIVEN a Step resource
- WHEN the resource is validated
- THEN `apiVersion` SHALL be `kfg.dev/v1alpha1`
- AND `kind` SHALL be `Step`
- AND `metadata.name` SHALL be present
- AND `spec.run` SHALL be present

#### Scenario: Optional needs field
- GIVEN a Step resource
- WHEN `spec.needs` is specified
- THEN it SHALL be a list of Step names
- AND each name SHALL reference an existing Step resource

#### Scenario: Optional when field
- GIVEN a Step resource
- WHEN `spec.when` is specified
- THEN it SHALL define a condition for step execution
- AND the condition SHALL reference only Step outputs

#### Scenario: Optional output field
- GIVEN a Step resource
- WHEN `spec.output` is specified
- THEN `spec.output.name` SHALL be present
- AND `spec.output.type` SHALL be `string`
- AND only one output SHALL be allowed per Step

#### Scenario: Optional failurePolicy field
- GIVEN a Step resource
- WHEN `spec.failurePolicy` is specified
- THEN it SHALL be either `Fail` or `Ignore`
- AND the default SHALL be `Fail`

### Requirement: Step Output Semantics

Step outputs MUST follow specific semantics.

#### Scenario: Output capture
- GIVEN a Step with `spec.output`
- WHEN the step's `run` executes successfully
- THEN stdout SHALL be captured as the output value
- AND stderr SHALL NOT be captured

#### Scenario: Output absence
- GIVEN a Step without `spec.output`
- WHEN the step executes
- THEN stdout is NOT stored as an output
- AND stdout passes through to the user

#### Scenario: Output scope
- GIVEN a Step output
- WHEN the output is produced
- THEN it SHALL be available only within the current command invocation
- AND the output SHALL be discarded when the command completes

### Requirement: When Condition Schema

The `when` field MUST support specific condition types.

#### Scenario: Output condition
- GIVEN a `when` clause
- WHEN referencing a Step output
- THEN `step` SHALL specify the Step name
- AND `name` SHALL specify the output name
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

### Requirement: Cmd Resource Schema

The Cmd resource MUST follow a defined schema.

#### Scenario: Required fields
- GIVEN a Cmd resource
- WHEN the resource is validated
- THEN `apiVersion` SHALL be `kfg.dev/v1alpha1`
- AND `kind` SHALL be `Cmd`
- AND `metadata.name` SHALL be present
- AND `spec.run` SHALL be present

#### Scenario: Optional before field
- GIVEN a Cmd resource
- WHEN `spec.before` is specified
- THEN it SHALL be a list of Step names
- AND each name SHALL reference an existing Step resource

#### Scenario: Optional after field
- GIVEN a Cmd resource
- WHEN `spec.after` is specified
- THEN it SHALL be a list of Step names
- AND each name SHALL reference an existing Step resource

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

### Requirement: Manifest Loading

Manifests MUST be loaded.

#### Scenario: Recursive loading
- GIVEN a manifest directory
- WHEN manifests are loaded
- THEN all `*.yaml` and `*.yml` files SHALL be loaded recursively
- AND files SHALL be processed in lexicographic order within each directory

#### Scenario: Non-existent directory
- GIVEN a path segment that does not exist
- WHEN manifests are loaded
- THEN the path SHALL be ignored
- AND no error SHALL be raised

#### Scenario: Multi-document files
- GIVEN a YAML file with multiple documents separated by `---`
- WHEN the file is loaded
- THEN each document SHALL be parsed as a separate resource

### Requirement: Manifest Merging

Manifests from multiple layers MUST be merged correctly.

#### Scenario: Override by precedence
- GIVEN the same resource identity exists in multiple layers
- WHEN manifests are merged
- THEN the resource from the rightmost (highest precedence) layer SHALL be used
- AND other versions SHALL be discarded

#### Scenario: Resource combination
- GIVEN different resources in different layers
- WHEN manifests are merged
- THEN all resources SHALL be combined
- AND the final set SHALL include resources from all layers

#### Scenario: Duplicate detection within layer
- GIVEN the same resource identity exists twice in the same layer
- WHEN manifests are loaded
- THEN an error SHALL be raised
- AND the error SHALL identify the duplicate resource

### Requirement: Reference Validation

Cross-resource references MUST be validated.

#### Scenario: Step reference in Cmd
- GIVEN a Cmd references a Step in `before` or `after`
- WHEN the manifest is validated
- THEN the Step SHALL exist in the merged manifest set
- AND missing references SHALL cause a validation error

#### Scenario: Cmd reference in CmdWorkflow
- GIVEN a CmdWorkflow references a Cmd in `cmds`
- WHEN the manifest is validated
- THEN the Cmd SHALL exist in the merged manifest set
- AND missing references SHALL cause a validation error

#### Scenario: Step reference in needs
- GIVEN a Step references another Step in `needs`
- WHEN the manifest is validated
- THEN the referenced Step SHALL exist
- AND missing references SHALL cause a validation error

#### Scenario: Output reference in when
- GIVEN a Step references an output in `when`
- WHEN the manifest is validated
- THEN the referenced Step SHALL exist
- AND the referenced Step SHALL declare the output
- AND missing declarations SHALL cause a validation error

### Requirement: Dependency Resolution

Step dependencies MUST be resolved correctly.

#### Scenario: Direct dependency
- GIVEN Step A has `needs: [B]`
- WHEN dependencies are resolved
- THEN Step B SHALL execute before Step A

#### Scenario: Transitive dependency
- GIVEN Step A has `needs: [B]` and Step B has `needs: [C]`
- WHEN dependencies are resolved
- THEN Step C SHALL execute before Step B
- AND Step B SHALL execute before Step A

#### Scenario: Circular dependency
- GIVEN Step A has `needs: [B]` and Step B has `needs: [A]`
- WHEN dependencies are resolved
- THEN a circular dependency error SHALL be raised
- AND the error SHALL identify the cycle

#### Scenario: Topological ordering
- GIVEN multiple Steps with `needs` relationships
- WHEN execution order is determined
- THEN the order SHALL be topologically sorted
- AND tie-breaking SHALL be consistent and deterministic