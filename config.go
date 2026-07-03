package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
)

type Theme struct {
	TitleBg       string `json:"titleBg"`
	TitleFg       string `json:"titleFg"`
	TabFg         string `json:"tabFg"`
	ActiveTabFg   string `json:"activeTabFg"`
	OverdueFg     string `json:"overdueFg"`
	CompletedFg   string `json:"completedFg"`
	ProgressBg    string `json:"progressBg"`
	ProgressFg    string `json:"progressFg"`
}

var defaultTheme = Theme{
	TitleBg:     "62",
	TitleFg:     "230",
	TabFg:       "240",
	ActiveTabFg: "255",
	OverdueFg:   "196",
	CompletedFg: "240",
	ProgressBg:  "240",
	ProgressFg:  "114",
}

var darkTheme = Theme{
	TitleBg:     "55",
	TitleFg:     "255",
	TabFg:       "243",
	ActiveTabFg: "255",
	OverdueFg:   "203",
	CompletedFg: "243",
	ProgressBg:  "237",
	ProgressFg:  "78",
}

var lightTheme = Theme{
	TitleBg:     "33",
	TitleFg:     "255",
	TabFg:       "246",
	ActiveTabFg: "255",
	OverdueFg:   "160",
	CompletedFg: "246",
	ProgressBg:  "250",
	ProgressFg:  "28",
}

var themePresets = map[string]Theme{
	"default": defaultTheme,
	"dark":    darkTheme,
	"light":   lightTheme,
}

type Config struct {
	DefaultCategory string `json:"defaultCategory"`
	ThemeName       string `json:"theme"`
}

func DefaultConfig() Config {
	return Config{
		DefaultCategory: "Home",
		ThemeName:       "default",
	}
}

func configFilePath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = "."
	}
	return filepath.Join(configDir, "my-tui-app", "config.json")
}

func LoadConfig(path string) (Config, error) {
	cfg := DefaultConfig()
	file, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, err
	}
	if err := json.Unmarshal(file, &cfg); err != nil {
		return DefaultConfig(), err
	}
	return cfg, nil
}

func resolveTheme(cfg Config) Theme {
	if t, ok := themePresets[cfg.ThemeName]; ok {
		return t
	}
	return defaultTheme
}

type styles struct {
	title      lipgloss.Style
	tab        lipgloss.Style
	activeTab  lipgloss.Style
	overdue    lipgloss.Style
	completed  lipgloss.Style
	progressBg lipgloss.Style
	progressFg lipgloss.Style
}

func buildStyles(t Theme) styles {
	return styles{
		title:      lipgloss.NewStyle().Background(lipgloss.Color(t.TitleBg)).Foreground(lipgloss.Color(t.TitleFg)).Padding(0, 1).Bold(true),
		tab:        lipgloss.NewStyle().Padding(0, 1).Foreground(lipgloss.Color(t.TabFg)),
		activeTab:  lipgloss.NewStyle().Padding(0, 1).Foreground(lipgloss.Color(t.ActiveTabFg)).Bold(true).Border(lipgloss.NormalBorder(), false, false, true, false),
		overdue:    lipgloss.NewStyle().Foreground(lipgloss.Color(t.OverdueFg)),
		completed:  lipgloss.NewStyle().Foreground(lipgloss.Color(t.CompletedFg)),
		progressBg: lipgloss.NewStyle().Foreground(lipgloss.Color(t.ProgressBg)),
		progressFg: lipgloss.NewStyle().Foreground(lipgloss.Color(t.ProgressFg)),
	}
}
