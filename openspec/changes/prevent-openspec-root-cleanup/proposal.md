## Why

Cached installation Steps can classify a newly created directory as an artifact and later register that directory for recursive session cleanup. When this happens to the repository `openspec/` root, `kfg.cleanup` deletes persistent specs and changes created after the cache entry was restored, causing irreversible planning data loss.

## What Changes

- Make automatic filesystem-diff artifact discovery cache concrete files and symbolic links without treating their parent directories as recursively disposable artifacts.
- Preserve explicitly declared directory artifacts so workflows can continue to mark generated scaffolds for whole-directory cleanup.
- Invalidate cache entries created under the unsafe artifact-discovery semantics so existing entries cannot restore and register `openspec/` as a directory artifact.
- Add cache miss and cache hit coverage proving that OpenSpec installation preserves an existing root and changes created during the agent session.
- Document that OpenSpec integration cleanup removes generated integration artifacts without removing the repository planning root.
- Non-goal: change the public manifest model, remove directory artifact support, or make kfg manage OpenSpec stores, registries, worksets, or synchronization.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `kfg-cache-step`: Automatic filesystem discovery records leaf artifacts rather than parent directories, and incompatible cache entries are not reused.
- `kfg-cache-sys-cache-command`: Cache storage and identity semantics prevent legacy directory artifacts from being restored into cleanup tracking.
- `domain-ai-agents-skill-installation-steps`: OpenSpec installation preserves the repository `openspec/` root across cache misses, cache hits, and final workflow cleanup.

## Impact

- Engine cache implementation under `src/internal/cache/` and generated cache snapshot shell code under `src/internal/generate/templates/` will change.
- Cache entries produced by the previous unsafe discovery semantics will be invalidated and rebuilt on their next use; this is an intentional one-time cache compatibility break with no persisted project-data migration.
- `kfg.cleanup` keeps its public artifact-driven contract and continues to remove explicitly registered directory artifacts recursively.
- The OpenSpec installation Step, AI-agent capability documentation, and domain regression tests will be updated; no OpenSpec CLI arguments, store configuration, public resource names, or manifest fields change.
- User-facing behavior changes only by preserving persistent files under `openspec/`; generated agent integration files remain eligible for cleanup.
