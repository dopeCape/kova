package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	White = lipgloss.Color("#ffffff")
	Black = lipgloss.Color("#000000")

	Gray100 = lipgloss.Color("#f5f5f5")
	Gray200 = lipgloss.Color("#e5e5e5")
	Gray300 = lipgloss.Color("#d4d4d4")
	Gray400 = lipgloss.Color("#a3a3a3")
	Gray500 = lipgloss.Color("#737373")
	Gray600 = lipgloss.Color("#525252")
	Gray700 = lipgloss.Color("#404040")
	Gray800 = lipgloss.Color("#262626")
	Gray900 = lipgloss.Color("#171717")

	Background = Black
	Surface    = Gray900
)

// Style definitions
var (
	// Base styles
	BaseStyle = lipgloss.NewStyle().
			Foreground(White).
			Background(Background)

	// Header styles
	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(White).
			Background(Surface).
			Padding(1, 2)

	// Title styles
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(White)

	// Menu item styles
	MenuItemStyle = lipgloss.NewStyle().
			Foreground(Gray300).
			Padding(0, 2)

	MenuItemSelectedStyle = lipgloss.NewStyle().
				Foreground(Black).
				Background(White).
				Bold(true).
				Padding(0, 2)

	// Button styles
	ButtonStyle = lipgloss.NewStyle().
			Foreground(Black).
			Background(White).
			Bold(true).
			Padding(0, 2).
			MarginRight(1)

	ButtonDisabledStyle = lipgloss.NewStyle().
				Foreground(Gray600).
				Background(Gray800).
				Padding(0, 2).
				MarginRight(1)

	// Input styles
	InputStyle = lipgloss.NewStyle().
			Foreground(White).
			Background(Surface).
			Padding(0, 1).
			Border(lipgloss.NormalBorder()).
			BorderForeground(Gray600).
			Width(40)

	InputFocusedStyle = lipgloss.NewStyle().
				Foreground(White).
				Background(Surface).
				Padding(0, 1).
				Border(lipgloss.NormalBorder()).
				BorderForeground(White).
				Width(40)

	// Label styles
	LabelStyle = lipgloss.NewStyle().
			Foreground(Gray300).
			Bold(true)

	// Error styles
	ErrorStyle = lipgloss.NewStyle().
			Foreground(White).
			Background(Gray800).
			Bold(true)

	// Success styles
	SuccessStyle = lipgloss.NewStyle().
			Foreground(White).
			Bold(true)

	// Warning styles
	WarningStyle = lipgloss.NewStyle().
			Foreground(Gray300).
			Bold(true)

	// Info styles
	InfoStyle = lipgloss.NewStyle().
			Foreground(Gray400)

	// Footer styles
	FooterStyle = lipgloss.NewStyle().
			Foreground(Gray500).
			Background(Surface).
			Padding(1, 2)

	// Container styles
	ContainerStyle = lipgloss.NewStyle().
			Padding(2, 4).
			Background(Background)

	// Content styles
	ContentStyle = lipgloss.NewStyle().
			Background(Background).
			Padding(1, 2)

	// Help text styles
	HelpStyle = lipgloss.NewStyle().
			Foreground(Gray500)
)

// Utility functions for common styling
func RenderTitle(text string) string {
	return TitleStyle.Render(text)
}

func RenderButton(text string, selected bool) string {
	if selected {
		return ButtonStyle.Render(text)
	}
	return ButtonDisabledStyle.Render(text)
}

func RenderError(text string) string {
	return ErrorStyle.Render("Error: " + text)
}

func RenderSuccess(text string) string {
	return SuccessStyle.Render("Success: " + text)
}

func RenderWarning(text string) string {
	return WarningStyle.Render("Warning: " + text)
}

func RenderInfo(text string) string {
	return InfoStyle.Render("Info: " + text)
}

