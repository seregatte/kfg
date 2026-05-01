# CLI Reference

## Commands

### `kfg build`

Build a kustomization and output the merged YAML.

```bash
kfg build path/to/kustomization          # Output to stdout
kfg build path/to/kustomization -o out.yaml  # Output to file
```

### `kfg apply`

Apply a kustomization and generate shell code.

```bash
kfg apply -k path/to/kustomization --workflow myworkflow           # Generate + source
kfg apply -k path/to/kustomization --workflow myworkflow --interactive  # Interactive shell
kfg apply -k path/to/kustomization --workflow myworkflow --cmds cmd1,cmd2  # Specific cmds
kfg apply -f manifest.yaml                               # From file
kfg apply -f -                                           # From stdin
```

### `kfg image` (alias: `img`)

Manage configuration images.

**Image operations**:
```bash
kfg image build --name myconfig:latest    # Build from Imagefile
kfg image build --push                     # Build + push (fails if exists)
kfg image build --push --keep-build        # Build + push, keep build dir
kfg image list [--json]                    # List images
kfg image inspect myconfig:latest          # Metadata
kfg image inspect myconfig:latest --files  # Artifact paths
kfg image inspect myconfig:latest --recipe # Imagefile only
kfg image push myconfig:latest             # Push existing image
kfg image remove myconfig:latest           # Remove image
```

### `kfg workspace` (alias: `ws`)

Manage workspace instances.
```bash
kfg workspace start myconfig:latest              # Materialize (default instance)
kfg workspace start myconfig:latest --name proj1 # Named instance
kfg workspace stop --name proj1                  # Restore + cleanup
```

### `kfg sys log`

Structured logging for use in shell scripts.

```bash
kfg sys log info "component" "message"
kfg sys log error "cmd:build" "failed to parse manifest"
kfg sys log debug "store:push" ""
kfg sys log --session-id "custom-123" info "component" "message"
```

Levels: `error`, `warn`, `info`, `detail`, `debug`. All levels persist to JSONL file regardless of verbosity.

### `kfg assets`

Transform Assets using Converters from build output.

Input source priority: `-f <file>` > stdin > `$KFG_BUILD_RESULT_FILE` > error.

```bash
kfg build path/to/kustomization | kfg assets convert --use myconverter
kfg assets convert --use myconverter -f build.yaml
kfg assets list -f build.yaml
```

### Command Aliases

Use `kfg apply -k` instead.

### Global Flags

| Flag | Description |
|------|-------------|
| `-d, --debug` | Enable debug mode (sets `KFG_VERBOSE=3`) |
| `-h, --help` | Show help |
| `--store <path>` | Override store directory (default: `$KFG_STORE_DIR`) |
| `--session-id <id>` | Override session ID for log correlation |

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `KFG_MANIFEST_PATH` | `~/.config/kfg/manifests:.kfg/manifests` | Colon-separated manifest paths (rightmost wins) |
| `KFG_VERBOSE` | `0` | `0`=quiet, `1`=error/warn/info, `2`=+detail, `3`=+debug |
| `KFG_STORE_DIR` | `~/.config/kfg/store` | Store directory |
| `KFG_LOG_FILE` | `$XDG_STATE_HOME/kfg/logs/kfg.jsonl` | Override JSONL log file path |
| `KFG_LOG_DIR` | `$XDG_STATE_HOME/kfg/logs` | Override log directory |
| `KFG_LOG_COLOR` | `auto` | `auto`/`always`/`never` |
| `KFG_DEBUG` | `false` | Set to `true` for debug mode |
| `KFG_SESSION_ID` | auto-generated | Session ID for log correlation |

## Exit Codes

| Code | Meaning |
|------|--------|
| `0` | Success |
| `1` | Error (manifest, validation, etc.) |
| `2` | Usage error (invalid flags, missing args) |

## Configuration Precedence

1. CLI flags
2. Environment variables
3. Default values
