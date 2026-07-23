## ADDED Requirements

### Requirement: OpenSpec common capability entrypoint

The AI agents domain SHALL expose common OpenSpec resources separately from project-specific planning mode configuration.

#### Scenario: Common OpenSpec composition
- **WHEN** a consumer references `packages/domains/ai-agents/manifests/openspec/`
- **THEN** the composition SHALL provide the documented OpenSpec command, installation, converter, and shared resources
- **AND** it SHALL NOT select a local root or external store pointer without an explicit configuration entrypoint

### Requirement: Repository-local planning configuration

The OpenSpec capability SHALL provide a `local` configuration Kustomization for repositories that keep planning artifacts beside code.

#### Scenario: Local configuration materialization
- **WHEN** a workflow composes the OpenSpec `local` configuration and materializes `openspec.assets.config.local`
- **THEN** it SHALL create or merge `openspec/config.yaml`
- **AND** the configuration SHALL support JSON Patch customization of schema, project context, and artifact rules
- **AND** OpenSpec commands SHALL write changes and specs in the repository's `openspec/` root

### Requirement: External store pointer configuration

The OpenSpec capability SHALL provide a `store` configuration Kustomization for repositories whose planning is fully externalized.

#### Scenario: Store pointer materialization
- **WHEN** a workflow composes the OpenSpec `store` configuration and materializes `openspec.assets.config.store`
- **THEN** `openspec/config.yaml` SHALL declare `store: <id>` using a documented JSON Patch target
- **AND** the project configuration SHALL NOT create `.openspec-store/store.yaml`
- **AND** the README SHALL state that the selected store must be registered on each machine

#### Scenario: Local planning conflicts with pointer mode
- **WHEN** a repository contains a qualifying local OpenSpec planning root and also declares `store: <id>`
- **THEN** the capability documentation MUST explain that native root resolution selects the local root and ignores the pointer
- **AND** local and store-pointer configuration entrypoints MUST be documented as mutually exclusive

### Requirement: Read-only store references configuration

The OpenSpec capability SHALL provide a `references` Kustomization for adding external planning context without changing the command write target.

#### Scenario: Reference configuration materialization
- **WHEN** a workflow composes and materializes `openspec.assets.config.references`
- **THEN** `openspec/config.yaml` SHALL contain a patchable `references` list
- **AND** each reference SHALL support an ID and optional remote
- **AND** OpenSpec commands SHALL continue writing to the resolved project root unless the user explicitly passes `--store`

#### Scenario: Local planning with references
- **WHEN** a consumer composes `local` and `references`
- **THEN** their Assets SHALL deep-merge into one `openspec/config.yaml`
- **AND** the referenced stores SHALL be treated as read-only working context

### Requirement: Store state remains explicit and local

The AI agents domain MUST NOT automatically manage OpenSpec's machine-local store registry, global default store, worksets, or source-control synchronization.

#### Scenario: Store prerequisite
- **WHEN** a generated configuration uses a store pointer or reference
- **THEN** the README and wizard SHALL provide the applicable explicit setup or registration command
- **AND** no devShell hook or cached installation Step SHALL register, clone, remove, push, or synchronize the store

### Requirement: OpenSpec configuration is patchable and documented

Each OpenSpec configuration entrypoint MUST expose stable Assets and documented JSON Patch paths.

#### Scenario: OpenSpec README
- **WHEN** a consumer reads the OpenSpec capability README
- **THEN** it SHALL explain local roots, standalone stores, pointers, references, native root precedence, versioned files, machine-local state, beta limitations, composition examples, JSON Patch examples, and validation with `openspec doctor`, `openspec context`, and `openspec validate`

#### Scenario: Configuration fixtures
- **WHEN** domain tests build OpenSpec configuration fixtures
- **THEN** local, pointer, and local-with-references compositions SHALL succeed
- **AND** the generated `openspec/config.yaml` structures SHALL match the documented modes
