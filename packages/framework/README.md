# kfg Framework Package

The framework package (`packages/framework/`) provides shared manifest primitives and reusable steps that domain packages can compose and extend.

## Public API

### Kustomization Entrypoint

**File:** `packages/framework/kustomization.yaml`

Consumers should reference this kustomization, not internal framework paths. The public kustomization composes all exported framework resources.

### Exported Steps

The framework exports five core reusable steps:

1. **`kfg.materialize`** - Generate shell code from manifests using converters
   - Supports per-item and aggregate modes
   - Handles asset-to-converter-to-output composition
   - Manages artifact registration

2. **`kfg.cleanup`** - Clean up generated artifacts
   - Removes all artifacts registered in the current session
   - Suitable for use as an after-step
   - Runs even if earlier steps failed

3. **`kfg.ensure-gitignore`** - Manage repository `.gitignore` entries
   - Adds entries to `.gitignore` without duplication
   - Creates the file if it doesn't exist
   - Prevents build artifacts from being tracked

4. **`kfg.copy-context`** - Copy context files into artifacts
   - Supports copying multiple files to a destination
   - Preserves directory structure
   - Useful for populating artifacts with supporting files

5. **`kfg.materialize-scaffold`** - Generate scaffolding from templates
   - Generates project/feature scaffolding structures
   - Supports template variable substitution
   - Typical usage: new project or feature setup

## Using Framework Steps

Domain packages compose the framework package and use its exported steps in their workflows.

### Example: Domain Composition

```yaml
apiVersion: kfg.dev/v1alpha1
kind: CmdWorkflow
metadata:
  name: my-domain-workflow
spec:
  cmds:
    - cmd1
  before:
    - step: kfg.materialize
      env:
        MODE: aggregate
        ASSETS: asset1:asset2
        CONVERTER: my-converter
        OUTPUTS: /tmp/output.yaml
```

### Shell Runtime API

Framework steps depend on the stable shell runtime API exported by the engine:

- `KFG_SESSION_ID` - Unique session identifier
- `KFG_WORKFLOW_NAME` - Current workflow name
- `KFG_ARTIFACTS` - Colon-separated artifact list
- `__kfg_add_artifact()` - Register artifacts
- `__kfg_build_result()` - Record build results
- `_kfg.log.*()` - Structured logging functions

Framework steps MUST NOT hardcode engine-specific paths or make assumptions about the engine's internal structure. They MUST use only the documented runtime API.

## Framework Specifications

Framework-specific capabilities are documented in `docs/context/framework/openspec/specs/`:

- **reusable-framework-steps** - Contract and behavior of exported steps
- **framework-package-contract** - Public API and stability guarantees

See `docs/context/framework/openspec/README.md` for the framework OpenSpec root.

## Framework Tests

Framework-specific Bats suites are in `packages/framework/tests/`:

- Step functionality tests
- Shell generation validation
- Artifact registration and cleanup
- Gitignore entry management

Run framework tests with:

```bash
bats packages/framework/tests/**/*.bats
```

Or use the canonical test target:

```bash
make test-bats  # Runs all tests including framework
```

## Framework Development

### Adding a New Framework Step

1. Create the step manifest in `packages/framework/manifests/steps/<name>.yaml`
2. Add it to `packages/framework/manifests/steps/kustomization.yaml`
3. Create a spec in `docs/context/framework/openspec/specs/` documenting the step
4. Add Bats tests in `packages/framework/tests/`
5. Update this documentation

### Updating Framework Behavior

Framework changes are tracked in `docs/context/framework/openspec/changes/`. When making changes that affect framework behavior:

1. Create a change proposal in `docs/context/framework/openspec/changes/`
2. Update relevant specs
3. Update this documentation
4. Add or update Bats tests
5. If cross-layer (affects engine or domains), create sibling changes with matching slugs

### Backward Compatibility

Framework maintains a backward-compatibility commitment:

- Exported step names and behavior remain stable
- New steps can be added without breaking existing consumers
- Breaking changes require a deprecation period and migration documentation
- The shell runtime API is stable (see shell-runtime-api spec)

## Integration with Domain Packages

Domain packages use the framework by:

1. Basing their kustomization on `packages/framework/`
2. Composing framework steps in their workflows
3. Extending framework behavior where needed
4. Testing against the framework API

See `packages/domains/ai-agents/kustomization.yaml` for an example.

## Further Reading

- **Shell Runtime API:** `docs/context/kfg/openspec/specs/shell-runtime-api/spec.md`
- **Framework Package Contract:** `docs/context/kfg/openspec/specs/framework-package-contract/spec.md`
- **Reusable Steps:** `docs/context/framework/openspec/specs/reusable-framework-steps/spec.md`
