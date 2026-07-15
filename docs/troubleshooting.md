# Troubleshooting

This guide helps diagnose and resolve common issues with KFG.

## Installation Issues

### `kfg: command not found`

**Cause**: The KFG binary is not in your PATH.

**Solutions**:

1. **Via Nix**: Add to current shell
   ```bash
   nix shell github:seregatte/kfg
   ```

2. **Build from source**: Add to PATH
   ```bash
   export PATH="$PWD/bin:$PATH"
   # Or copy to a PATH directory
   sudo cp bin/kfg /usr/local/bin/
   ```

3. **GOPATH**: Check if it's in your PATH
   ```bash
   export PATH="$GOPATH/bin:$PATH"
   ```

### Nix: `flakes experimental feature is disabled`

**Cause**: Flakes are not enabled in Nix.

**Solution**: Enable flakes in `~/.config/nix/nix.conf`:
```
experimental-features = nix-command flakes
```

Or use the `--extra-experimental-features` flag:
```bash
nix --extra-experimental-features "nix-command flakes" shell github:seregatte/kfg
```

### Build error: `go: command not found`

**Cause**: Go is not installed or not in your PATH.

**Solution**:
```bash
# Install Go (Linux)
sudo apt install golang-go

# Or via Nix
nix shell nixpkgs#go
```

Check the version (requires 1.21+):
```bash
go version
```

## Execution Issues

### `Error: manifest validation failed`

**Cause**: Invalid YAML manifest.

**Diagnosis**:
```bash
# Increase verbosity
KFG_VERBOSE=3 kfg apply -f manifest.yaml --workflow test

# Validate YAML
python3 -c "import yaml; yaml.safe_load(open('manifest.yaml'))"
```

**Common causes**:
- Incorrect `apiVersion` (must be `kfg.dev/v1alpha1`)
- Invalid `kind` (Cmd, CmdWorkflow, Step, etc.)
- Missing required fields
- Incorrect YAML indentation

**Correct example**:
```yaml
apiVersion: kfg.dev/v1alpha1  # Required
kind: Cmd                      # Required
metadata:
  name: myapp.cmd.example      # Required
  commandName: example         # Required for Cmd
spec:
  run: echo "hello"            # Required
```

### `Error: workflow not found`

**Cause**: The specified workflow does not exist in the manifest.

**Solution**:
```bash
# List available workflows
kfg build -k ./manifests | grep "kind: CmdWorkflow" -A 2

# Or check the exact name
grep "name:.*workflow" manifest.yaml
```

### `Error: circular dependency detected`

**Cause**: Steps form a cycle in the DAG.

**Problematic example**:
```yaml
before:
  - step: step-a
    when:
      output:
        step: step-b  # step-b depends on step-a
        name: STATUS
        equals: "ok"
  - step: step-b
    when:
      output:
        step: step-a  # step-a depends on step-b
        name: STATUS
        equals: "ok"
```

**Solution**: Reorganize dependencies to avoid cycles.

### Command not working after `kfg apply`

**Cause**: The generated code was not sourced into the current shell.

**Solution**:
```bash
# Option 1: Use --interactive to open an interactive shell
kfg apply -f manifest.yaml --workflow test --interactive

# Option 2: Source manually
eval "$(kfg apply -f manifest.yaml --workflow test --print)"

# Option 3: Add to .bashrc
echo 'eval "$(kfg apply -f manifest.yaml --workflow test --print)"' >> ~/.bashrc
```

## Cache Issues

### `Error: cache corrupted`

**Cause**: Corrupted cache files.

**Solution**:
```bash
# Clear all cache
rm -rf ~/.kfg/store/cache

# Or remove specific entries
kfg sys cache ls
kfg sys cache rm <id>

# Rebuild on next run
KFG_REFRESH=1 kfg apply -f manifest.yaml --workflow test
```

### Steps are not being cached

**Cause**: Step is not marked as cacheable.

**Solution**:
```yaml
kind: Step
metadata:
  name: myapp.steps.expensive
  annotations:
    kfg.dev/cacheable: "true"  # Add this annotation
spec:
  run: ...
```

### Cache not invalidating after manifest change

**Cause**: The step hash hasn't changed (same code).

**Solution**:
```bash
# Force invalidation
KFG_REFRESH=1 kfg apply -f manifest.yaml --workflow test

# Or remove cache manually
kfg sys cache rm <step-id>
```

## Logging Issues

### Logs are not appearing

**Cause**: Verbosity is too low.

**Solution**:
```bash
# Increase verbosity
KFG_VERBOSE=5 kfg apply -f manifest.yaml --workflow test
```

**Verbosity levels**:
- `0`: Quiet (no output)
- `1`: Error + Warn + Info (default)
- `2`: + Detail
- `3`: + Warn/Detail
- `4`: + Debug
- `5`: + Debug verbose

### Where are the logs?

**Default location**:
```bash
# Linux
~/.local/state/kfg/logs/kfg.jsonl

# macOS
~/Library/Application Support/kfg/logs/kfg.jsonl

# Or custom path
echo $KFG_LOG_FILE
```

**View logs in real time**:
```bash
tail -f ~/.local/state/kfg/logs/kfg.jsonl | jq .
```

### Logs are too verbose

**Solution**: Reduce verbosity
```bash
KFG_VERBOSE=1 kfg apply -f manifest.yaml --workflow test
```

## Performance Issues

### `kfg apply` is too slow

**Possible causes**:

1. **Too many manifests to process**
   ```bash
   # Diagnose
   KFG_VERBOSE=3 kfg apply -f manifest.yaml --workflow test 2>&1 | grep "Parsing"
   
   # Solution: Use kustomization instead of multiple -f flags
   kfg apply -k ./manifests --workflow test
   ```

2. **Steps are not cached**
   ```bash
   # Check cache
   kfg sys cache ls
   
   # Solution: Mark steps as cacheable
   annotations:
     kfg.dev/cacheable: "true"
   ```

3. **Kustomize processing complex overlays**
   ```bash
   # Diagnose
   KFG_VERBOSE=3 kfg build -k ./manifests
   
   # Solution: Simplify overlays or use base directly
   kfg apply -k ./base --workflow test
   ```

### Excessive memory usage

**Cause**: Very large manifests or too many resources.

**Solution**:
```bash
# Split into multiple applications
kfg apply -f part1.yaml --workflow test1
kfg apply -f part2.yaml --workflow test2

# Or use kustomization for modular composition
kfg apply -k ./manifests --workflow test
```

## Debug Issues

### How to view the generated shell code?

```bash
# Generate without executing
kfg build -k ./manifests -o generated.yaml

# Or to stdout
kfg build -k ./manifests

# View code for a specific workflow
kfg build -k ./manifests --workflow test
```

### How to debug a specific step?

```bash
# Increase verbosity
KFG_VERBOSE=5 kfg apply -f manifest.yaml --workflow test

# Run step manually
bash -x -c "$(kfg build -k ./manifests | grep -A 50 '_kfg.step.my_step')"

# View step outputs
echo $KFG_STEP_my_step_STATUS
```

### How to check environment variables?

```bash
# View all KFG variables
env | grep KFG

# View variables for a specific step
kfg sys cache inspect <step-id> | jq .env
```

## Common Kustomize Issues

### `Error: failed to load kustomization`

**Cause**: Invalid kustomization.yaml or missing resources.

**Diagnosis**:
```bash
# Validate kustomization
kustomize build ./manifests

# Check referenced resources
grep "resources:" -A 10 ./manifests/kustomization.yaml
```

**Solution**: Verify that all referenced files exist.

### Overlays are not being applied

**Cause**: Incorrect resource order or patches.

**Diagnosis**:
```bash
# View build result
kustomize build ./overlays/dev

# Check patches
cat ./overlays/dev/kustomization.yaml
```

**Solution**: Verify patch syntax and resource order.

### Duplicate resources

**Cause**: Same resource defined in both base and overlay.

**Solution**: Use patches instead of redefining resources:
```yaml
# overlays/dev/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - ../../base
patches:
  - target:
      kind: Cmd
      name: base.cmd.deploy
    patch: |
      - op: replace
        path: /spec/env/ENV
        value: development
```

## Runtime Errors

### `Error: placeholder resolution failed`

**Cause**: Environment variable doesn't exist.

**Solution**:
```yaml
# Bad: variable may not exist
env:
  API_KEY: "{env:API_KEY}"

# Good: with default value
env:
  API_KEY: "{env:API_KEY:-default_key}"
```

### `Error: step output not available`

**Cause**: Trying to use output from a step that hasn't executed yet.

**Solution**: Ensure correct workflow order:
```yaml
before:
  - step: step-a  # Executes first
  - step: step-b  # Can use step-a's output
    when:
      output:
        step: step-a
        name: STATUS
        equals: "ok"
```

### `Error: command execution failed`

**Cause**: Shell command failed (exit code != 0).

**Diagnosis**:
```bash
# Increase verbosity
KFG_VERBOSE=5 kfg apply -f manifest.yaml --workflow test

# Run manually
bash -x -c "failing_command"
```

**Solution**: Fix the command or add error handling:
```yaml
spec:
  run: |
    set -e  # Fail on any error
    command || { echo "Error: command failed"; exit 1; }
```

## How to Report Bugs

If you found a bug that's not in this guide:

1. **Check existing issues**: https://github.com/seregatte/kfg/issues
2. **Gather information**:
   ```bash
   kfg version
   KFG_VERBOSE=5 kfg apply -f manifest.yaml --workflow test 2>&1 | tee debug.log
   ```
3. **Open an issue**: https://github.com/seregatte/kfg/issues/new
   - Describe the problem
   - Include `debug.log`
   - Include the YAML manifest (remove sensitive data)

## Additional Resources

- **[Getting Started](getting-started.md)** - Step-by-step tutorial
- **[CLI Reference](cli-reference.md)** - Full reference
- **[Architecture](architecture.md)** - Internal architecture
- **[GitHub Discussions](https://github.com/seregatte/kfg/discussions)** - Questions and discussions
- **[GitHub Issues](https://github.com/seregatte/kfg/issues)** - Report bugs