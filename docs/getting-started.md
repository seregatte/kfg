# Getting Started with KFG

Este guia vai te levar do zero ao primeiro workflow funcional em 10 minutos.

## Pré-requisitos

Antes de começar, certifique-se de ter:

- **Nix** (recomendado) ou **Go 1.21+** para build do fonte
- Um editor de texto
- Terminal com bash ou zsh

## Instalação

### Opção 1: Via Nix (Recomendado)

Se você já tem o Nix instalado com flakes habilitados:

```bash
# Adicionar ao shell atual
nix shell github:seregatte/kfg

# Verificar instalação
kfg version
```

### Opção 2: Build do Fonte

```bash
git clone https://github.com/seregatte/kfg.git
cd kfg
make build

# Adicionar ao PATH
export PATH="$PWD/bin:$PATH"

# Verificar instalação
kfg version
```

## Primeiro Manifest

Vamos criar um comando simples que imprime "Hello, World!".

### 1. Crie o arquivo `hello.yaml`

```yaml
apiVersion: kfg.dev/v1alpha1
kind: Cmd
metadata:
  name: tutorial.cmd.hello
  commandName: hello
spec:
  run: echo "Hello, World!"
```

**Explicação:**
- `apiVersion`: Versão da API do KFG
- `kind`: Tipo de recurso (Cmd = função shell)
- `metadata.name`: Nome único do recurso
- `metadata.commandName`: Nome do comando que será gerado
- `spec.run`: Código shell a ser executado

### 2. Aplique o manifest

```bash
kfg apply -f hello.yaml --workflow default
```

O que acontece:
1. KFG lê o manifest YAML
2. Gera código shell
3. Disponibiliza o comando `hello` no seu shell

### 3. Execute o comando

```bash
hello
# Output: Hello, World!
```

## Adicionando Variáveis de Ambiente

Vamos modificar o comando para usar variáveis de ambiente:

```yaml
apiVersion: kfg.dev/v1alpha1
kind: Cmd
metadata:
  name: tutorial.cmd.greet
  commandName: greet
spec:
  env:
    NAME: "{env:USER:-Anonymous}"
    GREETING: "{env:GREETING:-Hello}"
  run: echo "$GREETING, $NAME!"
```

**Explicação:**
- `{env:VAR}`: Placeholder resolvido em tempo de geração
- `{env:VAR:-default}`: Com valor padrão
- Em tempo de execução, vira `$VAR`

### Teste com diferentes valores:

```bash
# Valor padrão
greet
# Output: Hello, Anonymous!

# Com variáveis de ambiente
USER=João GREETING=Olá greet
# Output: Olá, João!
```

## Trabalhando com Steps

Steps são unidades de trabalho reutilizáveis que podem ter outputs.

### 1. Crie um step

```yaml
apiVersion: kfg.dev/v1alpha1
kind: Step
metadata:
  name: tutorial.steps.check-file
spec:
  run: |
    if [ -f "config.yaml" ]; then
      echo "found"
    else
      echo "missing"
    fi
  output:
    name: FILE_STATUS
    type: string
```

### 2. Use o step em um workflow

```yaml
apiVersion: kfg.dev/v1alpha1
kind: Cmd
metadata:
  name: tutorial.cmd.deploy
  commandName: deploy
spec:
  run: echo "Deploying application..."

---
apiVersion: kfg.dev/v1alpha1
kind: CmdWorkflow
metadata:
  name: tutorial.workflow.deploy
spec:
  cmds: [tutorial.cmd.deploy]
  before:
    - step: tutorial.steps.check-file
```

### 3. Aplique e execute

```bash
kfg apply -f workflow.yaml --workflow deploy
deploy
```

O step `check-file` executa antes do comando `deploy`.

## Execução Condicional

Você pode executar steps condicionalmente baseado em outputs:

```yaml
apiVersion: kfg.dev/v1alpha1
kind: CmdWorkflow
metadata:
  name: tutorial.workflow.smart-deploy
spec:
  cmds: [tutorial.cmd.deploy]
  before:
    - step: tutorial.steps.check-file
    - step: tutorial.steps.validate-config
      when:
        output:
          step: tutorial.steps.check-file
          name: FILE_STATUS
          equals: "found"
```

**Operadores disponíveis:**
- `equals`: Igual a
- `in`: Em lista
- `contains`: Contém substring
- `matches`: Regex match
- `allOf`: Todas condições
- `anyOf`: Qualquer condição
- `not`: Negação

## Trabalhando com Kustomization

Kustomization permite compor manifests de múltiplas fontes.

### Estrutura de diretórios

```
myproject/
├── base/
│   ├── kustomization.yaml
│   ├── cmd-deploy.yaml
│   └── steps-validate.yaml
└── overlays/
    └── dev/
        ├── kustomization.yaml
        └── cmd-dev.yaml
```

### base/kustomization.yaml

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - cmd-deploy.yaml
  - steps-validate.yaml
```

### overlays/dev/kustomization.yaml

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - ../../base
  - cmd-dev.yaml
```

### Aplique o overlay

```bash
kfg apply -k overlays/dev --workflow deploy
```

## Cache de Steps

Steps podem ser cacheadas para evitar re-execuções:

```yaml
apiVersion: kfg.dev/v1alpha1
kind: Step
metadata:
  name: tutorial.steps.expensive
  annotations:
    kfg.dev/cacheable: "true"
spec:
  run: |
    echo "Running expensive operation..."
    sleep 5
    echo "done"
  output:
    name: RESULT
    type: string
```

### Gerenciando o cache

```bash
# Listar entradas do cache
kfg sys cache ls

# Inspecionar entrada específica
kfg sys cache inspect <id>

# Remover entradas antigas
kfg sys cache prune

# Invalidar cache na próxima execução
KFG_REFRESH=1 kfg apply -f workflow.yaml --workflow deploy
```

## Debug e Troubleshooting

### Aumentar verbosidade

```bash
# Níveis: 0 (quiet) a 5 (máximo debug)
KFG_VERBOSE=3 kfg apply -f workflow.yaml --workflow deploy
```

### Ver logs

```bash
# Localização padrão dos logs
ls ~/.local/state/kfg/logs/

# Ou caminho customizado
KFG_LOG_DIR=/tmp/kfg-logs kfg apply -f workflow.yaml --workflow deploy
```

### Ver código gerado

```bash
# Gerar sem executar
kfg build overlays/dev -o generated.yaml

# Ver conteúdo
cat generated.yaml
```

## Próximos Passos

Agora que você conhece o básico, explore:

- **[Manifest Model](manifest-model.md)** - Schema completo de manifests
- **[CLI Reference](cli-reference.md)** - Todos os comandos e flags
- **[Architecture](architecture.md)** - Como o KFG funciona internamente
- **[Exemplos](../packages/domains/ai-agents/manifests/)** - Manifests reais do projeto

## Exemplo Completo: Pipeline de CI/CD

Aqui está um exemplo real de pipeline de CI/CD:

```yaml
# steps.yaml
apiVersion: kfg.dev/v1alpha1
kind: Step
metadata:
  name: cicd.steps.checkout
spec:
  run: |
    git clone "$REPO_URL" /tmp/build
    echo "checked_out"
  output:
    name: STATUS
    type: string

---
apiVersion: kfg.dev/v1alpha1
kind: Step
metadata:
  name: cicd.steps.test
spec:
  run: |
    cd /tmp/build
    make test
    echo "tests_passed"
  output:
    name: TEST_STATUS
    type: string

---
apiVersion: kfg.dev/v1alpha1
kind: Step
metadata:
  name: cicd.steps.build
spec:
  run: |
    cd /tmp/build
    make build
    echo "build_complete"
  output:
    name: BUILD_STATUS
    type: string

---
# commands.yaml
apiVersion: kfg.dev/v1alpha1
kind: Cmd
metadata:
  name: cicd.cmd.deploy-staging
  commandName: deploy-staging
spec:
  env:
    ENV: "staging"
  run: |
    cd /tmp/build
    kubectl apply -f k8s/$ENV/

---
apiVersion: kfg.dev/v1alpha1
kind: Cmd
metadata:
  name: cicd.cmd.deploy-production
  commandName: deploy-production
spec:
  env:
    ENV: "production"
  run: |
    cd /tmp/build
    kubectl apply -f k8s/$ENV/

---
# workflow.yaml
apiVersion: kfg.dev/v1alpha1
kind: CmdWorkflow
metadata:
  name: cicd.workflow.full-pipeline
spec:
  cmds: [cicd.cmd.deploy-staging]
  before:
    - step: cicd.steps.checkout
      weight: -100
    - step: cicd.steps.test
      weight: -90
      when:
        output:
          step: cicd.steps.checkout
          name: STATUS
          equals: "checked_out"
    - step: cicd.steps.build
      weight: -80
      when:
        output:
          step: cicd.steps.test
          name: TEST_STATUS
          equals: "tests_passed"
  after:
    - step: kfg.cleanup
      weight: 100
```

### Use o pipeline

```bash
# Pipeline completo
REPO_URL=https://github.com/user/repo kfg apply -f cicd.yaml --workflow full-pipeline

# Executar deploy
deploy-staging

# Ou deploy de produção (com outro workflow)
deploy-production
```

Este pipeline:
1. Faz checkout do código
2. Roda testes (se checkout OK)
3. Builda (se testes passaram)
4. Deploy para staging
5. Limpa recursos temporários

---

**Dúvidas?** Veja [Troubleshooting](troubleshooting.md) ou abra uma [issue](https://github.com/seregatte/kfg/issues).
