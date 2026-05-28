## 1. Base Cmd

- [x] 1.1 Create `.manifests/base/cmds/openspec.yaml` with Cmd `kfg.agent.cmd.openspec` (commandName: openspec, run: `command openspec "$@"`)
- [x] 1.2 Update `.manifests/base/cmds/kustomization.yaml` to include `openspec.yaml`

## 2. Remove Overlay Dev Cmd

- [x] 2.1 Delete `.manifests/overlay/dev/cmds.yaml`
- [x] 2.2 Update `.manifests/overlay/dev/kustomization.yaml` to remove `cmds.yaml` from resources

## 3. Update Openspec Install Step

- [x] 3.1 Update `.manifests/base/extensions/openspec/steps/install.yaml` to copy both skills AND commands from `$AGENT_HOME/`

## 4. Integrate Openspec into Agents Workflow

- [x] 4.1 Add 4 openspec install steps to `.manifests/overlay/dev/agents-workflow.yaml` at weight -53 (per-agent: claude, opencode, gemini, pi)

## 5. Validation

- [x] 5.1 Run `make test` to verify Go unit tests pass
- [x] 5.2 Run `make test-bats` to verify integration tests pass
- [x] 5.3 Verify `kfg build` generates correct shell output with openspec cmd and install steps
