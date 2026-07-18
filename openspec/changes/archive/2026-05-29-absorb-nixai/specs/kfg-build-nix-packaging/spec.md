## MODIFIED Requirements

### Requirement: Version exposure

The flake SHALL expose the version string for downstream consumers.

#### Scenario: lib.version

- **WHEN** a consumer accesses `kfg.lib.version`
- **THEN** the value SHALL match the hardcoded `version` attribute in the flake

## REMOVED Requirements

### Requirement: nixai dependency via inputsFrom

**Reason**: kfg now declares all dev shell tooling inline. The `nixai` input and
the `inputsFrom` delegation are removed as part of the absorb-nixai change.

**Migration**: No migration needed. Developers run `nix develop` as before; the
shell is now self-contained. No reference to `github:seregatte/nixai` is required.
