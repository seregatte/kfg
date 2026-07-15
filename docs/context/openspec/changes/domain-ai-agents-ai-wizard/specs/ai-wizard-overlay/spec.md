## ADDED Requirements

### Requirement: WIZARD OVERLAY STRUCTURE

The wizard overlay SHALL be located at `packages/domains/ai-agents/overlays/ai/` and SHALL contain:
- `kustomization.yaml` — references framework base, domain manifests, and local resources
- `ai-workflow.yaml` — CmdWorkflow defining the wizard agent's lifecycle
- `assets/prompts/kfg-wizard.yaml` — the wizard skill prompt asset

#### Scenario: Overlay structure verified
- **WHEN** the directory `packages/domains/ai-agents/overlays/ai/` is listed
- **THEN** it SHALL contain `kustomization.yaml`, `ai-workflow.yaml`, and `assets/prompts/kfg-wizard.yaml`

### Requirement: KUSTOMIZATION ENTRYPOINT

The `kustomization.yaml` SHALL reference:
- Framework base at `../../../../framework/`
- Domain manifests at `../../manifests/`
- The wizard asset at `assets/prompts/kfg-wizard.yaml`
- The wizard workflow at `ai-workflow.yaml`

#### Scenario: Kustomization loaded
- **WHEN** `kfg build packages/domains/ai-agents/overlays/ai/` is run
- **THEN** it SHALL produce valid YAML output containing a CmdWorkflow resource with name `kfg.ai.workflow`

### Requirement: WIZARD WORKFLOW

The CmdWorkflow `kfg.ai.workflow` SHALL support agents `ai.pi.cmd.main` and `ai.opencode.cmd.main`.

The `before` phase SHALL include:
1. A `kfg.ensure-gitignore` step (weight: -90)
2. An `ai.steps.detect` step (weight: -70)
3. A `kfg.materialize-scaffold` step (weight: -65, per-agent) creating directories `.opencode/commands/` or `.pi/prompts/`
4. A `ctx7.steps.install` step (weight: -55, per-agent, conditional on detect-agent output)
5. A `kfg.materialize` step (weight: -45, per-agent, conditional on detect-agent output) materializing the `ai.prompts.kfg-wizard` asset using the agent's command converter

The `after` phase SHALL include a `kfg.cleanup` step.

Each per-agent step SHALL use `when.output.step: ai.detect-agent` with `name: AGENT` and `equals:` or `in:` to branch execution.

#### Scenario: Wizard workflow for opencode
- **WHEN** the wizard workflow is loaded
- **AND** the detected agent is `opencode`
- **THEN** the `ctx7.steps.install` step SHALL use `FLAGS: "--opencode --yes"` and `OUTPUT_DIR: ".opencode/skills/"`
- **AND** the wizard materialization SHALL output to `.opencode/commands/kfg-wizard.md`

#### Scenario: Wizard workflow for pi
- **WHEN** the wizard workflow is loaded
- **AND** the detected agent is `pi`
- **THEN** the `ctx7.steps.install` step SHALL use `FLAGS: "--opencode --yes"` and `OUTPUT_DIR: ".pi/skills"`
- **AND** the wizard materialization SHALL output to `.pi/prompts/kfg-wizard.md`
