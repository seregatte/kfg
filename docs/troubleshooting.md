# Troubleshooting

Este guia ajuda a diagnosticar e resolver problemas comuns com o KFG.

## Problemas de Instalação

### `kfg: command not found`

**Causa**: O binário do KFG não está no PATH.

**Soluções**:

1. **Via Nix**: Adicione ao shell atual
   ```bash
   nix shell github:seregatte/kfg
   ```

2. **Build do fonte**: Adicione ao PATH
   ```bash
   export PATH="$PWD/bin:$PATH"
   # Ou copie para um diretório no PATH
   sudo cp bin/kfg /usr/local/bin/
   ```

3. **GOPATH**: Verifique se está no PATH
   ```bash
   export PATH="$GOPATH/bin:$PATH"
   ```

### Nix: `flakes experimental feature is disabled`

**Causa**: Flakes não estão habilitados no Nix.

**Solução**: Habilite flakes em `~/.config/nix/nix.conf`:
```
experimental-features = nix-command flakes
```

Ou use o flag `--extra-experimental-features`:
```bash
nix --extra-experimental-features "nix-command flakes" shell github:seregatte/kfg
```

### Erro ao buildar: `go: command not found`

**Causa**: Go não está instalado ou não está no PATH.

**Solução**:
```bash
# Instalar Go (Linux)
sudo apt install golang-go

# Ou via Nix
nix shell nixpkgs#go
```

Verifique a versão (requer 1.21+):
```bash
go version
```

## Problemas de Execução

### `Error: manifest validation failed`

**Causa**: Manifest YAML inválido.

**Diagnóstico**:
```bash
# Aumentar verbosidade
KFG_VERBOSE=3 kfg apply -f manifest.yaml --workflow test

# Validar YAML
python3 -c "import yaml; yaml.safe_load(open('manifest.yaml'))"
```

**Causas comuns**:
- `apiVersion` incorreto (deve ser `kfg.dev/v1alpha1`)
- `kind` inválido (Cmd, CmdWorkflow, Step, etc.)
- Campos obrigatórios faltando
- Indentação YAML incorreta

**Exemplo correto**:
```yaml
apiVersion: kfg.dev/v1alpha1  # Obrigatório
kind: Cmd                      # Obrigatório
metadata:
  name: myapp.cmd.example      # Obrigatório
  commandName: example         # Obrigatório para Cmd
spec:
  run: echo "hello"            # Obrigatório
```

### `Error: workflow not found`

**Causa**: O workflow especificado não existe no manifest.

**Solução**:
```bash
# Listar workflows disponíveis
kfg build -k ./manifests | grep "kind: CmdWorkflow" -A 2

# Ou verificar o nome exato
grep "name:.*workflow" manifest.yaml
```

### `Error: circular dependency detected`

**Causa**: Steps formam um ciclo no DAG.

**Exemplo problemático**:
```yaml
before:
  - step: step-a
    when:
      output:
        step: step-b  # step-b depende de step-a
        name: STATUS
        equals: "ok"
  - step: step-b
    when:
      output:
        step: step-a  # step-a depende de step-b
        name: STATUS
        equals: "ok"
```

**Solução**: Reorganize as dependências para evitar ciclos.

### Comando não funciona após `kfg apply`

**Causa**: O código gerado não foi sourceado no shell atual.

**Solução**:
```bash
# Opção 1: Usar --interactive para abrir shell interativo
kfg apply -f manifest.yaml --workflow test --interactive

# Opção 2: Source manualmente
eval "$(kfg apply -f manifest.yaml --workflow test --print)"

# Opção 3: Adicionar ao .bashrc
echo 'eval "$(kfg apply -f manifest.yaml --workflow test --print)"' >> ~/.bashrc
```

## Problemas de Cache

### `Error: cache corrupted`

**Causa**: Arquivos de cache corrompidos.

**Solução**:
```bash
# Limpar todo o cache
rm -rf ~/.kfg/store/cache

# Ou apenas entradas específicas
kfg sys cache ls
kfg sys cache rm <id>

# Reconstruir na próxima execução
KFG_REFRESH=1 kfg apply -f manifest.yaml --workflow test
```

### Steps não estão sendo cacheadas

**Causa**: Step não está marcada como cacheable.

**Solução**:
```yaml
kind: Step
metadata:
  name: myapp.steps.expensive
  annotations:
    kfg.dev/cacheable: "true"  # Adicionar esta anotação
spec:
  run: ...
```

### Cache não invalida após mudar manifest

**Causa**: O hash do step não mudou (mesmo código).

**Solução**:
```bash
# Forçar invalidação
KFG_REFRESH=1 kfg apply -f manifest.yaml --workflow test

# Ou remover cache manualmente
kfg sys cache rm <step-id>
```

## Problemas de Logging

### Logs não aparecem

**Causa**: Verbosidade muito baixa.

**Solução**:
```bash
# Aumentar verbosidade
KFG_VERBOSE=5 kfg apply -f manifest.yaml --workflow test
```

**Níveis de verbosidade**:
- `0`: Quiet (sem output)
- `1`: Error + Warn + Info (padrão)
- `2`: + Detail
- `3`: + Warn/Detail
- `4`: + Debug
- `5`: + Debug verbose

### Não sei onde estão os logs

**Localização padrão**:
```bash
# Linux
~/.local/state/kfg/logs/kfg.jsonl

# macOS
~/Library/Application Support/kfg/logs/kfg.jsonl

# Ou caminho customizado
echo $KFG_LOG_FILE
```

**Ver logs em tempo real**:
```bash
tail -f ~/.local/state/kfg/logs/kfg.jsonl | jq .
```

### Logs muito verbosos

**Solução**: Reduza a verbosidade
```bash
KFG_VERBOSE=1 kfg apply -f manifest.yaml --workflow test
```

## Problemas de Performance

### `kfg apply` muito lento

**Causas possíveis**:

1. **Muitos manifests para processar**
   ```bash
   # Diagnosticar
   KFG_VERBOSE=3 kfg apply -f manifest.yaml --workflow test 2>&1 | grep "Parsing"
   
   # Solução: Usar kustomization ao invés de múltiplos -f
   kfg apply -k ./manifests --workflow test
   ```

2. **Steps não cacheadas**
   ```bash
   # Verificar cache
   kfg sys cache ls
   
   # Solução: Marcar steps como cacheable
   annotations:
     kfg.dev/cacheable: "true"
   ```

3. **Kustomize processando overlays complexos**
   ```bash
   # Diagnosticar
   KFG_VERBOSE=3 kfg build -k ./manifests
   
   # Solução: Simplificar overlays ou usar base diretamente
   kfg apply -k ./base --workflow test
   ```

### Uso excessivo de memória

**Causa**: Manifests muito grandes ou muitos recursos.

**Solução**:
```bash
# Dividir em múltiplas aplicações
kfg apply -f part1.yaml --workflow test1
kfg apply -f part2.yaml --workflow test2

# Ou usar kustomization para composição modular
kfg apply -k ./manifests --workflow test
```

## Problemas de Debug

### Como ver o código shell gerado?

```bash
# Gerar sem executar
kfg build -k ./manifests -o generated.yaml

# Ou para stdout
kfg build -k ./manifests

# Ver código de um workflow específico
kfg build -k ./manifests --workflow test
```

### Como debugar um step específico?

```bash
# Aumentar verbosidade
KFG_VERBOSE=5 kfg apply -f manifest.yaml --workflow test

# Executar step manualmente
bash -x -c "$(kfg build -k ./manifests | grep -A 50 '_kfg.step.my_step')"

# Ver outputs do step
echo $KFG_STEP_my_step_STATUS
```

### Como verificar variáveis de ambiente?

```bash
# Ver todas as variáveis KFG
env | grep KFG

# Ver variáveis de um step específico
kfg sys cache inspect <step-id> | jq .env
```

## Problemas Comuns com Kustomize

### `Error: failed to load kustomization`

**Causa**: kustomization.yaml inválido ou recursos faltando.

**Diagnóstico**:
```bash
# Validar kustomization
kustomize build ./manifests

# Ver recursos referenciados
grep "resources:" -A 10 ./manifests/kustomization.yaml
```

**Solução**: Verificar se todos os arquivos referenciados existem.

### Overlays não estão sendo aplicados

**Causa**: Ordem de recursos ou patches incorretos.

**Diagnóstico**:
```bash
# Ver resultado do build
kustomize build ./overlays/dev

# Verificar patches
cat ./overlays/dev/kustomization.yaml
```

**Solução**: Verificar sintaxe de patches e ordem de recursos.

### Recursos duplicados

**Causa**: Mesmo recurso definido em base e overlay.

**Solução**: Usar patches ao invés de redefinir recursos:
```yaml
# overlays/dev/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - ../../base
patches:
  - target:
      kind: Cmd
      name: base.cmd.deploy
    patch: |
      - op: replace
        path: /spec/env/ENV
        value: development
```

## Erros de Runtime

### `Error: placeholder resolution failed`

**Causa**: Variável de ambiente não existe.

**Solução**:
```yaml
# Ruim: variável pode não existir
env:
  API_KEY: "{env:API_KEY}"

# Bom: com valor padrão
env:
  API_KEY: "{env:API_KEY:-default_key}"
```

### `Error: step output not available`

**Causa**: Tentativa de usar output de step que ainda não executou.

**Solução**: Garantir ordem correta no workflow:
```yaml
before:
  - step: step-a  # Executa primeiro
  - step: step-b  # Pode usar output de step-a
    when:
      output:
        step: step-a
        name: STATUS
        equals: "ok"
```

### `Error: command execution failed`

**Causa**: Comando shell falhou (exit code != 0).

**Diagnóstico**:
```bash
# Aumentar verbosidade
KFG_VERBOSE=5 kfg apply -f manifest.yaml --workflow test

# Executar manualmente
bash -x -c "comando_que_falhou"
```

**Solução**: Corrigir o comando ou adicionar tratamento de erro:
```yaml
spec:
  run: |
    set -e  # Falhar em qualquer erro
    comando || { echo "Erro: comando falhou"; exit 1; }
```

## Como Reportar Bugs

Se você encontrou um bug que não está neste guia:

1. **Verifique issues existentes**: https://github.com/seregatte/kfg/issues
2. **Colete informações**:
   ```bash
   kfg version
   KFG_VERBOSE=5 kfg apply -f manifest.yaml --workflow test 2>&1 | tee debug.log
   ```
3. **Abra uma issue**: https://github.com/seregatte/kfg/issues/new
   - Descreva o problema
   - Inclua `debug.log`
   - Inclua o manifest YAML (remova dados sensíveis)

## Recursos Adicionais

- **[Getting Started](getting-started.md)** - Tutorial passo a passo
- **[CLI Reference](cli-reference.md)** - Referência completa
- **[Architecture](architecture.md)** - Arquitetura interna
- **[GitHub Discussions](https://github.com/seregatte/kfg/discussions)** - Perguntas e discussões
- **[GitHub Issues](https://github.com/seregatte/kfg/issues)** - Reportar bugs
