# Inspect Files List Specification

## Purpose

Provide machine-readable file listing for image inspection with support for both text and JSON output formats.

## Requirements

### Requirement: Files listing

The system SHALL list all image artifact paths as a simple newline-separated list.

#### Scenario: Files flag output format
- **WHEN** executing `inspect --files`
- **THEN** outputs one file path per line
- **AND** paths are sorted alphabetically
- **AND** output contains no metadata or headers

#### Scenario: Empty files list
- **WHEN** image has no files
- **THEN** outputs "No files" message

#### Scenario: Shell integration
- **WHEN** output is piped to shell commands
- **THEN** each line contains exactly one path
- **AND** paths are suitable for `while read` loops

### Requirement: JSON file listing

The system SHALL support JSON output format for `inspect --files` command.

#### Scenario: Files with JSON flag
- **WHEN** executing `kfg store image inspect myimage:v1 --files --json`
- **THEN** outputs JSON array of file paths
- **AND** paths are sorted alphabetically
- **AND** output is valid JSON parseable by jq or similar tools

#### Scenario: Empty files JSON output
- **WHEN** executing `inspect --files --json` on image with no files
- **THEN** outputs empty JSON array `[]`
- **AND** does NOT output "No files" text message

#### Scenario: Files JSON structure
- **WHEN** image has files CLAUDE.md, README.md, config.json
- **THEN** JSON output is `["CLAUDE.md", "README.md", "config.json"]`
- **AND** no metadata fields or extra structure included

#### Scenario: Files JSON parseable
- **WHEN** output is consumed by JSON parser
- **THEN** each path is a string element in array
- **AND** array is valid JSON for jq or similar tools