# Store & Images

kfg uses an image layer system for configuration management. Images are immutable snapshots of configuration files that get materialized into workspace instances.

## Imagefile

Dockerfile-like syntax parsed by `kfg internal package (github.com/seregatte/kfg)

```dockerfile
FROM base:v2 AS base
COPY config/ config/

FROM scratch
COPY --from=base config/ config/
COPY settings.json settings.json
TAG my-config:v1
```

**Instructions**: `FROM`, `COPY`, `COPY --from=<stage>`, `ENV`, `WORKDIR`, `RUN`, `TAG`.

**RUN limitation**: Use relative paths only. Stages are isolated — files from other stages must be explicitly copied with `COPY --from`.

## Image Operations

```bash
kfg image build --name myconfig:latest    # Build from Imagefile in cwd
kfg image build --push                     # Build + push (fails if exists)
kfg image build --push --keep-build        # Build + push, keep build dir
kfg image list                             # List all images
kfg image inspect myconfig:latest          # Metadata
kfg image inspect myconfig:latest --files  # Artifact paths
kfg image inspect myconfig:latest --recipe # Imagefile content
kfg image push myconfig:latest             # Push existing image
kfg image remove myconfig:latest           # Remove image
```

**`--push`**: Auto-pushes after build. Fails if image already exists (immutability).

## Workspace Instances

```bash
kfg workspace start myconfig:latest              # Materialize (default instance)
kfg workspace start myconfig:latest --name proj1 # Named instance
kfg workspace stop --name proj1                  # Restore + cleanup
```

## Artifact-Scoped Backup/Cleanup

Only files that **conflict** with image artifacts are backed up and restored:

- **On `start`**: Conflicting workspace files backed up to `$KFG_STORE_DIR/.workspace/<instance>/backup/`
- **On `stop`**: Only materialized paths removed. Backup restored. Unrelated files untouched.

Example:
```
# Workspace has: README.md, settings.json
# Image has: settings.json, config.json

# After start: settings.json backed up, image files extracted
# After stop: settings.json restored from backup, config.json removed
# README.md never touched
```

## Store Directory Structure

```
$KFG_STORE_DIR/
├── images/<name>/<tag>/
│   ├── Imagefile
│   ├── artifacts/         # Materialized files
│   └── metadata.json      # digest, created_at, artifacts list
└── .workspace/<instance>/
    ├── backup/data/       # Only conflicting files
    └── instance.json      # name, image_ref, materialized_paths
```