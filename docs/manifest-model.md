# Manifest Model

kfg uses YAML manifests to define resources. Resources follow a short, consistent naming convention: `<scope>.<kind>.<name>` (e.g., `myapp.cmd.deploy`, `ctx7.steps.install`, `ai.claude.asset.settings`).

## Resource Kinds

### Cmd

Pure shell function. No before/after — those belong in CmdWorkflow.

```yaml
apiVersion: kfg.dev/v1alpha1
kind: Cmd
metadata:
  name: myapp.cmd.deploy
  commandName: deploy   # Generated shell function name
spec:
  env:
    DEPLOY_TARGET: "{env:DEPLOY_TARGET}"
  run: |
    echo "Deploying to $DEPLOY_TARGET..."
    kubectl apply -f manifests/
```

### CmdWorkflow

Orchestration layer. Defines which cmds to include and global before/after steps.

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

Steps execute in **YAML order** (no weight sorting).

### Step

Reusable unit of work. Can produce outputs via stdout.

```yaml
apiVersion: kfg.dev/v1alpha1
kind: Step
metadata:
  name: myapp.steps.validate
spec:
  env:
    TIMEOUT: "30"
  run: |
    echo "Validating configuration..."
    if [ -f "config.yaml" ]; then
      echo "valid"
    else
      echo "missing"
    fi
  output:
    name: STATUS
    type: string
```

### Step References

When referencing a Step in a workflow, you can add:

- **`weight`**: Execution order (lower = earlier). No default for Cmd refs, default `1` for overlays.
- **`when`**: Conditional execution based on step outputs
- **`failurePolicy`**: `Fail` (default) or `Ignore`
- **`env`**: Subshell-scoped environment variable override

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
    OUTPUT_DIR: ".claude/skills/"
```

**When operators**: `equals`, `in`, `contains`, `matches`, `allOf`, `anyOf`, `not`.

### Placeholders

`{env:VAR_NAME}` in manifest data/env resolves to `$VAR` at generation time, then shell expands at runtime.

```yaml
env:
  API_KEY: "{env:ANTHROPIC_API_KEY}"  # → $ANTHROPIC_API_KEY
```

Two-stage resolution:
1. **Go stage**: `{env:VAR}` → `$VAR` (generation time)
2. **Shell stage**: `$VAR` → actual value (runtime)

Missing env vars resolve to empty string + warning at generation time.

### Source Layer (Schema / Assets / Converter)

**Schema** — JSON Schema draft-07 for validating Assets.

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
      items:
        type: object
        required: [enabled, type, command]
```

**Assets** — Typed data with schema validation.

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
        command: deploy
```

**Converter** — Transforms Assets using `template` (Go text/template) or `yq` engine.

```yaml
kind: Converter
metadata:
  name: myapp.conv.deploy-config
spec:
  input:
    schemaRef: schema://providers
  engine:
    type: template
    template: |
      {{- range .servers }}{{ if .enabled }}# {{ .type }}{{ end }}{{ end }}
  output:
    format: markdown
```

Converters are resolved from build output only (stdin, `-f` file, or `$KFG_BUILD_RESULT_FILE`).

```bash
kfg build path/to/kustomization | kfg assets convert --use myconverter
kfg assets list -f build.yaml
```

## Extensions

Extensions provide reusable MCP servers and skill installation steps. Each extension lives under `base/extensions/<name>/` and exposes:

- **MCP Assets**: `<ext>.assets.mcp` - Canonical MCP server definition (e.g., `ctx7.assets.mcp`)
- **Install Steps**: `<ext>.steps.install` - Skill installation step (e.g., `ctx7.steps.install`)

### Extension Structure

```
base/extensions/<name>/
├── kustomization.yaml
├── assets/
│   ├── kustomization.yaml
│   └── mcp.yaml         # kfg.extension.<name>.mcp
└── steps/
│   ├── kustomization.yaml
│   └── install.yaml     # kfg.extension.<name>.install
```

### MCP Assets

Extension MCP assets define the canonical MCP server configuration:

```yaml
apiVersion: kfg.dev/v1alpha1
kind: Assets
metadata:
  name: kfg.extension.playwright.mcp
spec:
  input:
    format: yaml
  data:
    name: playwright
    description: Playwright MCP server for end-to-end browser testing
    enabled: true
    server:
      command: npx
      args:
        - -y
        - "@playwright/mcp@latest"
        - "--extension"
      env: {}
```

### Install Steps

Extension install steps provide skill installation with explicit contracts:

```yaml
apiVersion: kfg.dev/v1alpha1
kind: Step
metadata:
  name: kfg.extension.ctx7.install
spec:
  run: |
    # Validate required env vars
    if [ -z "$FLAGS" ]; then
      _kfg.log.error "missing required env var: FLAGS"
      exit 1
    fi
    if [ -z "$OUTPUT_DIR" ]; then
      _kfg.log.error "missing required env var: OUTPUT_DIR"
      exit 1
    fi
    
    # Install and copy skills
    ctx7 setup --cli --project $FLAGS
    cp -Rf .claude/skills/* "$OUTPUT_DIR/"
  env:
    FLAGS: "--yes"
    OUTPUT_DIR: ""
  output:
    name: installed
    type: string
```

### Materialization

Workflows use `kfg.materialize` as the unified primitive for asset materialization:

#### Per-item Mode

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

#### Aggregate Mode

Merge multiple assets into a single output:

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

Overlays compose shared base manifests with project-specific resources. Each overlay lives under `.manifests/overlay/<name>/`:

```
.manifests/overlay/<name>/
├── kustomization.yaml     # References base + overlay resources
├── cmds.yaml              # Overlay-specific commands
├── agents-workflow.yaml   # Overlay-specific workflow
└── assets/
    ├── mcp.yaml           # Overlay-specific MCP assets
    └── ...                 # Other overlay assets
```

### Overlay Structure

```yaml
# overlay/<name>/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
metadata:
  name: my-overlay

resources:
  # Reference shared base (relative path)
  - ../../../kfg/.manifests/base
  # Overlay-specific resources
  - cmds.yaml
  - agents-workflow.yaml
```

### Overlay Workflow

Overlay workflows follow the normalized pattern with shared steps:

```yaml
apiVersion: kfg.dev/v1alpha1
kind: CmdWorkflow
metadata:
  name: myproject.workflow.main
spec:
  cmds:
    - myproject.cmd.foo
  before:
    # Phase 1 (-90): Gitignore
    - step: kfg.ensure-gitignore
      weight: -90
    # Phase 2 (-70): Detection
    - step: kfg.detect-agent
      weight: -70
    # Phase 3-7: Scaffold, Settings, Context, Extension installs, etc.
    # Phase 10 (-40): MCP aggregation with extension assets
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
# In generated shell (global scope)
export KFG_BUILD_RESULT_FILE=/tmp/kfg-build-XXXXXX.yaml

# In Steps (auto-detected)
kfg assets convert --use providers-to-$agent
```
