package ui

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/AlexandreSJ/aoi/internal/config"
	tea "github.com/charmbracelet/bubbletea"
)

type fileSelectModel struct {
	layout Layout
	cfg    *config.Config

	mode    gameMode
	files   []string
	cursor  int
	scroll  int
	isWords bool
}

func newFileSelectModel(cfg *config.Config, layout Layout, mode gameMode) fileSelectModel {
	f := fileSelectModel{
		cfg:    cfg,
		layout: layout,
		mode:   mode,
	}
	f.isWords = mode != modeQuote
	f.loadFiles()
	return f
}

func (f fileSelectModel) setSize(w, h int) fileSelectModel {
	f.layout = f.layout.SetSize(w, h)
	return f
}

func (f *fileSelectModel) loadFiles() {
	if f.isWords {
		dir, _ := config.ResolveDir(f.cfg.System.WordsDir, "~/.config/aoi/words")
		f.files = scanTxtFiles(dir)
	} else {
		dir, _ := config.ResolveDir(f.cfg.System.QuotesDir, "~/.config/aoi/quotes")
		f.files = scanTxtFiles(dir)
	}
	if len(f.files) == 0 {
		f.files = []string{"en"}
	}
}

func (f fileSelectModel) selectedFile() string {
	if f.cursor >= 0 && f.cursor < len(f.files) {
		return f.files[f.cursor]
	}
	return "en"
}

func (f fileSelectModel) Update(msg tea.Msg) (fileSelectModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if f.cursor > 0 {
				f.cursor--
				f = f.adjustScroll()
			}
		case "down", "j":
			if f.cursor < len(f.files)-1 {
				f.cursor++
				f = f.adjustScroll()
			}
		}
	}
	return f, nil
}

func (f fileSelectModel) adjustScroll() fileSelectModel {
	avail := f.layout.BodyHeight(f.footerSegments())
	if avail < 1 {
		avail = 1
	}
	if f.cursor < f.scroll {
		f.scroll = f.cursor
	}
	if f.cursor >= f.scroll+avail {
		f.scroll = f.cursor - avail + 1
	}
	if f.scroll < 0 {
		f.scroll = 0
	}
	return f
}

func (f fileSelectModel) footerSegments() []string {
	dir := "words"
	if !f.isWords {
		dir = "quotes"
	}
	return []string{
		fmt.Sprintf("mode: %s | %s", f.mode, dir),
		"\u2191/\u2193: pick",
		"enter: start",
		"esc: back",
	}
}

func (f fileSelectModel) View() string {
	if f.layout.Width == 0 {
		return ""
	}

	dir := "words"
	if !f.isWords {
		dir = "quotes"
	}

	footerSegs := f.footerSegments()
	bodyHeight := f.layout.BodyHeight(footerSegs)
	body := f.renderList(bodyHeight, dir)
	body = f.layout.CenterBody(body, bodyHeight)
	return f.layout.Render("A O I", footerSegs, body)
}

func (f fileSelectModel) renderList(maxLines int, dir string) string {
	avail := f.layout.BodyHeight(f.footerSegments())
	if avail < 1 {
		avail = 1
	}

	var b strings.Builder
	b.WriteString(f.layout.Styles.Subtitle.Render(fmt.Sprintf("  Select %s file:", dir)) + "\n\n")

	lineCount := 2
	for i := f.scroll; i < len(f.files) && lineCount < maxLines; i++ {
		name := f.files[i]
		if i == f.cursor {
			b.WriteString(f.layout.Styles.Marker.Render(fmt.Sprintf("  > %s", name)))
		} else {
			b.WriteString(f.layout.Styles.Dim.Render(fmt.Sprintf("    %s", name)))
		}
		b.WriteString("\n")
		lineCount++
	}

	return b.String()
}

func scanTxtFiles(dir string) []string {
	if dir == "" {
		return nil
	}
	entries, err := filepath.Glob(filepath.Join(dir, "*.txt"))
	if err != nil {
		return nil
	}
	var names []string
	for _, e := range entries {
		base := filepath.Base(e)
		name := strings.TrimSuffix(base, ".txt")
		names = append(names, name)
	}
	return names
}
