package config

import "github.com/charmbracelet/lipgloss"

type ColorConfig struct {
	Primary lipgloss.AdaptiveColor
	Text    lipgloss.AdaptiveColor
}

var (
	Colors = ColorConfig{
		Primary: lipgloss.AdaptiveColor{Light: "#2F51A2", Dark: "#2F51A2"},
		Text:    lipgloss.AdaptiveColor{Light: "#F5F5F5", Dark: "#F5F5F5"},
	}
)
