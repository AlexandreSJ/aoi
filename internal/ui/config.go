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
	pickerSections = append(pickerSections, pickerSection{standardStart, 16, 8})

	cubeStart := len(pickerColors)
	for r := 0; r < 6; r++ {
		for g := 0; g < 6; g++ {
			for b := 0; b < 6; b++ {
				pickerColors = append(pickerColors, fmt.Sprintf("%d", 16+r*36+g*6+b))
			}
		}
	}
	pickerSections = append(pickerSections, pickerSection{cubeStart, 216, 18})

	grayStart := len(pickerColors)
	for i := 232; i <= 255; i++ {
		pickerColors = append(pickerColors, fmt.Sprintf("%d", i))
	}
	pickerSections = append(pickerSections, pickerSection{grayStart, 24, 12})
}

type configModel struct {
	layout   Layout
	cfg      *config.Config
	keys     []string
	sections []config.Section
	cursor   int
	scroll   int

	items []displayItem

	editing    bool
	editKey    string
	input      string
	editCursor int
	err        string

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
	for _, s := range c.sections {
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
			fmt.Sprintf("Selected (%s): %s %3s", shortKey, swatch, selected),
			"\u2191\u2193\u2190\u2192: pick",
			"enter: confirm",
			"esc: cancel",
		}
	}
	if c.editing {
		return []string{"enter: confirm", "esc: cancel"}
	}
	return []string{"\u2191/\u2193: nav", "enter: edit", "c/esc: back"}
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
			c.editCursor = len([]rune(c.input))
			c.err = ""
		}
	}
	return c, nil
}

func (c configModel) handleEdit(msg tea.KeyMsg) (configModel, tea.Cmd) {
	runes := []rune(c.input)
	s := msg.String()

	switch {
	case s == "esc":
		c.editing = false
		c.input = ""
		c.editCursor = 0
		c.err = ""
	case s == "enter":
		input := strings.TrimSpace(c.input)
		c.cfg.Set(c.editKey, input)
		if err := config.Save(c.cfg); err != nil {
			c.err = fmt.Sprintf("save failed: %v", err)
			return c, nil
		}
		c.editing = false
		c.input = ""
		c.editCursor = 0
		c.err = ""
	case s == "backspace":
		if c.editCursor > 0 {
			c.input = string(runes[:c.editCursor-1]) + string(runes[c.editCursor:])
			c.editCursor--
		}
	case s == "delete":
		if c.editCursor < len(runes) {
			c.input = string(runes[:c.editCursor]) + string(runes[c.editCursor+1:])
		}
	case s == "left":
		if c.editCursor > 0 {
			c.editCursor--
		}
	case s == "right":
		if c.editCursor < len(runes) {
			c.editCursor++
		}
	case s == "home", s == "ctrl+a":
		c.editCursor = 0
	case s == "end", s == "ctrl+e":
		c.editCursor = len([]rune(c.input))
	case s == "ctrl+u":
		c.input = string(runes[c.editCursor:])
		c.editCursor = 0
	case s == "ctrl+k":
		c.input = string(runes[:c.editCursor])
	case s == "ctrl+w", s == "ctrl+backspace":
		pos := c.editCursor
		for pos > 0 && isEditWordSep(runes[pos-1]) {
			pos--
		}
		for pos > 0 && !isEditWordSep(runes[pos-1]) {
			pos--
		}
		c.input = string(runes[:pos]) + string(runes[c.editCursor:])
		c.editCursor = pos
	case s == "ctrl+delete":
		pos := c.editCursor
		for pos < len(runes) && !isEditWordSep(runes[pos]) {
			pos++
		}
		for pos < len(runes) && isEditWordSep(runes[pos]) {
			pos++
		}
		c.input = string(runes[:c.editCursor]) + string(runes[pos:])
	default:
		if msg.Type == tea.KeyRunes {
			inserted := string(msg.Runes)
			c.input = string(runes[:c.editCursor]) + inserted + string(runes[c.editCursor:])
			c.editCursor += len([]rune(inserted))
		}
	}
	return c, nil
}

func isEditWordSep(r rune) bool {
	return r == '/' || r == ' '
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
			if pl.idx >= 0 && pl.idx+pl.cellCount > c.pickerIdx {
				cursorLine = i
				break
			}
		}
	}
	if cursorLine < 0 {
		return c
	}

	avail := c.layout.BodyHeight(c.footerSegments())

	if cursorLine < c.pickerScroll {
		c.pickerScroll = cursorLine
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

// sectionCenteringOffset returns how many picker columns of left-padding
// centering adds to the given section. Each cell is 2 chars wide.
func (c configModel) sectionCenteringOffset(secIdx int) int {
	innerWidth := max(1, c.layout.Width-8) // matches CenterBody calculation
	sec := pickerSections[secIdx]
	secWidth := sec.cols * 2
	leftPad := (innerWidth - secWidth) / 2
	if leftPad < 0 {
		leftPad = 0
	}
	return leftPad / 2
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

	// Adjust column for visual centering difference
	visualPos := c.sectionCenteringOffset(secIdx) + col
	targetCol := visualPos - c.sectionCenteringOffset(secIdx-1)
	if targetCol < 0 {
		targetCol = 0
	}
	if targetCol >= prevSec.cols {
		targetCol = prevSec.cols - 1
	}

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

	// Adjust column for visual centering difference
	visualPos := c.sectionCenteringOffset(secIdx) + col
	targetCol := visualPos - c.sectionCenteringOffset(secIdx+1)
	if targetCol < 0 {
		targetCol = 0
	}
	if targetCol >= nextSec.cols {
		targetCol = nextSec.cols - 1
	}

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
		bodyContent = c.layout.CenterBody(bodyContent, bodyHeight)
	} else {
		bodyContent = c.renderList(c.layout.Width-4, bodyHeight)
		bodyContent = c.layout.CenterBody(bodyContent, bodyHeight)
	}

	return c.layout.Render("C F G", segments, bodyContent)
}

func editViewport(runes []rune, cursor, fieldWidth int) (start, end int, leftEllip, rightEllip bool) {
	n := len(runes)
	effLen := n
	if cursor >= n {
		effLen = n + 1
	}
	if effLen <= fieldWidth {
		return 0, n, false, false
	}

	useEllip := fieldWidth >= 5
	reserve := 0
	if useEllip {
		reserve = 2
	}
	if cursor >= n {
		reserve++
	}

	visWidth := fieldWidth - reserve
	if visWidth < 1 {
		visWidth = 1
	}

	target := cursor
	if target >= n {
		target = n - 1
	}
	if target < 0 {
		target = 0
	}

	start = target - visWidth/2
	if start < 0 {
		start = 0
	}
	end = start + visWidth
	if end > n {
		end = n
		start = end - visWidth
		if start < 0 {
			start = 0
		}
	}

	return start, end, useEllip && start > 0, useEllip && end < n
}

func renderEditValue(runes []rune, cursor, fieldWidth int, primaryColor, textColor string) string {
	n := len(runes)
	start, end, leftEllip, rightEllip := editViewport(runes, cursor, fieldWidth)

	cursorStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(primaryColor)).
		Bold(true)
	normalStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(textColor))

	var b strings.Builder

	if leftEllip {
		b.WriteString(normalStyle.Render("..."))
	}

	for i := start; i < end; i++ {
		ch := string(runes[i])
		if i == cursor {
			b.WriteString(cursorStyle.Render(ch))
		} else {
			b.WriteString(normalStyle.Render(ch))
		}
	}

	if cursor == n && n >= start {
		b.WriteString(cursorStyle.Render(""))
	}

	if rightEllip {
		b.WriteString(normalStyle.Render(""))
	}

	return b.String()
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
					Render("\u2588\u2588\u2588\u2588")
				labelWidth := lipgloss.Width(fmt.Sprintf("%s:%s %s", shortKey, value, swatch))
				padding := max(0, 20-labelWidth)
				line = fmt.Sprintf("%s %s: %s %s %s", marker, shortKey, strings.Repeat(" ", padding), value, swatch)
			} else {
				maxValW := 25

				if c.editing && c.editKey == key {
					fieldWidth := 25
					primary := c.cfg.Colors.Primary
					if primary == "" {
						primary = fallbackPrimary
					}
					text := c.cfg.Colors.Text
					if text == "" {
						text = fallbackText
					}
					renderedVal := renderEditValue([]rune(c.input), c.editCursor, fieldWidth, primary, text)

					labelWidth := lipgloss.Width(fmt.Sprintf("%s: [%s]", shortKey, renderedVal))
					padding := max(0, 35-labelWidth)
					prefix := fmt.Sprintf("%s %s: %s", marker, shortKey, strings.Repeat(" ", padding))
					line = fmt.Sprintf("%s[%s]", prefix, renderedVal)
				} else {
					displayVal := value
					if lipgloss.Width(displayVal) > maxValW {
						displayVal = truncatePrefix(displayVal, 25-1)
					}

					labelWidth := lipgloss.Width(fmt.Sprintf("%s: [%s]", shortKey, displayVal))
					padding := max(0, 35-labelWidth)
					prefix := fmt.Sprintf("%s %s: %s", marker, shortKey, strings.Repeat(" ", padding))
					line = fmt.Sprintf("%s[%s]", prefix, displayVal)
				}
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

func truncatePrefix(s string, maxW int) string {
	if lipgloss.Width(s) <= maxW {
		return s
	}
	// Keep "..." + as many trailing chars as fit
	if maxW <= 3 {
		return "..."
	}
	target := maxW - 3
	width := 0
	var b strings.Builder
	// Walk backwards to keep the tail
	runes := []rune(s)
	start := len(runes)
	for i := len(runes) - 1; i >= 0; i-- {
		w := lipgloss.Width(string(runes[i]))
		if width+w > target {
			break
		}
		width += w
		start = i
	}
	b.WriteString("...")
	for _, r := range runes[start:] {
		b.WriteRune(r)
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

	errLine := 0
	if c.err != "" {
		b.WriteString(c.layout.Styles.Error.Render(c.err) + "\n")
		errLine = 1
	}

	// Recalculate visible slice accounting for error line
	pickerAvail := avail - errLine
	if pickerAvail < 1 {
		pickerAvail = 1
	}
	end = scroll + pickerAvail
	if end > len(allLines) {
		end = len(allLines)
	}
	visible = allLines[scroll:end]

	for _, pl := range visible {
		if pl.text == "" {
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
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color(color)).
			Bold(true).
			Render("[]")
	}
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(color)).
		Render("\u2588\u2588")
}
