package config

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed default.yaml
var defaultConfig embed.FS

func Default() *Config {
	return &Config{
		Colors: Colors{
			Primary:   "32",
			Secondary: "27",
			Text:      "231",
			Dim:       "8",
			Title:     "232",
			Footer:    "232",
			Error:     "162",
			Success:   "75",
		},
		System: System{
			ConfigPath: "~/.config/aoi",
			WordsDir:   "~/.config/aoi/words",
			QuotesDir:  "~/.config/aoi/quotes",
		},
	}
}

type Colors struct {
	Primary   string `yaml:"primary"`
	Secondary string `yaml:"secondary"`
	Text      string `yaml:"text"`
	Dim       string `yaml:"dim"`
	Title     string `yaml:"title"`
	Footer    string `yaml:"footer"`
	Error     string `yaml:"error"`
	Success   string `yaml:"success"`
}

type System struct {
	ConfigPath string `yaml:"config"`
	WordsDir   string `yaml:"words"`
	QuotesDir  string `yaml:"quotes"`
}

type Config struct {
	Colors Colors `yaml:"colors"`
	System System `yaml:"system"`
}

type Section struct {
	Title string
	Keys  []string
}

func (c Config) Sections() []Section {
	return []Section{
		{
			Title: "Colors",
			Keys: []string{
				"colors.primary", "colors.secondary", "colors.text", "colors.dim",
				"colors.title", "colors.footer", "colors.error", "colors.success",
			},
		},
		{
			Title: "System",
			Keys:  []string{"system.config", "system.words", "system.quotes"},
		},
	}
}

func IsColorKey(key string) bool {
	return strings.HasPrefix(key, "colors.")
}

func IsInlineHintKey(key string) bool {
	return key == "system.config"
}

func InlineHint(key string) string {
	switch key {
	case "system.config":
		return "/config.yaml"
	default:
		return ""
	}
}

func (c Config) Get(key string) string {
	switch key {
	case "colors.primary":
		return c.Colors.Primary
	case "colors.secondary":
		return c.Colors.Secondary
	case "colors.text":
		return c.Colors.Text
	case "colors.dim":
		return c.Colors.Dim
	case "colors.title":
		return c.Colors.Title
	case "colors.footer":
		return c.Colors.Footer
	case "colors.error":
		return c.Colors.Error
	case "colors.success":
		return c.Colors.Success
	case "system.config":
		return c.System.ConfigPath
	case "system.words":
		return c.System.WordsDir
	case "system.quotes":
		return c.System.QuotesDir
	default:
		return ""
	}
}

func (c *Config) Set(key, value string) {
	switch key {
	case "colors.primary":
		c.Colors.Primary = value
	case "colors.secondary":
		c.Colors.Secondary = value
	case "colors.text":
		c.Colors.Text = value
	case "colors.dim":
		c.Colors.Dim = value
	case "colors.title":
		c.Colors.Title = value
	case "colors.footer":
		c.Colors.Footer = value
	case "colors.error":
		c.Colors.Error = value
	case "colors.success":
		c.Colors.Success = value
	case "system.config":
		c.System.ConfigPath = value
	case "system.words":
		c.System.WordsDir = value
	case "system.quotes":
		c.System.QuotesDir = value
	}
}

func expandPath(path string) (string, error) {
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

func resolveConfigPath(cfg *Config) (string, error) {
	if cfg.System.ConfigPath != "" && cfg.System.ConfigPath != "~/.config/aoi" {
		path, err := expandPath(cfg.System.ConfigPath)
		if err != nil {
			return "", err
		}
		return filepath.Join(path, "config.yaml"), nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home dir: %w", err)
	}
	return filepath.Join(home, ".config", "aoi", "config.yaml"), nil
}

func ResolveDir(raw string, fallback string) (string, error) {
	if raw == "" {
		raw = fallback
	}
	return expandPath(raw)
}

func Load() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("get home dir: %w", err)
	}
	defaultPath := filepath.Join(home, ".config", "aoi", "config.yaml")

	cfg, err := loadFrom(defaultPath)
	if err != nil {
		return nil, err
	}

	overridePath, err := resolveConfigPath(cfg)
	if err != nil {
		return cfg, nil
	}
	if overridePath != defaultPath {
		return loadFrom(overridePath)
	}

	return cfg, nil
}

func loadFrom(path string) (*Config, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return nil, fmt.Errorf("create config dir: %w", err)
		}
		data, err := defaultConfig.ReadFile("default.yaml")
		if err != nil {
			return nil, fmt.Errorf("read embedded default: %w", err)
		}
		if err := os.WriteFile(path, data, 0o644); err != nil {
			return nil, fmt.Errorf("write default config: %w", err)
		}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	return &cfg, nil
}

func Save(cfg *Config) error {
	path, err := resolveConfigPath(cfg)
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write config: %w", err)
	}
	return nil
}
