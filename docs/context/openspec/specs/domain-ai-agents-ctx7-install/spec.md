# Domain: AI Agents - ctx7 Install Specification

## Purpose

This specification defines the ctx7 CLI install step contract, including flag validation, output directory handling, and MCP asset aggregation.
## Requirements

### Requirement: ctx7 install step defaults

The step `kfg.extension.ctx7.install` SHALL use an explicit extension-install contract. The step SHALL validate that `FLAGS` and `OUTPUT_DIR` are set before execution. The step SHALL execute the ctx7 CLI setup flow requested by `FLAGS`, copy generated skills into `OUTPUT_DIR`, and expose the canonical Context7 MCP asset through `kfg.extension.ctx7.mcp` for overlays.

#### Scenario: Install with explicit output directory
- **WHEN** `FLAGS` is `--claude --yes` and `OUTPUT_DIR` is `.claude/skills/`
- **THEN** the step runs `ctx7 setup --cli --project --claude --yes`
- **AND** it copies generated skills into `.claude/skills/`

#### Scenario: Missing required env var
- **WHEN** `OUTPUT_DIR` is empty or unset
- **THEN** the step logs an error and exits with code `1`

#### Scenario: Opencode-style install copies from temporary home
- **WHEN** `FLAGS` contains `--opencode`
- **THEN** the step copies generated skills from the temporary install home into `OUTPUT_DIR`
- **AND** it removes the temporary install directory after copying

#### Scenario: Canonical MCP asset is aggregated
- **WHEN** an overlay aggregates Context7 MCP configuration
- **THEN** it SHALL use `kfg.extension.ctx7.mcp`
- **AND** it SHALL NOT require `kfg.extension.self.mcp.context7` for the normalized flow
