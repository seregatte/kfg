## MODIFIED Requirements

### Requirement: Cache Persistence Semantics
The runtime SHALL persist only the observable Step results needed for downstream reuse. Artifact detection uses a hybrid approach: registered artifact delta plus filesystem diff performed by the Go `store` subcommand. Automatic filesystem discovery SHALL persist leaf files and symbolic links without persisting their parent directories, while explicitly registered or declared directories SHALL remain valid artifacts.

#### Scenario: Declarative artifacts are cached
- **WHEN** a cacheable Step declares `spec.artifacts`
- **THEN** those artifact paths SHALL be included in the persisted cache entry for that Step invocation

#### Scenario: StepReference artifacts are cached
- **WHEN** a cacheable workflow Step reference declares `artifacts`
- **THEN** those artifact paths SHALL be included in the persisted cache entry for that Step invocation

#### Scenario: Runtime artifact registrations are cached
- **WHEN** a cacheable Step calls `__kfg_add_artifact` during execution
- **THEN** the runtime SHALL persist the artifact paths newly registered by that Step invocation
- **AND** SHALL NOT persist unrelated artifact paths that were already registered before the Step started
- **AND** the shell wrapper SHALL pass the before/after artifact lists to `kfg sys cache store` via JSON on stdin

#### Scenario: Filesystem diff detects unregistered artifacts
- **WHEN** a cacheable Step creates files in the workdir but does not call `__kfg_add_artifact` for them
- **THEN** the Go `store` subcommand SHALL detect those files via filesystem diff and include them in the cache entry

#### Scenario: Filesystem diff excludes parent directories
- **WHEN** a cacheable Step creates `openspec/config.yaml` under a new `openspec/` directory without registering either path
- **THEN** automatic filesystem discovery SHALL include `openspec/config.yaml` in the cache entry
- **AND** SHALL NOT include `openspec` as an artifact

#### Scenario: Explicit directory artifacts remain supported
- **WHEN** a cacheable Step explicitly declares or registers `.opencode/` as an artifact
- **THEN** the cache entry SHALL include `.opencode/` as a directory artifact
- **AND** cache restore SHALL retain its explicit directory cleanup semantics

#### Scenario: Step outputs are cached automatically
- **WHEN** a cacheable Step declares `spec.output`
- **THEN** the runtime SHALL persist the captured output value as part of the cache entry
- **AND** no extra manifest field SHALL be required to opt in output persistence

#### Scenario: Output-producing step caches output and runtime artifacts together
- **WHEN** a cacheable Step declares `spec.output` and registers artifacts during execution
- **THEN** the runtime SHALL persist both the captured output value and the Step-local runtime artifact registrations from the same invocation

## ADDED Requirements

### Requirement: Cache discovery compatibility
The runtime MUST NOT restore cache entries created under automatic directory-discovery semantics that can register parent directories for recursive cleanup.

#### Scenario: Legacy cache entry cannot produce a hit
- **WHEN** a cache entry was created in the previous unversioned identity namespace
- **THEN** the current runtime SHALL treat the corresponding StepReference as a cache miss
- **AND** SHALL rebuild the entry under the current cache identity namespace

#### Scenario: Legacy cache remains manageable
- **WHEN** an entry remains in a previous cache identity namespace
- **THEN** cache administration and garbage collection commands SHALL remain able to inspect or remove the stored entry
