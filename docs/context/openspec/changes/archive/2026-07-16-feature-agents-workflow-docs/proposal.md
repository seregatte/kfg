## Why

The development workflow used in practice differs from what's documented in AGENTS.md. The document lacks: the full OpenSpec cycle per feature, the rule that each change must live inside its worktree, PR target rules, and CLI caveats. Updating it ensures consistent behavior across future developments.

## What Changes

- **Add**: "Feature Development Workflow" section documenting the full OpenSpec cycle
- **Update**: Git Worktree Workflow — clarify PR base is the originating branch
- **Fix**: Stale `kfg-wizard.yaml` path reference → `wizard.yaml`
- **Add**: CLI caveat about `openspec new change` auto-detection

## Impact

- **Documentation only**: No code changes. All changes in `docs/AGENTS.md`.
