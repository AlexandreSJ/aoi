package ui

import (
	"fmt"
	"strconv"

	"github.com/AlexandreSJ/aoi/internal/config"
)

type gameMode int

const (
	modeInfinite gameMode = iota
	modeTimed
	modeCount
	modeQuote
	modeCountTotal
)

var modeNames = []string{
	"Zen",
	"Timed",
	"Count",
	"Quote",
}

var modeDescriptions = []string{
	"type infinitely, just practice at your own pace",
	"race against the clock, type as many words as you can",
	"type a fixed number of words as fast as you can",
	"type a random quote from start to finish",
}

var defaultTimedSeconds = 30
var defaultWordCount = 25

func (m gameMode) String() string {
	if int(m) < len(modeNames) {
		return modeNames[m]
	}
	return "?"
}

type homeModel struct {
	width       int
	height      int
	styles      Styles
	cfg         *config.Config
	modeIdx     int
	wordFiles   []string
	wordIdx     int

	timedSeconds int
	wordCount    int
	configInput  string
}

func newHomeModel(cfg *config.Config, styles Styles) homeModel {
	return homeModel{
		cfg:          cfg,
		styles:       styles,
		modeIdx:      int(modeInfinite),
		timedSeconds: defaultTimedSeconds,
		wordCount:    defaultWordCount,
	}
}

func (h homeModel) setSize(w, height int) homeModel {
	h.width = w
	h.height = height
	return h
}

func (h homeModel) selectedMode() gameMode {
	return gameMode(h.modeIdx)
}

func (h homeModel) modeLabel() string {
	switch h.selectedMode() {
	case modeTimed:
		return fmt.Sprintf("Timed %ds", h.timedSeconds)
	case modeCount:
		return fmt.Sprintf("Count %d", h.wordCount)
	default:
		return h.selectedMode().String()
	}
}

func (h homeModel) configValid() bool {
	if h.configInput == "" {
		return true
	}
	n, err := strconv.Atoi(h.configInput)
	if err != nil {
		return false
	}
	return n >= 5
}

func (h homeModel) View() string {
	currentMode := h.selectedMode()

	modeDisplay := ""
	for i, name := range modeNames {
		if i == h.modeIdx {
			label := h.modeLabel()
			if !h.configValid() {
				modeDisplay += h.styles.Error.Render(fmt.Sprintf(" [%s] ", label))
			} else {
				modeDisplay += h.styles.Marker.Render(fmt.Sprintf(" [%s] ", label))
			}
		} else {
			modeDisplay += h.styles.Dim.Render(fmt.Sprintf("  %s  ", name))
		}
	}

	help := ""
	switch currentMode {
	case modeTimed:
		help = fmt.Sprintf(
			"\nPress %s/%s to adjust time. Type a number to set custom time.",
			h.styles.Marker.Render("\u2191"),
			h.styles.Marker.Render("\u2193"),
		)
	case modeCount:
		help = fmt.Sprintf(
			"\nPress %s/%s to adjust count. Type a number to set custom count.",
			h.styles.Marker.Render("\u2191"),
			h.styles.Marker.Render("\u2193"),
		)
	}

	desc := ""
	if int(currentMode) < len(modeDescriptions) {
		desc = "\n" + h.styles.Dim.Render(modeDescriptions[currentMode])
	}

	body := fmt.Sprintf(
		"Mode:\n%s%s%s\n\nPress %s to start typing.\nPress %s to open config.\nPress %s to quit.",
		modeDisplay,
		help,
		desc,
		h.styles.Marker.Render("enter"),
		h.styles.Dim.Render("c"),
		h.styles.Dim.Render("q"),
	)

	footer := []string{
		footerVersion,
		fmt.Sprintf("mode: %s", h.modeLabel()),
		"\u2190/\u2192: mode",
		"enter: start",
		"c: config",
		"q: quit",
	}
	footerH := h.styles.RenderFooterHeight(footer, h.width)
	bodyHeight := BodyHeight(h.height, footerH)
	return h.styles.Layout(h.width, h.height, "A O I", footer, body, bodyHeight)
}
