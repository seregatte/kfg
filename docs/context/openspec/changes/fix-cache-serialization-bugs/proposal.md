## Why

Cache metadata serialization produces malformed YAML for multi-line base64 output values, preventing cache entry inspection and causing parse failures in `kfg sys gc` commands. Additionally, artifact paths contain double slashes and the `--store` flag is missing from CLI commands despite being specified in specs. These bugs were discovered during comprehensive cache testing and block production use of cache features for skill installation workflows.

## What Changes

- Fix YAML serialization to properly encode multi-line base64 values using block scalar syntax
- Normalize artifact path construction to eliminate double slashes
- Add `--store` flag to `run` and `apply` commands for isolated cache testing
- Ensure `metadata.yaml` follows YAML 1.2 specification for all output types
- Preserve backward compatibility with existing cache entries
- Fix typo in `kfg-cache-step` spec ("indentified" → "identified")
- Update `kfg-manifest-model` spec to remove "key" field acceptance
- Update test fixtures to remove obsolete `cache.key` usage
- Update `kfg-cache-sys-gc-command` spec implementation notes (remove key/script hash mentions)

## Capabilities

### New Capabilities

None - fixing existing capability bugs without introducing new features.

### Modified Capabilities

- `kfg-cache-step`: Fix typo "indentified" → "identified" in cache identity requirement; requirements otherwise unchanged
- `kfg-manifest-model`: Remove "key" from cache field specification (field removed in previous change but spec not updated)
- `kfg-cache-sys-gc-command`: Update implementation notes to reflect current cache identity (SHA256 of StepReference.name only, no key/script components)
- `kfg-cli-store-isolation`: Requirement implementation incomplete - adding `--store` flag to enable isolated testing per spec scenario

## Impact

**Affected Code:**
- `src/internal/generate/templates.go` - shell template generation for cache helpers
- `src/cmd/kfg/` - CLI command definitions (`run.go`, `apply.go`)
- Cache metadata serialization logic (likely in shell helpers or Go code)

**Affected APIs:**
- CLI flags: Adding `--store` to `kfg run` and `kfg apply` commands
- Shell helpers: `__kfg_cache_store()` metadata.yaml generation
- Cache commands: `kfg sys gc inspect` will work with fixed metadata

**Dependencies:**
- YAML serialization must handle multi-line base64 properly
- Path normalization in artifact registration
- Viper configuration for `--store` flag binding

**Systems:**
- Cache storage layer (metadata.yaml format)
- CLI command interface (flag parsing)
- Shell runtime (artifact path construction)
- Step cache helpers (metadata generation)

**Test Updates:**
- Remove obsolete `cache.key` from test fixtures in `tests/bats/workflows/step-cache-isolation.bats`
- Ensure fixtures reflect current manifest model (only `enabled` field)

**Breaking Changes:**
None - fixes are backward compatible with existing cache entries