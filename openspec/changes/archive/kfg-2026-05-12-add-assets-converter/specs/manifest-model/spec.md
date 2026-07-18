## MODIFIED Requirements

### Requirement: Resource Types

The manifest model MUST support 5 resource types.

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

#### Scenario: Assets resource
- GIVEN a manifest file
- WHEN a `kind: Assets` resource is defined
- THEN the resource declares a data payload
- AND the resource specifies its data format
- AND the resource is NOT part of shell generation

#### Scenario: Converter resource
- GIVEN a manifest file
- WHEN a `kind: Converter` resource is defined
- THEN the resource declares a data transformation
- AND the resource specifies input format, yq expression, and output format
- AND the resource is NOT part of shell generation

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
- THEN it SHALL be a valid namespace identifier
- AND it SHALL contain only lowercase alphanumeric characters, hyphens, and dots
- AND it SHALL not start with a digit
