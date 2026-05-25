## 1. YAML Serialization Fix

- [ ] 1.1 Locate cache metadata serialization logic in shell template generation
- [ ] 1.2 Identify where `valueEncoded` field is written to metadata.yaml
- [ ] 1.3 Implement YAML block scalar syntax (`|`) for multi-line base64 output values
- [ ] 1.4 Verify YAML parser can read new metadata format correctly

## 2. Path Normalization Fix

- [ ] 2.1 Locate artifact path construction in `__kfg_cache_store()` shell helper
- [ ] 2.2 Identify where double slashes occur (OUTPUT_DIR + artifact name concatenation)
- [ ] 2.3 Implement path normalization: `${path%%/}/${artifact#/}`
- [ ] 2.4 Verify artifact paths stored without double slashes

## 3. CLI Flag Implementation

- [ ] 3.1 Add `--store` flag to `kfg run` command in src/cmd/kfg/run.go
- [ ] 3.2 Add `--store` flag to `kfg apply` command in src/cmd/kfg/apply.go
- [ ] 3.3 Bind `--store` flag to Viper configuration with proper precedence
- [ ] 3.4 Update command help text to document `--store` flag usage

## 4. Spec Corrections

- [ ] 4.1 Fix typo in `kfg-cache-step` spec (line 10: "indentified" → "identified")
- [ ] 4.2 Update `kfg-manifest-model` spec (line 101: remove "and key" from cache field)
- [ ] 4.3 Update `kfg-cache-sys-gc-command` spec implementation notes (lines 100-103)
- [ ] 4.4 Update specs README to document no spec requirement changes (implementation fixes only)

## 5. Test Fixture Updates

- [ ] 5.1 Remove `cache.key` from test fixtures in `tests/bats/workflows/step-cache-isolation.bats`
- [ ] 5.2 Update fixture YAML to use only `cache.enabled` field
- [ ] 5.3 Verify bats tests still pass after fixture updates

## 6. Testing and Validation

- [ ] 6.1 Run cache test scenario 1 (initial execution) to verify metadata format
- [ ] 6.2 Run `kfg sys gc ls` to verify cache entries parseable
- [ ] 6.3 Run `kfg sys gc inspect <id>` to verify metadata inspection works
- [ ] 6.4 Test `--store` flag with isolated cache directory
- [ ] 6.5 Test backward compatibility with existing cache entries
- [ ] 6.6 Run full test suite: `make test-bats` to ensure no regressions

## 7. Documentation Updates

- [ ] 7.1 Update CLI help output to reflect `--store` flag addition
- [ ] 7.2 Verify cache specs remain accurate after corrections