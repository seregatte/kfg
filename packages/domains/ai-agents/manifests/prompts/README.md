# Prompts

Reusable prompt templates consumed by AI agent converters to generate
commands and skills.

## Resources

| Resource Name | Description |
|---------------|-------------|
| `ai.prompts.git-commit` | Generate a conventional commit message from staged changes |
| `ai.prompts.review-code` | Review code changes for quality, security, and best practices |
| `ai.prompts.review-search` | Search and review code using grep/ripgrep patterns |
| `ai.prompts.refactor-pure` | Refactor code with a pure functional approach |

## Selective Entrypoints

| Entrypoint | Description |
|------------|-------------|
| `git-commit/` | Git commit prompt only |
| `review-code/` | Code review prompt only |
| `review-search/` | Search review prompt only |
| `refactor-pure/` | Pure refactoring prompt only |

## Aggregate Entrypoint

| Entrypoint | Description |
|------------|-------------|
| `kustomization.yaml` | All prompts |

## Configuration

Prompts do not have configurable fields. Each prompt defines a fixed `data.name`,
`data.description`, and `data.prompt` that are passed through to the converter.

## Workflow Usage

Prompts are consumed by agent converters to produce agent-specific artifacts:

- `agents/opencode/converters/` — Converts prompts into OpenCode skills
- `agents/pi/converters/` — Converts prompts into Pi commands
