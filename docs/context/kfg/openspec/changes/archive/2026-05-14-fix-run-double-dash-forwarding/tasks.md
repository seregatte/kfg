## 1. Parser Update

- [x] 1.1 Update `parseLaunchArgs` in `src/cmd/kfg/run.go` to use Cobra's dash boundary instead of scanning `args` for `--`
- [x] 1.2 Preserve current behavior for no-arg and agent-only invocations while forwarding args after `--`

## 2. Test Coverage

- [x] 2.1 Update `src/cmd/kfg/run_test.go` to cover no args, agent only, agent with forwarded args, separator-only, and multiple forwarded args
- [x] 2.2 Add integration coverage for documented `kfg run ... -- ...` behavior if an existing test harness can exercise the command end-to-end
- [x] 2.3 Verify shell execution still forwards `extraArgs` unchanged into the generated agent invocation

## 3. Validation

- [x] 3.1 Run the relevant Go unit tests for `src/cmd/kfg`
- [x] 3.2 Run any targeted integration or shell validation needed to confirm the documented CLI UX matches actual behavior

## 4. Documentation

- [x] 4.1 Confirm CLI docs and examples for `kfg run` remain accurate after the parser fix
