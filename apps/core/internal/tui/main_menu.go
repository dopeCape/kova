package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MainMenuModel struct {
	cursor    int
	menuItems []MenuItem
}

type MenuItem struct {
	title       string
	description string
	action      string
}

func NewMainMenuModel() *MainMenuModel {
	return &MainMenuModel{
		cursor: 0,
		menuItems: []MenuItem{
			{
				title:       "Install Kova",
				description: "Install Kova deployment manager on local or remote machine",
				action:      "install",
			},

			{
				title:       "Exit",
				description: "Exit the interactive mode",
				action:      "exit",
			},
		},
	}
}

func (m *MainMenuModel) Init() tea.Cmd {
	return nil
}

func (m *MainMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.menuItems)-1 {
				m.cursor++
			}
		case "enter", " ":
			return m.handleSelection()
		}
	}

	return m, nil
}

func (m *MainMenuModel) handleSelection() (tea.Model, tea.Cmd) {
	selected := m.menuItems[m.cursor]

	switch selected.action {
	case "install":
		return m, func() tea.Msg {
			return SwitchViewMsg{View: InstallFormView}
		}
	case "exit":
		return m, tea.Quit
	default:
		// For now, just show a message for unimplemented features
		return m, func() tea.Msg {
			return StatusMsg{
				Type:    "info",
				Message: fmt.Sprintf("%s - Coming soon!", selected.title),
			}
		}
	}
}

func (m *MainMenuModel) View() string {
	var s strings.Builder

	title := RenderTitle("KOVA")
	s.WriteString(lipgloss.NewStyle().MarginBottom(1).Render(title))
	s.WriteString("\n")

	subtitle := lipgloss.NewStyle().
		Foreground(Gray500).
		Render("Select an option:")
	s.WriteString(subtitle)
	s.WriteString("\n\n")

	// Menu items
	for i, item := range m.menuItems {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		var itemStyle lipgloss.Style
		if m.cursor == i {
			itemStyle = MenuItemSelectedStyle
		} else {
			itemStyle = MenuItemStyle
		}

		// Render menu item
		itemText := fmt.Sprintf("%s %s", cursor, item.title)
		renderedItem := itemStyle.Render(itemText)
		s.WriteString(renderedItem)
		s.WriteString("\n")

		// Description (only for selected item)
		if m.cursor == i {
			description := lipgloss.NewStyle().
				Foreground(Gray500).
				MarginLeft(4).
				Render(item.description)
			s.WriteString(description)
			s.WriteString("\n")
		}
	}

	return s.String()
}

