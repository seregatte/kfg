## 1. Update Cmd References

- [x] 1.1 Replace `kfg.agent.cmd.claude` with `ai.claude.cmd.main`
- [x] 1.2 Replace `kfg.agent.cmd.gemini` with `ai.gemini.cmd.main`
- [x] 1.3 Replace `kfg.agent.cmd.opencode` with `ai.opencode.cmd.main`
- [x] 1.4 Replace `kfg.agent.cmd.pi` with `ai.pi.cmd.main`

## 2. Update Step References in When Conditions

- [x] 2.1 Replace all `kfg.detect-agent` references with `ai.steps.detect` in `when.output.step` fields (~50 occurrences)

## 3. Update Settings Steps

- [x] 3.1 Replace `kfg.agents.steps.settings` with `kfg.materialize` + params for claude
- [x] 3.2 Replace `kfg.agents.steps.settings` with `kfg.materialize` + params for gemini
- [x] 3.3 Replace `kfg.agents.steps.settings` with `kfg.materialize` + params for pi
- [x] 3.4 Replace `kfg.agents.steps.settings` with `kfg.materialize` + params for opencode

## 4. Update Command Conversion Steps

- [x] 4.1 Replace `kfg.convert` with `kfg.materialize` + params for all 3 commands × 4 agents (12 steps)
- [x] 4.2 Update converter references: `kfg.convert.self.command.*` → `ai.*.conv.command`
- [x] 4.3 Update parameter: `OUTPUT` → `OUTPUTS`, add `MODE: "per-item"`
- [x] 4.4 Replace `kfg.extension.self.commands.git-commit` with `ai.prompts.git-commit` in asset references

## 5. Update Subagent Conversion Steps

- [x] 5.1 Replace `kfg.convert` with `kfg.materialize` + params for elder subagent (2 steps)
- [x] 5.2 Update converter references: `kfg.convert.self.subagent.*` → `ai.*.conv.subagent`
- [x] 5.3 Update parameter: `OUTPUTS` → `OUTPUTS`, add `MODE: "per-item"`

## 6. Update MCP Aggregation Steps

- [x] 6.1 Replace `kfg.aggregate-mcp` with `kfg.materialize` + params for claude
- [x] 6.2 Replace `kfg.aggregate-mcp` with `kfg.materialize` + params for gemini
- [x] 6.3 Replace `kfg.aggregate-mcp` with `kfg.materialize` + params for opencode
- [x] 6.4 Update converter references: `kfg.convert.self.mcp.*` → `ai.*.conv.mcp`
- [x] 6.5 Update parameters: `TARGET` → `OUTPUTS`, add `MODE: "aggregate"`, add `WRAP_KEY`

## 7. Update Extension Install Steps

- [x] 7.1 Replace `kfg.extension.playwright.install` with `playwright.steps.install`
- [x] 7.2 Replace `kfg.extension.notebooklm.install` with `notebooklm.steps.install`
- [x] 7.3 Replace `kfg.extension.gws.install` with `gws.steps.install`

## 8. Validate Changes

- [x] 8.1 Build kfg binary: `cd ~/Sites/kfg && make build`
- [x] 8.2 Run `kfg apply -k .manifests/overlay/dev` in lifeos directory
- [x] 8.3 Verify exit code 0 and no "not found" errors
- [x] 8.4 Run syntax check on generated shell code: `bash -n <generated_file>`
- [x] 8.5 Verify all expected output files are generated (settings.json, commands, MCP configs)
