package ui

import (
	"fmt"
	"strings"

	"github.com/AlexandreSJ/aoi/internal/config"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var pickerColors []string

type pickerSection struct {
	name  string
	start int
	count int
	cols  int
}

var pickerSections []pickerSection

func init() {
	standardStart := len(pickerColors)
	for i := 0; i < 16; i++ {
		pickerColors = append(pickerColors, fmt.Sprintf("%d", i))
	}
	pickerSections = append(pickerSections, pickerSection{"Standard", standardStart, 16, 8})

	cubeStart := len(pickerColors)
	for r := 0; r < 6; r++ {
		for g := 0; g < 6; g++ {
			for b := 0; b < 6; b++ {
				pickerColors = append(pickerColors, fmt.Sprintf("%d", 16+r*36+g*6+b))
			}
		}
	}
	pickerSections = append(pickerSections, pickerSection{"Cube", cubeStart, 216, 18})

	grayStart := len(pickerColors)
	for i := 232; i <= 255; i++ {
		pickerColors = append(pickerColors, fmt.Sprintf("%d", i))
	}
	pickerSections = append(pickerSections, pickerSection{"Grayscale", grayStart, 24, 12})
}

type configModel struct {
	layout Layout
	cfg    *config.Config
	keys   []string
	sections []config.Section
	cursor   int
	scroll   int

	items []displayItem

	editing bool
	editKey string
	input   string
	err     string

	picker       bool
	pickerIdx    int
	pickerScroll int
}

func newConfigModel(cfg *config.Config, layout Layout) configModel {
	sections := cfg.Sections()
	var keys []string
	for _, s := range sections {
		keys = append(keys, s.Keys...)
	}
	c := configModel{
		layout:   layout,
		cfg:      cfg,
		sections: sections,
		keys:     keys,
	}
	c.items = c.buildDisplayItems()
	return c
}

func (c configModel) setSize(w, h int) configModel {
	c.layout = c.layout.SetSize(w, h)
	return c
}

type displayItem struct {
	kind string
	text string
	key  string
	hint string
}

func (c configModel) buildDisplayItems() []displayItem {
	var items []displayItem
	for si, s := range c.sections {
		if si > 0 {
			items = append(items, displayItem{kind: "spacing"})
		}
		items = append(items, displayItem{kind: "subtitle", text: s.Title})
		for _, k := range s.Keys {
			hint := ""
			if config.IsInlineHintKey(k) {
				hint = config.InlineHint(k)
			}
			items = append(items, displayItem{kind: "item", key: k, hint: hint})
		}
	}
	return items
}

func (c configModel) cursorToDisplayIdx() int {
	keyCount := 0
	for i, item := range c.items {
		if item.kind == "item" {
			if keyCount == c.cursor {
				return i
			}
			keyCount++
		}
	}
	return 0
}

func (c configModel) footerSegments() []string {
	if c.picker {
		selected := pickerColors[c.pickerIdx]
		shortKey := c.editKey[strings.LastIndex(c.editKey, ".")+1:]
		swatch := FooterSwatch(selected, c.cfg.Colors.Footer)
		return []string{
			footerVersion,
			fmt.Sprintf("Selected (%s): %s %3s", shortKey, swatch, selected),
			"\u2191\u2193\u2190\u2192: pick",
			"enter: confirm",
			"esc: cancel",
		}
	}
	if c.editing {
		return []string{footerVersion, "enter: confirm", "esc: cancel"}
	}
	return []string{footerVersion, "\u2191/\u2193: nav", "enter: edit", "c/esc: back"}
}

func (c configModel) Update(msg tea.Msg) (configModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if c.picker {
			return c.handlePicker(msg)
		}
		if c.editing {
			return c.handleEdit(msg)
		}
		return c.handleNav(msg)
	}
	return c, nil
}

func (c configModel) handleNav(msg tea.KeyMsg) (configModel, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if c.cursor > 0 {
			c.cursor--
			c = c.adjustScroll()
		}
	case "down", "j":
		if c.cursor < len(c.keys)-1 {
			c.cursor++
			c = c.adjustScroll()
		}
	case "enter":
		c.editKey = c.keys[c.cursor]
		if config.IsColorKey(c.editKey) {
			c.picker = true
			c.pickerIdx = 0
			c.pickerScroll = 0
			c.err = ""
		} else {
			c.editing = true
			c.input = c.cfg.Get(c.editKey)
			c.err = ""
		}
	}
	return c, nil
}

func (c configModel) handleEdit(msg tea.KeyMsg) (configModel, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		c.editing = false
		c.input = ""
		c.err = ""
	case tea.KeyEnter:
		input := strings.TrimSpace(c.input)
		c.cfg.Set(c.editKey, input)
		if err := config.Save(c.cfg); err != nil {
			c.err = fmt.Sprintf("save failed: %v", err)
			return c, nil
		}
		c.editing = false
		c.input = ""
		c.err = ""
	case tea.KeyBackspace:
		if len(c.input) > 0 {
			c.input = c.input[:len(c.input)-1]
		}
	case tea.KeyRunes:
		c.input += string(msg.Runes)
	}
	return c, nil
}

func (c configModel) handlePicker(msg tea.KeyMsg) (configModel, tea.Cmd) {
	total := len(pickerColors)

	switch msg.String() {
	case "esc":
		c.picker = false
		c.err = ""
	case "enter":
		color := pickerColors[c.pickerIdx]
		c.cfg.Set(c.editKey, color)
		if err := config.Save(c.cfg); err != nil {
			c.err = fmt.Sprintf("save failed: %v", err)
			return c, nil
		}
		c.picker = false
		c.err = ""
	case "up":
		c.pickerIdx = c.pickerUp()
		c = c.adjustPickerScroll()
	case "down":
		c.pickerIdx = c.pickerDown()
		c = c.adjustPickerScroll()
	case "left", "h":
		if c.pickerIdx > 0 {
			c.pickerIdx--
			c = c.adjustPickerScroll()
		}
	case "right", "l":
		if c.pickerIdx < total-1 {
			c.pickerIdx++
			c = c.adjustPickerScroll()
		}
	}
	return c, nil
}

func (c configModel) adjustPickerScroll() configModel {
	lines := c.buildPickerLines()
	cursorLine := -1
	for i, pl := range lines {
		if pl.idx == c.pickerIdx {
			cursorLine = i
			break
		}
	}
	if cursorLine < 0 {
		for i, pl := range lines {
			if pl.idx != -1 && pl.idx+pl.cellCount > c.pickerIdx {
				cursorLine = i
				break
			}
		}
	}
	if cursorLine < 0 {
		return c
	}

	avail := c.layout.BodyHeight(c.footerSegments())

	targetLine := cursorLine
	if cursorLine > 0 && lines[cursorLine-1].idx == -1 {
		targetLine = cursorLine - 1
	}

	if targetLine < c.pickerScroll {
		c.pickerScroll = targetLine
	}
	if cursorLine >= c.pickerScroll+avail {
		c.pickerScroll = cursorLine - avail + 1
	}

	if c.pickerScroll < 0 {
		c.pickerScroll = 0
	}
	return c
}

type pickerLine struct {
	text      string
	idx       int
	cellCount int
}

func (c configModel) buildPickerLines() []pickerLine {
	var lines []pickerLine
	for _, sec := range pickerSections {
		lines = append(lines, pickerLine{text: sec.name + ":", idx: -1, cellCount: 0})
		var row strings.Builder
		row.WriteString("  ")
		rowCells := 0
		for i := 0; i < sec.count; i++ {
			if i > 0 && i%sec.cols == 0 {
				lines = append(lines, pickerLine{text: row.String(), idx: sec.start + i - sec.cols, cellCount: sec.cols})
				row.Reset()
				row.WriteString("  ")
				rowCells = 0
			}
			idx := sec.start + i
			row.WriteString(c.renderPickerCell(idx))
			rowCells++
		}
		lines = append(lines, pickerLine{text: row.String(), idx: sec.start + ((sec.count-1)/sec.cols)*sec.cols, cellCount: rowCells})
		lines = append(lines, pickerLine{text: "", idx: -1, cellCount: 0})
	}
	return lines
}

func (c configModel) pickerSectionIdx() int {
	for i, sec := range pickerSections {
		if c.pickerIdx >= sec.start && c.pickerIdx < sec.start+sec.count {
			return i
		}
	}
	return 0
}

func (c configModel) pickerColInSection(secIdx int) int {
	sec := pickerSections[secIdx]
	return (c.pickerIdx - sec.start) % sec.cols
}

func (c configModel) pickerUp() int {
	secIdx := c.pickerSectionIdx()
	sec := pickerSections[secIdx]
	col := c.pickerColInSection(secIdx)

	if c.pickerIdx-sec.start >= sec.cols {
		return c.pickerIdx - sec.cols
	}
	if secIdx == 0 {
		return c.pickerIdx
	}
	prevSec := pickerSections[secIdx-1]
	targetCol := min(col, prevSec.cols-1)
	lastRowStart := prevSec.start + ((prevSec.count-1)/prevSec.cols)*prevSec.cols
	idx := lastRowStart + targetCol
	if idx >= prevSec.start+prevSec.count {
		idx = prevSec.start + prevSec.count - 1
	}
	return idx
}

func (c configModel) pickerDown() int {
	secIdx := c.pickerSectionIdx()
	sec := pickerSections[secIdx]
	col := c.pickerColInSection(secIdx)
	relIdx := c.pickerIdx - sec.start

	if relIdx+sec.cols < sec.count {
		return c.pickerIdx + sec.cols
	}
	if secIdx == len(pickerSections)-1 {
		return c.pickerIdx
	}
	nextSec := pickerSections[secIdx+1]
	targetCol := min(col, nextSec.cols-1)
	idx := nextSec.start + targetCol
	if idx >= nextSec.start+nextSec.count {
		idx = nextSec.start + nextSec.count - 1
	}
	return idx
}

func (c configModel) adjustScroll() configModel {
	dIdx := c.cursorToDisplayIdx()
	avail := c.layout.BodyHeight(c.footerSegments())

	linesAfterCursor := 1
	if c.items[dIdx].hint != "" {
		linesAfterCursor++
	}

	targetLine := dIdx
	if dIdx > 0 && c.items[dIdx-1].kind == "subtitle" {
		targetLine = dIdx - 1
	}

	if targetLine < c.scroll {
		c.scroll = targetLine
	}
	if dIdx+linesAfterCursor-1 >= c.scroll+avail {
		c.scroll = dIdx + linesAfterCursor - avail
	}

	if c.scroll < 0 {
		c.scroll = 0
	}
	return c
}

func (c configModel) View() string {
	segments := c.footerSegments()
	bodyHeight := c.layout.BodyHeight(segments)

	var bodyContent string
	if c.picker {
		bodyContent = c.renderPicker(bodyHeight)
	} else {
		bodyContent = c.renderList(c.layout.Width-4, bodyHeight)
	}

	return c.layout.Render("Config", segments, bodyContent)
}

func (c configModel) renderList(availWidth, maxLines int) string {
	dIdx := c.cursorToDisplayIdx()

	var b strings.Builder
	lineCount := 0
	for i := 0; i < len(c.items) && lineCount < maxLines; i++ {
		if i < c.scroll {
			continue
		}
		item := c.items[i]
		switch item.kind {
		case "spacing":
			b.WriteString("\n")
			lineCount++
		case "subtitle":
			b.WriteString(c.layout.Styles.Subtitle.Render(fmt.Sprintf("  %s", item.text)) + "\n")
			lineCount++
		case "item":
			key := item.key
			value := c.cfg.Get(key)
			shortKey := key[strings.LastIndex(key, ".")+1:]

			marker := "  "
			if i == dIdx {
				marker = c.layout.Styles.Marker.Render(" >")
			}

			var line string
			if config.IsColorKey(key) {
				swatch := lipgloss.NewStyle().
					Foreground(lipgloss.Color(value)).
					Render("\u2588\u2588")
				line = fmt.Sprintf("%s %s %s: %s", marker, swatch, shortKey, value)
			} else {
				line = fmt.Sprintf("%s   %s: %s", marker, shortKey, value)
				if item.hint != "" && !(c.editing && c.editKey == key) {
					line += c.layout.Styles.Dim.Render(fmt.Sprintf("  %s", item.hint))
				}
			}

			if c.editing && c.editKey == key {
				line = fmt.Sprintf("%s   %s: %s\u2588", marker, shortKey, c.input)
			}

			if availWidth > 0 && lipgloss.Width(line) > availWidth {
				line = truncateLine(line, availWidth)
			}
			b.WriteString(line + "\n")
			lineCount++
		}
	}

	if c.err != "" && lineCount < maxLines {
		b.WriteString("\n" + c.layout.Styles.Error.Render(c.err))
	}

	return b.String()
}

func truncateLine(line string, maxWidth int) string {
	width := 0
	inEscape := false
	var b strings.Builder
	for _, r := range line {
		if r == '\x1b' {
			inEscape = true
			b.WriteRune(r)
			continue
		}
		if inEscape {
			b.WriteRune(r)
			if r == 'm' {
				inEscape = false
			}
			continue
		}
		if width >= maxWidth {
			break
		}
		b.WriteRune(r)
		width++
	}
	return b.String()
}

func (c configModel) renderPicker(avail int) string {
	allLines := c.buildPickerLines()
	if avail < 1 {
		avail = 1
	}

	scroll := c.pickerScroll
	if scroll < 0 {
		scroll = 0
	}
	maxScroll := len(allLines) - avail
	if maxScroll < 0 {
		maxScroll = 0
	}
	if scroll > maxScroll {
		scroll = maxScroll
	}

	end := scroll + avail
	if end > len(allLines) {
		end = len(allLines)
	}
	visible := allLines[scroll:end]

	var b strings.Builder

	if c.err != "" {
		b.WriteString(c.layout.Styles.Error.Render(c.err) + "\n")
	}

	for i, pl := range visible {
		if pl.text == "" {
			if i < len(visible)-1 && visible[i+1].text != "" {
				b.WriteString("\n")
			}
			continue
		}
		if pl.idx == -1 {
			b.WriteString(c.layout.Styles.Subtitle.Render("  "+pl.text) + "\n")
			continue
		}
		b.WriteString(pl.text + "\n")
	}
	return b.String()
}

func (c configModel) renderPickerCell(idx int) string {
	if idx >= len(pickerColors) {
		return "  "
	}
	color := pickerColors[idx]
	if idx == c.pickerIdx {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(color)).
			Background(lipgloss.Color("#FFFFFF")).
			Render("\u2593\u2593")
	}
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(color)).
		Render("\u2588\u2588")
}
