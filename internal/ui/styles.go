package ui

import (
	"fmt"
	"strings"

	"github.com/AlexandreSJ/aoi/internal/config"
	"github.com/charmbracelet/lipgloss"
)

const fallbackPrimary = "32"
const fallbackSecondary = "57"
const fallbackText = "231"
const fallbackDim = "0"
const fallbackTitle = "27"
const fallbackFooter = "16"
const fallbackError = "162"
const fallbackSuccess = "39"

const footerVersion = "Aoi v0.1.0"

const layoutOverhead = 3

type Styles struct {
	Border   lipgloss.Style
	Footer   lipgloss.Style
	Title    lipgloss.Style
	Subtitle lipgloss.Style
	Cursor   lipgloss.Style
	Error    lipgloss.Style
	Success  lipgloss.Style
	Dim      lipgloss.Style
	Text     lipgloss.Style
	Marker   lipgloss.Style
}

func colorOr(cfg *config.Config, value, fallback string) string {
	if value != "" {
		return value
	}
	return fallback
}

func StylesFromConfig(cfg *config.Config) Styles {
	primary := colorOr(cfg, cfg.Colors.Primary, fallbackPrimary)
	footer := colorOr(cfg, cfg.Colors.Footer, fallbackFooter)
	titleColor := colorOr(cfg, cfg.Colors.Title, fallbackTitle)
	errorColor := colorOr(cfg, cfg.Colors.Error, fallbackError)
	successColor := colorOr(cfg, cfg.Colors.Success, fallbackSuccess)
	dimColor := colorOr(cfg, cfg.Colors.Dim, fallbackDim)
	textColor := colorOr(cfg, cfg.Colors.Text, fallbackText)

	return Styles{
		Border: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(primary)).
			Padding(1, 2),

		Footer: lipgloss.NewStyle().
			Foreground(lipgloss.Color(footer)).
			Background(lipgloss.Color(primary)).
			Padding(0, 1),

		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(titleColor)).
			Background(lipgloss.Color(primary)).
			Align(lipgloss.Center),

		Subtitle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(primary)).
			Bold(true),

		Cursor: lipgloss.NewStyle().
			Foreground(lipgloss.Color(primary)).
			Bold(true),

		Error: lipgloss.NewStyle().
			Foreground(lipgloss.Color(errorColor)),

		Success: lipgloss.NewStyle().
			Foreground(lipgloss.Color(successColor)),

		Dim: lipgloss.NewStyle().
			Foreground(lipgloss.Color(dimColor)),

		Text: lipgloss.NewStyle().
			Foreground(lipgloss.Color(textColor)),

		Marker: lipgloss.NewStyle().
			Foreground(lipgloss.Color(primary)).
			Bold(true),
	}
}

func BodyHeight(totalHeight, footerHeight int) int {
	h := totalHeight - 1 - footerHeight - layoutOverhead
	if h < 1 {
		return 1
	}
	return h
}

type Layout struct {
	Width  int
	Height int
	Styles Styles
}

func (l Layout) SetSize(w, h int) Layout {
	l.Width = w
	l.Height = h
	return l
}

func (l Layout) BodyHeight(footerSegments []string) int {
	footerH := l.Styles.RenderFooterHeight(footerSegments, l.Width)
	return BodyHeight(l.Height, footerH)
}

func (l Layout) Render(title string, footerSegments []string, bodyContent string) string {
	bodyHeight := l.BodyHeight(footerSegments)
	return l.Styles.Layout(l.Width, l.Height, title, footerSegments, bodyContent, bodyHeight)
}

func (l Layout) CenterBody(bodyContent string, bodyHeight int) string {
	innerWidth := maxWidth(1, l.Width-8) // match typing screen: width minus border+padding
	bodyContent = strings.TrimRight(bodyContent, "\n")
	lines := strings.Split(bodyContent, "\n")

	// Horizontal centering per line
	var centered strings.Builder
	for _, line := range lines {
		centered.WriteString(lipgloss.NewStyle().Width(innerWidth).Align(lipgloss.Center).Render(line))
		centered.WriteString("\n")
	}

	// Vertical centering
	contentLines := len(lines)
	padAbove := 0
	if bodyHeight > contentLines {
		padAbove = (bodyHeight - contentLines) / 2
	}

	var b strings.Builder
	for i := 0; i < padAbove; i++ {
		b.WriteString("\n")
	}
	b.WriteString(centered.String())
	return b.String()
}

func maxWidth(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (s Styles) RenderFooter(segments []string, width int) string {
	if width < 1 {
		width = 1
	}
	sep := " \u2502 "
	availWidth := width - 2

	widths := make([]int, len(segments))
	for i, seg := range segments {
		widths[i] = lipgloss.Width(seg)
	}

	totalWidth := lipgloss.Width(strings.Join(segments, sep))
	if totalWidth <= availWidth {
		return s.Footer.Width(width).Render(strings.Join(segments, sep))
	}

	sepW := lipgloss.Width(sep)
	line1W := 0
	splitIdx := 0
	for i, w := range widths {
		need := w
		if i > 0 {
			need += sepW
		}
		if line1W+need > availWidth {
			break
		}
		line1W += need
		splitIdx = i + 1
	}
	if splitIdx == 0 {
		splitIdx = 1
	}
	content := strings.Join(segments[:splitIdx], sep) + "\n" + strings.Join(segments[splitIdx:], sep)
	return s.Footer.Width(width).Render(content)
}

func (s Styles) RenderFooterHeight(segments []string, width int) int {
	return lipgloss.Height(s.RenderFooter(segments, width))
}

func FooterSwatch(color256, footerFg string) string {
	if footerFg == "" {
		footerFg = "16"
	}
	return fmt.Sprintf("\x1b[38;5;%sm\u2588\u2588\u2588\u2588\x1b[38;5;%sm", color256, footerFg)
}

func (s Styles) Layout(width, height int, titleText string, footerSegments []string, bodyContent string, bodyHeight int) string {
	borderStyle := s.Border.Padding(0, 2)

	titleView := s.Title.Width(width).Render(titleText)
	footerView := s.RenderFooter(footerSegments, width)

	body := borderStyle.
		Width(max(1, width-2)).
		Height(bodyHeight).
		Render(bodyContent)

	return strings.Join([]string{titleView, body, footerView}, "\n")
}
