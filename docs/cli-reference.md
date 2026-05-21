# CLI Reference

## Commands

### `kfg build`

Build a kustomization and output the merged YAML.

The path argument is optional when `KFG_KPATH` is set. GitHub URLs are supported.

```bash
kfg build path/to/kustomization          # Output to stdout
kfg build path/to/kustomization -o out.yaml  # Output to file
kfg build https://github.com/owner/repo//path  # From GitHub URL
kfg build https://github.com/owner/repo//path?ref=v1.0.0  # With tag reference
KFG_KPATH=./manifests kfg build         # Using env var
```

### `kfg apply`

Apply a kustomization and generate shell code.

The `-k` flag is optional when `KFG_KPATH` is set. GitHub URLs are supported.

```bash
kfg apply -k path/to/kustomization --workflow myworkflow           # Generate + source
kfg apply -k path/to/kustomization --workflow myworkflow --interactive  # Interactive shell
kfg apply -k path/to/kustomization --workflow myworkflow --cmds cmd1,cmd2  # Specific cmds
kfg apply -k https://github.com/owner/repo//path                    # From GitHub URL
kfg apply -f manifest.yaml                               # From file
kfg apply -f -                                           # From stdin
KFG_KPATH=./manifests kfg apply                         # Using env var
```

### `kfg sys gc`

Garbage collection commands for managing Step cache entries.

```bash
kfg sys gc ls                    # List cache entries with metadata
kfg sys gc inspect <id>          # Show detailed metadata for an entry
kfg sys gc rm <id> [<id>...]     # Remove specific cache entries
kfg sys gc prune                 # Remove entries older than 30 days
kfg sys gc du                    # Show disk usage of cache entries
```

Cache entries are stored under `KFG_STORE_DIR/cache` (defaults to `~/.kfg/store/cache`).

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
| `KFG_KPATH` | (empty) | Default kustomization source path or GitHub URL |
| `KFG_REFRESH` | (empty) | Set to "1" to force refresh of cached Steps (bypasses cache) |
| `KFG_VERBOSE` | `0` | `0`=quiet, `1`=error/warn/info, `2`=+detail, `3`=+debug |
| `KFG_STORE_DIR` | `~/.config/kfg/store` | Store directory for cache entries |
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
