package models

// RailpackOutput represents the output from railpack analyzer
type RailpackOutput struct {
	Plan struct {
		Steps []struct {
			Name   string `json:"name"`
			Inputs []struct {
				Image string `json:"image"`
			} `json:"inputs"`
			Commands []struct {
				Path       string `json:"path,omitempty"`
				Name       string `json:"name,omitempty"`
				CustomName string `json:"customName,omitempty"`
				Cmd        string `json:"cmd,omitempty"`
			} `json:"commands"`
		} `json:"steps"`
		Deploy struct {
			Base struct {
				Step string `json:"step"`
			} `json:"base"`
			Inputs []struct {
				Step    string   `json:"step"`
				Include []string `json:"include"`
			} `json:"inputs"`
			StartCommand string `json:"startCommand"`
		} `json:"deploy"`
	} `json:"plan"`
	Success bool `json:"success"`
	Logs    []struct {
		Level string `json:"level"`
		Msg   string `json:"msg"`
	} `json:"logs"`
}

// RepositoryAnalysis represents the parsed analysis result
type RepositoryAnalysis struct {
	Install []string `json:"install"`
	Build   []string `json:"build"`
	Deploy  string   `json:"deploy"`
	Success bool     `json:"success"`
}

// AnalyzeRepositoryRequest represents a request to analyze a repository
type AnalyzeRepositoryRequest struct {
	RepoURL   string `json:"repo_url" validate:"required,url"`
	Branch    string `json:"branch" validate:"omitempty"`
	RepoID    int64  `json:"repo_id" validate:"required"`
	RepoName  string `json:"repo_name" validate:"required"`
	RepoOwner string `json:"repo_owner" validate:"required"`
}
