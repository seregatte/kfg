## ADDED Requirements

### Requirement: Ensure gitignore entries exist

The step `kfg.ensure-gitignore` SHALL add specified entries to `.gitignore`, creating the file if it does not exist. The step MUST be idempotent — entries already present SHALL NOT be duplicated.

#### Scenario: Add new entries to empty gitignore
- **WHEN** `GITIGNORE_ENTRIES` is set to `"/.opencode/ /.gemini/"` and `.gitignore` does not exist
- **THEN** the step creates `.gitignore` with both entries, one per line

#### Scenario: Skip existing entries
- **WHEN** `GITIGNORE_ENTRIES` is set to `"/.opencode/ /.gemini/"` and `.gitignore` already contains `"/.opencode/"`
- **THEN** the step adds only `"/.gemini/"` to `.gitignore`

#### Scenario: No entries configured
- **WHEN** `GITIGNORE_ENTRIES` is empty or unset
- **THEN** the step exits with success and makes no changes

#### Scenario: Gitignore file already complete
- **WHEN** `GITIGNORE_ENTRIES` is set and all entries already exist in `.gitignore`
- **THEN** the step exits with success and makes no changes

### Requirement: Configurable gitignore file path

The step SHALL use `GITIGNORE_FILE` env var to determine the target file path, defaulting to `.gitignore`.

#### Scenario: Custom gitignore path
- **WHEN** `GITIGNORE_FILE` is set to `.gitignore.dev` and `GITIGNORE_ENTRIES` is `"/.opencode/"`
- **THEN** the step creates or updates `.gitignore.dev` with the entry
