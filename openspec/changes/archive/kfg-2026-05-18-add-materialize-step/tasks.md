## 1. Shared Step Contract

- [x] 1.1 Add `.manifests/base/steps/materialize.yaml` with required `MODE`, `ASSETS`, `CONVERTER`, and `OUTPUTS` inputs plus optional `WRAP_KEY`
- [x] 1.2 Implement `MODE=per-item` with positional `ASSETS[i] -> OUTPUTS[i]` mapping and strict count validation
- [x] 1.3 Implement `MODE=aggregate` with multi-document deep merge, optional `WRAP_KEY`, and default merge-with-existing output behavior
- [x] 1.4 Register every generated output path in `KFG_ARTIFACTS` and fail fast on missing or invalid configuration

## 2. Workflow Migration

- [x] 2.1 Replace `kfg.agents.steps.settings` usage in `.manifests/overlay/dev/agents-workflow.yaml` with `kfg.materialize` in `per-item` mode
- [x] 2.2 Replace `kfg.convert` usage in `.manifests/overlay/dev/agents-workflow.yaml` with `kfg.materialize` grouped by agent and type
- [x] 2.3 Replace `kfg.aggregate-mcp` usage in `.manifests/overlay/dev/agents-workflow.yaml` with `kfg.materialize` in `aggregate` mode
- [x] 2.4 Remove obsolete specialized step manifests or update their references so the shared base exposes only the new stable primitive

## 3. Validation

- [x] 3.1 Add or update Bats coverage for `kfg.materialize` per-item mode with single-item and multi-item conversions
- [x] 3.2 Add or update Bats coverage for aggregate mode, including merge, wrap, and merge-with-existing output behavior
- [x] 3.3 Add negative Bats coverage for missing required vars, mismatched `ASSETS`/`OUTPUTS`, invalid mode combinations, and unsupported aggregate output cardinality
- [x] 3.4 Add unit or integration coverage where needed for any Go-side assumptions affected by the step migration (N/A - no Go-side changes required, migration is manifest/shell-level only)
- [x] 3.5 Run the relevant Bats and Go test suites to confirm deterministic workflow behavior after migration

## 4. Documentation and Specs

- [x] 4.1 Add a new `materialize-step` spec describing the stable shared step contract and validation rules
- [x] 4.2 Update `dev-workflow` spec and any workflow documentation examples to use `kfg.materialize`
- [x] 4.3 Update manifest or developer docs that currently describe `kfg.convert`, `kfg.aggregate-mcp`, or settings materialization as separate step contracts
