## ADDED Requirements

### Requirement: LifeOS overlay SHALL use the normalized extension workflow pattern

The LifeOS repository SHALL define its agent overlay under `~/Sites/lifeos/.manifests/overlay/lifeos/` and SHALL compose shared `kfg` steps, extension assets, and converters instead of referencing deprecated `kfg.core.steps.*` or `kfg.extension.steps.*` resources.

#### Scenario: LifeOS overlay entrypoint
- **WHEN** the LifeOS manifests are organized for the normalized migration
- **THEN** `~/Sites/lifeos/.manifests/overlay/lifeos/kustomization.yaml` SHALL reference the shared `kfg` base manifests and the overlay-local resources

#### Scenario: LifeOS workflow uses shared step names
- **WHEN** the LifeOS workflow invokes reusable behavior
- **THEN** it SHALL reference shared resources such as `kfg.ensure-gitignore`, `kfg.detect-agent`, `kfg.copy-context`, `kfg.convert`, `kfg.aggregate-mcp`, and `kfg.cleanup`
- **AND** it SHALL NOT reference `kfg.core.steps.*` or `kfg.extension.steps.*`

### Requirement: LifeOS-specific resources SHALL remain overlay-local assets

The LifeOS overlay SHALL model project-specific commands, subagents, and MCP definitions as local assets inside the overlay rather than embedding LifeOS-specific behavior into shared `kfg` base extensions.

#### Scenario: LifeOS commands are local assets
- **WHEN** the overlay needs the `prepare-discourse` or `wpp-elder` command definitions
- **THEN** it SHALL define them as overlay-local assets and convert them with shared `kfg.convert.self.command.*` converters

#### Scenario: LifeOS subagent is local asset
- **WHEN** the overlay needs the `elder` subagent definition
- **THEN** it SHALL define the subagent as an overlay-local asset and convert it with shared `kfg.convert.self.subagent.*` converters

#### Scenario: LifeOS MCP is local asset
- **WHEN** the overlay needs the `eld` MCP definition
- **THEN** it SHALL define the MCP as an overlay-local asset and aggregate it with shared MCP assets from normalized extensions

### Requirement: LifeOS overlay SHALL aggregate normalized extension MCP assets

The LifeOS overlay SHALL aggregate extension-owned MCP assets for `ctx7`, `chrome-devtools`, and `playwright`, plus the local `eld` MCP asset, for supported agents.

#### Scenario: Claude receives all LifeOS MCP assets
- **WHEN** the detected agent is `claude`
- **THEN** the overlay SHALL aggregate `lifeos.mcp.eld`, `kfg.extension.ctx7.mcp`, `kfg.extension.chrome-devtools.mcp`, and `kfg.extension.playwright.mcp` into the Claude MCP target

#### Scenario: Pi excludes MCP aggregation
- **WHEN** the detected agent is `pi`
- **THEN** the overlay SHALL NOT aggregate MCP assets for Pi unless the overlay explicitly defines a Pi-compatible MCP target
