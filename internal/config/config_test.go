package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadCreatesDefault(t *testing.T) {
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	path := filepath.Join(tmpDir, ".config", "aoi", "config.yaml")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("config file not created")
	}

	if cfg.Colors.Primary != "32" {
		t.Fatalf("expected primary 32, got %s", cfg.Colors.Primary)
	}
}

func TestSaveAndReload(t *testing.T) {
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	cfg, _ := Load()
	cfg.Colors.Primary = "#FF0000"
	if err := Save(cfg); err != nil {
		t.Fatalf("Save: %v", err)
	}

	reload, err := Load()
	if err != nil {
		t.Fatalf("Load after save: %v", err)
	}
	if reload.Colors.Primary != "#FF0000" {
		t.Fatalf("expected #FF0000, got %s", reload.Colors.Primary)
	}
}

func TestGetSet(t *testing.T) {
	cfg := &Config{}
	cfg.Set("colors.primary", "#123456")
	cfg.Set("system.config", "/tmp/aoi")

	if got := cfg.Get("colors.primary"); got != "#123456" {
		t.Fatalf("expected #123456, got %s", got)
	}
	if got := cfg.Get("system.config"); got != "/tmp/aoi" {
		t.Fatalf("expected /tmp/aoi, got %s", got)
	}
}

func TestSections(t *testing.T) {
	cfg := &Config{}
	sections := cfg.Sections()
	if len(sections) < 2 {
		t.Fatalf("expected at least 2 sections, got %d", len(sections))
	}
	if sections[0].Title != "Colors" {
		t.Fatalf("expected first section 'Colors', got %s", sections[0].Title)
	}
	if sections[1].Title != "System" {
		t.Fatalf("expected second section 'System', got %s", sections[1].Title)
	}
}
