## Context

The cache implementation successfully stores and restores Step artifacts and outputs, but metadata serialization has critical bugs discovered during comprehensive testing. The current implementation generates malformed YAML for multi-line base64-encoded output values, uses incorrect path construction causing double slashes, and lacks the `--store` CLI flag despite being specified in specs.

**Current State:**
- Cache metadata stored in `metadata.yaml` within cache entry directories
- Output values encoded as base64 in `valueEncoded` field
- Multi-line base64 values written without YAML block scalar syntax (causes parse failures)
- Artifact paths constructed with double slashes (`.pi/skills//openspec-*`)
- `kfg sys gc inspect` fails on malformed cache entries
- `kfg run` and `kfg apply` commands have no `--store` flag for isolated testing

**Stakeholders:**
- Cache implementation (shell helpers in generated code)
- CLI commands (`run`, `apply`, `sys gc`)
- Framework and domain packages using cacheable Steps (skill installation)
- Testing workflows requiring isolated cache stores

**Constraints:**
- Must maintain backward compatibility with existing cache entries
- YAML format must follow YAML 1.2 specification
- Shell helper changes must work across bash environments
- CLI flag addition must integrate with Viper configuration

## Goals / Non-Goals

**Goals:**
- Fix YAML serialization to properly encode multi-line base64 using block scalar syntax (`|`)
- Normalize artifact paths to eliminate double slashes
- Add `--store` flag to `run` and `apply` commands
- Ensure `kfg sys gc inspect` works with all cache entries
- Maintain backward compatibility with existing cache entries

**Non-Goals:**
- Changing cache identity computation (SHA256(StepReference.name))
- Modifying cache restore/store logic beyond metadata serialization
- Adding new cache features or capabilities
- Refactoring cache storage structure
- Changing YAML schema for metadata (only fixing encoding)

## Decisions

### Decision 1: YAML Block Scalar Syntax for Multi-line Base64

**Choice:** Use YAML literal block scalar (`|`) for multi-line base64 output values

**Rationale:**
- YAML block scalar (`|`) preserves line breaks exactly as written
- Base64 encoding can span multiple lines when output is large
- Literal block scalar ensures YAML parsers can read the value correctly
- Alternative: Use folded block scalar (`>`), but this converts line breaks to spaces (incorrect for base64)
- Alternative: Keep base64 as single line, but this fails for large outputs and is not what current implementation produces

**Implementation:**
```yaml
output:
  name: installed
  valueEncoded: |
    ZmluZC1kb2NzCmZpbmQtZG9jcwpvcGVuc3BlYy...
```

### Decision 2: Path Normalization in Artifact Registration

**Choice:** Normalize paths by removing redundant slashes before storing

**Rationale:**
- Double slashes come from path construction: `OUTPUT_DIR` + `/` + artifact name
- If `OUTPUT_DIR` ends with `/` and artifact name starts with `/`, result is `//`
- Shell path normalization: `${path%%/}/${artifact#/}` removes trailing/leading slashes
- Alternative: Fix in Go code, but artifact paths are constructed in shell helpers
- Alternative: Require OUTPUT_DIR to not end with `/`, but this is a convention change

**Implementation:**
```bash
# Normalize path before storing
local normalized_path="${OUTPUT_DIR%%/}/${artifact#/}"
```

### Decision 3: CLI Flag Binding for `--store`

**Choice:** Add `--store` flag to `run` and `apply` commands with Viper binding

**Rationale:**
- Spec `kfg-cli-store-isolation` requires `--store <path>` flag
- Viper configuration already supports KFG_STORE_DIR environment variable
- Flag should override environment variable, following existing patterns
- Alternative: Add global `--store` flag, but specs mention it for specific commands
- Alternative: Only use environment variable, but flag provides better test isolation UX

**Implementation:**
- Add flag to Cobra command definitions
- Bind flag to Viper configuration key
- Override KFG_STORE_DIR when flag is provided

## Risks / Trade-offs

**Risk: Breaking existing cache entries with new YAML format**
→ Mitigation: Use backward-compatible YAML syntax that parsers can read. Existing entries with malformed YAML will remain unparseable (can't fix without recreation). New entries will be parseable.

**Risk: Path normalization changes artifact storage locations**
→ Mitigation: Normalization only affects new cache entries. Existing entries with double slashes will restore correctly (paths stored as-is). No migration needed.

**Risk: Flag binding conflicts with existing configuration**
→ Mitigation: Follow existing pattern: flag overrides env var. Test with both flag and env var set to ensure precedence.

**Risk: Shell helper changes affect multiple packages**
→ Mitigation: Changes are in generated shell code templates, tested through existing Bats tests. All packages using cache helpers will benefit from fix.

**Trade-off: Can't fix existing malformed cache entries**
→ Acceptance: Malformed entries require recreation via `--refresh`. This is acceptable as cache is ephemeral and rebuildable.

**Trade-off: Path normalization adds processing overhead**
→ Acceptance: Minimal overhead (string manipulation). Worth avoiding double slashes in stored metadata.

## Migration Plan

**Phase 1: Fix Implementation**
- Implement YAML block scalar syntax for metadata serialization
- Implement path normalization in artifact registration
- Add `--store` flag to CLI commands

**Phase 2: Testing**
- Run cache tests from testing report to verify fixes
- Test with existing cache entries (backward compatibility)
- Test with new cache entries (YAML parseable)
- Test with both `--store` flag and KFG_STORE_DIR env var

**Phase 3: Documentation**
- Update cache spec if needed (requirements unchanged, implementation fixed)
- No user-facing changes beyond `--store` flag addition

**Rollback Strategy:**
- If YAML format causes issues, revert to old format (but entries already malformed)
- If path normalization causes issues, disable normalization (double slashes tolerated)
- If `--store` flag causes issues, remove flag and use env var only

## Open Questions

None - all implementation details are clear from testing report and existing specs.