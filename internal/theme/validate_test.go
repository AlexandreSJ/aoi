package theme

import "testing"

func TestIsValidColor(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		// Valid hex
		{"#FFF", true},
		{"#7D56F4", true},
		{"#abc", true},
		// Valid named
		{"red", true},
		{"brightgreen", true},
		{"CYAN", true}, // case insensitive
		// Valid ANSI
		{"0", true},
		{"255", true},
		{"128", true},
		// Invalid
		{"#12", false},
		{"#GGGGGG", false},
		{"purple", false},
		{"256", false},
		{"-1", false},
		{"", false},
	}

	for _, tt := range tests {
		got := IsValidColor(tt.input)
		if got != tt.want {
			t.Errorf("IsValidColor(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}
