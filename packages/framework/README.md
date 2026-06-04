# kfg Framework Package

The framework package (`packages/framework/`) provides shared manifest primitives and reusable steps for domain packages to compose and extend.

## Public API

### Kustomization Entrypoint

**File:** `packages/framework/kustomization.yaml`

Consumers reference this kustomization, not internal framework paths.

### Exported Steps

| Step | Purpose |
|------|---------|
| `kfg.materialize` | Generate shell code from manifests using converters (per-item or aggregate modes) |
| `kfg.cleanup` | Clean up all artifacts registered in current session |
| `kfg.ensure-gitignore` | Add entries to `.gitignore` without duplication |
| `kfg.copy-context` | Copy context files into artifacts, preserving structure |
| `kfg.materialize-scaffold` | Generate scaffolding from templates |

## Using Framework Steps

Domain packages compose the framework and use its exported steps in workflows:

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

## Shell Runtime API

Framework steps depend on the stable shell runtime API:

- `KFG_SESSION_ID` - Unique session identifier
- `KFG_WORKFLOW_NAME` - Current workflow name
- `KFG_ARTIFACTS` - Colon-separated artifact list
- `__kfg_add_artifact()` - Register artifacts
- `__kfg_build_result()` - Record build results
- `__kfg_log_*()` - Structured logging functions

Framework steps MUST use only documented runtime API. See `docs/context/openspec/specs/kfg-shell-runtime-api/spec.md`.

## Framework Tests

Bats suites in `packages/framework/tests/`:

```bash
make test-bats  # Canonical entrypoint (all tests including framework)
```

Tests cover: step functionality, shell generation, artifact registration, gitignore management.

## Framework Development

### Adding a Step

1. Create manifest in `packages/framework/manifests/steps/<name>.yaml`
2. Add to `packages/framework/manifests/steps/kustomization.yaml`
3. Create spec: `docs/context/openspec/specs/framework-<name>/spec.md`
4. Add Bats tests in `packages/framework/tests/`
5. Update this documentation

### Backward Compatibility

Framework maintains stability:
- Exported step names and behavior remain stable
- New steps added without breaking existing consumers
- Breaking changes require deprecation period + migration docs
- Shell runtime API is stable (see spec)

## Integration with Domain Packages

Domain packages use framework by:

1. Basing kustomization on `packages/framework/`
2. Composing framework steps in workflows
3. Extending framework behavior where needed
4. Testing against framework API

See `packages/domains/ai-agents/kustomization.yaml` for example.

## Further Reading

- **Shell Runtime API:** `docs/context/openspec/specs/kfg-shell-runtime-api/spec.md`
- **Framework Package Contract:** `docs/context/openspec/specs/kfg-framework-package-contract/spec.md`
- **Reusable Steps:** `docs/context/openspec/specs/framework-reusable-framework-steps/spec.md`
