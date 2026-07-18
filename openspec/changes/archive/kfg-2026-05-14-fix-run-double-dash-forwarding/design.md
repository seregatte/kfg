## Context

`kfg run` documents support for forwarding extra agent arguments after `--`, but the implementation searches the `args` slice for a literal `--`. Cobra strips the separator before `RunE` receives positional arguments and exposes the split point through command metadata instead. Because the parser ignores that metadata, documented extra arguments are dropped before the generated agent function runs.

The bug is isolated to command parsing in `src/cmd/kfg/run.go`, but it directly affects the user-facing CLI contract and should be locked down with unit coverage.

## Goals / Non-Goals

**Goals:**
- Make `kfg run agent -- <args...>` forward every argument after `--` to the agent.
- Keep `kfg run agent` behavior unchanged when no `--` is present.
- Preserve current agent discovery behavior when no agent name is supplied.
- Add targeted tests for Cobra dash-boundary cases.

**Non-Goals:**
- Change shell script execution or temp file lifecycle.
- Introduce custom flag parsing or disable Cobra flag parsing globally.
- Change how agents are selected or how workflows are resolved.

## Decisions

### Use `cmd.ArgsLenAtDash()` as the source of truth

`parseLaunchArgs` should derive the split point from Cobra's recorded dash boundary instead of searching for a literal separator in `args`. When `ArgsLenAtDash()` returns `-1`, the command behaves as it does today: first positional arg is the agent and there are no forwarded extras. When it returns `0`, there is no agent name before the separator and the remaining args are extra args only.

Alternative considered: enable custom parsing or `DisableFlagParsing`. Rejected because Cobra already provides the needed boundary information and changing global parsing behavior would be a much riskier fix.

### Keep forwarding behavior localized to `parseLaunchArgs`

The smallest correct fix is to keep all dash-boundary logic in `parseLaunchArgs` and leave `executeAgent` unchanged. `executeAgent` already forwards `extraArgs` with `"$@"`; the bug is that `extraArgs` is computed incorrectly.

Alternative considered: reconstruct the argument split later in `executeAgent`. Rejected because that duplicates parsing concerns and would make tests less direct.

### Add unit tests for dash-boundary edge cases

Tests should cover:
- no args
- agent only
- agent with `--` and extra args
- separator with no agent before it
- multiple extra args after the separator

Alternative considered: rely on integration tests only. Rejected because this behavior is deterministic and easier to specify at the parser boundary.

## Risks / Trade-offs

- [Tests may not reflect Cobra's real dash metadata] -> Build test cases with a command instance that explicitly sets the dash boundary state used by `parseLaunchArgs`.
- [Parser changes could affect agent discovery with no args] -> Keep the no-arg and no-dash paths covered in unit tests.
- [Future commands may copy the old parsing pattern] -> Keep the behavior documented in `run-command` and `cli-conventions` so regressions are easier to spot.

## Migration Plan

1. Update `parseLaunchArgs` to use `cmd.ArgsLenAtDash()`.
2. Add or update unit tests for dash-boundary scenarios.
3. Verify `kfg run` documentation and specs still match the implemented behavior.

## Open Questions

- Whether any existing tests already create Cobra commands with a dash boundary that can be reused.
