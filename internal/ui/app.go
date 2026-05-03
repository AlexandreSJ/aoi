package ui

import (
	"fmt"
	"strconv"

	"github.com/AlexandreSJ/aoi/internal/config"
	"github.com/AlexandreSJ/aoi/internal/quotes"
	"github.com/AlexandreSJ/aoi/internal/words"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	minTermWidth  = 64
	minTermHeight = 15
)

type screen int

const (
	screenHome screen = iota
	screenConfig
	screenFileSelect
	screenTyping
)

type App struct {
	width      int
	height     int
	screen     screen
	cfg        *config.Config
	styles     Styles
	home       homeModel
	config     configModel
	fileSelect fileSelectModel
	typing     typingModel
}

func NewApp() App {
	cfg, err := config.Load()
	if err != nil {
		cfg = config.Default()
	}

	bootstrapDirs(cfg)

	styles := StylesFromConfig(cfg)

	return App{
		cfg:    cfg,
		styles: styles,
		home:   newHomeModel(cfg, styles),
		config: newConfigModel(cfg, styles),
		typing: newTypingModel(cfg, styles),
	}
}

func bootstrapDirs(cfg *config.Config) {
	wordsDir, _ := config.ResolveDir(cfg.System.WordsDir, "~/.config/aoi/words")
	quotesDir, _ := config.ResolveDir(cfg.System.QuotesDir, "~/.config/aoi/quotes")

	words.EnsureUserDir(wordsDir)
	words.CopyEmbeddedToUser("en", wordsDir)

	quotes.EnsureUserDir(quotesDir)
	quotes.CopyEmbeddedToUser("en", quotesDir)
}

func (a App) Init() tea.Cmd {
	return nil
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		a.home = a.home.setSize(msg.Width, msg.Height)
		a.config = a.config.setSize(msg.Width, msg.Height)
		a.fileSelect = a.fileSelect.setSize(msg.Width, msg.Height)
		a.typing = a.typing.setSize(msg.Width, msg.Height)
		return a, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return a, tea.Quit
		}
	}

	switch a.screen {
	case screenHome:
		return a.updateHome(msg)
	case screenConfig:
		return a.updateConfig(msg)
	case screenFileSelect:
		return a.updateFileSelect(msg)
	case screenTyping:
		return a.updateTyping(msg)
	}
	return a, nil
}

func (a App) updateHome(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return a, tea.Quit
		case "c":
			a.screen = screenConfig
			a.config = newConfigModel(a.cfg, a.styles)
			a.config = a.config.setSize(a.width, a.height)
			return a, nil
		case "left", "h":
			if a.home.modeIdx > 0 {
				a.home.modeIdx--
			}
			a.home.configInput = ""
			return a, nil
		case "right", "l":
			if a.home.modeIdx < int(modeCountTotal)-1 {
				a.home.modeIdx++
			}
			a.home.configInput = ""
			return a, nil
		case "up", "k":
			a.home.configInput = ""
			switch a.home.selectedMode() {
			case modeTimed:
				a.home.timedSeconds += 5
				if a.home.timedSeconds > 300 {
					a.home.timedSeconds = 300
				}
			case modeCount:
				a.home.wordCount += 5
				if a.home.wordCount > 500 {
					a.home.wordCount = 500
				}
			}
			return a, nil
		case "down", "j":
			a.home.configInput = ""
			switch a.home.selectedMode() {
			case modeTimed:
				a.home.timedSeconds -= 5
				if a.home.timedSeconds < 5 {
					a.home.timedSeconds = 5
				}
			case modeCount:
				a.home.wordCount -= 5
				if a.home.wordCount < 5 {
					a.home.wordCount = 5
				}
			}
			return a, nil
		case "backspace":
			if a.home.configInput != "" {
				a.home.configInput = a.home.configInput[:len(a.home.configInput)-1]
				if a.home.configInput == "" {
					switch a.home.selectedMode() {
					case modeTimed:
						a.home.timedSeconds = defaultTimedSeconds
					case modeCount:
						a.home.wordCount = defaultWordCount
					}
				} else {
					a.applyHomeInput()
				}
			}
			return a, nil
		case "enter":
			if !a.home.configValid() {
				return a, nil
			}
			a.fileSelect = newFileSelectModel(a.cfg, a.styles, a.home.selectedMode())
			a.fileSelect = a.fileSelect.setSize(a.width, a.height)
			a.screen = screenFileSelect
			return a, nil
		default:
			if msg.Type == tea.KeyRunes {
				a.home = a.handleHomeNumberInput(msg.String())
			}
		}
	}
	return a, nil
}

func (a App) handleHomeNumberInput(input string) homeModel {
	if a.home.selectedMode() != modeTimed && a.home.selectedMode() != modeCount {
		return a.home
	}
	a.home.configInput += input
	a.applyHomeInput()
	return a.home
}

func (a *App) applyHomeInput() {
	n, err := strconv.Atoi(a.home.configInput)
	if err != nil {
		return
	}
	switch a.home.selectedMode() {
	case modeTimed:
		if n > 300 {
			n = 300
			a.home.configInput = "300"
		}
		if n >= 5 {
			a.home.timedSeconds = n
		}
	case modeCount:
		if n > 500 {
			n = 500
			a.home.configInput = "500"
		}
		if n >= 5 {
			a.home.wordCount = n
		}
	}
}

func (a App) updateConfig(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if !a.config.editing && !a.config.picker {
			switch keyMsg.String() {
			case "esc", "c":
				a.cfg = a.config.cfg
				a.styles = StylesFromConfig(a.cfg)
				a.home = newHomeModel(a.cfg, a.styles)
				a.home = a.home.setSize(a.width, a.height)
				a.typing = newTypingModel(a.cfg, a.styles)
				a.typing = a.typing.setSize(a.width, a.height)
				a.screen = screenHome
				return a, nil
			}
		}
	}

	var cmd tea.Cmd
	a.config, cmd = a.config.Update(msg)
	a.cfg = a.config.cfg
	return a, cmd
}

func (a App) updateFileSelect(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			a.screen = screenHome
			return a, nil
		case "enter":
			selectedFile := a.fileSelect.selectedFile()
			a.typing = newTypingModel(a.cfg, a.styles)
			a.typing.mode = a.home.selectedMode()
			a.typing.timedSeconds = a.home.timedSeconds
			a.typing.wordCountTarget = a.home.wordCount
			a.typing.wordListName = selectedFile
			a.typing = a.typing.setSize(a.width, a.height)
			if a.typing.mode == modeQuote {
				a.typing = a.typing.loadQuote()
			} else {
				a.typing = a.typing.loadWords()
			}
			a.screen = screenTyping
			return a, nil
		}
	}

	var cmd tea.Cmd
	a.fileSelect, cmd = a.fileSelect.Update(msg)
	return a, cmd
}

func (a App) updateTyping(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "esc" {
			a.screen = screenHome
			return a, nil
		}
	}

	var cmd tea.Cmd
	a.typing, cmd = a.typing.Update(msg)
	return a, cmd
}

func (a App) View() string {
	if a.width == 0 {
		return "Loading..."
	}

	if a.width < minTermWidth || a.height < minTermHeight {
		return fmt.Sprintf("Terminal too small (min %dx%d)", minTermWidth, minTermHeight)
	}

	switch a.screen {
	case screenConfig:
		return a.config.View()
	case screenFileSelect:
		return a.fileSelect.View()
	case screenTyping:
		return a.typing.View()
	default:
		return a.home.View()
	}
}
