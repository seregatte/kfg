## Context

The `nixai-freeze-v1` project had a layered merge system for agent configs: base templates →
feature definitions → user local overrides → companion layers. KFG's materialization pipeline
(`kfg.materialize`, `kfg.copy-context`) can generate files but has no mechanism for users to
overlay their own files on top without modifying manifests.

## Architecture

### Step: `kfg.merge-local`

A framework Step that merges files from a source directory into a target directory.

```
SRC directory/                    DEST directory/
├── opencode.json        ──deep──▶  opencode.json (merged)
├── skills/                       skills/
│   └── my-skill/                  my-skill/
│       └── SKILL.md    ──copy──▶    SKILL.md (new)
└── mcp.json            ──skip──▶  mcp.json (already exists, skipped)
```

**Parameters (env vars):**

| Var | Default | Description |
|-----|---------|-------------|
| `SRC` | (required) | Source directory with user's custom files |
| `DEST` | `.` | Target directory (where agent configs live) |
| `MODE` | `deep` | `deep` = JSON/YAML deep merge, `overwrite` = blind copy, `skip` = preserve existing |
| `FILTER` | `""` | Optional glob filter (e.g., `*.json`, `*.yaml`) |

**Behavior:**
- If `SRC` is empty or doesn't exist: logs warning, returns 0 (idempotent skip)
- Walk `SRC` recursively with optional `FILTER`
- For each file, compute relative path and target
- `MODE=deep`: for `.json` files uses `jq -s '.[0] * .[1]'`; for `.yaml`/`.yml` uses `yq eval-all 'select(fileIndex == 0) * select(fileIndex == 1)'`
- `MODE=overwrite`: `cp` unconditionally
- `MODE=skip`: skip if target exists
- Register copied/merged files as artifacts

### Wizard Integration

The wizard prompt (`wizard.yaml`) is extended with a **Phase 3.5** after the user confirms
the proposal but before generation:

1. Ask: "Do you want to add a local merge layer for agent-specific overrides?"
2. If yes, ask for `SRC` path (default: `.kfg/local/<agent>/`) and merge `MODE`
3. Include `kfg.merge-local` step in the generated workflow with user's choices

### Workflow Placement

The merge step runs at weight `-60` (between scaffold at `-65` and skill install at `-55`),
per-agent conditional on the detected agent:

```yaml
- name: ai.merge-local.opencode
  step: kfg.merge-local
  weight: -60
  when:
    output:
      step: ai.detect-agent
      name: AGENT
      equals: "opencode"
  env:
    SRC: ".kfg/local/opencode/"
    DEST: ".opencode/"
    MODE: "deep"
```

## Non-Goals

- Not modifying the Go engine (pure Step + wizard prompt change)
- Not auto-including the step in all workflows (opt-in)
- Not handling git-aware merges or conflict resolution
