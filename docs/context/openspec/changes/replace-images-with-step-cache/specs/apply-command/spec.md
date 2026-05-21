## ADDED Requirements

### Requirement: Apply Refresh Propagation
The apply command SHALL generate shell code that can force cacheable Steps to refresh.

#### Scenario: Apply with refresh flag
- **WHEN** user runs `kfg apply -k path --refresh`
- **THEN** the generated shell code SHALL export or embed refresh state equivalent to `KFG_REFRESH`
- **AND** cacheable Steps in that shell SHALL bypass matching cache entries when executed

#### Scenario: Apply without refresh flag
- **WHEN** user runs `kfg apply -k path` without `--refresh`
- **THEN** the generated shell code SHALL use cached Step entries when available
