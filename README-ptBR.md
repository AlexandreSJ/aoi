<div align="center">
  <img width=100% src="https://capsule-render.vercel.app/api?type=waving&height=200&color=0:00aaff,100:00aaff&text=A%20O%20I&fontColor=ffffff&fontSize=50&fontAlignY=40" alt="AOI">
</div>

<h1 align="center">🔹 あおい 🔹</h1>

<p align="center"> 
  Um teste de digitação no terminal. 
  <br>
  Pratique sua digitação, relaxe e curta com aoi.
</p>

<div align="center">

  <img src="https://img.shields.io/badge/Go-1.24-00add8?style=flat-square&logo=go" alt="Go" href="https://go.dev/doc/install">
  <img src="https://img.shields.io/badge/Bubble_Tea-1.3-ff69b4?style=flat-square" alt="Bubble Tea" href="https://github.com/charmbracelet/bubbletea">
  <img src="https://img.shields.io/badge/Lipgloss-1.1-7d56f4?style=flat-square" alt="Lipgloss" href="https://github.com/charmbracelet/lipgloss">
  <img src="https://img.shields.io/badge/Licença-MIT-00aaff?style=flat-square" alt="Licença" href="LICENSE">
  <br>
  <a href="https://www.buymeacoffee.com/aelxand" target="_blank">
    <img width=120px src="assets/bmc/bmc.png" alt="Buy Me A Coffee">
  </a>
  <br>
  <a href="README-ptBR.md"><img src="https://img.shields.io/badge/ptBR🇧🇷-README-0af?style=flat-square" alt="License" href="LICENSE"></a>
  <a href="README.md"><img src="https://img.shields.io/badge/en🇬🇧-README-fff?style=flat-square" alt="License" href="LICENSE"></a>
  <a href="README-es.md"><img src="https://img.shields.io/badge/es🇪🇸-README-fff?style=flat-square" alt="License" href="LICENSE"></a>
</div>

## Sumário

- [O que é o Aoi?](#o-que-é-o-aoi)
- [Instalação](#instalação)
- [Uso](#uso)
- [Recursos](#recursos)
- [Configuração](#configuração)
- [Solução de Problemas](#solução-de-problemas)
- [Licença](#licença)

## O que é o Aoi?

Comecei a gostar de fazer testes de digitação como hobby e para manter minhas habilidades de digitação afiadas, mas sempre quis isso em uma TUI. Então criei o AOI!

Escolha entre 4 modos diferentes de prática de digitação no Aoi:
- Zen: Digite infinitamente no seu próprio ritmo
- Temporizado: Corra contra o relógio
- Contado: Digite um número fixo de palavras
- Citação: Digite uma citação aleatória

Configure as cores do jeito que quiser. Você também pode adicionar mais palavras ou citações, escalável para usar qualquer idioma!

<div align="center">
  <img width=100% src="assets/prints/typing.png" alt="Digitando">
</div>

## Instalação

### Pré-requisitos

- Go 1.24+ (necessário para compilar a partir do código-fonte)
- Emulador de terminal com suporte a cores ANSI e Unicode

### Métodos de Instalação

#### Método 1: Instalar usando o Go

```bash
# Certifique-se de ter o Go configurado no seu ~/.zshrc ou ~/.bashrc
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin

# Feche o terminal ou execute para aplicar
source ~/.zshrc 
# Ou
source ~/.bashrc

# Por fim, instale diretamente do GitHub
go install -a github.com/AlexandreSJ/aoi/cmd/aoi@latest
```

#### Método 2: Compilar a partir do Código-Fonte (para desenvolvedores)

```bash
# Clone o repositório
git clone https://github.com/AlexandreSJ/aoi.git
cd aoi

# Compile a aplicação
make build
```

### Início Rápido

Após a instalação, simplesmente execute:

```bash
aoi
```

### Comandos de Build

Se você tem o repositório git instalado, em /aoi você pode executar:

```bash
make clean  # Remove o diretório /build
make build  # Compila o binário
make run    # Compila e executa imediatamente
```

### Funcionalidades

- **Leve e rápido** - Leve e veloz como um ouriço
- **Feedback de digitação em tempo real** - Veja sua precisão e velocidade enquanto digita
- **Suporte a Unicode** - Funciona com vários conjuntos de caracteres
- **Design responsivo** - Adapta-se a diferentes tamanhos de terminal

### Requisitos do Sistema

- **Sistema Operacional**: Linux, macOS ou Windows (com WSL)
- **Terminal**: Qualquer emulador de terminal moderno (Terminal, iTerm2, Alacritty, Windows Terminal, etc.)
- **Espaço em Disco**: ~5MB para o binário

### Solução de Problemas

**P: Recebo "command not found: aoi"**
R: Certifique-se de que o diretório GOPATH/bin está no seu PATH, ou use o caminho completo para o binário.

**P: As cores parem estranhas no meu terminal**
R: Tente definir `TERM=xterm-256color` ou use um terminal com suporte a true color.

**P: A aplicação não inicia**
R: Certifique-se de ter o Go 1.24+ instalado e que seu terminal suporta caracteres Unicode.

**P: Estou tendo problemas para instalar/atualizar o aoi para a versão mais recente**
R: Se você já tem o Go instalado, execute o seguinte comando para evitar o proxy.golang.org e usar a tag `-a` para forçar a recompilação:
`GOPROXY=direct go install -a github.com/AlexandreSJ/aoi/cmd/aoi@latest`

## Roadmap

> Plano de desenvolvimento em fases. Cada fase mistura trabalho de múltiplas áreas.

### Fase 1 — Persistência de Configuração e Temas `v0.3`


**Gerenciamento de Configuração** (Concluído para v0.3.1)
- Melhorias no padding e espaçamento da UI de entrada
- Correções na entrada de configuração

**Persistência de configuração** (Em andamento)
- Lembrar último modo usado, duração temporizada, contagem de palavras, último arquivo por modo, preferência de tamanho de texto
- Campo de versão na config com migração automática em mudanças de schema
- Opção cursor ligado/desligado
- Opção flash de tecla ligado/desligado
- Refatorar tratamento de erros do `config.go`

**Sistema de temas**
- 4+ temas predefinidos (Aoi padrão, Monokai, Dracula, Solarized, Catppuccin)
- Tela de seleção de tema acessível pela home ou configuração
- Pré-visualização de tema antes de aplicar
- Temas personalizados de `~/.config/aoi/themes/` como arquivos YAML
- Importar/exportar arquivos de tema para compartilhar

### Fase 2 — Telemetria e Estatísticas `v0.4`

**Armazenamento de telemetria (JSONL)**
- Uma linha por sessão de teste em `~/.config/aoi/telemetry.jsonl`
- Por sessão: modo, WPM, caracteres digitados, palavras completas, ok/erros, precisão, duração, timestamp
- Timestamps por caractere para análise de velocidade
- Somente adição, sem dependências externas

**Tela de estatísticas**
- Nova tela acessível pela home (tecla `s`)
- Média de WPM por modo
- Lista de histórico de sessões — navegável por data/modo/WPM
- Mapa de calor por tecla — teclas mais lentas e com mais erros
- Melhores resultados por modo — exibidos na tela home
- Rastreamento de sequência — dias consecutivos de prática, exibidos na home
- Exportar telemetria como CSV

**Gráficos e visualizações**
- Gráfico de WPM ao longo do tempo (sparkline ou gráfico de barras)
- Tendência de precisão
- Barras de desempenho por tecla
- Tudo renderizado com Lipgloss/Bubble Tea (sem lib de gráficos externa)

### Fase 3 — UX/UI e Som `v0.5`

**Barra de progresso**
- Barra de progresso estilizada para modo Temporizado (tempo restante) e modo Contado (palavras digitadas/restantes)
- Renderizada acima ou abaixo da área de digitação

**Refatoração da tela de digitação**
- Experiência de leitura mais suave — melhor espaçamento entre linhas, renderização de espaçamento entre palavras
- Indicador de WPM em tempo real durante a digitação
- Refatorar exibição do timer (modo Temporizado) para mostrar na área da barra de progresso em vez do rodapé
- Repetição de palavras com erro — após o teste, opção de repraticar apenas palavras com erro

**Sistema de som**
- `gopxl/beep` para reprodução de áudio
- 3+ pacotes de sons embutidos (toque mecânico, clique suave, etc.)
- Controles separados: som de digitação ok, som de digitação erro, jingle de finalização
- Sons personalizados de `~/.config/aoi/sounds/`

**Refatoração do rodapé**
- Mover informações de tempo/modo para outros locais (área da barra de progresso)
- Rodapé mais limpo com apenas atalhos de teclado

## Contribuindo

Contribuições são bem-vindas! Sinta-se à vontade para enviar um Pull Request.

## Licença

MIT

<div align="center">
  <a href="https://git.io/typing-svg">
    <img src="https://readme-typing-svg.herokuapp.com?font=Fira+Code&duration=1&color=00AAFF&center=true&vCenter=true&repeat=false&width=435&lines=stay+blue+%3C3" alt="Typing SVG" />
  </a>
</div>

<img width=100% src="https://capsule-render.vercel.app/api?type=slice&height=300&color=00aaff&text=AOI&section=footer&fontAlign=22&fontAlignY=69&rotate=19&fontSize=50&fontColor=ffffff&desc=あおい&descAlignY=80&descAlign=22" alt="AOI">
