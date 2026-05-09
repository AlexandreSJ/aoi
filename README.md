<div align="center">
  <img width=100% src="https://capsule-render.vercel.app/api?type=waving&height=200&color=0:00aaff,100:00aaff&text=A%20O%20I&fontColor=ffffff&fontSize=50&fontAlignY=40" alt="AOI">
</div>

<h1 align="center">🔹 あおい 🔹</h1>

<p align="center"> 
  A terminal-based typing test. 
  <br>
  Practice your typing, relax and vibe with aoi.
</p>

<div align="center">

  <img src="https://img.shields.io/badge/Go-1.24-00add8?style=flat-square&logo=go" alt="Go" href="https://go.dev/doc/install">
  <img src="https://img.shields.io/badge/Bubble_Tea-1.3-ff69b4?style=flat-square" alt="Bubble Tea" href="https://github.com/charmbracelet/bubbletea">
  <img src="https://img.shields.io/badge/Lipgloss-1.1-7d56f4?style=flat-square" alt="Lipgloss" href="https://github.com/charmbracelet/lipgloss">
  <img src="https://img.shields.io/badge/License-MIT-00aaff?style=flat-square" alt="License" href="LICENSE">

  <br>

  <a href="https://www.buymeacoffee.com/aelxand" target="_blank">
    <img width=120px src="assets/bmc/bmc.png" alt="Buy Me A Coffee">
  </a>

  <br>

  <a href="README-ptBR.md">
    (README 🇧🇷)
  </a>
</div>

## What is Aoi?

I started to like doing typing tests for hobby and keeping my digitation skills sharp, but always wanted it in a TUI. So i made AOI!

Choose 4 different modes of typing practice in Aoi:
- Zen: Type infinitely at your own pace
- Timed: Race against the clock
- Count: Type a fixed number of words
- Quote: Type a random quote

Configure the colors anyway you like. You can also add more words or quotes, scalable to use any language you want!

<div align="center">
  <img width=100% src="assets/prints/typing.png" alt="Typing">
</div>

## Installation

### Prerequisites

- Go 1.24+ (required for building from source)
- Terminal emulator that supports ANSI colors and Unicode

### Installation Methods

#### Method 1: Install using Go

```bash
# Make sure you have Go configured to your ~/.zshrc or ~/.bashrc
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin

# Close the terminal or run to apply
source ~/.zshrc 
# Or
source ~/.bashrc

# Finally, install directly from GitHub
go install -a github.com/AlexandreSJ/aoi/cmd/aoi@latest
```

#### Method 2: Build from Source (for developers)

```bash
# Clone the repository
git clone https://github.com/AlexandreSJ/aoi.git
cd aoi

# Build the application
make build
```

### Quick Start

After installation, simply run:

```bash
aoi
```

### Build Commands

If you have the git repo installed, ond /aoi you can run:

```bash
make clean  # Remove /build directory
make build  # Compile the binary
make run    # Build and run immediately
```

### Features

- **Lightweight and fast** - Light and quick as a hedgehog
- **Real-time typing feedback** - See your accuracy and speed as you type
- **Unicode support** - Works with various character sets
- **Responsive design** - Adapts to different terminal sizes

### System Requirements

- **Operating System**: Linux, macOS, or Windows (with WSL)
- **Terminal**: Any modern terminal emulator (Terminal, iTerm2, Alacritty, Windows Terminal, etc.)
- **Disk Space**: ~5MB for the binary

### Troubleshooting

**Q: I get "command not found: aoi"**
A: Make sure your GOPATH/bin directory is in your PATH, or use the full path to the binary.

**Q: The colors look strange in my terminal**
A: Try setting `TERM=xterm-256color` or use a terminal that supports true color.

**Q: The application won't start**
A: Ensure you have Go 1.24+ installed and that your terminal supports Unicode characters.

**Q: I am having trouble installing/updating aoi to the latest version**
A: If you already have Go installed, run the following command to avoid the proxy.golang.org and use `-a` tag to force rebuild:
`GOPROXY=direct go install -a github.com/AlexandreSJ/aoi/cmd/aoi@latest`

<div align="center">
  <a href="https://git.io/typing-svg">
    <img src="https://readme-typing-svg.herokuapp.com?font=Fira+Code&duration=1&color=00AAFF&center=true&vCenter=true&repeat=false&width=435&lines=stay+blue+%3C3" alt="Typing SVG" />
  </a>
</div>

<img width=100% src="https://capsule-render.vercel.app/api?type=slice&height=300&color=00aaff&text=AOI&section=footer&fontAlign=22&fontAlignY=69&rotate=19&fontSize=50&fontColor=ffffff&desc=あおい&descAlignY=80&descAlign=22" alt="AOI">