package words

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

type WordList struct {
	Name  string
	Words []string
}

func LoadList(name string, userDir string) (*WordList, error) {
	if userDir != "" {
		path := filepath.Join(userDir, name+".txt")
		if wl, err := loadFromFile(name, path); err == nil {
			return wl, nil
		}
	}

	data, err := embeddedFiles.ReadFile("embedded/" + name + ".txt")
	if err == nil {
		return parseWordList(name, data), nil
	}

	return nil, fmt.Errorf("word list %q not found", name)
}

func loadFromFile(name, path string) (*WordList, error) {
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
		return nil, fmt.Errorf("empty word list")
	}
	return &WordList{Name: name, Words: lines}, nil
}

func parseWordList(name string, data []byte) *WordList {
	var words []string
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			words = append(words, line)
		}
	}
	return &WordList{Name: name, Words: words}
}

func AvailableLists(userDir string) ([]string, error) {
	seen := map[string]bool{}
	var names []string

	entries, err := embeddedFiles.ReadDir("embedded")
	if err == nil {
		for _, e := range entries {
			name := strings.TrimSuffix(e.Name(), ".txt")
			if !seen[name] {
				seen[name] = true
				names = append(names, name)
			}
		}
	}

	if userDir != "" {
		entries, err = os.ReadDir(userDir)
		if err == nil {
			for _, e := range entries {
				name := strings.TrimSuffix(e.Name(), ".txt")
				if !seen[name] {
					seen[name] = true
					names = append(names, name)
				}
			}
		}
	}

	return names, nil
}

func (wl *WordList) Sample(n int) []string {
	result := make([]string, n)
	for i := range n {
		result[i] = wl.Words[rand.IntN(len(wl.Words))]
	}
	return result
}

func (wl *WordList) Infinite() *InfiniteGenerator {
	return &InfiniteGenerator{words: wl.Words}
}

type InfiniteGenerator struct {
	words []string
}

func (g *InfiniteGenerator) Next() string {
	return g.words[rand.IntN(len(g.words))]
}

func EnsureUserDir(dir string) error {
	if dir == "" {
		return nil
	}
	expanded, err := expandHome(dir)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(expanded, 0o755); err != nil {
		return fmt.Errorf("create dir %s: %w", expanded, err)
	}
	return nil
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
