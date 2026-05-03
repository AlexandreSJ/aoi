package theme

import (
	"regexp"
	"strconv"
	"strings"
)

var hexShortRe = regexp.MustCompile(`^#[0-9a-fA-F]{3}$`)
var hexLongRe = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)

var namedColors = map[string]bool{
	"black":        true,
	"red":          true,
	"green":        true,
	"yellow":       true,
	"blue":         true,
	"magenta":      true,
	"cyan":         true,
	"white":        true,
	"brightblack":  true,
	"brightred":    true,
	"brightgreen":  true,
	"brightyellow": true,
	"brightblue":   true,
	"brightmagenta": true,
	"brightcyan":   true,
	"brightwhite":  true,
}

func IsValidColor(s string) bool {
	if hexShortRe.MatchString(s) || hexLongRe.MatchString(s) {
		return true
	}
	if namedColors[strings.ToLower(s)] {
		return true
	}
	if n, err := strconv.Atoi(s); err == nil && n >= 0 && n <= 255 {
		return true
	}
	return false
}
