## 1. Engine execution fix

- [ ] 1.1 Replace command-substitution execution for `spec.output` Steps with parent-shell stdout capture.
- [ ] 1.2 Preserve output storage and artifact registration semantics in the updated wrapper path.

## 2. Regression tests

- [ ] 2.1 Add generator or golden tests covering output-producing Steps that mutate runtime state.
- [ ] 2.2 Add integration coverage for cache store/restore of both output values and runtime artifacts from the same Step invocation.
