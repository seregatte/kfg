# Manifest Model Specification

## Purpose

kfg uses YAML manifests to define 5 resource kinds organized in two layers. These manifests are loaded from configured directories, merged according to precedence rules, and compiled into shell functions. This specification defines the structure and semantics of the manifest model.

**Execution Layer** (orchestration):
- Step: Reusable units of work
- Cmd: Shell function definitions
- CmdWorkflow: Entry point for shell generation

**Source Layer** (data declarations):
- Assets: Structured data payloads
- Converter: Data transformation rules via yq-go expressions

## ADDED Requirements

### Requirement: Source Layer Resource Types

The manifest model SHALL support 2 source layer resource types that are excluded from the execution pipeline but accessible via the build result.

#### Scenario: Assets resource
- **GIVEN** a manifest file
- **WHEN** a `kind: Assets` resource is defined
- **THEN** the resource declares structured data with format metadata
- **AND** the resource SHALL be excluded from the execution index
- **AND** the resource SHALL be included in `$KFG_BUILD_RESULT_FILE`

#### Scenario: Converter resource
- **GIVEN** a manifest file
- **WHEN** a `kind: Converter` resource is defined
- **THEN** the resource declares a yq-go transformation expression
- **AND** the resource SHALL be excluded from the execution index
- **AND** the resource SHALL be included in `$KFG_BUILD_RESULT_FILE`

### Requirement: Source Layer Resolution Behavior

Source layer resources SHALL be explicitly skipped during resolution.

#### Scenario: Resolution skips source kinds
- **WHEN** the resolution pipeline indexes parsed resources
- **THEN** Assets and Converter resources SHALL be excluded from Step, Cmd, and CmdWorkflow indexes
- **AND** a debug log entry SHALL be emitted for each skipped resource
- **AND** the resources SHALL remain available in the YAML build result

### Requirement: Resource Kind Count

The manifest model SHALL support exactly 5 resource kinds.

#### Scenario: All resource kinds
- **WHEN** the manifest model is queried for supported kinds
- **THEN** it SHALL return: Step, Cmd, CmdWorkflow, Assets, Converter
- **AND** exactly 3 SHALL be execution kinds (Step, Cmd, CmdWorkflow)
- **AND** exactly 2 SHALL be source kinds (Assets, Converter)
