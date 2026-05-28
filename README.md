# KFG - Declarative Shell Compiler

KFG is a declarative shell compiler that transforms YAML manifests into bash functions. It allows you to define shell commands, their dependencies, and execution steps in YAML manifests, then generates shell integration code that can be sourced or used interactively.

## Installation

### Install via Nix (Recommended)

The easiest way to get kfg is via Nix. Pre-built binaries are available from GitHub Releases:

```bash
# Build from GitHub Releases
nix build github:seregatte/kfg

# Run directly without installing
nix run github:seregatte/kfg -- --help

# Add to current shell temporarily
nix shell github:seregatte/kfg
```

This works on Linux and macOS (x86_64 and ARM64).

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

### Apply a Kustomization

```bash
# Apply a kustomization directory
kfg apply -k path/to/kustomization --workflow myworkflow

# Apply with explicit file
kfg apply -f manifest.yaml --workflow myworkflow

# Apply from stdin
kfg apply -f - --workflow myworkflow
```

### Run an Agent

```bash
# Run a specific agent
kfg run -k path/to/kustomization myagent

# List available agents
kfg run -k path/to/kustomization

# Run with arguments
kfg run -k path/to/kustomization myagent -- --option value
```

## Command Reference

| Command | Alias | Description |
|---------|-------|-------------|
| `kfg apply` | | Apply a kustomization or manifest file |
| `kfg run` | | Run an agent one-shot |
| `kfg build` | | Build kustomization to YAML |
| `kfg sys log` | | System logging (internal) |
| `kfg sys cache` | | Step cache management |
| `kfg version` | | Show version information |

## Environment Variables

| Variable | Description |
|----------|-------------|
| `KFG_VERBOSE` | Verbosity level (0-5) |
| `KFG_STORE_DIR` | Store directory (default: ~/.kfg/store) |
| `KFG_LOG_FILE` | Log file path |
| `KFG_LOG_DIR` | Log directory |
| `KFG_LOG_COLOR` | Log color mode (auto/always/never) |

## API Version

KFG uses the `kfg.dev/v1alpha1` API version for manifests:

```yaml
apiVersion: kfg.dev/v1alpha1
kind: Cmd
metadata:
  name: example
spec:
  run: echo "Hello, World!"
```

## Development

Enter the Nix dev shell as your development entrypoint. It provides Go, Node.js, bats, and all AI agent tools:

```bash
nix develop

# Once inside the dev shell:
make build        # Build the binary → ./bin/kfg
make test         # Go unit tests
make test-bats    # Bats integration tests
make fmt          # Format code
make lint         # Run linter
```

## Repository Structure

KFG uses a package-oriented architecture:

```
├── src/                          # Engine implementation (Go)
│   ├── cmd/kfg/                  # CLI commands
│   └── internal/                 # Internal packages
├── packages/
│   ├── framework/                # Shared manifest primitives
│   │   ├── manifests/            # Reusable steps (materialize, cleanup, etc.)
│   │   └── tests/                # Framework test suite
│   └── domains/
│       └── ai-agents/            # AI agents domain package
│           ├── manifests/        # AI agent resources
│           ├── overlays/dev/     # Development overlay
│           └── tests/            # Domain test suite
├── docs/
│   ├── AGENTS.md                 # AI agent operating context
│   └── context/
│       └── openspec/             # Unified OpenSpec root
├── tests/
│   └── bats/                     # Engine and integration tests
│       ├── cli/                  # CLI command tests
│       ├── workflows/            # Runtime workflow tests
│       └── helpers/              # Shared test helpers
└── Makefile                      # Build and test targets
```

**Public Entrypoints:**
- Engine CLI: `./bin/kfg`
- Framework package: `packages/framework/kustomization.yaml`
- AI agents domain: `packages/domains/ai-agents/kustomization.yaml`
- Domain overlay (dev): `packages/domains/ai-agents/overlays/dev/`

For more details, see [AGENTS.md](AGENTS.md).

## License

MIT License
