## MODIFIED Requirements

### Requirement: KFG_VERBOSE control

The system SHALL use `KFG_VERBOSE` as the sole variable controlling human output visibility.

#### Scenario: Default value
- **WHEN** `KFG_VERBOSE` is not set
- **THEN** default value is `1`
- **AND** error events produce human stderr output

#### Scenario: Verbose=0 (silent)
- **WHEN** `KFG_VERBOSE=0`
- **THEN** all log events persist to JSONL
- **AND** no human output appears in stderr

#### Scenario: Verbose=1 (error only)
- **WHEN** `KFG_VERBOSE=1`
- **THEN** error events produce human stderr output
- **AND** warn, info, detail, debug do not produce human output

#### Scenario: Verbose=2 (error + warn + info)
- **WHEN** `KFG_VERBOSE=2`
- **THEN** error, warn, info events produce human stderr output
- **AND** detail, debug do not produce human output

#### Scenario: Verbose=3 (all levels)
- **WHEN** `KFG_VERBOSE=3`
- **THEN** all levels produce human stderr output

#### Scenario: Nested internal kfg subprocesses stay quiet
- **WHEN** runtime wrappers invoke an internal `kfg` subprocess while the parent invocation has `KFG_VERBOSE=3`
- **THEN** the child subprocess SHALL suppress its human log output by forcing child-scoped `KFG_VERBOSE=0`
- **AND** the parent invocation SHALL continue using its original verbosity setting
