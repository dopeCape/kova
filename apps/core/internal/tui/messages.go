package tui

// StatusMsg represents status messages (success, error, info, warning)
type StatusMsg struct {
	Type    string // "success", "error", "info", "warning"
	Message string
}

// LoadingMsg represents loading state changes
type LoadingMsg struct {
	Loading bool
	Message string
}

// InstallCompleteMsg is sent when installation is complete
type InstallCompleteMsg struct {
	Success bool
	Message string
}

