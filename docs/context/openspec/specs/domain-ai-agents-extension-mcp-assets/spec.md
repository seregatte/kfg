# Domain: AI Agents - Extension MCP Assets Specification

## Purpose

This specification defines extension-owned MCP asset contracts, including asset naming, server commands, and aggregation conventions.
## Requirements

### Requirement: Extension-owned MCP assets are canonical

Each extension that exposes an MCP server SHALL define its canonical MCP asset under `base/extensions/<name>/assets/` using the manifest name `kfg.extension.<name>.mcp`. Overlays SHALL aggregate extension-owned MCP assets directly and SHALL NOT rely on transitional `kfg.extension.self.*` assets for normalized extensions.

#### Scenario: Playwright MCP asset is aggregated by extension namespace
- **WHEN** an overlay needs the Playwright MCP server definition
- **THEN** it aggregates `kfg.extension.playwright.mcp`
- **AND** it SHALL NOT aggregate `kfg.extension.self.mcp.playwright`

#### Scenario: Chrome DevTools MCP asset is aggregated by extension namespace
- **WHEN** an overlay needs the Chrome DevTools MCP server definition
- **THEN** it aggregates `kfg.extension.chrome-devtools.mcp`
- **AND** it SHALL NOT aggregate `kfg.extension.self.mcp.chrome-devtools`

### Requirement: MCP assets SHALL expose a stable aggregation contract

Each extension-owned MCP asset SHALL provide the fields required by `kfg.aggregate-mcp` and the MCP converters: `name`, `description`, `enabled`, `server.command`, `server.args`, and `server.env`.

#### Scenario: Asset can be converted for Claude
- **WHEN** `kfg.aggregate-mcp` converts `kfg.extension.chrome-devtools.mcp` with `kfg.convert.self.mcp.claude`
- **THEN** the conversion SHALL succeed using the asset's `name`, `server.command`, `server.args`, and `server.env`

#### Scenario: Asset can be converted for OpenCode
- **WHEN** `kfg.aggregate-mcp` converts `kfg.extension.playwright.mcp` with `kfg.convert.self.mcp.opencode`
- **THEN** the conversion SHALL succeed using the asset's `name`, `server.command`, `server.args`, and `server.env`

### Requirement: MCP assets SHALL preserve extension-specific server commands

The canonical MCP asset for each extension SHALL preserve the working server command observed in the normalized extension contract.

#### Scenario: Chrome DevTools server command
- **WHEN** the repository defines `kfg.extension.chrome-devtools.mcp`
- **THEN** the asset SHALL describe the local server command `npx -y chrome-devtools-mcp@latest`

#### Scenario: Playwright server command
- **WHEN** the repository defines `kfg.extension.playwright.mcp`
- **THEN** the asset SHALL describe the local server command `npx -y @playwright/mcp@latest --extension`

#### Scenario: Context7 server command
- **WHEN** the repository defines `kfg.extension.ctx7.mcp`
- **THEN** the asset SHALL describe the local server command for the Context7 MCP package used by the repository
