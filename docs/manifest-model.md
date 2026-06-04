# Manifest Model

kfg uses YAML manifests with short, consistent naming: `<scope>.<kind>.<name>` (e.g., `myapp.cmd.deploy`, `ctx7.steps.install`).

## Resource Kinds

### Cmd — Shell Function

Pure shell function (no before/after — those go in CmdWorkflow).

```yaml
apiVersion: kfg.dev/v1alpha1
kind: Cmd
metadata:
  name: myapp.cmd.deploy
  commandName: deploy
spec:
  env:
    DEPLOY_TARGET: "{env:DEPLOY_TARGET}"
  run: kubectl apply -f manifests/
```

### CmdWorkflow — Orchestration

Defines which cmds to include and global before/after steps. Steps execute in **YAML order** (no weight sorting).

```yaml
apiVersion: kfg.dev/v1alpha1
kind: CmdWorkflow
metadata:
  name: myapp.workflow.main
  shell: bash
spec:
  cmds: [myapp.cmd.deploy, myapp.cmd.test]
  before:
    - step: myapp.steps.validate
    - step: ctx7.steps.install
      failurePolicy: Ignore
  after:
    - step: kfg.cleanup
```

### Step — Reusable Work Unit

Produces outputs via stdout.

```yaml
apiVersion: kfg.dev/v1alpha1
kind: Step
metadata:
  name: myapp.steps.validate
spec:
  env:
    TIMEOUT: "30"
  run: |
    [ -f "config.yaml" ] && echo "valid" || echo "missing"
  output:
    name: STATUS
    type: string
```

### Step References

When referencing a Step in workflow:

| Property | Purpose |
|----------|---------|
| `weight` | Execution order (lower = earlier). Default: `1` for overlays |
| `when` | Conditional execution based on step outputs |
| `failurePolicy` | `Fail` (default) or `Ignore` |
| `env` | Subshell-scoped environment override |

```yaml
- step: ctx7.steps.install
  weight: -55
  when:
    output:
      step: myapp.steps.validate
      name: STATUS
      equals: "valid"
  failurePolicy: Ignore
  env:
    FLAGS: "--yes"
```

**When operators:** `equals`, `in`, `contains`, `matches`, `allOf`, `anyOf`, `not`

### Placeholders

`{env:VAR_NAME}` in manifest data/env resolves to `$VAR` at generation time.

```yaml
env:
  API_KEY: "{env:ANTHROPIC_API_KEY}"  # → $ANTHROPIC_API_KEY
```

Two-stage resolution:
1. **Go stage (generation):** `{env:VAR}` → `$VAR`
2. **Shell stage (runtime):** `$VAR` → actual value

Missing vars resolve to empty string + warning.

## Source Layer: Schema / Assets / Converter

### Schema — JSON Schema Draft-07

```yaml
kind: Schema
metadata:
  name: providers
spec:
  type: object
  required: [servers]
  properties:
    servers:
      type: array
      items: {type: object}
```

### Assets — Typed Data

```yaml
kind: Assets
metadata:
  name: myapp.assets.providers
spec:
  schemaRef: schema://providers
  data:
    servers:
      - enabled: true
        type: production
```

### Converter — Transform Assets

```yaml
kind: Converter
metadata:
  name: myapp.conv.deploy-config
spec:
  input:
    schemaRef: schema://providers
  engine:
    type: template
    template: "{{- range .servers }}{{ if .enabled }}# {{ .type }}{{ end }}{{ end }}"
  output:
    format: markdown
```

Converters resolved from build output only:

```bash
kfg build path/to/kustomization | kfg assets convert --use myconverter
kfg assets list -f build.yaml
```

## Extensions

Extensions provide reusable MCP servers and skill installation steps. Each at `base/extensions/<name>/`:

| Asset | Purpose | Name |
|-------|---------|------|
| **MCP Assets** | Canonical MCP server config | `<ext>.assets.mcp` |
| **Install Steps** | Skill installation | `<ext>.steps.install` |

### Extension Structure

```
base/extensions/<name>/
├── kustomization.yaml
├── assets/
│   └── mcp.yaml         # kfg.extension.<name>.mcp
└── steps/
    └── install.yaml     # kfg.extension.<name>.install
```

### MCP Assets Example

```yaml
apiVersion: kfg.dev/v1alpha1
kind: Assets
metadata:
  name: kfg.extension.playwright.mcp
spec:
  data:
    name: playwright
    enabled: true
    server:
      command: npx
      args: [-y, "@playwright/mcp@latest", "--extension"]
```

### Install Steps Example

```yaml
apiVersion: kfg.dev/v1alpha1
kind: Step
metadata:
  name: kfg.extension.ctx7.install
spec:
  run: |
    if [ -z "$FLAGS" ] || [ -z "$OUTPUT_DIR" ]; then
      _kfg.log.error "missing required env vars"
      exit 1
    fi
    ctx7 setup --cli --project $FLAGS
    cp -Rf .claude/skills/* "$OUTPUT_DIR/"
  env:
    FLAGS: "--yes"
    OUTPUT_DIR: ""
```

## Materialization

Unified primitive for asset materialization using `kfg.materialize`:

### Per-Item Mode

Convert one or more assets to individual outputs:

```yaml
- step: kfg.materialize
  weight: -45
  env:
    MODE: "per-item"
    ASSETS: "ai.prompts.git-commit"
    CONVERTER: "ai.claude.conv.command"
    OUTPUTS: ".claude/commands/git-commit.md"
```

### Aggregate Mode

Merge multiple assets into single output:

```yaml
- step: kfg.materialize
  weight: -40
  env:
    MODE: "aggregate"
    ASSETS: "ctx7.assets.mcp:chrome.assets.mcp:playwright.assets.mcp"
    CONVERTER: "ai.claude.conv.mcp"
    OUTPUTS: ".mcp.json"
    WRAP_KEY: "mcpServers"
```

## Overlays

Overlays compose shared base manifests with project-specific resources under domain packages:

```
packages/domains/<domain>/overlays/<name>/
├── kustomization.yaml     # References domain + overlay resources
└── agents-workflow.yaml    # Overlay-specific workflow
```

### Overlay Workflow Pattern

```yaml
apiVersion: kfg.dev/v1alpha1
kind: CmdWorkflow
metadata:
  name: myproject.workflow.main
spec:
  cmds:
    - myproject.cmd.foo
  before:
    - step: kfg.ensure-gitignore
      weight: -90
    - step: kfg.detect-agent
      weight: -70
    - step: kfg.materialize
      weight: -40
      env:
        MODE: "aggregate"
        ASSETS: "myproject.mcp.local:kfg.extension.ctx7.mcp"
        CONVERTER: "kfg.convert.self.mcp.claude"
        OUTPUTS: ".mcp.json"
        WRAP_KEY: "mcpServers"
```

## Build Result

During shell generation, a base64-encoded build result file is created at **global scope** (once per workflow). All Cmds and Steps share `$KFG_BUILD_RESULT_FILE`.

```bash
export KFG_BUILD_RESULT_FILE=/tmp/kfg-build-XXXXXX.yaml

# In Steps (auto-detected)
kfg assets convert --use providers-to-$agent
```
