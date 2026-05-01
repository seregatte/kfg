# Store Image Persistence Specification

## Purpose

Specifies image persistence operations: pushing built candidates to the store, listing and inspecting images with metadata, and removing stored images with safety checks.

## Requirements

### Requirement: Image push to store
The system SHALL persist built image candidates to the store with metadata.

#### Scenario: Successful push
- **WHEN** executing `kfg store image push /path/to/candidate`
- **THEN** system verifies metadata.json exists, writes to store, and cleans up build directory
- **AND** metadata was created by build command

#### Scenario: Push with existing image name:tag
- **WHEN** pushing image with name and tag that already exist in store
- **THEN** system fails with error (immutable images; prevents overwrites)

#### Scenario: Push preserves candidate on flag
- **WHEN** executing `kfg store image push <build-dir> --keep-build`
- **THEN** system persists to store but leaves build directory intact

#### Scenario: Candidate structure required
- **WHEN** pushing directory missing metadata.json
- **THEN** system fails with error "metadata.json not found"
- **AND** suggests running build first

#### Scenario: Build → push workflow
- **WHEN** user runs `kfg store image build --name myconfig:latest`
- **THEN** build creates candidate with files AND metadata.json
- **AND** push succeeds immediately: `kfg store image push <candidate>`

#### Scenario: Metadata pre-existing in candidate
- **WHEN** push is executed on valid candidate directory
- **THEN** system finds metadata.json created by build command
- **AND** no metadata creation happens during push

### Requirement: Image listing
The system SHALL enable querying stored images with optional filtering.

#### Scenario: List all images
- **WHEN** executing `kfg store image list`
- **THEN** system outputs human-readable table with columns: NAME, TAG, DIGEST (first 12 chars), CREATED

#### Scenario: List with JSON output
- **WHEN** executing `kfg store image list --json`
- **THEN** system outputs JSON array of image objects with full metadata
- **AND** empty store outputs `[]` (valid JSON)

#### Scenario: List with no images
- **WHEN** store is empty
- **THEN** system outputs "No images found" in human format
- **OR** outputs empty JSON array `[]` when `--json` flag is set

#### Scenario: List sorting
- **WHEN** multiple images exist
- **THEN** system sorts by name, then by tag

### Requirement: Image inspection
The system SHALL provide detailed metadata access for debugging.

#### Scenario: Inspect image metadata
- **WHEN** executing `kfg store image inspect my-config:v1`
- **THEN** system outputs metadata (name, tag, digest, created date, source images)

#### Scenario: Inspect with JSON output
- **WHEN** executing `kfg store image inspect my-config:v1 --json`
- **THEN** system outputs full metadata JSON structure

#### Scenario: Inspect with recipe
- **WHEN** executing `kfg store image inspect my-config:v1 --recipe`
- **THEN** system outputs original Imagefile content alongside metadata

#### Scenario: Inspect non-existent image
- **WHEN** inspecting image that doesn't exist
- **THEN** system fails with helpful error (suggests similar names)

#### Scenario: Image reference resolution
- **WHEN** specifying image as `name:tag` or `name:latest`
- **THEN** system resolves to exact match or defaults to `:latest` if tag omitted

### Requirement: Flag conflict validation
The system SHALL reject conflicting output format flags with clear error message.

#### Scenario: Recipe and files conflict
- **WHEN** executing `kfg store image inspect myimage:v1 --recipe --files`
- **THEN** system fails with error "flags --recipe and --files are mutually exclusive"
- **AND** exit code is non-zero

#### Scenario: Recipe and json conflict
- **WHEN** executing `kfg store image inspect myimage:v1 --recipe --json`
- **THEN** system fails with error "flags --recipe and --json are mutually exclusive"
- **AND** exit code is non-zero

#### Scenario: Files and json allowed
- **WHEN** executing `kfg store image inspect myimage:v1 --files --json`
- **THEN** system succeeds and outputs JSON file array
- **AND** no validation error occurs

#### Scenario: Single output flag succeeds
- **WHEN** executing with only one output flag (--recipe, --files, or --json)
- **THEN** system succeeds with expected output format
- **AND** no validation error occurs

### Requirement: Image removal
The system SHALL allow deletion of stored images with safety.

#### Scenario: Remove image
- **WHEN** executing `kfg store image remove my-config:v1`
- **THEN** system deletes image and confirms removal

#### Scenario: Remove non-existent image
- **WHEN** trying to remove image that doesn't exist
- **THEN** system fails with error (no silent success)

#### Scenario: Remove image in use (warning)
- **WHEN** removing image that is referenced by active workspace instance
- **THEN** system warns user before proceeding (does not prevent removal)

#### Scenario: Alias command
- **WHEN** executing `kfg store image rm <name:tag>` (short form)
- **THEN** system treats as equivalent to `remove`

### Requirement: Store organization
The system SHALL maintain well-defined directory structure for images.

#### Scenario: Image storage location
- **WHEN** image is pushed
- **THEN** system stores at `$KFG_STORE_DIR/images/<name>/<tag>/` with files and metadata

#### Scenario: Metadata file structure
- **WHEN** querying stored image
- **THEN** system finds `metadata.json` alongside image files

### Requirement: Tag resolution
The system SHALL handle image references with reasonable defaults.

#### Scenario: Explicit tag
- **WHEN** specifying `name:tag`
- **THEN** system uses exact tag

#### Scenario: Implicit latest tag
- **WHEN** specifying `name` (no tag)
- **THEN** system defaults to `name:latest`

#### Scenario: Latest tag semantics
- **WHEN** no image with `:latest` tag exists
- **THEN** system treats `:latest` as simple default string (no special resolution logic)

### Requirement: Store consistency
The system SHALL maintain metadata consistency and validity.

#### Scenario: Metadata validation on list/inspect
- **WHEN** accessing stored images
- **THEN** system validates metadata format; reports corruption if found

#### Scenario: File manifest accuracy
- **WHEN** inspecting image
- **THEN** system reports accurate file count and sizes