# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run

```bash
make build          # build binary to ./build/aoi
make run            # build + run
make clean          # remove build/
```

Run tests:

```bash
go test ./...                        # all tests
go test ./internal/config/           # single package
go test ./internal/words/            # single package
go test -run TestLoadCreatesDefault  # single test
```

No lint configuration exists. Use `go vet ./...` for static analysis.

## Architecture

TUI typing test built with Bubble Tea (Elm Architecture: Model-Update-View) + Lipgloss for styling. Entry point: `cmd/aoi/main.go`.

### Screen system

`internal/ui/app.go` — top-level `App` struct delegates to screen-specific models:
- `screenHome` → `homeModel` (`home.go`) — mode selection with Left/Right, Enter to start
- `screenConfig` → `configModel` (`config.go`) — config editor with inline text editing and ANSI 256-color picker
- `screenFileSelect` → `fileSelectModel` (`fileselect.go`) — word/quote file picker before starting test
- `screenTyping` → `typingModel` (`typing.go`) — typing test engine with character-level feedback

Screen routing happens in `App.Update()` via a `screen` enum. Each sub-model has its own `Update`/`View` methods and `setSize` for resize propagation. All models use value receivers (return modified copy — Bubble Tea pattern).

### Layout pattern

Every screen renders three vertical sections: title bar, bordered body, footer via `Styles.Layout()`. Footer uses `Styles.RenderFooter()` which auto-wraps segments across lines when terminal is narrow. Minimum terminal size: 64x15. Body height computed by `BodyHeight()` helper.

### Styling

`internal/ui/styles.go` — `Styles` struct holds all pre-built lipgloss styles (Border, Footer, Title, Subtitle, Cursor, Error, Success, Dim, Text, Marker). Built once from config via `StylesFromConfig()`. Fallback colors are centralized constants.

### Typing engine

`internal/ui/typing.go` — Characters stored as `[]typedChar` with states: Pending, Correct, Error. Cursor tracked via `cursor bool`. Word-boundary wrapping via `computeLines()` ensures words never split across lines. Infinite mode appends words dynamically. Auto-scrolling via `adjustScroll()` keeps cursor visible.

### Mode system

`internal/ui/home.go` — Four modes: Infinite, Timed (configurable via Up/Down or typing a number), Count (configurable), Quote. `homeModel` stores `timedSeconds` and `wordCount` for configuration.

### Word and quote systems

`internal/words/words.go` — Loads word lists from embedded files (`internal/words/embedded/*.txt`) or user directory (`~/.config/aoi/words/`). Supports random sampling and infinite generation.

`internal/quotes/quotes.go` — Same pattern as words but for quote lists (`internal/quotes/embedded/*.txt`, `~/.config/aoi/quotes/`). Each line = one quote. `QuoteList.Random()` picks one.

### Config system

`internal/config/` — YAML config loaded from `~/.config/aoi/config.yaml` (or custom path via `system.config` key). Default config embedded via `//go:embed default.yaml`. Config mutated in-place by the config screen, saved with `config.Save()`.

Config keys use dotted notation (`colors.primary`, `system.config`). `IsColorKey()` distinguishes color keys (show picker) from text keys (show inline editor).

### Theme validation

`internal/theme/validate.go` — `IsValidColor()` accepts hex (`#RGB`, `#RRGGBB`), ANSI 256 (0-255), and CSS named colors (16 standard + 8 bright).

## Dependencies

- `github.com/charmbracelet/bubbletea` — TUI framework
- `github.com/charmbracelet/lipgloss` — styling/layout
- `gopkg.in/yaml.v3` — config parsing

Go 1.24+. Module path: `github.com/AlexandreSJ/aoi`.
