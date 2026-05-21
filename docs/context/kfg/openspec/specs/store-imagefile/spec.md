# Store Imagefile Specification

## Purpose

Define Imagefile format and multi-stage composition for building configuration images.

## Requirements

### Requirement: Imagefile parsing

The system SHALL parse Imagefile manifests with declarative instructions.

#### Scenario: Valid Imagefile
- **WHEN** parsing valid Imagefile with `FROM`, `COPY`, `ENV`, `RUN`, `TAG`
- **THEN** produces abstract syntax tree without errors

#### Scenario: Invalid instruction
- **WHEN** parsing unknown instruction (e.g., `INVALID command`)
- **THEN** rejects and reports line number

#### Scenario: Case-insensitive keywords
- **WHEN** parsing `from`, `From`, `FROM`
- **THEN** normalizes and accepts all variants

### Requirement: FROM instruction

The system SHALL support `FROM` for base images or stages.

#### Scenario: FROM with image reference
- **WHEN** parsing `FROM claude-base:v2`
- **THEN** records stage name and image reference

#### Scenario: FROM scratch
- **WHEN** parsing `FROM scratch`
- **THEN** creates empty stage

#### Scenario: FROM with alias
- **WHEN** parsing `FROM claude-base:v2 AS claude`
- **THEN** assigns stage name `claude`

#### Scenario: Single FROM per stage
- **WHEN** multiple `FROM` in same stage
- **THEN** rejects Imagefile

### Requirement: COPY instruction

The system SHALL support `COPY` for file transfer.

#### Scenario: COPY from workspace
- **WHEN** parsing `COPY docs/AGENTS.md AGENTS.md`
- **THEN** records source and destination

#### Scenario: COPY from prior stage
- **WHEN** parsing `COPY --from=base .claude/ .claude/`
- **THEN** records stage reference

#### Scenario: COPY with wildcard
- **WHEN** parsing `COPY docs/*.md ./`
- **THEN** records glob pattern

#### Scenario: COPY directory
- **WHEN** parsing `COPY .claude/ .claude/`
- **THEN** recursively copies directory

#### Scenario: Duplicate COPY
- **WHEN** two COPY instructions to same path
- **THEN** accepts both (last write wins)

### Requirement: COPY destination semantics

The system SHALL follow Docker-compatible semantics for COPY destination paths.

#### Scenario: COPY to current directory (dot)
- **WHEN** parsing `COPY file.txt .`
- **THEN** copies `file.txt` to current working directory as `file.txt`
- **AND** preserves source filename

#### Scenario: COPY to directory (trailing slash)
- **WHEN** parsing `COPY file.txt target/`
- **THEN** copies `file.txt` into `target/` directory as `target/file.txt`
- **AND** trailing `/` marks destination as directory

#### Scenario: COPY with rename (no trailing slash)
- **WHEN** parsing `COPY file.txt newname.txt`
- **THEN** copies `file.txt` as `newname.txt` (rename operation)
- **AND** destination becomes final filename

#### Scenario: COPY multiple sources to directory
- **WHEN** parsing `COPY a.txt b.txt target/`
- **THEN** copies both files into `target/` directory
- **AND** destination MUST end with `/` (directory required)

#### Scenario: COPY explicit current directory
- **WHEN** parsing `COPY file.txt ./`
- **THEN** copies `file.txt` to current directory as `file.txt`
- **AND** `./` explicitly marks current directory

#### Scenario: COPY with --from preserves semantics
- **WHEN** parsing `COPY --from=base file.txt .`
- **THEN** copies from stage with same destination semantics
- **AND** `.` resolves to current directory in final stage

### Requirement: ENV instruction

The system SHALL support `ENV` for environment variables.

#### Scenario: ENV key=value
- **WHEN** parsing `ENV DEBUG=true TARGET=AGENTS.md`
- **THEN** records key-value pairs

#### Scenario: ENV with spaces
- **WHEN** parsing `ENV MESSAGE="Hello World"`
- **THEN** parses quoted values

#### Scenario: Multiple ENV
- **WHEN** parsing multiple ENV instructions
- **THEN** merges all variables

### Requirement: RUN instruction

The system SHALL support `RUN` for build-time commands.

#### Scenario: RUN with shell
- **WHEN** parsing `RUN sh -c 'cat base.md override.md > AGENTS.md'`
- **THEN** records shell command

#### Scenario: RUN inherits ENV
- **WHEN** executing RUN after ENV
- **THEN** exports ENV variables to subshell

#### Scenario: RUN multiline
- **WHEN** parsing RUN with backslash continuation
- **THEN** joins lines into single command

#### Scenario: RUN in WORKDIR context
- **WHEN** executing RUN after WORKDIR instruction
- **THEN** executes command in working directory context
- **AND** uses WORKDIR as command's working directory

### Requirement: WORKDIR instruction

The system SHALL support `WORKDIR` for setting working directory.

#### Scenario: WORKDIR with absolute path
- **WHEN** parsing `WORKDIR /app`
- **THEN** sets working directory to `/app` (normalized to `app`)
- **AND** creates directory if it doesn't exist

#### Scenario: WORKDIR with relative path
- **WHEN** parsing `WORKDIR subdir`
- **THEN** sets working directory relative to current WORKDIR
- **AND** chains with previous WORKDIR if set

#### Scenario: WORKDIR affects COPY destination
- **WHEN** COPY after WORKDIR with `.` destination
- **THEN** resolves `.` against WORKDIR
- **AND** copies file into working directory

#### Scenario: WORKDIR chained
- **WHEN** multiple WORKDIR instructions
- **THEN** each relative WORKDIR chains on previous
- **AND** absolute WORKDIR resets working directory

### Requirement: Glob pattern expansion

The system SHALL expand glob patterns in COPY sources.

#### Scenario: Star pattern expansion
- **WHEN** parsing `COPY *.txt ./`
- **THEN** expands to all matching `.txt` files
- **AND** copies each file individually

#### Scenario: Question mark pattern
- **WHEN** parsing `COPY file?.txt ./`
- **THEN** expands to files matching single character

#### Scenario: Directory glob
- **WHEN** parsing `COPY docs/*.md ./`
- **THEN** expands files in directory matching pattern

#### Scenario: Glob zero match error
- **WHEN** glob pattern matches zero files
- **THEN** fails build with error
- **AND** reports pattern that matched nothing

#### Scenario: Multiple expanded sources
- **WHEN** glob expands to multiple files
- **THEN** destination treated as directory
- **AND** each file copied with preserved basename

### Requirement: TAG instruction

The system SHALL support `TAG` for name and tag assignment.

#### Scenario: TAG assignment
- **WHEN** parsing `TAG my-image:v1.0`
- **THEN** records image name and tag

#### Scenario: TAG format
- **WHEN** parsing `TAG <name>:<tag>`
- **THEN** validates format and splits

#### Scenario: TAG only in final stage
- **WHEN** parsing TAG in intermediate stage
- **THEN** rejects Imagefile

### Requirement: Multi-stage composition

The system SHALL support multi-stage builds.

#### Scenario: Multiple FROM stages
- **WHEN** Imagefile contains multiple FROM
- **THEN** builds each stage independently

#### Scenario: COPY --from stage reference
- **WHEN** stage uses `COPY --from=<stage>`
- **THEN** uses files from specified stage

#### Scenario: Last write wins
- **WHEN** multiple stages write to same path
- **THEN** final stage's version used

#### Scenario: Directory merge
- **WHEN** stages copy different files to same directory
- **THEN** all files coexist

### Requirement: Environment scope

The system SHALL manage ENV across stages.

#### Scenario: ENV per stage
- **WHEN** stage declares ENV
- **THEN** available to RUN in that stage only

#### Scenario: ENV not inherited
- **WHEN** stage A declares ENV, stage B starts fresh
- **THEN** stage B does not inherit

### Requirement: Intermediate stage handling

The system SHALL NOT persist intermediate stages.

#### Scenario: Only final image persisted
- **WHEN** multi-stage build completes
- **THEN** only final image pushed to store

#### Scenario: Intermediate discarded
- **WHEN** final stage copies from intermediate
- **THEN** intermediate stages discarded after build

### Requirement: Source image resolution

The system SHALL load parent images from store.

#### Scenario: Image lookup
- **WHEN** FROM specifies `image:tag`
- **THEN** looks up in store and loads files

#### Scenario: Parent not found
- **WHEN** parent image doesn't exist
- **THEN** fails with error

#### Scenario: Digest tracking
- **WHEN** stage loads parent image
- **THEN** records resolved digest in metadata