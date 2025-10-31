package models

import (
	"time"
)

type EnvironmentVariable struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Project struct {
	ID               string                `json:"id"`
	Name             string                `json:"name"`
	UserID           string                `json:"user_id"`
	RepoID           int64                 `json:"repo_id"`
	RepoName         string                `json:"repo_name"`
	RepoFullName     string                `json:"repo_full_name"`
	RepoURL          string                `json:"repo_url"`
	RepoBranch       string                `json:"repo_branch"`
	Status           string                `json:"status"`
	EnvVariables     []EnvironmentVariable `json:"env_variables"`
	DeploymentStatus string                `json:"deployment_status"`
	Domain           string                `json:"domain,omitempty"`
	Port             int                   `json:"port,omitempty"`
	CreatedAt        time.Time             `json:"created_at"`
	UpdatedAt        time.Time             `json:"updated_at"`
}

type CreateProjectRequest struct {
	Name         string                `json:"name" validate:"required,min=1,max=50"`
	Domain       string                `json:"domain" validate:"required,min=3"`
	RepoID       int64                 `json:"repo_id" validate:"required,min=1"`
	RepoName     string                `json:"repo_name" validate:"required,min=1,max=255"`
	RepoFullName string                `json:"repo_full_name" validate:"required,min=1,max=255"`
	RepoURL      string                `json:"repo_url" validate:"required,url"`
	RepoBranch   string                `json:"repo_branch" validate:"omitempty,min=1,max=255"`
	EnvVariables []EnvironmentVariable `json:"env_variables" validate:"omitempty,dive"`
}

type UpdateProjectRequest struct {
	Name       string `json:"name" validate:"omitempty,min=1,max=50,alphanum_dash"`
	RepoBranch string `json:"repo_branch" validate:"omitempty,min=1,max=255"`
	Status     string `json:"status" validate:"omitempty,oneof=active inactive archived"`
	Domain     string `json:"domain" validate:"omitempty,min=3"`
}

type ProjectWithRepository struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	UserID           string    `json:"user_id"`
	RepoID           int64     `json:"repo_id"`
	RepoName         string    `json:"repo_name"`
	RepoFullName     string    `json:"repo_full_name"`
	RepoURL          string    `json:"repo_url"`
	RepoBranch       string    `json:"repo_branch"`
	Status           string    `json:"status"`
	DeploymentStatus string    `json:"deployment_status"`
	Domain           string    `json:"domain,omitempty"`
	Port             int       `json:"port,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// IsOwnedBy checks if the project belongs to the specified user
func (p *Project) IsOwnedBy(userID string) bool {
	return p.UserID == userID
}

// IsActive checks if the project is in active status
func (p *Project) IsActive() bool {
	return p.Status == "active"
}

// ToPublic returns a copy of the Project with all fields visible
func (p *Project) ToPublic() *Project {
	return &Project{
		ID:               p.ID,
		Name:             p.Name,
		UserID:           p.UserID,
		RepoID:           p.RepoID,
		RepoName:         p.RepoName,
		RepoFullName:     p.RepoFullName,
		RepoURL:          p.RepoURL,
		RepoBranch:       p.RepoBranch,
		Status:           p.Status,
		DeploymentStatus: p.DeploymentStatus,
		Domain:           p.Domain,
		Port:             p.Port,
		EnvVariables:     p.EnvVariables,
		CreatedAt:        p.CreatedAt,
		UpdatedAt:        p.UpdatedAt,
	}
}

// ValidateStatus checks if the status is valid
func (p *Project) ValidateStatus() bool {
	validStatuses := []string{"active", "inactive", "archived"}
	for _, status := range validStatuses {
		if p.Status == status {
			return true
		}
	}
	return false
}

// SetDefaults sets default values for optional fields
func (req *CreateProjectRequest) SetDefaults() {
	if req.RepoBranch == "" {
		req.RepoBranch = "main"
	}
}
