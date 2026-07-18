## ADDED Requirements

### Requirement: Public capability documentation

Domain packages MUST document every independently consumable public capability adjacent to its Kustomization entrypoint.

#### Scenario: Capability usage contract
- **WHEN** a domain exposes a selective capability entrypoint
- **THEN** an adjacent README SHALL document composition, exported resources, configurable fields, JSON Patch examples, artifacts, operational behavior, limitations, and validation
- **AND** the domain root README SHALL link to that capability documentation

#### Scenario: Capability behavior changes
- **WHEN** a change modifies public usage of a domain capability
- **THEN** its adjacent README MUST be updated in the same change
- **AND** automated domain tests SHALL validate documentation coverage and advertised entrypoints

### Requirement: Selective domain entrypoints

Domain packages SHALL support documented selective entrypoints for independently useful capabilities in addition to the stable complete-package entrypoint.

#### Scenario: Selective capability consumption
- **WHEN** a consumer needs one independently useful domain capability
- **THEN** the consumer MAY reference the capability's documented Kustomization
- **AND** that Kustomization SHALL build without unrelated capability resources
- **AND** aggregate domain entrypoints SHALL remain available unless a breaking migration is explicitly specified
