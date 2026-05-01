# Template Engine Specification

## Purpose

kfg uses Go's text/template package for generating shell code from YAML manifests. This specification defines how templates are loaded, executed, and how they produce deterministic, valid shell output.

## Requirements

### Requirement: Template-based Shell Generation

The CLI MUST use text/template for generating shell code.

#### Scenario: Template embedding with Go embed package
- **WHEN** the CLI is built
- **THEN** template files are embedded in the binary using Go 1.16+ embed package
- **AND** //go:embed directive is used to embed template files
- **AND** templates are included in the binary for single-file distribution

#### Scenario: Template loading with ParseFS
- **WHEN** the CLI initializes the generator
- **THEN** template.ParseFS() is used to parse embedded templates from embed.FS
- **AND** templates are parsed and validated at startup
- **AND** template parsing errors are reported immediately
- **AND** ParseFS() is more efficient than ParseFiles() for embedded content

#### Scenario: Template execution
- **WHEN** shell code is generated for a CmdWorkflow
- **THEN** templates are executed with the resolved CmdWorkflow data
- **AND** template output is written to a buffer
- **AND** the buffer is returned as the final shell code

#### Scenario: Template organization
- **WHEN** templates are organized
- **THEN** there is a base template for the shell header
- **AND** there are templates for each shell function type
- **AND** templates can include other templates using {{ template "name" . }}
- **AND** {{ define "name" }} blocks create reusable template components

#### Scenario: Parse once at initialization
- **WHEN** templates are loaded
- **THEN** templates are parsed ONCE during initialization
- **AND** parsed templates are reused for every generation request
- **AND** templates are NOT parsed on every request (performance optimization)
- **AND** template parsing overhead is minimized

### Requirement: Template Data Structure

Template data MUST be well-structured and documented.

#### Scenario: Data structure definition
- **WHEN** templates are designed
- **THEN** a clear data structure is defined for template execution
- **AND** the data structure includes all necessary fields (CmdWorkflow name, commands, steps, etc.)
- **AND** the data structure is documented

#### Scenario: Data structure population
- **WHEN** a resolved CmdWorkflow is converted to template data
- **THEN** all required fields are populated
- **AND** the data structure matches the template's expectations
- **AND** missing fields result in clear errors

### Requirement: Template Functions

Templates MUST have access to custom functions for shell-specific operations.

#### Scenario: Shell escaping function
- **WHEN** a template needs to escape a string for shell
- **THEN** a `shellEscape` function is available
- **AND** the function properly escapes special characters
- **AND** the escaped string is safe for shell execution

#### Scenario: String manipulation functions
- **WHEN** a template needs to manipulate strings
- **THEN** standard Go template functions are available (upper, lower, trim, etc.)
- **AND** custom functions can be added as needed

#### Scenario: Conditional logic
- **WHEN** a template needs conditional output
- **THEN** standard template conditionals are available (if, else, with)
- **AND** conditions are evaluated correctly
- **AND** template readability is maintained

### Requirement: FuncMap Registration

Custom template functions MUST be registered before template parsing.

#### Scenario: FuncMap definition
- **WHEN** custom functions are defined for templates
- **THEN** they are created as template.FuncMap (map[string]interface{})
- **AND** each function is named with a string key
- **AND** function names are used in templates: {{ functionName .Data }}

#### Scenario: FuncMap registration before Parse
- **GIVEN** a FuncMap is defined
- **WHEN** templates are parsed
- **THEN** the FuncMap is registered BEFORE calling Parse()
- **AND** template.New("name").Funcs(funcMap).Parse(content) is used
- **AND** functions registered after Parse() will not be available
- **AND** registration order is critical

#### Scenario: Function return values
- **WHEN** a custom function is defined
- **THEN** the function MUST return either one value or two values
- **AND** single return value: the result
- **AND** two return values: (result, error) - error halts template execution
- **AND** functions returning more than two values cause runtime error

#### Scenario: Shell escaping function signature
- **GIVEN** the shellEscape function is defined
- **WHEN** used in templates
- **THEN** it returns (string, error) for proper error handling
- **AND** invalid input returns an error
- **AND** valid input returns escaped string with nil error
- **AND** template execution halts on error

### Requirement: Template Context Navigation

Templates MUST use proper context navigation with dot (.) and context-changing actions.

#### Scenario: Dot (.) represents current context
- **WHEN** a template is executed with data
- **THEN** the dot (.) represents the current data context
- **AND** {{ .FieldName }} accesses fields of the current data
- **AND** the dot starts at the root data passed to Execute()

#### Scenario: {{ with .Field }} changes context
- **GIVEN** a nested data structure with .CmdWorkflow.Cmds
- **WHEN** {{ with .CmdWorkflow }} is used
- **THEN** the dot (.) context changes to .CmdWorkflow
- **AND** inside the with block, {{ .Cmds }} accesses Cmds field
- **AND** {{ end }} restores the original context
- **AND** context changes simplify nested data access

#### Scenario: {{ range .Items }} iterates and changes context
- **GIVEN** a data structure with .Commands as a slice
- **WHEN** {{ range .Commands }} is used
- **THEN** the dot (.) context changes to each item in the slice
- **AND** inside the range block, {{ .Name }} accesses each command's Name
- **AND** {{ end }} restores the original context after iteration
- **AND** iteration provides clean template code for lists

#### Scenario: Context preservation for nested operations
- **GIVEN** complex nested data structures
- **WHEN** multiple context changes are needed
- **THEN** $ variable can save the root context: {{ $root := . }}
- **AND** $root can be referenced inside with/range blocks
- **AND** this preserves access to parent data in nested contexts
- **AND** template readability is maintained

### Requirement: Template Readability

Templates MUST be readable and maintainable.

#### Scenario: Template formatting
- **WHEN** templates are written
- **THEN** they use consistent indentation
- **AND** comments are used to explain complex logic
- **AND** templates are organized logically

#### Scenario: Template testing
- **WHEN** templates are modified
- **THEN** they can be tested independently
- **AND** unit tests exist for template functions
- **AND** integration tests verify template output

### Requirement: Template Output Quality

Generated shell code MUST be valid and deterministic.

#### Scenario: Valid shell output
- **WHEN** templates are executed
- **THEN** the output is valid bash code
- **AND** the code can be sourced or eval'd without errors
- **AND** syntax errors are detected before output

#### Scenario: Deterministic output
- **GIVEN** the same input data
- **WHEN** templates are executed multiple times
- **THEN** the output is identical each time
- **AND** no timestamps or random values are included
- **AND** output order is consistent

#### Scenario: Output formatting
- **WHEN** shell code is generated
- **THEN** the output has consistent formatting
- **AND** indentation is correct
- **AND** line breaks are appropriate for readability

### Requirement: Template Error Handling

Template errors MUST be clear and actionable.

#### Scenario: Template parsing error
- **GIVEN** a template has syntax errors
- **WHEN** templates are loaded
- **THEN** a clear error message is returned
- **AND** the error indicates the template name and line number
- **AND** the error describes the syntax issue

#### Scenario: Template execution error
- **GIVEN** a template accesses a missing field
- **WHEN** the template is executed
- **THEN** a clear error message is returned
- **AND** the error indicates which field is missing
- **AND** the error includes context about the operation being performed

#### Scenario: Missing template
- **GIVEN** a required template file is missing
- **WHEN** the CLI starts
- **THEN** an error is returned
- **AND** the error indicates which template is missing
- **AND** the CLI fails fast with a clear message

### Requirement: Template Extensibility

Templates MUST support easy extension for future shells.

#### Scenario: Adding new shell template
- **WHEN** a developer adds support for a new shell (e.g., zsh, fish)
- **THEN** new template files are added for the shell
- **AND** the generator detects and uses the new templates
- **AND** no changes to core generation logic are needed

#### Scenario: Template overrides
- **GIVEN** a user wants to customize shell output
- **WHEN** custom templates are provided (future feature)
- **THEN** custom templates override default templates
- **AND** the override mechanism is documented
- **AND** fallback to default templates is supported

### Requirement: Migration from String Concatenation

Existing string concatenation code MUST be replaced with templates.

#### Scenario: Function generation migration
- **GIVEN** code uses fmt.Sprintf or string concatenation to build shell functions
- **WHEN** migrated to templates
- **THEN** a template is created for the function structure
- **AND** the template produces identical output
- **AND** the code is more readable and maintainable

#### Scenario: Header generation migration
- **GIVEN** code uses string concatenation to build the shell header
- **WHEN** migrated to templates
- **THEN** a header template is created
- **AND** the template includes all header components
- **AND** header content is centralized in one place

#### Scenario: Helper function generation
- **GIVEN** code generates helper functions (e.g., __kfg_* functions)
- **WHEN** migrated to templates
- **THEN** templates are created for each helper function type
- **AND** helper function logic is clear and testable
- **AND** helper functions can be easily modified