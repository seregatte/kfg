## MODIFIED Requirements

### Requirement: Cache store subcommand
The `kfg sys cache store` subcommand SHALL persist Step execution results (artifacts and output) to the cache, reading structured input from stdin. Automatically detected filesystem artifacts SHALL exclude directory entries, while explicit artifact inputs SHALL preserve directory entries.

#### Scenario: Store cache entry with artifacts and output
- **WHEN** user runs `kfg sys cache store <step-ref> --workdir <path>` with JSON on stdin containing `before`, `after`, `declarative`, and `output` fields
- **THEN** the subcommand SHALL compute the artifact delta (after - before), merge with declarative artifacts, perform fs diff on workdir to detect unregistered leaf artifacts, copy all artifacts to the cache entry directory, write `metadata.yaml`, and commit atomically via rename

#### Scenario: Store reads JSON from stdin
- **WHEN** `kfg sys cache store` is invoked
- **THEN** it SHALL read a JSON object from stdin with fields: `before` (string array), `after` (string array), `declarative` (string array), and `output` (object with `name` string and `value` base64 string)
- **AND** all arrays contain relative paths from the workdir

#### Scenario: Store computes artifact delta
- **WHEN** `before` contains `["a.txt"]` and `after` contains `["a.txt", "b.txt"]`
- **THEN** the delta SHALL be `["b.txt"]`
- **AND** only `b.txt` plus any declarative artifacts and fs-diff-detected leaf artifacts SHALL be cached

#### Scenario: Store performs fs diff on workdir
- **WHEN** the step creates files in workdir that are not in the `after` list
- **THEN** the subcommand SHALL detect those files via filesystem diff and include them in the cache entry
- **AND** SHALL NOT include newly created parent directories unless they were provided as explicit artifacts

#### Scenario: Store uses versioned StepReference identity
- **WHEN** user runs `kfg sys cache store ctx7.steps.install`
- **THEN** the cache entry SHALL be written to `$KFG_STORE_DIR/cache/<sha256(cache-version + ctx7.steps.install)>/`
- **AND** the `metadata.yaml` SHALL contain `stepRefName: ctx7.steps.install`
- **AND** changing the cache version SHALL prevent entries from an incompatible identity namespace from producing cache hits

#### Scenario: Store with no output
- **WHEN** the JSON has no `output` field or `output` is null
- **THEN** the `metadata.yaml` SHALL NOT contain an `output` section

#### Scenario: Store with spaces in artifact paths
- **WHEN** an artifact path contains spaces (e.g., `my dir/file.txt`)
- **THEN** the path SHALL be preserved correctly through the store cycle
