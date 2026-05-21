## REMOVED Requirements

### Requirement: Image build command behavior
**Reason**: The `kfg image` command family is removed along with the image storage model.
**Migration**: Move startup optimization use cases to cacheable Steps and manage persisted cache entries through `kfg sys gc`.
