## Why

`kfg run` documents support for forwarding arguments after `--`, but the current parser looks for the separator in Cobra's already-processed `args` slice. As a result, extra agent arguments are silently dropped and the documented CLI contract is broken.

## What Changes

- Fix `kfg run` argument parsing to use Cobra's dash boundary information instead of searching `args` for a literal `--`.
- Ensure `kfg run agent -- <args...>` forwards only the arguments after `--` to the agent function.
- Add unit coverage for the dash boundary cases and no-dash cases.
- Keep the documented `kfg run` UX aligned with implementation and specs.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `run-command`: require `kfg run` to forward arguments after `--` to the selected agent while consuming the separator itself.
- `cli-conventions`: keep the one-shot run command examples and guarantees aligned with `kfg run` behavior.

## Impact

- Affected code: `src/cmd/kfg/run.go`
- Affected tests: `src/cmd/kfg/run_test.go`
- Affected specs: `run-command`, `cli-conventions`
- User-facing CLI UX: documented `kfg run ... -- ...` behavior starts working as described
