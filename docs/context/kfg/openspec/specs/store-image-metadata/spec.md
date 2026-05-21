# Store Image Metadata Specification

## Purpose

Define image metadata structure and digest computation for stored images.

## Requirements

### Requirement: Image digest computation

The system SHALL compute deterministic SHA256 digest for each image.

#### Scenario: Digest from recipe and files
- **WHEN** image is built
- **THEN** system computes digest as SHA256 of (recipe + sorted files)

#### Scenario: Digest determinism
- **WHEN** same Imagefile and sources are built twice
- **THEN** digests are identical

#### Scenario: Digest uniqueness
- **WHEN** any file or recipe changes
- **THEN** digest changes

#### Scenario: Nix32 encoding
- **WHEN** encoding hash bytes to string
- **THEN** encoding uses alphabet `0123456789abcdfghijklmnpqrsvwxyz`
- **AND** output is 32 characters

### Requirement: Metadata file structure

The system SHALL persist metadata alongside image files.

#### Scenario: Metadata location
- **WHEN** image build completes
- **THEN** system creates `metadata.json` in candidate directory
- **AND** after push, metadata exists at `$KFG_STORE_DIR/images/<name>/<tag>/metadata.json`

#### Scenario: Metadata content
- **WHEN** querying image metadata
- **THEN** includes: name, tag, image_digest, created_at, source_images[], recipe{}, files{}

#### Scenario: Metadata from BuildResult
- **WHEN** creating metadata
- **THEN** system derives metadata from BuildResult object
- **AND** includes: name, tag, digest, recipe, files, source_images

#### Scenario: Required fields
- **WHEN** storing metadata
- **THEN** validates: name, tag, image_digest, source_images, recipe

### Requirement: Recipe persistence

The system SHALL store original Imagefile for reproducibility.

#### Scenario: Recipe storage
- **WHEN** build completes
- **THEN** records original Imagefile content in metadata

#### Scenario: Recipe format
- **WHEN** retrieving recipe
- **THEN** includes: source_path, content, format ("imagefile.v1")

#### Scenario: Recipe retrieval
- **WHEN** executing `inspect --recipe`
- **THEN** outputs original Imagefile alongside metadata

### Requirement: Source image tracking

The system SHALL record parent image references and digests.

#### Scenario: Source image array
- **WHEN** image is built from stages
- **THEN** metadata includes `source_images` array

#### Scenario: Source entry structure
- **WHEN** stage references parent image
- **THEN** entry includes: ref, resolved_digest

#### Scenario: Multiple source tracking
- **WHEN** image uses multiple `FROM` stages
- **THEN** each source tracked independently

### Requirement: File manifest

The system SHALL provide detailed file listing.

#### Scenario: File manifest generation
- **WHEN** image is built
- **THEN** records all files in `files` object

#### Scenario: File count reporting
- **WHEN** inspecting image
- **THEN** reports total file count and aggregate size

### Requirement: Timestamp recording

The system SHALL record image creation time.

#### Scenario: Created timestamp
- **WHEN** build completes
- **THEN** records `created_at` as ISO 8601 timestamp

### Requirement: Human-readable inspection

The system SHALL present metadata in accessible formats.

#### Scenario: Inspect table output
- **WHEN** executing `inspect` (text mode)
- **THEN** outputs formatted table

#### Scenario: JSON inspection
- **WHEN** executing `inspect --json`
- **THEN** outputs full metadata JSON

#### Scenario: Recipe display
- **WHEN** executing `inspect --recipe`
- **THEN** outputs only Imagefile content without metadata

#### Scenario: Files listing
- **WHEN** executing `inspect --files`
- **THEN** outputs one file path per line, sorted alphabetically
- **AND** output contains no metadata or headers

### Requirement: Cmd derivation hash

The system SHALL compute hash for resolved Cmd derivations (v1 store).

#### Scenario: Hash includes cmd name
- **WHEN** resolver processes Cmd
- **THEN** hash input includes cmd name

#### Scenario: Hash includes step run scripts
- **WHEN** Cmd depends on Steps
- **THEN** hash input includes step run scripts

#### Scenario: Hash excludes runtime values
- **WHEN** computing derivation hash
- **THEN** excludes: step outputs, env vars, file contents

#### Scenario: Hash deterministic
- **WHEN** same manifest resolved multiple times
- **THEN** hash identical each time