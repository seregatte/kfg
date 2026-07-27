# Restore Wizard Skill Materialization Steps

## What Changed

Commit `5a9c1d1` ("feat: rename /kfg-wizard command to /wizard") accidentally
removed the entire Phase 5b (weight -44) wizard skill materialization steps
from `packages/domains/ai-agents/overlays/ai/ai-workflow.yaml`.

The steps for `ai.wizard.skill.opencode` and `ai.wizard.skill.pi` were deleted
along with their `when` condition blocks.

## Impact

- The `kfg ai` wizard no longer materializes `.opencode/skills/wizard/SKILL.md`
  and `.pi/skills/wizard/SKILL.md`
- The README still documents these artifacts, but they are not generated
- The wizard skill prompt (`assets/prompts/wizard.yaml`) is not used as a skill

## Fix

Restore the Phase 5b steps in their last-known-good form (from commit `d136ba3`).
