# KFG - Declarative Shell Compiler [DEPRECATED]

> ⚠️ **This project is no longer maintained.**
>
> KFG served as an experimental declarative shell compiler and AI agent configuration manager from 2024–2026.
> Its agent configuration manifests have been **materialized directly into projects** (`.pi/` and `.opencode/`),
> removing the runtime dependency on KFG. The `openspec` tool is now consumed via `npx @fission-ai/openspec`.
>
> This repository is **archived and read-only**. See individual project repositories for current agent configurations.

[![GitHub release (latest by date)](https://img.shields.io/github/v/release/seregatte/kfg)](https://github.com/seregatte/kfg/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Nix](https://img.shields.io/badge/Nix-5277C3?logo=nixos&logoColor=white)](https://nixos.org)

**KFG** (Key Function Generator) is a declarative shell compiler that transforms YAML manifests into bash functions. Define commands, dependencies, and execution steps in YAML, and KFG generates shell code that can be sourced or used interactively.

## Why KFG?

- **Declarative**: Define what to do in YAML, not how to do it in bash
- **Dependencies**: KFG manages execution order automatically via DAG
- **Cache**: Steps are cached to avoid unnecessary re-executions
- **Reusable**: Create modular and reusable manifests
- **Versionable**: Your shell workflows can now be versioned as code

## Use Cases

- **Application deployment**: Define declarative deploy pipelines
- **Environment setup**: Automate project configuration with dependencies
- **AI agent workflows**: Generate commands for Claude, Copilot, etc.
- **CI/CD**: Standardize build and deploy processes
- **MCP server management**: Configure and manage MCP servers declaratively

## Installation

### Prerequisites

To install via Nix (recommended):
- [Nix](https://nixos.org/download.html) with flakes enabled

To build from source:
- Go 1.21+
- Make

### Via Nix (Recommended)

```bash
# Build and install
nix build github:seregatte/kfg

# Run without installing
nix run github:seregatte/kfg -- --help

# Add to current shell
nix shell github:seregatte/kfg
```

Supports Linux and macOS (x86_64 and ARM64).

### Build from Source

```bash
git clone https://github.com/seregatte/kfg.git
cd kfg
make build
```

The binary will be placed in `./bin/kfg`.

### Install to GOPATH

```bash
make install
```

## Quick Start

### 1. Create your first manifest

Create a file `hello.yaml`:

```yaml
apiVersion: kfg.dev/v1alpha1
kind: Cmd
metadata:
  name: myapp.cmd.hello
  commandName: hello
spec:
  run: echo "Hello from KFG!"
```

### 2. Apply the manifest

```bash
kfg apply -f hello.yaml --workflow default
```

This generates and executes the shell code. Now you have a `hello` command available!

### 3. Run the command

```bash
hello
# Output: Hello from KFG!
```

### Complete Example: Deploy Pipeline

Create `deploy.yaml`:

```yaml
apiVersion: kfg.dev/v1alpha1
kind: Cmd
metadata:
  name: myapp.cmd.deploy
  commandName: deploy
spec:
  env:
    DEPLOY_TARGET: "{env:DEPLOY_TARGET:-production}"
  run: |
    echo "Deploying to $DEPLOY_TARGET..."
    kubectl apply -f manifests/

---
apiVersion: kfg.dev/v1alpha1
kind: Step
metadata:
  name: myapp.steps.validate
spec:
  run: |
    [ -f "config.yaml" ] && echo "Config found" || exit 1
  output:
    name: STATUS
    type: string

---
apiVersion: kfg.dev/v1alpha1
kind: CmdWorkflow
metadata:
  name: myapp.workflow.deploy
spec:
  cmds: [myapp.cmd.deploy]
  before:
    - step: myapp.steps.validate
```

Apply with:

```bash
DEPLOY_TARGET=staging kfg apply -f deploy.yaml --workflow deploy
```

📖 **More examples**: See [docs/getting-started.md](docs/getting-started.md) for a full tutorial.

## Documentation

- **[Getting Started](docs/getting-started.md)** - Step-by-step tutorial
- **[CLI Reference](docs/cli-reference.md)** - Complete CLI reference
- **[Manifest Model](docs/manifest-model.md)** - Schema and manifest types
- **[Architecture](docs/architecture.md)** - KFG internal architecture
- **[Troubleshooting](docs/troubleshooting.md)** - Common issues and solutions
- **[Contributing](CONTRIBUTING.md)** - How to contribute

## Command Reference

| Command | Description | Example |
|---------|-------------|---------|
| `kfg apply` | Apply kustomization/manifest and generate shell code | `kfg apply -f manifest.yaml --workflow main` |
| `kfg run` | Run an agent one-shot | `kfg run -k ./manifests myagent` |
| `kfg build` | Build kustomization to YAML | `kfg build ./manifests -o output.yaml` |
| `kfg sys cache` | Step cache management | `kfg sys cache ls` |
| `kfg sys log` | Structured logging for scripts | `kfg sys log info "component" "message"` |
| `kfg version` | Show version information | `kfg version` |

📖 **Full reference**: See [docs/cli-reference.md](docs/cli-reference.md) for all commands, flags, and environment variables.

## Comparison with Alternatives

| Feature | KFG | Make | Just | Task |
|---------|-----|------|------|------|
| **Syntax** | Declarative YAML | Makefile | Justfile | YAML |
| **Dependencies** | Automatic via DAG | Manual | Manual | Manual |
| **Step caching** | ✅ Native | ❌ | ❌ | ❌ |
| **Modular composition** | ✅ Kustomize | ❌ | ❌ | Limited |
| **Code generation** | ✅ Shell functions | ❌ | ❌ | ❌ |
| **Placeholders** | ✅ `{env:VAR}` | ❌ | ❌ | Limited |
| **Versionable** | ✅ Yes | ✅ Yes | ✅ Yes | ✅ Yes |

**When to use KFG?**
- When you need **automatic dependencies** between tasks
- When you want **smart caching** to avoid re-executions
- When you need to **compose manifests** from different sources (Kustomize)
- When you want to **generate reusable shell functions**
- When you're working with **AI agents** that need structured commands

## Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `KFG_VERBOSE` | Verbosity level (0-5) | `KFG_VERBOSE=3` |
| `KFG_STORE_DIR` | Store directory (cache) | `~/.kfg/store` |
| `KFG_LOG_FILE` | Log file path | `/tmp/kfg.log` |
| `KFG_LOG_DIR` | Log directory | `~/.local/state/kfg/logs` |
| `KFG_LOG_COLOR` | Color mode (auto/always/never) | `auto` |
| `KFG_KPATH` | Default kustomization path | `./manifests` |
| `KFG_REFRESH` | Invalidate cache (set to "1") | `1` |

📖 **Full reference**: See [docs/cli-reference.md](docs/cli-reference.md#environment-variables).

## API Version

KFG uses `kfg.dev/v1alpha1` as the API version for manifests:

```yaml
apiVersion: kfg.dev/v1alpha1
kind: Cmd
metadata:
  name: example
spec:
  run: echo "Hello, World!"
```

📖 **Full schema**: See [docs/manifest-model.md](docs/manifest-model.md) for all resource types.

## Real-World Examples

KFG is used in this repository itself to manage AI agent workflows:

```bash
# Apply development overlay
kfg apply -k packages/domains/ai-agents/overlays/dev --workflow agents

# Run a specific agent
kfg run -k packages/domains/ai-agents/overlays/dev openspec
```

See more examples in `packages/domains/ai-agents/manifests/`.

## Development

### DevShells

KFG provides three devShells via Nix flakes:

| Shell | Usage | Description |
|-------|-------|-------------|
| `default` | `nix develop` | **Consumer shell** — tools for using KFG |
| `dev` | `nix develop .#dev` | **Development shell** — full development environment |
| `ci` | `nix develop .#ci` | **CI shell** — minimal for build and tests |

#### Preferring System-Installed Tools

By default, every devShell prefers your host-installed tools (Homebrew, system profiles, `/usr/local/bin`, etc.) over the Nix devShell versions for a preconfigured list of commands. When a tool is not found on the host, the Nix version is used as a fallback.

**Default allowlist:** `go`, `node`, `npm`, `npx`, `corepack`, `uv`, `uvx`, `bats`, `openspec`, `pi`, `ctx7`, `chrome-devtools-mcp`, `gws`, `notebooklm`, `nblm`, `opencode`, `playwright`

| Variable | Description | Example |
|----------|-------------|---------|
| `KFG_PREFER_SYSTEM` | Set to `0` to disable and use Nix versions only | `KFG_PREFER_SYSTEM=0 nix develop` |
| `KFG_PREFER_SYSTEM_COMMANDS` | Space-separated list of commands to prefer from host | `KFG_PREFER_SYSTEM_COMMANDS="go opencode" nix develop` |

**Note:** This changes which executable is resolved at runtime, but does not prevent Nix from realizing the declared packages. The Nix versions remain available as immediate fallback on PATH.

### Building

```bash
# Using the dev shell
nix develop .#dev --command make build        # → ./bin/kfg
nix develop .#dev --command make test         # Go unit tests
nix develop .#dev --command make test-bats    # Bats integration tests
```

### Repository Structure

```
├── src/                          # Engine implementation (Go)
│   ├── cmd/kfg/                  # CLI commands
│   └── internal/                 # Internal packages
├── packages/
│   ├── framework/                # Shared manifest primitives
│   │   ├── manifests/            # Reusable steps
│   │   └── tests/                # Framework tests
│   └── domains/
│       └── ai-agents/            # AI agents domain package
│           ├── manifests/        # AI agent resources
│           ├── overlays/dev/     # Development overlay
│           └── tests/            # Domain tests
├── docs/
│   ├── AGENTS.md                 # AI agent operating context
│   └── context/
│       └── openspec/             # OpenSpec specifications
├── tests/
│   └── bats/                     # Engine and integration tests
└── Makefile                      # Build and test targets
```

📖 **Detailed architecture**: See [docs/architecture.md](docs/architecture.md).

## License

MIT License — see [LICENSE](LICENSE) for details.

## Links

- **Repository**: https://github.com/seregatte/kfg
- **Releases**: https://github.com/seregatte/kfg/releases
- **Issues**: https://github.com/seregatte/kfg/issues
- **Discussions**: https://github.com/seregatte/kfg/discussions
