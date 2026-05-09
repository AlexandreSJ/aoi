package ui

import (
	"fmt"
	"math/rand/v2"
	"strings"
	"time"

	"github.com/AlexandreSJ/aoi/internal/config"
	"github.com/AlexandreSJ/aoi/internal/quotes"
	"github.com/AlexandreSJ/aoi/internal/words"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type charState int

const (
	charPending charState = iota
	charCorrect
	charError
)

type typedChar struct {
	char  string
	state charState
}

type textLine struct {
	start int
	end   int
}

type tickMsg time.Time

type flashMsg struct{}

func flashCmd() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(time.Time) tea.Msg { return flashMsg{} })
}

type typingModel struct {
	layout Layout
	cfg    *config.Config

	mode            gameMode
	wordListName    string
	wordList        *words.WordList
	quoteList       *quotes.QuoteList
	timedSeconds    int
	wordCountTarget int

	chars        []typedChar
	cursor       int
	scrollOffset int
	finished     bool
	err          string

	lastWord     string
	lastKey      string
	lastKeyError bool
	lastKeyFlash bool
	errorCount   int
	trimmedOk    int
	trimmedErr   int

	timeRemaining int
	timerRunning  bool

	charTimestamps []time.Time
	startTime      time.Time
	showStats      bool
	charLimit      int
}

func newTypingModel(cfg *config.Config, layout Layout) typingModel {
	return typingModel{
		cfg:          cfg,
		layout:       layout,
		wordListName: "en",
		showStats:    true,
		charLimit:    100,
	}
}

func (t typingModel) setSize(w, h int) typingModel {
	t.layout = t.layout.SetSize(w, h)
	return t
}

func (t typingModel) Update(msg tea.Msg) (typingModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		if !t.timerRunning {
			return t, nil
		}
		t.timeRemaining--
		if t.timeRemaining <= 0 {
			t.timeRemaining = 0
			t.finished = true
			t.timerRunning = false
			return t, nil
		}
		return t, tea.Tick(time.Second, func(t time.Time) tea.Msg {
			return tickMsg(t)
		})
	case flashMsg:
		t.lastKeyFlash = false
		return t, nil
	case tea.KeyMsg:
		if t.finished {
			if msg.String() == "enter" {
				return t.restart()
			}
			return t, nil
		}

		if len(t.chars) == 0 {
			return t, nil
		}

		if msg.String() == "tab" {
			t.showStats = !t.showStats
			return t, nil
		}

		t.lastKeyFlash = true

		switch msg.Type {
		case tea.KeyEnter:
			t.lastKey = "enter"
			t.lastKeyError = false
			return t.restart()
		case tea.KeyBackspace:
			t.lastKey = "backspace"
			t.lastKeyError = false
			return t.handleBackspace(), flashCmd()
		case tea.KeyDown:
			t.lastKey = "\u2193"
			t.lastKeyError = false
			return t.handleSkipRow()
		case tea.KeyUp:
			t.lastKey = "\u2191"
			t.lastKeyError = false
			return t, flashCmd()
		case tea.KeyLeft:
			t.lastKey = "\u2190"
			t.lastKeyError = false
			t.charLimit = max(25, t.charLimit-25)
			return t, flashCmd()
		case tea.KeyRight:
			t.lastKey = "\u2192"
			t.lastKeyError = false
			t.charLimit = min(500, t.charLimit+25)
			return t, flashCmd()
		case tea.KeyRunes:
			if msg.String() == "j" {
				return t.handleCharOrSkip(msg.String())
			}
			return t.handleChar(msg.String())
		case tea.KeySpace:
			return t.handleChar(" ")
		}
	}

	return t, nil
}

func (t typingModel) handleCharOrSkip(input string) (typingModel, tea.Cmd) {
	lines := t.computeLines(t.lineWidth())
	if len(lines) == 0 {
		return t.handleChar(input)
	}

	cursorLine := t.cursorLineIdx(lines)
	if t.cursor == lines[cursorLine].start {
		return t.handleSkipRow()
	}

	return t.handleChar(input)
}

func (t typingModel) restart() (typingModel, tea.Cmd) {
	savedMode := t.mode
	savedName := t.wordListName
	savedTimed := t.timedSeconds
	savedCount := t.wordCountTarget
	savedW := t.layout.Width
	savedH := t.layout.Height
	savedLastKey := t.lastKey
	savedLastKeyError := t.lastKeyError
	savedFlash := t.lastKeyFlash
	savedShowStats := t.showStats
	savedCharLimit := t.charLimit

	t = newTypingModel(t.cfg, t.layout)
	t.mode = savedMode
	t.wordListName = savedName
	t.timedSeconds = savedTimed
	t.wordCountTarget = savedCount
	t.layout.Width = savedW
	t.layout.Height = savedH
	t.timeRemaining = savedTimed
	t.lastKey = savedLastKey
	t.lastKeyError = savedLastKeyError
	t.lastKeyFlash = savedFlash
	t.showStats = savedShowStats
	t.charLimit = savedCharLimit

	if savedMode == modeQuote {
		t = t.loadQuote()
	} else {
		t = t.loadWords()
	}

	return t, nil
}

func (t typingModel) resolveWordsDir() string {
	dir, _ := config.ResolveDir(t.cfg.System.WordsDir, "~/.config/aoi/words")
	return dir
}

func (t typingModel) resolveQuotesDir() string {
	dir, _ := config.ResolveDir(t.cfg.System.QuotesDir, "~/.config/aoi/quotes")
	return dir
}

func (t typingModel) loadWords() typingModel {
	wl, err := words.LoadList(t.wordListName, t.resolveWordsDir())
	if err != nil {
		t.err = fmt.Sprintf("cannot load %q: %v", t.wordListName, err)
		return t
	}
	t.wordList = wl

	switch t.mode {
	case modeCount:
		ws := t.sampleWords(t.wordCountTarget)
		t.setChars(strings.Join(ws, " "))
	default:
		t.initInfiniteText()
	}
	return t
}

func (t typingModel) loadQuote() typingModel {
	ql, err := quotes.LoadList(t.wordListName, t.resolveQuotesDir())
	if err != nil {
		t.err = fmt.Sprintf("cannot load quotes %q: %v", t.wordListName, err)
		return t
	}
	t.quoteList = ql

	raw := ql.Random()
	t.setChars(raw)
	return t
}

func (t *typingModel) setChars(raw string) {
	t.chars = make([]typedChar, 0, len(raw))
	for _, c := range raw {
		t.chars = append(t.chars, typedChar{char: string(c)})
	}
	t.cursor = 0
	t.scrollOffset = 0
}

func (t *typingModel) sampleWords(n int) []string {
	all := t.wordList.Words
	result := make([]string, 0, n)
	for len(result) < n {
		w := all[rand.IntN(len(all))]
		if len(result) == 0 || w != result[len(result)-1] {
			result = append(result, w)
		}
	}
	t.lastWord = result[len(result)-1]
	return result
}

func (t *typingModel) wordsPerRow() int {
	availWidth := max(1, t.layout.Width-8)
	return max(1, availWidth/5)
}

func (t *typingModel) lineWidth() int {
	availWidth := max(1, t.layout.Width-8)
	return min(t.charLimit, availWidth)
}

func (t *typingModel) initInfiniteText() {
	ws := t.sampleWords(t.wordsPerRow() * 4)
	t.setChars(strings.Join(ws, " "))
}

func (t *typingModel) ensureBufferRows() {
	lines := t.computeLines(t.lineWidth())
	if len(lines) == 0 {
		return
	}

	cursorLine := t.cursorLineIdx(lines)
	linesBelowCursor := len(lines) - cursorLine - 1

	if linesBelowCursor < 2 {
		need := t.wordsPerRow() * 3
		ws := t.sampleWords(need)
		additional := " " + strings.Join(ws, " ")
		for _, c := range additional {
			t.chars = append(t.chars, typedChar{char: string(c)})
		}
	}
}

func (t typingModel) handleBackspace() typingModel {
	// If cursor is on an error char (missed, not yet advanced), clear it in place
	if t.cursor < len(t.chars) && t.chars[t.cursor].state == charError {
		t.chars[t.cursor].state = charPending
		t.errorCount--
		return t.adjustScroll()
	}
	if t.cursor <= 0 {
		return t
	}
	t.cursor--
	if t.chars[t.cursor].state == charError {
		t.errorCount--
	}
	t.chars[t.cursor].state = charPending
	return t.adjustScroll()
}

func (t typingModel) handleSkipRow() (typingModel, tea.Cmd) {
	lines := t.computeLines(t.lineWidth())
	if len(lines) == 0 {
		return t, nil
	}

	cursorLine := t.cursorLineIdx(lines)
	if cursorLine+1 < len(lines) {
		t.cursor = lines[cursorLine+1].start
	} else {
		t.cursor = lines[cursorLine].end
	}

	if t.mode == modeInfinite || t.mode == modeTimed {
		t.ensureBufferRows()
	}

	t = t.adjustScroll()
	t.trimCompleted()
	return t, flashCmd()
}

func (t typingModel) handleChar(input string) (typingModel, tea.Cmd) {
	if t.cursor >= len(t.chars) {
		return t, nil
	}

	t.lastKey = input

	var cmd tea.Cmd
	if t.mode == modeTimed && !t.timerRunning && !t.finished {
		t.timerRunning = true
		t.startTime = time.Now()
		cmd = tea.Tick(time.Second, func(t time.Time) tea.Msg {
			return tickMsg(t)
		})
	}

	if t.startTime.IsZero() {
		t.startTime = time.Now()
	}

	expected := t.chars[t.cursor].char

	if input == expected {
		if t.chars[t.cursor].state != charError {
			t.chars[t.cursor].state = charCorrect
			t.charTimestamps = append(t.charTimestamps, time.Now())
		}
		t.lastKeyError = false
		t.cursor++
		if t.mode == modeInfinite || t.mode == modeTimed {
			t.ensureBufferRows()
		}
		if t.mode == modeInfinite && t.cursor >= len(t.chars) {
			t.finished = true
		}
		if t.mode == modeCount && t.cursor >= len(t.chars) {
			t.finished = true
		}
		if t.mode == modeQuote && t.cursor >= len(t.chars) {
			t.finished = true
		}
	} else {
		if t.chars[t.cursor].state != charError {
			t.errorCount++
		}
		t.chars[t.cursor].state = charError
		t.lastKeyError = true
	}

	t = t.adjustScroll()
	t.trimCompleted()
	if cmd != nil {
		return t, tea.Batch(cmd, flashCmd())
	}
	return t, flashCmd()
}

func (t *typingModel) trimCompleted() {
	lines := t.computeLines(t.lineWidth())
	if len(lines) < 3 {
		return
	}

	cursorLine := t.cursorLineIdx(lines)
	if cursorLine < 2 {
		return
	}

	trimLineIdx := cursorLine - 1
	trimCount := lines[trimLineIdx].start
	if trimCount <= 0 {
		return
	}

	for i := 0; i < trimCount; i++ {
		switch t.chars[i].state {
		case charCorrect:
			t.trimmedOk++
		case charError:
			t.trimmedErr++
		}
	}

	t.chars = t.chars[trimCount:]
	t.cursor -= trimCount
	t.scrollOffset = 0
}

func (t typingModel) computeLines(availWidth int) []textLine {
	if availWidth < 1 || len(t.chars) == 0 {
		return nil
	}

	type wordBound struct {
		start int
		end   int
	}

	var boundaries []wordBound
	i := 0
	for i < len(t.chars) {
		start := i
		for i < len(t.chars) && t.chars[i].char != " " {
			i++
		}
		end := i
		if i < len(t.chars) && t.chars[i].char == " " {
			end = i + 1
			i++
		}
		boundaries = append(boundaries, wordBound{start: start, end: end})
	}

	var lines []textLine
	lineStart := 0
	lineWidth := 0

	for _, w := range boundaries {
		wordWidth := w.end - w.start
		if lineWidth > 0 && lineWidth+wordWidth > availWidth {
			lines = append(lines, textLine{start: lineStart, end: w.start})
			lineStart = w.start
			lineWidth = 0
		}
		lineWidth += wordWidth
	}

	if lineStart < len(t.chars) {
		lines = append(lines, textLine{start: lineStart, end: len(t.chars)})
	}

	return lines
}

func (t typingModel) cursorLineIdx(lines []textLine) int {
	for i, line := range lines {
		if t.cursor >= line.start && t.cursor < line.end {
			return i
		}
	}
	if len(lines) > 0 {
		return len(lines) - 1
	}
	return 0
}

func (t typingModel) adjustScroll() typingModel {
	lines := t.computeLines(t.lineWidth())
	if len(lines) == 0 {
		return t
	}

	cursorLine := t.cursorLineIdx(lines)

	t.scrollOffset = cursorLine
	if cursorLine > 0 {
		t.scrollOffset = cursorLine - 1
	}

	if t.scrollOffset < 0 {
		t.scrollOffset = 0
	}

	return t
}

func (t typingModel) correctCount() int {
	n := t.trimmedOk
	for _, c := range t.chars {
		if c.state == charCorrect {
			n++
		}
	}
	return n
}

func (t typingModel) wpm() float64 {
	if t.startTime.IsZero() {
		return 0
	}

	var elapsed time.Duration
	if t.mode == modeTimed {
		elapsed = time.Duration(t.timedSeconds-t.timeRemaining) * time.Second
	} else {
		elapsed = time.Since(t.startTime)
		if t.finished {
			elapsed = time.Duration(t.correctCount()+t.errorCount) * time.Second / time.Duration(max(1, t.correctCount()+t.errorCount))
		}
	}

	if elapsed < time.Second {
		return 0
	}

	seconds := elapsed.Seconds()
	if seconds <= 0 {
		return 0
	}

	var charsInWindow int
	if t.mode == modeTimed {
		charsInWindow = t.correctCount()
	} else {
		cutoff := time.Now().Add(-5 * time.Second)
		for _, ts := range t.charTimestamps {
			if ts.After(cutoff) {
				charsInWindow++
			}
		}
		if charsInWindow == 0 {
			return 0
		}
		seconds = 5.0
	}

	wpm := (float64(charsInWindow) / 5.0) * (60.0 / seconds)
	if wpm >= 10000 {
		return 9999.99
	}
	return wpm
}

func (t typingModel) precision() float64 {
	ok := t.correctCount()
	total := ok + t.trimmedErr + t.errorCount
	if total == 0 {
		return 0
	}
	p := float64(ok) / float64(total) * 100
	if p >= 10000 {
		return 9999.99
	}
	return p
}

func (t typingModel) footerSegments() []string {
	status := "ready"
	if t.cursor > 0 {
		status = "typing"
	}
	if t.finished {
		status = "done"
	}

	modeLabel := t.mode.String()
	switch t.mode {
	case modeTimed:
		remaining := t.timeRemaining
		if !t.timerRunning && !t.finished {
			remaining = t.timedSeconds
		}
		modeLabel = fmt.Sprintf("Timed %ds", remaining)
	case modeCount:
		modeLabel = fmt.Sprintf("Count %d", t.wordCountTarget)
	}

	ok := t.correctCount()
	errCount := t.trimmedErr + t.errorCount
	okStr := fmt.Sprintf("%d", ok)
	errStr := fmt.Sprintf("%d", errCount)
	if ok >= 10000 {
		okStr = "9999+"
	}
	if errCount >= 10000 {
		errStr = "9999+"
	}

	tabLabel := "tab: hide"
	if !t.showStats {
		tabLabel = "tab: show"
	}
	return []string{
		footerVersion,
		fmt.Sprintf("%s | %s", modeLabel, t.wordListName),
		fmt.Sprintf("%s ok / %s err", okStr, errStr),
		status,
		tabLabel,
		"enter: restart",
		"esc: back",
	}
}

func (t typingModel) View() string {
	if t.layout.Width == 0 {
		return ""
	}

	if t.err != "" {
		footer := []string{footerVersion, "esc: back"}
		body := t.layout.Styles.Error.Render(t.err)
		return t.layout.Render("A O I", footer, body)
	}

	footerSegs := t.footerSegments()
	bodyHeight := t.layout.BodyHeight(footerSegs)
	bodyContent := t.renderBody(bodyHeight)
	return t.layout.Render("A O I", footerSegs, bodyContent)
}

func (t typingModel) renderBody(bodyHeight int) string {
	text := t.renderText()
	textLines := strings.Count(text, "\n")
	if text == "" {
		textLines = 0
	}

	successColor := t.cfg.Colors.Success
	if successColor == "" {
		successColor = fallbackSuccess
	}
	errorColor := t.cfg.Colors.Error
	if errorColor == "" {
		errorColor = fallbackError
	}
	primaryColor := t.cfg.Colors.Primary
	if primaryColor == "" {
		primaryColor = fallbackPrimary
	}

	display := t.formatLastKey()
	var lastKeyLine string
	if t.lastKeyFlash {
		if t.lastKeyError {
			lastKeyLine = lipgloss.NewStyle().
				Foreground(lipgloss.Color(errorColor)).
				Bold(true).
				Render(display)
		} else if t.lastKey != "" {
			lastKeyLine = lipgloss.NewStyle().
				Foreground(lipgloss.Color(successColor)).
				Bold(true).
				Render(display)
		} else {
			lastKeyLine = lipgloss.NewStyle().Render(" ")
		}
	} else if t.lastKey != "" {
		lastKeyLine = lipgloss.NewStyle().
			Foreground(lipgloss.Color(primaryColor)).
			Render(display)
	} else {
		lastKeyLine = lipgloss.NewStyle().Render(" ")
	}

	innerWidth := max(1, t.layout.Width-8)
	centered := lipgloss.NewStyle().
		Width(innerWidth).
		Align(lipgloss.Center).
		Render(lastKeyLine)

	bottomSection := 3
	if t.showStats {
		bottomSection = 4
	}
	totalContentLines := textLines + bottomSection

	padAbove := 0
	if bodyHeight > totalContentLines {
		padAbove = (bodyHeight - totalContentLines) / 2
	}

	var b strings.Builder
	for i := 0; i < padAbove; i++ {
		b.WriteString("\n")
	}

	b.WriteString(text)

	b.WriteString("\n\n")
	b.WriteString(centered)

	if t.showStats {
		wpmVal := t.wpm()
		precVal := t.precision()
		statsLine := fmt.Sprintf("%.2f wpm | %.2f%% prec | %d chars", wpmVal, precVal, t.charLimit)
		centeredStats := lipgloss.NewStyle().
			Width(innerWidth).
			Align(lipgloss.Center).
			Foreground(lipgloss.Color(primaryColor)).
			Render(statsLine)
		b.WriteString("\n")
		b.WriteString(centeredStats)
	}

	return b.String()
}

func (t typingModel) formatLastKey() string {
	if t.lastKey == "" {
		return " "
	}
	switch t.lastKey {
	case " ":
		return "space"
	default:
		return t.lastKey
	}
}

func (t typingModel) renderText() string {
	if len(t.chars) == 0 {
		return ""
	}

	lines := t.computeLines(t.lineWidth())
	if len(lines) == 0 {
		return ""
	}

	primary := t.cfg.Colors.Primary
	if primary == "" {
		primary = fallbackPrimary
	}

	successColor := t.cfg.Colors.Success
	if successColor == "" {
		successColor = fallbackSuccess
	}

	errorColor := t.cfg.Colors.Error
	if errorColor == "" {
		errorColor = fallbackError
	}

	const textWindow = 3
	startLine := t.scrollOffset
	if startLine < 0 {
		startLine = 0
	}

	endLine := startLine + textWindow
	if endLine > len(lines) {
		endLine = len(lines)
	}

	cursorLine := t.cursorLineIdx(lines)

	successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(successColor)).Bold(true)
	successDimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(successColor)).Faint(true).Bold(true)
	errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(errorColor)).Bold(true)
	errorDimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(errorColor)).Faint(true).Bold(true)
	errorCursorStyle := lipgloss.NewStyle().Background(lipgloss.Color(errorColor)).Bold(true)
	successCursorStyle := lipgloss.NewStyle().Background(lipgloss.Color(successColor)).Bold(true)
	pendingCursorStyle := lipgloss.NewStyle().Background(lipgloss.Color(primary)).Bold(true)
	pendingStyle := t.layout.Styles.Dim.Bold(true)

	innerWidth := max(1, t.layout.Width-8)
	var b strings.Builder

	for li := startLine; li < endLine; li++ {
		line := lines[li]
		var lineBuf strings.Builder

		for i := line.start; i < line.end && i < len(t.chars); i++ {
			c := t.chars[i]
			isCursor := (i == t.cursor)

			display := c.char
			if c.char == " " && c.state == charError {
				display = "\u00b7"
			}

			switch {
			case isCursor && c.state == charError:
				lineBuf.WriteString(errorCursorStyle.Render(display))
			case isCursor && c.state == charCorrect:
				lineBuf.WriteString(successCursorStyle.Render(display))
			case isCursor:
				lineBuf.WriteString(pendingCursorStyle.Render(display))
			case c.state == charCorrect:
				if li < cursorLine {
					lineBuf.WriteString(successDimStyle.Render(display))
				} else {
					lineBuf.WriteString(successStyle.Render(display))
				}
			case c.state == charError:
				if li < cursorLine {
					lineBuf.WriteString(errorDimStyle.Render(display))
				} else {
					lineBuf.WriteString(errorStyle.Render(display))
				}
			default:
				if li <= cursorLine+1 {
					lineBuf.WriteString(pendingStyle.Render(display))
				}
			}
		}
		lineStr := lineBuf.String()
		centered := lipgloss.NewStyle().Width(innerWidth).Align(lipgloss.Center).Render(lineStr)
		b.WriteString(centered)
		b.WriteString("\n")
	}

	return b.String()
}
