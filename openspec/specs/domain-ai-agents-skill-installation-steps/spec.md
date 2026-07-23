## MODIFIED Requirements

### Requirement: Generic skill installation Step

The system SHALL provide a Step resource for each extension that installs agent skills via external CLIs. Each Step MUST receive agent-specific configuration through `spec.env`, validate required inputs, create required output directories, preserve unrelated project files, and document generated artifacts and cleanup behavior. Steps MUST NOT select an agent through internal `case` or equivalent branching when workflow conditions can select the agent-specific invocation.

#### Scenario: ctx7 skill installation for OpenCode
- **WHEN** a workflow invokes `ctx7.steps.install` with `FLAGS="--opencode --yes"` and `OUTPUT_DIR=".opencode/skills/"`
- **THEN** the Step SHALL execute the documented Context7 project setup
- **AND** it SHALL copy generated skills into `.opencode/skills/`

#### Scenario: Chrome DevTools skill installation
- **WHEN** a workflow invokes `chrome-devtools.steps.install` with its documented skill, agent, home, and output inputs
- **THEN** the Step SHALL install the Chrome DevTools skill for the selected agent
- **AND** it SHALL copy generated files into the requested output directory

#### Scenario: Playwright skill installation
- **WHEN** a workflow invokes `playwright.steps.install` with its documented output inputs
- **THEN** the Step SHALL install the Playwright skill for the selected agent
- **AND** cleanup SHALL NOT remove pre-existing agent configuration outside the Step's registered artifacts

#### Scenario: Google Workspace skill installation
- **WHEN** a workflow invokes `gws.steps.install` with its documented skill and agent inputs
- **THEN** the Step SHALL install the Google Workspace skill for the selected agent
- **AND** generated files SHALL be registered as artifacts

#### Scenario: NotebookLM skill installation
- **WHEN** a workflow invokes `notebooklm.steps.install` with its documented output inputs
- **THEN** the Step SHALL install NotebookLM skill files into the requested location
- **AND** it SHALL NOT assume an undocumented agent home

#### Scenario: OpenSpec agent integration installation
- **WHEN** a workflow invokes `openspec.steps.install` with a documented OpenSpec tool selection and output locations
- **THEN** the Step SHALL generate and copy the selected agent's OpenSpec integration files
- **AND** it MUST preserve an existing `openspec/` root
- **AND** it MUST NOT create and remove a disposable root that can shadow canonical project planning data

#### Scenario: Install Step logging attribution
- **WHEN** an installation Step emits a runtime log event through `__kfg_log_*`
- **THEN** the event SHALL rely on runtime-provided `step_name` attribution
- **AND** the Step SHALL NOT need to encode its Step identity inside the component string
