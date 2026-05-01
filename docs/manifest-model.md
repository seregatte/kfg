# Manifest Model

kfg uses YAML manifests to define resources. All resources use namespace-style names: `kfg.core.steps.detect-provider`.

## Resource Kinds

### Cmd

Pure shell function. No before/after — those belong in CmdWorkflow.

```yaml
apiVersion: kfg.dev/v1alpha1
kind: Cmd
metadata:
  name: claude
  commandName: claude   # Generated shell function name
spec:
  env:
    ANTHROPIC_API_KEY: "{env:ANTHROPIC_API_KEY}"
  run: command claude "$@"
```

### CmdWorkflow

Orchestration layer. Defines which cmds to include and global before/after steps.

```yaml
apiVersion: kfg.dev/v1alpha1
kind: CmdWorkflow
metadata:
  name: dev
  shell: bash
spec:
  cmds: [claude, gemini]
  before:
    - step: kfg.core.steps.detect-provider
    - step: kfg.core.steps.install-ctx7
      failurePolicy: Ignore
  after:
    - step: kfg.core.steps.cleanup-temp
```

Steps execute in **YAML order** (no weight sorting).

### Step

Reusable unit of work. Can produce outputs via stdout.

```yaml
apiVersion: kfg.dev/v1alpha1
kind: Step
metadata:
  name: kfg.core.steps.detect-provider
spec:
  env:
    TIMEOUT: "30"
  run: |
    if command -v claude >/dev/null 2>&1; then
      echo "anthropic"
    else
      echo "unknown"
    fi
  output:
    name: provider
    type: string
```

### Step References

When referencing a Step in a workflow, you can add:

- **`weight`**: Execution order (lower = earlier). No default for Cmd refs, default `1` for overlays.
- **`when`**: Conditional execution based on step outputs
- **`failurePolicy`**: `Fail` (default) or `Ignore`
- **`env`**: Subshell-scoped environment variable override

```yaml
- step: kfg.core.steps.install-ctx7
  weight: -50
  when:
    output:
      step: kfg.core.steps.detect-provider
      name: provider
      equals: "anthropic"
  failurePolicy: Ignore
  env:
    TIMEOUT: "60"
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
  name: providers
spec:
  schemaRef: schema://providers
  data:
    servers:
      - enabled: true
        type: anthropic
        command: claude
```

**Converter** — Transforms Assets using `template` (Go text/template) or `yq` engine.

```yaml
kind: Converter
metadata:
  name: providers-to-claude
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

## Build Result

During shell generation, a base64-encoded build result file is created at **global scope** (once per workflow). All Cmds and Steps share `$KFG_BUILD_RESULT_FILE`.

```bash
# In generated shell (global scope)
export KFG_BUILD_RESULT_FILE=/tmp/kfg-build-XXXXXX.yaml

# In Steps (auto-detected)
kfg assets convert --use providers-to-$agent
```
