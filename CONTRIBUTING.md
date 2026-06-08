# Contributing to KFG

Obrigado pelo interesse em contribuir com o KFG! Este guia vai te ajudar a entender como contribuir de forma eficaz.

## Código de Conduta

Ao participar deste projeto, você concorda em manter um ambiente respeitoso e inclusivo. Seja respeitoso com todos os contribuidores, independente de experiência, gênero, orientação sexual, deficiência, etnia, religião ou qualquer outra característica pessoal.

## Como Contribuir

### Reportando Bugs

Antes de reportar um bug:

1. **Verifique issues existentes**: https://github.com/seregatte/kfg/issues
2. **Consulte o troubleshooting**: [docs/troubleshooting.md](docs/troubleshooting.md)

Se o bug não foi reportado, [abra uma issue](https://github.com/seregatte/kfg/issues/new) incluindo:

- **Descrição clara** do problema
- **Passos para reproduzir** o problema
- **Comportamento esperado** vs **comportamento atual**
- **Ambiente**: SO, versão do KFG, versão do Go/Nix
- **Logs**: Execute com `KFG_VERBOSE=5` e inclua os logs
- **Manifests**: Inclua o YAML mínimo que reproduz o problema (remova dados sensíveis)

### Sugerindo Melhorias

Ideias e sugestões são bem-vindas! Abra uma issue com a label `enhancement` descrevendo:

- **Caso de uso**: Qual problema essa melhoria resolve?
- **Proposta**: Como você imagina essa funcionalidade?
- **Alternativas**: Outras soluções consideradas?

### Contribuindo com Código

#### Setup do Ambiente de Desenvolvimento

1. **Fork o repositório** no GitHub

2. **Clone seu fork**:
   ```bash
   git clone https://github.com/YOUR_USERNAME/kfg.git
   cd kfg
   ```

3. **Configure o ambiente de desenvolvimento**:
   
   **Opção A: Nix (Recomendado)**
   ```bash
   nix develop .#dev
   ```
   
   **Opção B: Go nativo**
   ```bash
   # Requer Go 1.21+
   go mod download
   ```

4. **Adicione o upstream**:
   ```bash
   git remote add upstream https://github.com/seregatte/kfg.git
   ```

#### Workflow de Desenvolvimento

1. **Sincronize com upstream**:
   ```bash
   git checkout main
   git pull upstream main
   ```

2. **Crie uma branch**:
   ```bash
   git checkout -b feature/nome-da-sua-feature
   # ou
   git checkout -b fix/descricao-do-bug
   ```
   
   **Convenções de nomenclatura**:
   - `feature/*` - Novas funcionalidades
   - `fix/*` - Correções de bugs
   - `docs/*` - Melhorias na documentação
   - `refactor/*` - Refatorações
   - `test/*` - Adição de testes

3. **Faça suas alterações**:
   
   Siga os padrões do projeto:
   - Código em Go: siga o estilo padrão do Go (`gofmt`)
   - Commits: siga [Conventional Commits](https://www.conventionalcommits.org/)
   - Testes: adicione testes para novas funcionalidades
   
   **Estrutura de commits**:
   ```
   feat: adicionar suporte a placeholders aninhados
   fix: corrigir cache invalidation em steps condicionais
   docs: melhorar exemplos no getting-started
   refactor: simplificar dependency resolver
   test: adicionar testes para kustomize overlays
   ```

4. **Execute os testes**:
   ```bash
   # Testes unitários
   make test
   
   # Testes de integração
   make test-bats
   
   # Linting e formatação
   make fmt lint vet
   ```

5. **Commit suas alterações**:
   ```bash
   git add .
   git commit -m "feat: descrição clara da mudança"
   ```

6. **Push para seu fork**:
   ```bash
   git push origin feature/nome-da-sua-feature
   ```

7. **Abra um Pull Request** em https://github.com/seregatte/kfg/pulls

#### Guias de Estilo

**Go**:
- Siga [Effective Go](https://go.dev/doc/effective-go)
- Use `gofmt` para formatação
- Execute `go vet` antes de commitar
- Nomeie variáveis e funções de forma descritiva

**YAML/Manifests**:
- Use indentação de 2 espaços
- Nomes de recursos: `<scope>.<kind>.<name>`
- Comentários em inglês

**Documentação**:
- Markdown com linhas de no máximo 100 caracteres
- Exemplos de código sempre testados
- Links relativos dentro do repositório

**Commits**:
- Use [Conventional Commits](https://www.conventionalcommits.org/)
- Primeira linha: máximo 72 caracteres
- Corpo do commit: explique o "porquê", não o "o quê"
- Inclua "BREAKING CHANGE:" no footer para mudanças incompatíveis

#### Testes

**Tipos de testes**:

1. **Unitários** (`src/internal/*_test.go`):
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

2. **Integração** (`tests/bats/`):
   ```bash
   @test "apply generates shell code" {
       run kfg apply -f test.yaml --workflow test
       [ "$status" -eq 0 ]
       [[ "$output" =~ "test_function()" ]]
   }
   ```

3. **E2E** (`packages/*/tests/`):
   Testes completos com manifests reais.

**Cobertura**:
- Novas funcionalidades: mínimo 80% de cobertura
- Bug fixes: adicione teste que reproduz o bug

#### Documentação

Ao contribuir com código, atualize a documentação:

- **README.md**: Se mudou uso básico
- **docs/getting-started.md**: Se mudou fluxo inicial
- **docs/cli-reference.md**: Se mudou comandos/flags
- **docs/manifest-model.md**: Se mudou schema
- **docs/architecture.md**: Se mudou arquitetura interna
- **Comentários no código**: Para lógica complexa

### Contribuindo com Documentação

Documentação é tão importante quanto código! Você pode:

- Corrigir typos e erros gramaticais
- Melhorar explicações confusas
- Adicionar exemplos
- Traduzir documentação
- Criar tutoriais

**Processo**:
1. Siga o mesmo workflow de código
2. Use branch `docs/*`
3. Build local para verificar formatação:
   ```bash
   # Se usar MkDocs ou similar
   make docs-serve
   ```

### Revisando Pull Requests

Revisões são bem-vindas! Ao revisar:

- Seja respeitoso e construtivo
- Foque no código, não na pessoa
- Explique o "porquê" das sugestões
- Teste as mudanças localmente se possível
- Use labels apropriados (`needs-changes`, `approved`, etc.)

## Política de Versão

KFG segue [Semantic Versioning](https://semver.org/):

- **MAJOR** (X.0.0): Mudanças incompatíveis
- **MINOR** (0.X.0): Novas funcionalidades (backward compatible)
- **PATCH** (0.0.X): Correções de bugs

**Importante**: Mudanças de versão são feitas apenas em branches `release/*` pelo maintainers.

## Processo de Release

Releases são feitas pelos maintainers:

1. Branch `release/vX.Y.Z` é criada
2. Version bump no `flake.nix`
3. Tag criada: `git tag vX.Y.Z`
4. CI builda e publica
5. PR para `main`

Contribuidores não devem fazer version bumps em branches de feature.

## Comunidade

- **GitHub Discussions**: https://github.com/seregatte/kfg/discussions
- **Issues**: https://github.com/seregatte/kfg/issues
- **PRs**: https://github.com/seregatte/kfg/pulls

## Reconhecimento

Contribuidores são reconhecidos no README e nas release notes. Obrigado por ajudar a tornar o KFG melhor!

## Dúvidas?

Se tiver dúvidas sobre como contribuir:

1. Consulte este guia
2. Verifique issues existentes
3. Abra uma issue com a label `question`
4. Participe das discussions

---

**Resumo do workflow**:
```
1. Fork → 2. Branch → 3. Code → 4. Test → 5. Commit → 6. Push → 7. PR
```

Obrigado por contribuir! 🎉
