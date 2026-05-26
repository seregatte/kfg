## ADDED Requirements

### Requirement: Cache store subcommand
The `kfg sys cache store` subcommand SHALL persist Step execution results (artifacts and output) to the cache, reading structured input from stdin.

#### Scenario: Store cache entry with artifacts and output
- **WHEN** user runs `kfg sys cache store <step-ref> --workdir <path>` with JSON on stdin containing `before`, `after`, `declarative`, and `output` fields
- **THEN** the subcommand SHALL compute the artifact delta (after - before), merge with declarative artifacts, perform fs diff on workdir to detect unregistered artifacts, copy all artifacts to the cache entry directory, write `metadata.yaml`, and commit atomically via rename

#### Scenario: Store reads JSON from stdin
- **WHEN** `kfg sys cache store` is invoked
- **THEN** it SHALL read a JSON object from stdin with fields: `before` (string array), `after` (string array), `declarative` (string array), and `output` (object with `name` string and `value` base64 string)
- **AND** all arrays contain relative paths from the workdir

#### Scenario: Store computes artifact delta
- **WHEN** `before` contains `["a.txt"]` and `after` contains `["a.txt", "b.txt"]`
- **THEN** the delta SHALL be `["b.txt"]`
- **AND** only `b.txt` plus any declarative artifacts and fs-diff-detected artifacts SHALL be cached

#### Scenario: Store performs fs diff on workdir
- **WHEN** the step creates files in workdir that are not in the `after` list
- **THEN** the subcommand SHALL detect those files via filesystem diff and include them in the cache entry

#### Scenario: Store uses StepReference.name as directory key
- **WHEN** `kfg sys cache store ctx7.steps.install` is invoked
- **THEN** the cache entry SHALL be written to `$KFG_STORE_DIR/cache/<sha256("ctx7.steps.install")>/`
- **AND** the `metadata.yaml` SHALL contain `stepRefName: ctx7.steps.install`

#### Scenario: Store with no output
- **WHEN** the JSON has no `output` field or `output` is null
- **THEN** the `metadata.yaml` SHALL NOT contain an `output` section

#### Scenario: Store with spaces in artifact paths
- **WHEN** an artifact path contains spaces (e.g., `my dir/file.txt`)
- **THEN** the path SHALL be preserved correctly through the store cycle

### Requirement: Cache restore subcommand
The `kfg sys cache restore` subcommand SHALL restore cached Step results and emit shell eval-safe output.

#### Scenario: Restore emits shell eval-safe output
- **WHEN** user runs `kfg sys cache restore <step-ref> --workdir <path>`
- **THEN** the subcommand SHALL emit lines to stdout that, when eval'd in bash, call `__kfg_add_artifact` for each cached artifact and `__kfg_output_set` for the cached output

#### Scenario: Restore restores artifacts to original relative paths
- **WHEN** a cached artifact has relative path `.pi/skills/test/file.txt`
- **THEN** the emitted `__kfg_add_artifact` call SHALL use that same relative path
- **AND** the file SHALL be copied from cache to `<workdir>/.pi/skills/test/file.txt`

#### Scenario: Restore restores output value
- **WHEN** a cached entry has output with name `ctx7_context` and base64-encoded value
- **THEN** the emitted `__kfg_output_set` call SHALL use the decoded value

#### Scenario: Restore with no output
- **WHEN** the cached entry has no output
- **THEN** the emitted lines SHALL contain only `__kfg_add_artifact` calls

#### Scenario: Restore with spaces in paths
- **WHEN** cached artifact paths contain spaces
- **THEN** the emitted `__kfg_add_artifact` calls SHALL quote the paths correctly

### Requirement: Cache exists subcommand
The `kfg sys cache exists` subcommand SHALL check for a valid cache entry.

#### Scenario: Cache hit
- **WHEN** a valid cache entry exists for the given StepReference.name (directory and `metadata.yaml` present)
- **THEN** the subcommand SHALL exit with code 0

#### Scenario: Cache miss
- **WHEN** no valid cache entry exists for the given StepReference.name
- **THEN** the subcommand SHALL exit with code 1

### Requirement: Cache list subcommand
The `kfg sys cache ls` subcommand SHALL list all cache entries with metadata.

#### Scenario: List entries in table format
- **WHEN** user runs `kfg sys cache ls`
- **THEN** the subcommand SHALL display each entry with STEP REF NAME, TIMESTAMP, and SIZE columns

#### Scenario: List entries in JSON format
- **WHEN** user runs `kfg sys cache ls --json`
- **THEN** the subcommand SHALL output a JSON array of objects with `stepRef`, `timestamp`, and `size` fields

#### Scenario: List entries in YAML format
- **WHEN** user runs `kfg sys cache ls --yaml`
- **THEN** the subcommand SHALL output a YAML list of objects with `stepRef`, `timestamp`, and `size` fields

### Requirement: Cache inspect subcommand
The `kfg sys cache inspect` subcommand SHALL show detailed metadata for a cache entry.

#### Scenario: Inspect entry by StepReference.name
- **WHEN** user runs `kfg sys cache inspect <step-ref>`
- **THEN** the subcommand SHALL display the entry identified by the StepReference.name (not a hash ID)

#### Scenario: Inspect output in YAML format
- **WHEN** user runs `kfg sys cache inspect <step-ref>`
- **THEN** the default output SHALL be YAML containing `stepRef`, `timestamp`, `size`, `artifacts` (list), and `output` (with `name` and `value`)
- **AND** the output value SHALL be displayed completely (no truncation)

#### Scenario: Inspect output in JSON format
- **WHEN** user runs `kfg sys cache inspect <step-ref> --json`
- **THEN** the subcommand SHALL output JSON with the same fields

#### Scenario: Inspect nonexistent entry
- **WHEN** user runs `kfg sys cache inspect <step-ref>` for an entry that does not exist
- **THEN** the subcommand SHALL exit with a non-zero code and display an error message

### Requirement: Cache remove subcommand
The `kfg sys cache rm` subcommand SHALL remove cache entries by StepReference.name.

#### Scenario: Remove single entry
- **WHEN** user runs `kfg sys cache rm <step-ref>`
- **THEN** the subcommand SHALL remove the cache entry for that StepReference.name

#### Scenario: Remove multiple entries
- **WHEN** user runs `kfg sys cache rm <ref1> <ref2> ...`
- **THEN** the subcommand SHALL remove all specified entries

#### Scenario: Remove nonexistent entry
- **WHEN** user runs `kfg sys cache rm <step-ref>` for an entry that does not exist
- **THEN** the subcommand SHALL display a warning and continue with remaining entries

### Requirement: Cache prune subcommand
The `kfg sys cache prune` subcommand SHALL remove old cache entries.

#### Scenario: Prune entries older than 30 days
- **WHEN** user runs `kfg sys cache prune`
- **THEN** the subcommand SHALL remove entries with timestamps older than 30 days
- **AND** SHALL display which entries were pruned and total freed bytes

#### Scenario: Prune with JSON output
- **WHEN** user runs `kfg sys cache prune --json`
- **THEN** the subcommand SHALL output JSON with `pruned` (list of step refs), `count`, and `freedBytes` fields

### Requirement: Cache disk usage subcommand
The `kfg sys cache du` subcommand SHALL report disk usage for cache entries.

#### Scenario: Show disk usage in table format
- **WHEN** user runs `kfg sys cache du`
- **THEN** the subcommand SHALL display the cache directory path, per-entry sizes, and total

#### Scenario: Show disk usage in JSON format
- **WHEN** user runs `kfg sys cache du --json`
- **THEN** the subcommand SHALL output JSON with `cacheDir`, `entries` (list with `stepRef` and `size`), and `totalBytes` fields

### Requirement: Output format flags
All cache admin subcommands SHALL support structured output flags.

#### Scenario: Default output is human-readable
- **WHEN** user runs any cache admin subcommand without format flags
- **THEN** the output SHALL be human-readable text (table or formatted text)

#### Scenario: JSON flag
- **WHEN** user passes `--json` to any cache admin subcommand
- **THEN** the output SHALL be valid JSON

#### Scenario: YAML flag
- **WHEN** user passes `--yaml` to any cache admin subcommand
- **THEN** the output SHALL be valid YAML
