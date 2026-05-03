package quotes

import (
	"bufio"
	"embed"
	"fmt"
	"math/rand/v2"
	"os"
	"path/filepath"
	"strings"
)

//go:embed embedded/*.txt
var embeddedFiles embed.FS

type QuoteList struct {
	Name   string
	Quotes []string
}

func LoadList(name string, userDir string) (*QuoteList, error) {
	if userDir != "" {
		path := filepath.Join(userDir, name+".txt")
		if ql, err := loadFromFile(name, path); err == nil {
			return ql, nil
		}
	}

	data, err := embeddedFiles.ReadFile("embedded/" + name + ".txt")
	if err == nil {
		return parseQuoteList(name, data), nil
	}

	return nil, fmt.Errorf("quote list %q not found", name)
}

func loadFromFile(name, path string) (*QuoteList, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			lines = append(lines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if len(lines) == 0 {
		return nil, fmt.Errorf("empty quote list")
	}
	return &QuoteList{Name: name, Quotes: lines}, nil
}

func parseQuoteList(name string, data []byte) *QuoteList {
	var quotes []string
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			quotes = append(quotes, line)
		}
	}
	return &QuoteList{Name: name, Quotes: quotes}
}

func (ql *QuoteList) Random() string {
	return ql.Quotes[rand.IntN(len(ql.Quotes))]
}

func EnsureUserDir(dir string) error {
	if dir == "" {
		return nil
	}
	expanded, err := expandHome(dir)
	if err != nil {
		return err
	}
	return os.MkdirAll(expanded, 0o755)
}

func CopyEmbeddedToUser(name string, userDir string) error {
	if userDir == "" {
		return nil
	}
	expanded, err := expandHome(userDir)
	if err != nil {
		return err
	}
	dest := filepath.Join(expanded, name+".txt")
	if _, err := os.Stat(dest); err == nil {
		return nil
	}
	data, err := embeddedFiles.ReadFile("embedded/" + name + ".txt")
	if err != nil {
		return nil
	}
	return os.WriteFile(dest, data, 0o644)
}

func expandHome(path string) (string, error) {
	path = os.ExpandEnv(path)
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = filepath.Join(home, path[2:])
	}
	return path, nil
}
