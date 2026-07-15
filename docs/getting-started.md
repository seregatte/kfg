# Getting Started with KFG

This guide will take you from zero to your first working workflow in 10 minutes.

## Prerequisites

Before you start, make sure you have:

- **Nix** (recommended) or **Go 1.21+** for building from source
- A text editor
- Terminal with bash or zsh

## Installation

### Option 1: Via Nix (Recommended)

If you already have Nix installed with flakes enabled:

```bash
# Add to current shell
nix shell github:seregatte/kfg

# Verify installation
kfg version
```

### Option 2: Build from Source

```bash
git clone https://github.com/seregatte/kfg.git
cd kfg
make build

# Add to PATH
export PATH="$PWD/bin:$PATH"

# Verify installation
kfg version
```

## First Manifest

Let's create a simple command that prints "Hello, World!".

### 1. Create the file `hello.yaml`

```yaml
apiVersion: kfg.dev/v1alpha1
kind: Cmd
metadata:
  name: tutorial.cmd.hello
  commandName: hello
spec:
  run: echo "Hello, World!"
```

**Explanation**:
- `apiVersion`: KFG API version
- `kind`: Resource type (Cmd = shell function)
- `metadata.name`: Unique resource name
- `metadata.commandName`: Name of the generated command
- `spec.run`: Shell code to execute

### 2. Apply the manifest

```bash
kfg apply -f hello.yaml --workflow default
```

What happens:
1. KFG reads the YAML manifest
2. Generates shell code
3. Makes the `hello` command available in your shell

### 3. Run the command

```bash
hello
# Output: Hello, World!
```

## Adding Environment Variables

Let's modify the command to use environment variables:

```yaml
apiVersion: kfg.dev/v1alpha1
kind: Cmd
metadata:
  name: tutorial.cmd.greet
  commandName: greet
spec:
  env:
    NAME: "{env:USER:-Anonymous}"
    GREETING: "{env:GREETING:-Hello}"
  run: echo "$GREETING, $NAME!"
```

**Explanation**:
- `{env:VAR}`: Placeholder resolved at generation time
- `{env:VAR:-default}`: With default value
- At runtime, becomes `$VAR`

### Test with different values:

```bash
# Default values
greet
# Output: Hello, Anonymous!

# With environment variables
USER=John GREETING=Hi greet
# Output: Hi, John!
```

## Working with Steps

Steps are reusable units of work that can have outputs.

### 1. Create a step

```yaml
apiVersion: kfg.dev/v1alpha1
kind: Step
metadata:
  name: tutorial.steps.check-file
spec:
  run: |
    if [ -f "config.yaml" ]; then
      echo "found"
    else
      echo "missing"
    fi
  output:
    name: FILE_STATUS
    type: string
```

### 2. Use the step in a workflow

```yaml
apiVersion: kfg.dev/v1alpha1
kind: Cmd
metadata:
  name: tutorial.cmd.deploy
  commandName: deploy
spec:
  run: echo "Deploying application..."

---
apiVersion: kfg.dev/v1alpha1
kind: CmdWorkflow
metadata:
  name: tutorial.workflow.deploy
spec:
  cmds: [tutorial.cmd.deploy]
  before:
    - step: tutorial.steps.check-file
```

### 3. Apply and run

```bash
kfg apply -f workflow.yaml --workflow deploy
deploy
```

The `check-file` step runs before the `deploy` command.

## Conditional Execution

You can run steps conditionally based on outputs:

```yaml
apiVersion: kfg.dev/v1alpha1
kind: CmdWorkflow
metadata:
  name: tutorial.workflow.smart-deploy
spec:
  cmds: [tutorial.cmd.deploy]
  before:
    - step: tutorial.steps.check-file
    - step: tutorial.steps.validate-config
      when:
        output:
          step: tutorial.steps.check-file
          name: FILE_STATUS
          equals: "found"
```

**Available operators**:
- `equals`: Equals value
- `in`: In list
- `contains`: Contains substring
- `matches`: Regex match
- `allOf`: All conditions
- `anyOf`: Any condition
- `not`: Negation

## Working with Kustomization

Kustomization allows composing manifests from multiple sources.

### Directory structure

```
myproject/
├── base/
│   ├── kustomization.yaml
│   ├── cmd-deploy.yaml
│   └── steps-validate.yaml
└── overlays/
    └── dev/
        ├── kustomization.yaml
        └── cmd-dev.yaml
```

### base/kustomization.yaml

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - cmd-deploy.yaml
  - steps-validate.yaml
```

### overlays/dev/kustomization.yaml

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - ../../base
  - cmd-dev.yaml
```

### Apply the overlay

```bash
kfg apply -k overlays/dev --workflow deploy
```

## Step Caching

Steps can be cached to avoid re-executions:

```yaml
apiVersion: kfg.dev/v1alpha1
kind: Step
metadata:
  name: tutorial.steps.expensive
  annotations:
    kfg.dev/cacheable: "true"
spec:
  run: |
    echo "Running expensive operation..."
    sleep 5
    echo "done"
  output:
    name: RESULT
    type: string
```

### Managing the cache

```bash
# List cache entries
kfg sys cache ls

# Inspect specific entry
kfg sys cache inspect <id>

# Remove old entries
kfg sys cache prune

# Invalidate cache on next run
KFG_REFRESH=1 kfg apply -f workflow.yaml --workflow deploy
```

## Debug and Troubleshooting

### Increase verbosity

```bash
# Levels: 0 (quiet) to 5 (maximum debug)
KFG_VERBOSE=3 kfg apply -f workflow.yaml --workflow deploy
```

### View logs

```bash
# Default log location
ls ~/.local/state/kfg/logs/

# Or custom path
KFG_LOG_DIR=/tmp/kfg-logs kfg apply -f workflow.yaml --workflow deploy
```

### View generated code

```bash
# Generate without executing
kfg build overlays/dev -o generated.yaml

# View content
cat generated.yaml
```

## Next Steps

Now that you know the basics, explore:

- **[Manifest Model](manifest-model.md)** - Full manifest schema
- **[CLI Reference](cli-reference.md)** - All commands and flags
- **[Architecture](architecture.md)** - How KFG works internally
- **[Examples](../packages/domains/ai-agents/manifests/)** - Real manifests from the project

## Complete Example: CI/CD Pipeline

Here's a real-world CI/CD pipeline example:

```yaml
# steps.yaml
apiVersion: kfg.dev/v1alpha1
kind: Step
metadata:
  name: cicd.steps.checkout
spec:
  run: |
    git clone "$REPO_URL" /tmp/build
    echo "checked_out"
  output:
    name: STATUS
    type: string

---
apiVersion: kfg.dev/v1alpha1
kind: Step
metadata:
  name: cicd.steps.test
spec:
  run: |
    cd /tmp/build
    make test
    echo "tests_passed"
  output:
    name: TEST_STATUS
    type: string

---
apiVersion: kfg.dev/v1alpha1
kind: Step
metadata:
  name: cicd.steps.build
spec:
  run: |
    cd /tmp/build
    make build
    echo "build_complete"
  output:
    name: BUILD_STATUS
    type: string

---
# commands.yaml
apiVersion: kfg.dev/v1alpha1
kind: Cmd
metadata:
  name: cicd.cmd.deploy-staging
  commandName: deploy-staging
spec:
  env:
    ENV: "staging"
  run: |
    cd /tmp/build
    kubectl apply -f k8s/$ENV/

---
apiVersion: kfg.dev/v1alpha1
kind: Cmd
metadata:
  name: cicd.cmd.deploy-production
  commandName: deploy-production
spec:
  env:
    ENV: "production"
  run: |
    cd /tmp/build
    kubectl apply -f k8s/$ENV/

---
# workflow.yaml
apiVersion: kfg.dev/v1alpha1
kind: CmdWorkflow
metadata:
  name: cicd.workflow.full-pipeline
spec:
  cmds: [cicd.cmd.deploy-staging]
  before:
    - step: cicd.steps.checkout
      weight: -100
    - step: cicd.steps.test
      weight: -90
      when:
        output:
          step: cicd.steps.checkout
          name: STATUS
          equals: "checked_out"
    - step: cicd.steps.build
      weight: -80
      when:
        output:
          step: cicd.steps.test
          name: TEST_STATUS
          equals: "tests_passed"
  after:
    - step: kfg.cleanup
      weight: 100
```

### Use the pipeline

```bash
# Full pipeline
REPO_URL=https://github.com/user/repo kfg apply -f cicd.yaml --workflow full-pipeline

# Run deploy
deploy-staging

# Or production deploy (with another workflow)
deploy-production
```

This pipeline:
1. Checks out the code
2. Runs tests (if checkout OK)
3. Builds (if tests passed)
4. Deploys to staging
5. Cleans up temporary resources

---

**Questions?** See [Troubleshooting](troubleshooting.md) or open an [issue](https://github.com/seregatte/kfg/issues).
