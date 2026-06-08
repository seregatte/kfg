# KFG - Declarative Shell Compiler

[![GitHub release (latest by date)](https://img.shields.io/github/v/release/seregatte/kfg)](https://github.com/seregatte/kfg/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Nix](https://img.shields.io/badge/Nix-5277C3?logo=nixos&logoColor=white)](https://nixos.org)

**KFG** (Key Function Generator) é um compilador shell declarativo que transforma manifests YAML em funções bash. Defina comandos, dependências e passos de execução em YAML, e o KFG gera código shell que pode ser "sourceado" ou usado interativamente.

## Por que usar KFG?

- **Declarativo**: Defina o que fazer em YAML, não como fazer em bash
- **Dependências**: O KFG gerencia a ordem de execução automaticamente
- **Cache**: Steps são cacheadas para evitar re-execuções desnecessárias
- **Reutilizável**: Crie manifests modulares e reutilizáveis
- **Versionável**: Seus workflows de shell agora podem ser versionados como código

## Casos de Uso

- **Deploy de aplicações**: Defina pipelines de deploy declarativos
- **Setup de ambientes**: Automatize configuração de projetos com dependências
- **Workflows de AI agents**: Gere comandos para Claude, Copilot, etc.
- **CI/CD**: Padronize processos de build e deploy
- **Gerenciamento de MCP servers**: Configure e gerencie MCP servers declarativamente

## Instalação

### Pré-requisitos

Para instalar via Nix (recomendado):
- [Nix](https://nixos.org/download.html) com flakes habilitados

Para buildar do fonte:
- Go 1.21+
- Make

### Via Nix (Recomendado)

```bash
# Build e instalar
nix build github:seregatte/kfg

# Executar sem instalar
nix run github:seregatte/kfg -- --help

# Adicionar ao shell atual
nix shell github:seregatte/kfg
```

Suporta Linux e macOS (x86_64 e ARM64).

### Build do Fonte

```bash
git clone https://github.com/seregatte/kfg.git
cd kfg
make build
```

O binário será colocado em `./bin/kfg`.

### Instalar no GOPATH

```bash
make install
```

## Quick Start

### 1. Crie seu primeiro manifest

Crie um arquivo `hello.yaml`:

```yaml
apiVersion: kfg.dev/v1alpha1
kind: Cmd
metadata:
  name: myapp.cmd.hello
  commandName: hello
spec:
  run: echo "Hello from KFG!"
```

### 2. Aplique o manifest

```bash
kfg apply -f hello.yaml --workflow default
```

Isso gera e executa o código shell. Agora você tem um comando `hello` disponível!

### 3. Execute o comando

```bash
hello
# Output: Hello from KFG!
```

### Exemplo Completo: Pipeline de Deploy

Crie `deploy.yaml`:

```yaml
apiVersion: kfg.dev/v1alpha1
kind: Cmd
metadata:
  name: myapp.cmd.deploy
  commandName: deploy
spec:
  env:
    DEPLOY_TARGET: "{env:DEPLOY_TARGET:-production}"
  run: |
    echo "Deploying to $DEPLOY_TARGET..."
    kubectl apply -f manifests/

---
apiVersion: kfg.dev/v1alpha1
kind: Step
metadata:
  name: myapp.steps.validate
spec:
  run: |
    [ -f "config.yaml" ] && echo "Config found" || exit 1
  output:
    name: STATUS
    type: string

---
apiVersion: kfg.dev/v1alpha1
kind: CmdWorkflow
metadata:
  name: myapp.workflow.deploy
spec:
  cmds: [myapp.cmd.deploy]
  before:
    - step: myapp.steps.validate
```

Aplique com:

```bash
DEPLOY_TARGET=staging kfg apply -f deploy.yaml --workflow deploy
```

📖 **Mais exemplos**: Veja [docs/getting-started.md](docs/getting-started.md) para um tutorial completo.

## Documentação

- **[Getting Started](docs/getting-started.md)** - Tutorial passo a passo
- **[CLI Reference](docs/cli-reference.md)** - Referência completa da CLI
- **[Manifest Model](docs/manifest-model.md)** - Schema e tipos de manifests
- **[Architecture](docs/architecture.md)** - Arquitetura interna do KFG
- **[Troubleshooting](docs/troubleshooting.md)** - Problemas comuns e soluções
- **[Contributing](CONTRIBUTING.md)** - Como contribuir

## Command Reference

| Comando | Descrição | Exemplo |
|---------|-----------|---------|
| `kfg apply` | Aplica kustomization/manifest e gera shell code | `kfg apply -f manifest.yaml --workflow main` |
| `kfg run` | Executa um agent one-shot | `kfg run -k ./manifests myagent` |
| `kfg build` | Build kustomization para YAML | `kfg build ./manifests -o output.yaml` |
| `kfg sys cache` | Gerenciamento de cache de steps | `kfg sys cache ls` |
| `kfg sys log` | Logging estruturado para scripts | `kfg sys log info "component" "message"` |
| `kfg version` | Mostra informações de versão | `kfg version` |

📖 **Referência completa**: Veja [docs/cli-reference.md](docs/cli-reference.md) para todos os comandos, flags e variáveis de ambiente.

## Comparação com Alternativas

| Ferramenta | KFG | Make | Just | Task |
|------------|-----|------|------|------|
| **Sintaxe** | YAML declarativo | Makefile | Justfile | YAML |
| **Dependências** | Automáticas via DAG | Manual | Manual | Manual |
| **Cache de steps** | ✅ Nativo | ❌ | ❌ | ❌ |
| **Composição modular** | ✅ Kustomize | ❌ | ❌ | Limitado |
| **Geração de código** | ✅ Shell functions | ❌ | ❌ | ❌ |
| **Placeholders** | ✅ `{env:VAR}` | ❌ | ❌ | Limitado |
| **Versionável** | ✅ Sim | ✅ Sim | ✅ Sim | ✅ Sim |

**Quando usar KFG?**
- Quando você precisa de **dependências automáticas** entre tarefas
- Quando você quer **cache inteligente** para evitar re-execuções
- Quando você precisa **compor manifests** de diferentes fontes (Kustomize)
- Quando você quer **gerar funções shell** reutilizáveis
- Quando você está trabalhando com **AI agents** que precisam de comandos estruturados

## Environment Variables

| Variável | Descrição | Exemplo |
|----------|-----------|---------|
| `KFG_VERBOSE` | Nível de verbosidade (0-5) | `KFG_VERBOSE=3` |
| `KFG_STORE_DIR` | Diretório do store (cache) | `~/.kfg/store` |
| `KFG_LOG_FILE` | Caminho do arquivo de log | `/tmp/kfg.log` |
| `KFG_LOG_DIR` | Diretório de logs | `~/.local/state/kfg/logs` |
| `KFG_LOG_COLOR` | Modo de cor (auto/always/never) | `auto` |
| `KFG_KPATH` | Caminho padrão para kustomization | `./manifests` |
| `KFG_REFRESH` | Invalida cache (set to "1") | `1` |

📖 **Referência completa**: Veja [docs/cli-reference.md](docs/cli-reference.md#environment-variables).

## API Version

KFG usa `kfg.dev/v1alpha1` como versão de API para manifests:

```yaml
apiVersion: kfg.dev/v1alpha1
kind: Cmd
metadata:
  name: example
spec:
  run: echo "Hello, World!"
```

📖 **Schema completo**: Veja [docs/manifest-model.md](docs/manifest-model.md) para todos os tipos de recursos.

## Exemplos Reais

O KFG é usado no próprio repositório para gerenciar workflows de AI agents:

```bash
# Aplicar overlay de desenvolvimento
kfg apply -k packages/domains/ai-agents/overlays/dev --workflow agents

# Executar agent específico
kfg run -k packages/domains/ai-agents/overlays/dev openspec
```

Veja mais exemplos em `packages/domains/ai-agents/manifests/`.

## Development

### DevShells

KFG fornece três devShells via Nix flakes:

| Shell | Uso | Descrição |
|-------|-----|-----------|
| `default` | `nix develop` | **Consumer shell** — ferramentas para usar KFG |
| `dev` | `nix develop .#dev` | **Development shell** — ambiente completo para desenvolvimento |
| `ci` | `nix develop .#ci` | **CI shell** — mínimo para build e testes |

### Building

```bash
# Usando o dev shell
nix develop .#dev --command make build        # → ./bin/kfg
nix develop .#dev --command make test         # Testes unitários Go
nix develop .#dev --command make test-bats    # Testes de integração Bats
```

### Repository Structure

```
├── src/                          # Implementação do engine (Go)
│   ├── cmd/kfg/                  # Comandos CLI
│   └── internal/                 # Pacotes internos
├── packages/
│   ├── framework/                # Primitivas de manifest compartilhadas
│   │   ├── manifests/            # Steps reutilizáveis
│   │   └── tests/                # Testes do framework
│   └── domains/
│       └── ai-agents/            # Pacote de domínio AI agents
│           ├── manifests/        # Recursos de AI agents
│           ├── overlays/dev/     # Overlay de desenvolvimento
│           └── tests/            # Testes do domínio
├── docs/
│   ├── AGENTS.md                 # Contexto para AI agents
│   └── context/
│       └── openspec/             # Especificações OpenSpec
├── tests/
│   └── bats/                     # Testes de engine e integração
└── Makefile                      # Targets de build e teste
```

📖 **Arquitetura detalhada**: Veja [docs/architecture.md](docs/architecture.md).

## License

MIT License - veja [LICENSE](LICENSE) para detalhes.

## Links

- **Repositório**: https://github.com/seregatte/kfg
- **Releases**: https://github.com/seregatte/kfg/releases
- **Issues**: https://github.com/seregatte/kfg/issues
- **Discussions**: https://github.com/seregatte/kfg/discussions
