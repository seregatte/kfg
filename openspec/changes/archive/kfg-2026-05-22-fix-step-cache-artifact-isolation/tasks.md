## 1. Artifact isolation in generated runtime

- [x] 1.1 Update generated command wrappers so Step artifacts are not pre-registered globally before the Step executes.
- [x] 1.2 Add before/after artifact snapshots around cacheable Step execution and compute the Step-local artifact delta.
- [x] 1.3 Persist only Step-local artifacts plus the relevant declarative Step and StepReference artifact paths.

## 2. Path-preserving cache restore

- [x] 2.1 Update cache storage helpers to preserve each artifact's original relative path inside the cache entry.
- [x] 2.2 Update restore helpers to recreate parent directories and restore artifacts back to their original relative paths.
- [x] 2.3 Add regression tests for nested paths and duplicate basenames in different directories.

## 3. Cache diagnostics and regression coverage

- [x] 3.1 Add detail logs for refresh bypass, cache hit, cache miss, store start/end, and restore start/end.
- [x] 3.2 Add debug logs for cache path/identity and the artifact paths being stored or restored.
- [x] 3.3 Add unit, golden, and integration/Bats tests that verify artifact isolation, path-preserving restore, output restore, and cache logging behavior.
