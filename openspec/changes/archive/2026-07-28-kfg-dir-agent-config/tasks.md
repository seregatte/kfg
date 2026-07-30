## 1. Update Pi Cmd

- [x] 1.1 Update `manifests/cmds/agents/pi.yaml` — add KFG_DIR setup bash preamble
      to `spec.run` and export `PI_CODING_AGENT_DIR`
- [x] 1.2 Verify trap cleanup only fires when KFG created the temp dir

## 2. Update OpenCode Cmd

- [x] 2.1 Update `manifests/cmds/agents/opencode.yaml` — add KFG_DIR setup bash preamble
      to `spec.run` and export `OPENCODE_CONFIG_DIR`
- [x] 2.2 Verify trap cleanup only fires when KFG created the temp dir

## 3. Update Documentation

- [x] 3.1 Update `manifests/cmds/README.md` — add KFG_DIR usage description

## 4. Validation

- [x] 4.1 Run `go test ./internal/generate/ -count=1` — all tests pass
- [x] 4.2 Verify generated shell output includes the KFG_DIR preamble
