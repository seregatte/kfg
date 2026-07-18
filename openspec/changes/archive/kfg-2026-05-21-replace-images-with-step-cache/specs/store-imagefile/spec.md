## REMOVED Requirements

### Requirement: Imagefile-based image composition
**Reason**: The engine no longer supports image-based configuration composition and instead uses Step cache for persisted runtime artifacts.
**Migration**: Replace image/workspace-based flows with cacheable Steps and `kfg sys gc` operational management where persistence inspection or cleanup is required.
