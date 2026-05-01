# Store Image Build Specification

## Purpose

Specifies the image build process that creates candidate directories from Imagefiles, resolves multi-stage builds, and computes deterministic digests.

## Requirements

### Requirement: Build candidate creation
The system SHALL create a local candidate image directory from an Imagefile without modifying the store.

#### Scenario: Build with default Imagefile
- **WHEN** executing `kfg store image build` in directory with `./Imagefile`
- **THEN** system parses Imagefile, resolves stages, copies files, and outputs candidate path (e.g., `$TMPDIR/<name>/<tag>`)
- **AND** removes any existing build directory before creating new one

#### Scenario: Build with custom Imagefile path
- **WHEN** executing `kfg store image build -f ./Custom.dockerfile`
- **THEN** system uses specified file instead of default `./Imagefile`

#### Scenario: Build output in specified directory
- **WHEN** executing `kfg store image build --output /tmp/mybuild`
- **THEN** system creates candidate at `/tmp/mybuild` instead of default temp directory
- **AND** removes existing `/tmp/mybuild` if present before build

#### Scenario: Build with custom root
- **WHEN** executing `kfg store image build --root /path/to/project`
- **THEN** system resolves file paths relative to specified root instead of current directory

#### Scenario: Build directory cleanup
- **WHEN** build directory already exists from previous build
- **THEN** system removes entire directory tree before creating new build
- **AND** ensures deterministic builds without stale file interference

### Requirement: Stage resolution
The system SHALL resolve multi-stage builds by loading referenced images and composing files.

#### Scenario: FROM with store image reference
- **WHEN** build encounters `FROM claude-base:v2`
- **THEN** system looks up `claude-base:v2` in store, reads its files, and uses as stage context

#### Scenario: Stage composition from multiple images
- **WHEN** Imagefile references multiple `FROM` stages
- **THEN** system loads each image independently and composes files in order

#### Scenario: FROM scratch (empty stage)
- **WHEN** encountering `FROM scratch`
- **THEN** system creates empty stage (no files, no environment)

#### Scenario: Image not found in store
- **WHEN** `FROM` references non-existent image
- **THEN** system fails with error identifying missing image and suggesting available images

### Requirement: File composition and conflict handling
The system SHALL combine files from multiple stages into single candidate directory.

#### Scenario: Last write wins for duplicate paths
- **WHEN** multiple stages copy to same destination path
- **THEN** final stage's version is used in candidate

#### Scenario: Directory merge
- **WHEN** stages copy different files to same directory
- **THEN** all files are merged (non-conflicting files coexist)

#### Scenario: Recursive directory copy
- **WHEN** copying directory with subdirectories and files
- **THEN** entire tree is preserved in candidate

### Requirement: Build-time command execution
The system SHALL execute `RUN` instructions with environment variables and file context.

#### Scenario: RUN with environment variables
- **WHEN** executing `RUN` after prior `ENV` declarations
- **THEN** system exports ENV variables to RUN subshell and executes command

#### Scenario: RUN failure handling
- **WHEN** `RUN` command exits with non-zero code
- **THEN** system fails build and leaves candidate directory for inspection

#### Scenario: RUN success
- **WHEN** `RUN` command completes successfully
- **THEN** system continues to next instruction

#### Scenario: RUN with file access
- **WHEN** `RUN` command reads/writes files in candidate
- **THEN** system has full filesystem access to build context

### Requirement: RUN output logging
The system SHALL capture RUN command stdout and stderr via structured logging.

#### Scenario: RUN stdout logging
- **WHEN** RUN command produces stdout output
- **THEN** system captures each line via logger with component `build:run:stdout`
- **AND** uses log level `detail` (visible with KFG_VERBOSE=2)

#### Scenario: RUN stderr logging
- **WHEN** RUN command produces stderr output
- **THEN** system captures each line via logger with component `build:run:stderr`
- **AND** uses log level `warn`

#### Scenario: RUN failure logging
- **WHEN** RUN command fails with non-zero exit code
- **THEN** system logs error with component `build:run`
- **AND** includes exit code and command in message

#### Scenario: JSONL format compliance
- **WHEN** logging RUN output
- **THEN** each line is valid JSONL with fields: ts, level, component, msg, source, pid
- **AND** source="go" for all build operations

#### Scenario: Verbose level control
- **WHEN** KFG_VERBOSE=0 (default)
- **THEN** RUN stdout/stderr detail logs are suppressed in console output
- **AND** logs are still written to JSONL file if KFG_LOG_FILE is set

### Requirement: Build output structure
The system SHALL produce a well-defined directory structure for candidates.

#### Scenario: Candidate directory layout
- **WHEN** build succeeds
- **THEN** candidate contains all composed files in `artifacts/` subdirectory
- **AND** includes `metadata.json` at candidate root
- **AND** maintains `stages/` directory for stage isolation

#### Scenario: File permissions preserved
- **WHEN** copying files from stages
- **THEN** execute bits and other permissions are preserved in candidate

### Requirement: Digest computation
The system SHALL compute a deterministic SHA256 digest of the build result.

#### Scenario: Digest calculation
- **WHEN** build completes
- **THEN** system computes SHA256 of sorted(recipe + all files) and records as `image_digest`

#### Scenario: Digest reproducibility
- **WHEN** same Imagefile and sources build twice
- **THEN** resulting digests are identical
- **AND** cleanup of existing directories ensures no stale files affect hash

### Requirement: Stage isolation
The system SHALL maintain isolation between stages in multi-stage builds.

#### Scenario: Stage directories are isolated
- **WHEN** multi-stage build processes multiple stages
- **THEN** each stage has its own isolated directory
- **AND** files from previous stages do NOT leak into current stage

#### Scenario: COPY --from uses source stage directory
- **WHEN** `COPY --from=base` copies files from base stage
- **THEN** system uses base stage's isolated directory as source
- **AND** files are NOT copied from workspace or other stages

#### Scenario: FROM scratch creates empty stage
- **WHEN** stage starts with `FROM scratch`
- **THEN** stage directory contains no files
- **AND** workspace files do NOT appear in stage

#### Scenario: Final stage artifacts only
- **WHEN** build completes
- **THEN** artifacts directory contains only final stage files
- **AND** intermediate stage files are discarded

### Requirement: Build status and logging
The system SHALL provide clear feedback during build execution.

#### Scenario: Build progress reporting
- **WHEN** executing build
- **THEN** system logs each stage resolution and file composition step

#### Scenario: Build error reporting
- **WHEN** build fails
- **THEN** system reports specific error (e.g., file not found, command failed) with context

#### Scenario: Build completion
- **WHEN** build succeeds
- **THEN** system reports candidate path for user review before push

### Requirement: Structured logging format
The system SHALL use JSONL structured logging for build operations.

#### Scenario: JSONL log format
- **WHEN** logging build events
- **THEN** outputs JSONL with fields: ts, level, component, msg, source, pid
- **AND** Go code sets source="go" and component prefixed with "core:""

#### Scenario: Build stage logging
- **WHEN** processing stage
- **THEN** logs with component "build:stage" and stage details

#### Scenario: COPY logging
- **WHEN** copying files
- **THEN** logs with component "build:copy" and file details

#### Scenario: WORKDIR logging
- **WHEN** setting working directory
- **THEN** logs with component "build:workdir" and path details

#### Scenario: RUN logging
- **WHEN** executing RUN command
- **THEN** logs with component "build:run" and command details

#### Scenario: Workspace operations logging
- **WHEN** materializing workspace
- **THEN** logs with component "workspace:start" or "workspace:stop"
- **AND** includes file counts and backup status

#### Scenario: Verbose level control
- **WHEN** KFG_VERBOSE environment variable is set
- **THEN** controls which log levels are displayed in human output
- **AND** JSONL file always receives all events

### Requirement: Metadata creation at build completion

The system SHALL create metadata.json in the candidate directory after successful build.

#### Scenario: Metadata creation from BuildResult
- **WHEN** build succeeds
- **THEN** system creates ImageMetadata from BuildResult
- **AND** includes: name, tag, digest, files, recipe

#### Scenario: Metadata validation before save
- **WHEN** creating metadata
- **THEN** system validates metadata fields
- **AND** fails build if validation fails

#### Scenario: Metadata saved to candidate
- **WHEN** metadata validation passes
- **THEN** system writes metadata.json to candidate directory root
- **AND** push command can immediately consume candidate

#### Scenario: Build → push workflow
- **WHEN** user runs build then push
- **THEN** metadata.json exists from build
- **AND** push succeeds without metadata creation

### Requirement: Auto push after build
The system SHALL support automatic push to store after successful build via `--push` flag.

#### Scenario: Build with push flag
- **WHEN** executing `kfg store image build --push`
- **THEN** system builds image and automatically pushes to store if build succeeds
- **AND** reports both build success and push success

#### Scenario: Build with push fails on existing image
- **WHEN** executing `kfg store image build --push` for image that already exists in store
- **THEN** system fails push phase with immutability error
- **AND** leaves build candidate intact for inspection

#### Scenario: Build with push and keep-build
- **WHEN** executing `kfg store image build --push --keep-build`
- **THEN** system pushes to store and preserves build candidate directory
- **AND** reports candidate path for user reference

#### Scenario: Build failure before push
- **WHEN** executing `kfg store image build --push` and build fails
- **THEN** system reports build error and does not attempt push
- **AND** leaves failed build candidate for debugging