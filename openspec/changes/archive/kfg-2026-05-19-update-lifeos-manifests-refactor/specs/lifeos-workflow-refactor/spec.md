## ADDED Requirements

### Requirement: Workflow uses new kfg naming conventions
The lifeos CmdWorkflow SHALL reference all kfg-provided steps, assets, and converters using the new naming conventions established in the kfg manifest refactoring (commits `e41ab21`, `44f186a`, `19fa1e8`).

#### Scenario: Cmd references are updated
- **WHEN** the workflow defines cmds
- **THEN** all cmd references use `ai.*.cmd.main` format instead of `kfg.agent.cmd.*`

#### Scenario: Step references are updated
- **WHEN** the workflow references steps in before/after blocks
- **THEN** all step references use new naming (`ai.steps.detect`, `<ext>.steps.install`, etc.) instead of old naming (`kfg.detect-agent`, `kfg.extension.*.install`, etc.)

#### Scenario: When conditions reference new step names
- **WHEN** a step has a `when.output.step` condition
- **THEN** the referenced step name matches the new naming convention (`ai.steps.detect`)

### Requirement: Workflow uses kfg.materialize for conversions
The workflow SHALL use `kfg.materialize` step instead of deprecated `kfg.convert` and `kfg.aggregate-mcp` steps.

#### Scenario: Per-item conversion for commands
- **WHEN** converting a command asset to agent-specific format
- **THEN** the step uses `kfg.materialize` with `MODE: "per-item"`, `ASSETS`, `CONVERTER`, and `OUTPUTS` parameters

#### Scenario: Aggregate conversion for MCP configs
- **WHEN** aggregating multiple MCP assets into a single config file
- **THEN** the step uses `kfg.materialize` with `MODE: "aggregate"`, `ASSETS`, `CONVERTER`, `OUTPUTS`, and `WRAP_KEY` parameters

#### Scenario: Settings materialization
- **WHEN** generating agent settings files
- **THEN** the step uses `kfg.materialize` with `MODE: "per-item"` instead of `kfg.agents.steps.settings`

### Requirement: Converter references follow new pattern
The workflow SHALL reference converters using the `ai.<agent>.conv.<type>` pattern instead of the old `kfg.convert.self.<type>.<agent>` pattern.

#### Scenario: Command converters
- **WHEN** a converter transforms a prompt/command asset
- **THEN** the converter name follows `ai.claude.conv.command`, `ai.gemini.conv.command`, `ai.opencode.conv.command`, or `ai.pi.conv.command` format

#### Scenario: MCP converters
- **WHEN** a converter transforms an MCP asset
- **THEN** the converter name follows `ai.claude.conv.mcp`, `ai.gemini.conv.mcp`, or `ai.opencode.conv.mcp` format

#### Scenario: Subagent converters
- **WHEN** a converter transforms a subagent asset
- **THEN** the converter name follows `ai.claude.conv.subagent` or `ai.opencode.conv.subagent` format

### Requirement: Local lifeos assets maintain current naming
LifeOS-specific assets defined in `lifeos/.manifests/base/lifeos/` SHALL maintain their current naming convention (`lifeos.commands.*`, `lifeos.mcp.*`, `lifeos.subagents.*`).

#### Scenario: Command assets are referenced by original name
- **WHEN** the workflow references lifeos commands
- **THEN** the asset names remain `lifeos.commands.prepare-discourse`, `lifeos.commands.wpp-elder`, and `lifeos.commands.hourglass-sync`

#### Scenario: MCP assets are referenced by original name
- **WHEN** the workflow references lifeos MCPs
- **THEN** the asset names remain `lifeos.mcp.blender`, `lifeos.mcp.playwright`, and `lifeos.mcp.eld`

#### Scenario: Subagent assets are referenced by original name
- **WHEN** the workflow references lifeos subagents
- **THEN** the asset name remains `lifeos.subagents.elder`

### Requirement: kfg apply validates successfully
After all reference updates, running `kfg apply -k .manifests/overlay/dev` SHALL succeed without errors.

#### Scenario: kfg apply completes without errors
- **WHEN** `kfg apply` is executed with the lifeos dev overlay
- **THEN** the exit code is 0 and no "not found" or "invalid reference" errors appear in output

#### Scenario: Generated shell code passes syntax check
- **WHEN** the generated shell code is validated with `bash -n`
- **THEN** the syntax check passes with exit code 0
