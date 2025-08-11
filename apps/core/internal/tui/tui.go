package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func StartTUI() error {
	m := NewMainModel()
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to start TUI: %w", err)
	}
	return nil
}

type App struct {
	currentView View
	mainMenu    *MainMenuModel
	installForm *InstallFormModel
	width       int
	height      int
}

type View int

const (
	MainMenuView View = iota
	InstallFormView
	LoadingView
)

func NewMainModel() *App {
	return &App{
		currentView: MainMenuView,
		mainMenu:    NewMainMenuModel(),
		installForm: NewInstallFormModel(),
	}
}

func (a *App) Init() tea.Cmd {
	return nil
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		return a, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return a, tea.Quit
		}
	case SwitchViewMsg:
		a.currentView = msg.View
		return a, nil
	}

	switch a.currentView {
	case MainMenuView:
		newModel, cmd := a.mainMenu.Update(msg)
		a.mainMenu = newModel.(*MainMenuModel)
		return a, cmd
	case InstallFormView:
		newModel, cmd := a.installForm.Update(msg)
		a.installForm = newModel.(*InstallFormModel)
		return a, cmd
	}

	return a, nil
}

// View implements tea.Model
func (a *App) View() string {
	if a.width == 0 || a.height == 0 {
		return "Loading..."
	}

	// Header
	header := a.renderHeader()

	// Footer
	footer := a.renderFooter()

	// Calculate content height (total height - header - footer)
	headerHeight := lipgloss.Height(header)
	footerHeight := lipgloss.Height(footer)
	contentHeight := a.height - headerHeight - footerHeight

	// Current view content
	var content string
	switch a.currentView {
	case MainMenuView:
		content = a.mainMenu.View()
	case InstallFormView:
		content = a.installForm.View()
	default:
		content = "Unknown view"
	}

	// Style the content to fill the available space
	styledContent := ContentStyle.
		Width(a.width).
		Height(contentHeight).
		Render(content)

	// Combine all parts with proper sizing
	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		styledContent,
		footer,
	)
}

// SwitchViewMsg is sent when switching between views
type SwitchViewMsg struct {
	View View
}

