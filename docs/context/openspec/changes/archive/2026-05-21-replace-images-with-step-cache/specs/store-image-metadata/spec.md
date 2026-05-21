## REMOVED Requirements

### Requirement: Stored image metadata inspection
**Reason**: Image metadata and inspection APIs are removed with the image feature set.
**Migration**: Use `kfg sys gc inspect <id>` to inspect persisted runtime cache metadata for cacheable Steps.
