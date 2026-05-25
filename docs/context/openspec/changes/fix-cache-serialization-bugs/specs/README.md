# Spec Changes Summary

This change updates existing specs to reflect current implementation state after `kfg-simplify-step-cache-identity` change was archived.

## Changes

### kfg-cache-step
- Fix typo: "indentified" → "identified" in cache identity requirement
- No requirement changes - only typo correction

### kfg-manifest-model  
- Remove "key" from cache field specification
- Schema now accepts only `enabled` field (key removed in previous change)

### kfg-cache-sys-gc-command
- Update implementation notes to reflect current cache identity
- Remove mention of key/script hash components (no longer used)
- Cache identity: Step reference name only

## Rationale

Previous change `kfg-simplify-step-cache-identity` removed `cache.key` from manifest model but did not update related specs. This change corrects those oversights to maintain spec-implementation alignment.