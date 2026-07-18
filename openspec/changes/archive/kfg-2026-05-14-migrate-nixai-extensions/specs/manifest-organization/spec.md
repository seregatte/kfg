## ADDED Requirements

### Requirement: Hierarchical kustomization files

Each directory under `manifests/base/` SHALL contain a `kustomization.yaml` file that references its immediate children (files and subdirectories).

#### Scenario: Root kustomization references top-level directories
- **WHEN** kustomize loads `manifests/base/kustomization.yaml`
- **THEN** it references exactly 4 directories: `agents`, `cmds`, `extensions`, `steps`

#### Scenario: agents kustomization references its resources
- **WHEN** kustomize loads `manifests/base/agents/kustomization.yaml`
- **THEN** it references `claude.yaml`, `gemini.yaml`, `opencode.yaml`, `pi.yaml`, `converters/`, and `steps/`

#### Scenario: extensions kustomization references all extension directories
- **WHEN** kustomize loads `manifests/base/extensions/kustomization.yaml`
- **THEN** it references `self`, `ctx7`, `chrome-devtools`, `playwright`, `gws`, `notebooklm`, and `openspec`

#### Scenario: self extension kustomization references assets and converters
- **WHEN** kustomize loads `manifests/base/extensions/self/kustomization.yaml`
- **THEN** it references `assets` and `converters` subdirectories

#### Scenario: Extension with no resources
- **WHEN** kustomize loads `manifests/base/extensions/gws/kustomization.yaml`
- **THEN** it has an empty `resources` list (no error)

### Requirement: Backward compatibility

The reorganized kustomization structure MUST produce the same merged output as the current flat structure.

#### Scenario: Same manifest output after reorganization
- **WHEN** `kfg build` is run with the new hierarchical kustomization
- **THEN** the generated shell output is identical to the output from the current flat kustomization

### Requirement: Extension install Steps registered in kustomization

Each extension's `kustomization.yaml` SHALL reference its `steps/` directory after the install Steps are created.

#### Scenario: ctx7 kustomization includes steps
- **WHEN** kustomize loads `manifests/base/extensions/ctx7/kustomization.yaml`
- **THEN** it references both `assets` and `steps` subdirectories

#### Scenario: gws kustomization includes steps
- **WHEN** kustomize loads `manifests/base/extensions/gws/kustomization.yaml`
- **THEN** it references the `steps` subdirectory
