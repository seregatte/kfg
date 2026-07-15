# Contributing to KFG

Thank you for your interest in contributing to KFG! This guide will help you understand how to contribute effectively.

## Code of Conduct

By participating in this project, you agree to maintain a respectful and inclusive environment. Be respectful to all contributors regardless of experience, gender, sexual orientation, disability, ethnicity, religion, or any other personal characteristic.

## How to Contribute

### Reporting Bugs

Before reporting a bug:

1. **Check existing issues**: https://github.com/seregatte/kfg/issues
2. **Consult the troubleshooting guide**: [docs/troubleshooting.md](docs/troubleshooting.md)

If the bug hasn't been reported, [open an issue](https://github.com/seregatte/kfg/issues/new) including:

- **Clear description** of the problem
- **Steps to reproduce** the issue
- **Expected behavior** vs **actual behavior**
- **Environment**: OS, KFG version, Go/Nix version
- **Logs**: Run with `KFG_VERBOSE=5` and include the logs
- **Manifests**: Include the minimal YAML that reproduces the issue (remove sensitive data)

### Suggesting Enhancements

Ideas and suggestions are welcome! Open an issue with the `enhancement` label describing:

- **Use case**: What problem does this enhancement solve?
- **Proposal**: How do you envision this feature?
- **Alternatives**: Other solutions considered?

### Contributing Code

#### Development Environment Setup

1. **Fork the repository** on GitHub

2. **Clone your fork**:
   ```bash
   git clone https://github.com/YOUR_USERNAME/kfg.git
   cd kfg
   ```

3. **Set up the development environment**:
   
   **Option A: Nix (Recommended)**
   ```bash
   nix develop .#dev
   ```
   
   **Option B: Native Go**
   ```bash
   # Requires Go 1.21+
   go mod download
   ```

4. **Add upstream**:
   ```bash
   git remote add upstream https://github.com/seregatte/kfg.git
   ```

#### Development Workflow

1. **Sync with upstream**:
   ```bash
   git checkout main
   git pull upstream main
   ```

2. **Create a branch**:
   ```bash
   git checkout -b feature/your-feature-name
   # or
   git checkout -b fix/bug-description
   ```
   
   **Naming conventions**:
   - `feature/*` - New features
   - `fix/*` - Bug fixes
   - `docs/*` - Documentation improvements
   - `refactor/*` - Refactoring
   - `test/*` - Test additions

3. **Make your changes**:
   
   Follow project standards:
   - Go code: follow standard Go style (`gofmt`)
   - Commits: follow [Conventional Commits](https://www.conventionalcommits.org/)
   - Tests: add tests for new features
   
   **Commit structure**:
   ```
   feat: add support for nested placeholders
   fix: fix cache invalidation in conditional steps
   docs: improve examples in getting-started
   refactor: simplify dependency resolver
   test: add tests for kustomize overlays
   ```

4. **Run tests**:
   ```bash
   # Unit tests
   make test
   
   # Integration tests
   make test-bats
   
   # Linting and formatting
   make fmt lint vet
   ```

5. **Commit your changes**:
   ```bash
   git add .
   git commit -m "feat: clear description of the change"
   ```

6. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```

7. **Open a Pull Request** at https://github.com/seregatte/kfg/pulls

#### Style Guides

**Go**:
- Follow [Effective Go](https://go.dev/doc/effective-go)
- Use `gofmt` for formatting
- Run `go vet` before committing
- Name variables and functions descriptively

**YAML/Manifests**:
- Use 2-space indentation
- Resource names: `<scope>.<kind>.<name>`
- Comments in English

**Documentation**:
- Markdown with maximum 100 characters per line
- Code examples should always be tested
- Use relative links within the repository

**Commits**:
- Use [Conventional Commits](https://www.conventionalcommits.org/)
- First line: maximum 72 characters
- Commit body: explain the "why", not the "what"
- Include "BREAKING CHANGE:" in the footer for incompatible changes

#### Tests

**Test types**:

1. **Unit tests** (`src/internal/*_test.go`):
   ```go
   func TestParser_ValidManifest(t *testing.T) {
       // Arrange
       manifest := `...`
       
       // Act
       result, err := Parse(manifest)
       
       // Assert
       assert.NoError(t, err)
       assert.Equal(t, expected, result)
   }
   ```

2. **Integration tests** (`tests/bats/`):
   ```bash
   @test "apply generates shell code" {
       run kfg apply -f test.yaml --workflow test
       [ "$status" -eq 0 ]
       [[ "$output" =~ "test_function()" ]]
   }
   ```

3. **E2E tests** (`packages/*/tests/`):
   Full tests with real manifests.

**Coverage**:
- New features: minimum 80% coverage
- Bug fixes: add a test that reproduces the bug

#### Documentation

When contributing code, update the documentation:

- **README.md**: If basic usage changed
- **docs/getting-started.md**: If the initial flow changed
- **docs/cli-reference.md**: If commands/flags changed
- **docs/manifest-model.md**: If schema changed
- **docs/architecture.md**: If internal architecture changed
- **Code comments**: For complex logic

### Contributing Documentation

Documentation is as important as code! You can:

- Fix typos and grammatical errors
- Improve confusing explanations
- Add examples
- Translate documentation
- Create tutorials

**Process**:
1. Follow the same code workflow
2. Use `docs/*` branch
3. Build locally to verify formatting:
   ```bash
   # If using MkDocs or similar
   make docs-serve
   ```

### Reviewing Pull Requests

Reviews are welcome! When reviewing:

- Be respectful and constructive
- Focus on the code, not the person
- Explain the "why" behind suggestions
- Test changes locally if possible
- Use appropriate labels (`needs-changes`, `approved`, etc.)

## Versioning Policy

KFG follows [Semantic Versioning](https://semver.org/):

- **MAJOR** (X.0.0): Incompatible changes
- **MINOR** (0.X.0): New features (backward compatible)
- **PATCH** (0.0.X): Bug fixes

**Important**: Version changes are only made in `release/*` branches by maintainers.

## Release Process

Releases are made by maintainers:

1. `release/vX.Y.Z` branch is created
2. Version bump in `flake.nix`
3. Tag created: `git tag vX.Y.Z`
4. CI builds and publishes
5. PR to `main`

Contributors should not make version bumps in feature branches.

## Community

- **GitHub Discussions**: https://github.com/seregatte/kfg/discussions
- **Issues**: https://github.com/seregatte/kfg/issues
- **PRs**: https://github.com/seregatte/kfg/pulls

## Recognition

Contributors are recognized in the README and release notes. Thank you for helping make KFG better!

## Questions?

If you have questions about contributing:

1. Consult this guide
2. Check existing issues
3. Open an issue with the `question` label
4. Participate in discussions

---

**Workflow summary**:
```
1. Fork → 2. Branch → 3. Code → 4. Test → 5. Commit → 6. Push → 7. PR
```

Thank you for contributing! 🎉