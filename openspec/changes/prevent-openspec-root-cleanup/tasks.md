## 1. Cache Identity Safety

- [x] 1.1 Add an engine-controlled cache schema namespace to `ComputeIdentity` so current entries cannot collide with unsafe unversioned entries
- [x] 1.2 Add Go unit tests proving identity remains deterministic within a version and differs from the legacy StepReference-only hash
- [x] 1.3 Verify cache store, exists, restore, inspection, and garbage collection continue to resolve or enumerate entries through their intended identity paths

## 2. Leaf-Only Filesystem Discovery

- [x] 2.1 Update the Go filesystem snapshot to include files and symbolic links while excluding directory entries
- [x] 2.2 Update the generated Bash pre-Step snapshot to use the same leaf-only path semantics
- [x] 2.3 Add Go unit tests for nested files, symbolic links, excluded parent directories, and deterministic snapshot ordering
- [x] 2.4 Add cache store tests proving automatic discovery excludes `openspec/` while caching `openspec/config.yaml`, and explicit directory artifacts remain supported
- [x] 2.5 Add shell generation tests asserting cacheable Steps emit a leaf-only pre-execution snapshot command

## 3. OpenSpec Cleanup Regression Coverage

- [x] 3.1 Add a domain Bats regression that exercises an OpenSpec installation cache miss, creates session planning data under `openspec/changes/`, runs final cleanup, and verifies the root and planning files remain
- [x] 3.2 Extend the domain regression through a cache hit and verify restore registers concrete integration files without registering the `openspec/` root
- [x] 3.3 Verify generated OpenSpec integration artifacts remain eligible for cleanup and explicitly declared disposable directory artifacts retain recursive cleanup behavior
- [x] 3.4 Run the AI-agent overlay build tests to confirm existing workflow composition, conditions, and public resource names are unchanged

## 4. Documentation

- [x] 4.1 Update the adjacent OpenSpec capability README generated-artifact and cache/cleanup sections to describe persistent planning-root preservation on cache miss and cache hit
- [x] 4.2 Verify the README retains every capability documentation contract section and accurately describes unchanged store ownership and cleanup limitations

## 5. Validation

- [x] 5.1 Run focused Go cache and shell-generation unit tests through `nix develop .#dev`
- [x] 5.2 Run focused framework and AI-agent Bats tests, including the session cleanup regression, through `nix develop .#dev`
- [x] 5.3 Run `nix develop .#dev --command make fmt lint vet test` and `nix develop .#dev --command make test-bats`
- [x] 5.4 Run `openspec validate --change prevent-openspec-root-cleanup` and resolve every proposal, design, spec, or task validation error
