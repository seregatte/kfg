## REMOVED Requirements

### Requirement: Workspace materialization and restore
**Reason**: Workspace start/stop behavior is removed with the image/workspace system.
**Migration**: Replace workspace-based agent bootstrapping with cacheable Steps that restore generated artifacts directly into the working tree.
