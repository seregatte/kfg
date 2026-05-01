# KFG - Declarative Shell Compiler

KFG is a declarative shell compiler that transforms YAML manifests into bash functions. It allows you to define shell commands, their dependencies, and execution steps in YAML manifests, then generates shell integration code that can be sourced or used interactively.

## Installation

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
kfg apply .nixai/overlay/dev

# Apply with explicit file
kfg apply -f manifest.yaml

# Apply with specific workflow
kfg apply .nixai/overlay/dev --workflow dev
```

### Run an Agent

```bash
# Run a specific agent
kfg run .nixai/overlay/dev claude

# List available agents
kfg run .nixai/overlay/dev

# Run with arguments
kfg run .nixai/overlay/dev claude -- --model gpt-4
```

### Build an Image

```bash
# Build an image
kfg image build --name myconfig:latest

# Build and push
kfg image build --name myconfig:latest --push

# List images (alias: kfg img ls)
kfg image ls
```

### Manage Workspaces

```bash
# Start a workspace
kfg workspace start myconfig:latest

# Stop a workspace
kfg workspace stop --name myinstance

# List workspaces (alias: kfg ws ls)
kfg workspace ls
```

## Command Reference

| Command | Alias | Description |
|---------|-------|-------------|
| `kfg apply` | | Apply a kustomization or manifest file |
| `kfg run` | | Run an agent one-shot |
| `kfg build` | | Build kustomization to YAML |
| `kfg image` | `img` | Image management commands |
| `kfg workspace` | `ws` | Workspace management commands |
| `kfg sys log` | | System logging (internal) |
| `kfg version` | | Show version information |

## Environment Variables

| Variable | Description |
|----------|-------------|
| `KFG_VERBOSE` | Verbosity level (0-3) |
| `KFG_STORE_DIR` | Store directory (default: ~/.config/kfg/store) |
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

```bash
# Run tests
make test

# Run tests with coverage
make test-coverage

# Format code
make fmt

# Run linter
make lint
```

## License

MIT License

## Migration from NixAI

For users migrating from NixAI, see the [Migration Guide](docs/migration.md).

Key changes:
- Command structure: `nixai store image build` → `kfg image build` (or `kfg img build`)
- Environment variables: `NIXAI_*` → `KFG_*`
- API version: `nixai.dev/v1alpha1` → `kfg.dev/v1alpha1`
- Config paths: `~/.config/nixai/` → `~/.config/kfg/`# kfg
