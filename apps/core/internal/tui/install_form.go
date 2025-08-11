package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	install "github.com/dopeCape/kova/scripts"
)

type InstallStep int

const (
	StepSelectType InstallStep = iota
	StepLocalForm
	StepRemoteForm
	StepInstalling
	StepComplete
)

type InstallFormModel struct {
	cursor      int
	step        InstallStep
	installType string // "local" or "remote"
	inputs      []InputField
	loading     bool
	errorMsg    string
	successMsg  string
}

type InputField struct {
	label       string
	value       string
	placeholder string
	isPassword  bool
	required    bool
	focused     bool
}

func NewInstallFormModel() *InstallFormModel {
	return &InstallFormModel{
		cursor:      0,
		step:        StepSelectType,
		installType: "",
		inputs:      []InputField{},
	}
}

func (m *InstallFormModel) Init() tea.Cmd {
	return nil
}

func (m *InstallFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.step == StepSelectType {
				return m, func() tea.Msg {
					return SwitchViewMsg{View: MainMenuView}
				}
			} else {
				// Go back to previous step
				m.step = StepSelectType
				m.cursor = 0
				m.errorMsg = ""
			}
		case "tab", "down":
			m.nextInput()
		case "shift+tab", "up":
			m.prevInput()
		case "enter":
			return m.handleEnter()
		case "ctrl+c":
			return m, tea.Quit
		default:
			// Handle input for current field
			if m.step == StepLocalForm || m.step == StepRemoteForm {
				if m.cursor < len(m.inputs) {
					m.handleInput(msg.String())
				}
			}
		}

	case StatusMsg:
		switch msg.Type {
		case "error":
			m.errorMsg = msg.Message
			m.loading = false
		case "success":
			m.successMsg = msg.Message
			m.loading = false
			m.step = StepComplete
		}
	}

	return m, nil
}

func (m *InstallFormModel) handleEnter() (tea.Model, tea.Cmd) {
	switch m.step {
	case StepSelectType:
		if m.cursor == 0 {
			// Local installation
			m.installType = "local"
			m.setupLocalForm()
			m.step = StepLocalForm
			m.cursor = 0
		} else if m.cursor == 1 {
			// Remote installation - disabled for now
			return m, func() tea.Msg {
				return StatusMsg{
					Type:    "error",
					Message: "Remote installation is currently disabled",
				}
			}
		} else if m.cursor == 2 {
			// Back button
			return m, func() tea.Msg {
				return SwitchViewMsg{View: MainMenuView}
			}
		}
	case StepLocalForm:
		if m.cursor == len(m.inputs) {
			// Install button pressed
			return m.handleInstall()
		} else if m.cursor == len(m.inputs)+1 {
			// Back button pressed
			m.step = StepSelectType
			m.cursor = 0
			m.errorMsg = ""
		}
	case StepRemoteForm:
		// Similar to local form but with remote-specific fields
		if m.cursor == len(m.inputs) {
			return m.handleInstall()
		} else if m.cursor == len(m.inputs)+1 {
			m.step = StepSelectType
			m.cursor = 0
			m.errorMsg = ""
		}
	case StepComplete:
		// Any key goes back to main menu
		return m, func() tea.Msg {
			return SwitchViewMsg{View: MainMenuView}
		}
	}
	return m, nil
}

func (m *InstallFormModel) setupLocalForm() {
	m.inputs = []InputField{
		{
			label:       "Admin Email",
			placeholder: "admin@example.com",
			required:    true,
		},
		{
			label:       "Admin Password",
			placeholder: "Enter secure password",
			isPassword:  true,
			required:    true,
		},
		{
			label:       "Admin Username",
			value:       "admin",
			placeholder: "admin",
			required:    false,
		},
		{
			label:       "Domain (Optional)",
			placeholder: "your-domain.com",
			required:    false,
		},
	}
}

func (m *InstallFormModel) setupRemoteForm() {
	m.inputs = []InputField{
		{
			label:       "Remote IP Address",
			placeholder: "192.168.1.100",
			required:    true,
		},
		{
			label:       "SSH Username",
			placeholder: "root",
			required:    true,
		},
		{
			label:       "SSH Key Path",
			placeholder: "/path/to/private/key",
			required:    false,
		},
		{
			label:       "Admin Email",
			placeholder: "admin@example.com",
			required:    true,
		},
		{
			label:       "Admin Password",
			placeholder: "Enter secure password",
			isPassword:  true,
			required:    true,
		},
		{
			label:       "Admin Username",
			value:       "admin",
			placeholder: "admin",
			required:    false,
		},
	}
}

func (m *InstallFormModel) handleInput(key string) {
	if m.cursor >= len(m.inputs) {
		return
	}

	currentInput := &m.inputs[m.cursor]

	switch key {
	case "backspace":
		if len(currentInput.value) > 0 {
			currentInput.value = currentInput.value[:len(currentInput.value)-1]
		}
	case "space":
		currentInput.value += " "
	default:
		// Only add printable characters
		if len(key) == 1 && key >= " " && key <= "~" {
			currentInput.value += key
		}
	}

	// Clear error when user starts typing
	if m.errorMsg != "" {
		m.errorMsg = ""
	}
}

func (m *InstallFormModel) nextInput() {
	maxCursor := 0
	switch m.step {
	case StepSelectType:
		maxCursor = 2 // Local, Remote (disabled), Back
	case StepLocalForm, StepRemoteForm:
		maxCursor = len(m.inputs) + 1 // inputs + Install + Back
	}

	m.cursor++
	if m.cursor > maxCursor {
		m.cursor = 0
	}
}

func (m *InstallFormModel) prevInput() {
	maxCursor := 0
	switch m.step {
	case StepSelectType:
		maxCursor = 2
	case StepLocalForm, StepRemoteForm:
		maxCursor = len(m.inputs) + 1
	}

	m.cursor--
	if m.cursor < 0 {
		m.cursor = maxCursor
	}
}

func (m *InstallFormModel) handleInstall() (tea.Model, tea.Cmd) {
	// Validate required fields
	for _, input := range m.inputs {
		if input.required && strings.TrimSpace(input.value) == "" {
			return m, func() tea.Msg {
				return StatusMsg{
					Type:    "error",
					Message: fmt.Sprintf("%s is required", input.label),
				}
			}
		}
	}

	m.loading = true
	m.errorMsg = ""
	m.successMsg = ""
	m.step = StepInstalling

	return m, func() tea.Msg {
		// Create install configuration from form inputs
		config := &install.InstallConfig{
			InstallType: m.installType,
			DataDir:     "/data/kova",
			InstallDir:  "/opt/kova",
		}

		// Map form inputs to config
		for _, input := range m.inputs {
			switch input.label {
			case "Admin Email":
				config.AdminEmail = input.value
			case "Admin Password":
				config.AdminPassword = input.value
			case "Admin Username":
				config.AdminUsername = input.value
			case "Domain (Optional)":
				config.Domain = input.value
			}
		}

		if config.AdminUsername == "" {
			config.AdminUsername = "admin"
		}

		if err := install.Install(config); err != nil {
			return StatusMsg{
				Type:    "error",
				Message: fmt.Sprintf("Installation failed: %v", err),
			}
		}

		return StatusMsg{
			Type:    "success",
			Message: "Installation completed successfully!",
		}
	}
}

func (m *InstallFormModel) View() string {
	switch m.step {
	case StepSelectType:
		return m.renderTypeSelection()
	case StepLocalForm:
		return m.renderLocalForm()
	case StepRemoteForm:
		return m.renderRemoteForm()
	case StepInstalling:
		return m.renderInstalling()
	case StepComplete:
		return m.renderComplete()
	}

	return "Unknown step"
}

func (m *InstallFormModel) renderTypeSelection() string {
	var s strings.Builder

	// Title
	title := RenderTitle("Install Kova")
	s.WriteString(lipgloss.NewStyle().MarginBottom(1).Render(title))
	s.WriteString("\n")

	// Subtitle
	subtitle := lipgloss.NewStyle().
		Foreground(Gray500).
		Render("Select installation type:")
	s.WriteString(subtitle)
	s.WriteString("\n\n")

	// Options
	options := []string{"Local Installation", "Remote Installation (Disabled)", "Back"}
	for i, option := range options {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		var optionStyle lipgloss.Style
		if i == 1 {
			// Remote option - disabled
			optionStyle = lipgloss.NewStyle().
				Foreground(Gray600).
				Padding(0, 2)
		} else if m.cursor == i {
			optionStyle = MenuItemSelectedStyle
		} else {
			optionStyle = MenuItemStyle
		}

		optionText := fmt.Sprintf("%s %s", cursor, option)
		s.WriteString(optionStyle.Render(optionText))
		s.WriteString("\n")
	}

	s.WriteString("\n")
	helpText := HelpStyle.Render("Up/Down: Navigate • Enter: Select • Esc: Back")
	s.WriteString(helpText)

	return s.String()
}

func (m *InstallFormModel) renderLocalForm() string {
	var s strings.Builder

	// Title
	title := RenderTitle("Local Installation")
	s.WriteString(lipgloss.NewStyle().MarginBottom(1).Render(title))
	s.WriteString("\n")

	// Error message
	if m.errorMsg != "" {
		s.WriteString(RenderError(m.errorMsg))
		s.WriteString("\n\n")
	}

	// Form fields
	for i, input := range m.inputs {
		// Label
		labelStyle := LabelStyle
		if i == m.cursor {
			labelStyle = labelStyle.Foreground(White)
		}
		s.WriteString(labelStyle.Render(input.label))
		if input.required {
			s.WriteString(ErrorStyle.Render(" *"))
		}
		s.WriteString("\n")

		// Input field
		value := input.value
		if input.isPassword && value != "" {
			value = strings.Repeat("*", len(value))
		}

		if value == "" && input.placeholder != "" {
			value = lipgloss.NewStyle().Foreground(Gray600).Render(input.placeholder)
		}

		inputStyle := InputStyle
		if i == m.cursor {
			inputStyle = InputFocusedStyle
			value += "█" // cursor
		}

		s.WriteString(inputStyle.Render(value))
		s.WriteString("\n\n")
	}

	// Buttons
	installBtn := "Install"
	if m.loading {
		installBtn = "Installing..."
	}

	installButton := RenderButton(installBtn, m.cursor == len(m.inputs))
	backButton := RenderButton("Back", m.cursor == len(m.inputs)+1)

	s.WriteString(lipgloss.JoinHorizontal(lipgloss.Left, installButton, backButton))
	s.WriteString("\n\n")

	// Help text
	helpText := HelpStyle.Render("Tab: Next • Shift+Tab: Previous • Enter: Select • Esc: Back")
	s.WriteString(helpText)

	return s.String()
}

func (m *InstallFormModel) renderRemoteForm() string {
	// Similar to local form but with different title
	var s strings.Builder

	title := RenderTitle("Remote Installation")
	s.WriteString(lipgloss.NewStyle().MarginBottom(1).Render(title))
	s.WriteString("\n")

	// Rest is similar to local form...
	// (Implementation would be similar to renderLocalForm)

	return s.String()
}

func (m *InstallFormModel) renderInstalling() string {
	var s strings.Builder

	title := RenderTitle("Installing...")
	s.WriteString(lipgloss.NewStyle().MarginBottom(1).Render(title))
	s.WriteString("\n")

	s.WriteString(lipgloss.NewStyle().
		Foreground(Gray400).
		Render("Please wait while Kova is being installed..."))

	return s.String()
}

func (m *InstallFormModel) renderComplete() string {
	var s strings.Builder

	title := RenderTitle("Installation Complete")
	s.WriteString(lipgloss.NewStyle().MarginBottom(1).Render(title))
	s.WriteString("\n")

	if m.successMsg != "" {
		s.WriteString(RenderSuccess(m.successMsg))
		s.WriteString("\n\n")
	}

	helpText := HelpStyle.Render("Press any key to return to main menu")
	s.WriteString(helpText)

	return s.String()
}

