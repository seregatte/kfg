# Design: Restore Wizard Skill Materialization Steps

## Changes to `ai-workflow.yaml`

Add back the Phase 5b block after the existing Phase 5a block (after `ai.wizard.pi` at weight -45):

```yaml
    # Phase 5b (-44): Wizard skill materialize (per-agent)
    - name: ai.wizard.skill.opencode
      step: kfg.materialize
      weight: -44
      when:
        output:
          step: ai.detect-agent
          name: AGENT
          equals: "opencode"
      env:
        MODE: "per-item"
        ASSETS: "ai.prompts.wizard"
        CONVERTER: "ai.opencode.conv.skill"
        OUTPUTS: ".opencode/skills/wizard/SKILL.md"
    - name: ai.wizard.skill.pi
      step: kfg.materialize
      weight: -44
      when:
        output:
          step: ai.detect-agent
          name: AGENT
          equals: "pi"
      env:
        MODE: "per-item"
        ASSETS: "ai.prompts.wizard"
        CONVERTER: "ai.pi.conv.command"
        OUTPUTS: ".pi/skills/wizard/SKILL.md"
```

### Rationale

- `ai.prompts.wizard` — the full skill prompt (exists in `assets/prompts/wizard.yaml`)
- `ai.opencode.conv.skill` — dedicated skill converter for opencode (exists in manifests)
- `ai.pi.conv.command` — pi-only converter (pi has no dedicated skill converter)
- Output paths match what the README documents as generated artifacts
- README is already correct; no changes needed

## Files Changed

1. `packages/domains/ai-agents/overlays/ai/ai-workflow.yaml` — Add Phase 5b block
