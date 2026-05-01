# AGENTS.md

This file provides guidance to AI agents when working with code in the kfg repository.

kfg is a standalone CLI for processing YAML manifests into shell functions. It's written in Go using Cobra for CLI framework and Viper for configuration.

## Quick Reference

```bash
# Build
make build         # Builds to ./bin/kfg

# Test
make test          # Go unit tests: cd src && go test ./...
make test-bats     # Bats integration tests: bats tests/bats

# Development
make fmt           # Format code
make lint          # Run linter (requires golangci-lint)
make vet           # Run go vet
```

## Architecture

### Pipeline

```
YAML manifests → Load → Kustomize Merge → Validate → Generate Shell (Go templates)
```

### CLI Structure

```
src/cmd/kfg/       # Cobra CLI commands
├── main.go        # Entry point
├── root.go        # Root command
├── build.go       # kfg build
├── apply.go       # kfg apply
├── image.go       # kfg image (img alias)
├── workspace.go   # kfg workspace (ws alias)
├── run.go         # kfg run
├── sys.go         # kfg sys (log subcommand)
└── version.go     # kfg --version
```

### Internal Packages

```
src/internal/
├── config/          # Viper configuration management
├── generate/        # Shell code generation + embedded templates/
├── image/           # Image builder, materializer, metadata
├── imagefile/       # Dockerfile-like parser (FROM, COPY, RUN, TAG)
├── kustomize/       # Loader, adapter, schema validation
├── logger/          # Zerolog structured logging
├── manifest/        # YAML parser + types (Cmd, Step, CmdWorkflow)
├── resolve/         # Step dependency resolution
├── resolver/        # {env:VAR} placeholder resolution
├── store/           # Local store (images, workspace)
```

### Manifest Model

**Resource Kinds** (apiVersion: `kfg.dev/v1alpha1`):
- **Cmd**: Pure shell function (`spec.run`, `spec.env`, `metadata.commandName`)
- **CmdWorkflow**: Orchestration — lists cmds, `before`/`after` steps
- **Step**: Reusable unit of work (`spec.run`, `spec.env`, optional `spec.output`)
- **Schema / Assets / Converter**: Source-layer data and transformations

**CmdWorkflow structure**:
```yaml
apiVersion: kfg.dev/v1alpha1
kind: CmdWorkflow
metadata:
  name: workflow-name
  shell: bash
spec:
  cmds: [cmd1, cmd2]
  before:
    - step: setup-step
      failurePolicy: Ignore  # or Fail (default)
  after:
    - step: cleanup-step
```

**Step reference fields**: `step` (required), `weight`, `when`, `failurePolicy`, `env`.

**Placeholders**: `{env:VAR_NAME}` resolves to `$VAR` at generation time, then shell expands at runtime.

### Store-Centric Architecture

`kfg image` and `kfg workspace` manage immutable configuration images:

```
$KFG_STORE_DIR/
├── images/<name>/<tag>/
│   ├── Imagefile
│   ├── artifacts/
│   └── metadata.json
└── .workspace/<instance>/
    ├── backup/
    └── instance.json
```

**Imagefile**: `FROM`, `COPY`, `COPY --from=<stage>`, `ENV`, `WORKDIR`, `RUN`, `TAG`.

**Artifact-scoped backup/cleanup**: Only conflicting files are backed up/restored.

## Testing

### Go Unit Tests

Tests are in `src/internal/*_test.go` files. Run with `make test`.

Key test areas:
- Manifest parsing and validation
- Step dependency resolution
- Placeholder resolution
- Template generation
- Imagefile parsing
- Store operations

### Bats Integration Tests

Integration tests are in `tests/bats/`. Run with `make test-bats`.

Test files:
- `cli.bats`: CLI command tests
- `image_layer_system.bats`: Image build/push/start/stop
- `imagefile_features.bats`: Imagefile parsing
- `store_workspace.bats`: Workspace operations

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `KFG_MANIFEST_PATH` | `~/.config/kfg/manifests:.kfg/manifests` | Manifest paths (rightmost wins) |
| `KFG_VERBOSE` | `0` | 0=quiet, 1=error/warn/info, 2=+detail, 3=+debug |
| `KFG_STORE_DIR` | `~/.config/kfg/store` | Store directory |
| `KFG_LOG_FILE` | `$XDG_STATE_HOME/kfg/logs/kfg.log` | Log file path |
| `KFG_LOG_DIR` | `$XDG_STATE_HOME/kfg/logs` | Log directory |
| `KFG_LOG_COLOR` | `auto` | auto/always/never |
| `KFG_SESSION_ID` | auto-generated | Session ID for correlation |

## Code Conventions

### Error Messages

Use structured error format with component prefix:

```go
// Good: structured error with component
logger.Error("core:cli", "invalid argument", "arg", arg)

// Bad: unstructured error
log.Println("invalid argument: " + arg)
```

### CLI Flags

Follow Cobra conventions:
- Use `-k` for `--kustomize` (short flag first)
- Use `-f` for `--file`
- Use `-w` for `--workflow`
- Use `-o` for `--output`

### Shell Generation

Generated shell code must:
- Start with `#!/bin/bash`
- Use `__kfg_` prefix for internal helpers
- Export metadata env vars (`KFG_*`)
- Include session ID generation
- Support before/after steps

## Local State & Gotchas

- Version is injected via ldflags at build time (see Makefile)
- The binary must be built before running Bats tests
- Store directory is created on first use
- Log files use `.log` extension (not `.jsonl`)