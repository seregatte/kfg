## Implementation Notes

- Each cache entry is identified by a stable identifier derived from Step reference name
- The cache identity uses `StepReference.name` only (no additional components)
- The `metadata.yaml` file contains:
  - `stepRefName`: The workflow step reference name
  - `timestamp`: When the cache entry was created
  - `output`: Optional output metadata (name, valueEncoded)
  - `artifacts`: List of cached artifact relative paths
- The prune policy currently removes entries older than 30 days