package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

func (a *App) renderHeader() string {
	// Logo/Title
	logo := lipgloss.NewStyle().
		Bold(true).
		Foreground(White).
		Render("KOVA")

	// Current view indicator
	var viewName string
	switch a.currentView {
	case MainMenuView:
		viewName = "Main Menu"
	case InstallFormView:
		viewName = "Installation"
	default:
		viewName = "Unknown"
	}

	viewIndicator := lipgloss.NewStyle().
		Foreground(Gray500).
		Render(fmt.Sprintf("» %s", viewName))

	// Time
	timeStr := lipgloss.NewStyle().
		Foreground(Gray500).
		Render(time.Now().Format("15:04:05"))

	// Join left parts
	leftSide := strings.Join([]string{logo, viewIndicator}, " ")

	// Calculate spacing for right alignment
	headerContent := leftSide
	if a.width > 0 {
		leftWidth := lipgloss.Width(leftSide)
		timeWidth := lipgloss.Width(timeStr)
		availableSpace := a.width - leftWidth - timeWidth - 4 // account for padding

		if availableSpace > 0 {
			spacing := strings.Repeat(" ", availableSpace)
			headerContent = leftSide + spacing + timeStr
		}
	}

	return HeaderStyle.Width(a.width).Render(headerContent)
}

func (a *App) renderFooter() string {
	var helpText string

	switch a.currentView {
	case MainMenuView:
		helpText = "↑/↓: Navigate • Enter: Select • Q: Quit"
	case InstallFormView:
		helpText = "Tab: Next • Shift+Tab: Previous • Enter: Confirm • Esc: Back • Q: Quit"
	default:
		helpText = "Q: Quit"
	}

	return FooterStyle.Width(a.width).Render(helpText)
}

