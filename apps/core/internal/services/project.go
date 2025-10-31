package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dopeCape/kova/internal/models"
	"github.com/dopeCape/kova/internal/store"
	"github.com/go-playground/validator/v10"
)

type ProjectService struct {
	store        store.Store
	validator    *validator.Validate
	buildService *BuildService
}

func NewProjectService(store store.Store, buildService *BuildService) *ProjectService {
	return &ProjectService{
		store:        store,
		validator:    validator.New(),
		buildService: buildService,
	}
}

func (s *ProjectService) CreateProject(ctx context.Context, userID string, req *models.CreateProjectRequest) (*models.Project, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	_, err := s.store.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	req.SetDefaults()

	if req.EnvVariables == nil {
		req.EnvVariables = []models.EnvironmentVariable{}
	}

	exists, err := s.store.ProjectExistsByUserIDAndName(ctx, userID, req.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to check project name existence: %w", err)
	}
	if exists {
		return nil, errors.New("project with this name already exists for this user")
	}

	project := &models.Project{
		Name:             req.Name,
		UserID:           userID,
		RepoID:           req.RepoID,
		RepoName:         req.RepoName,
		RepoFullName:     req.RepoFullName,
		RepoURL:          req.RepoURL,
		RepoBranch:       req.RepoBranch,
		Status:           "active",
		DeploymentStatus: "pending",
		Domain:           req.Domain,
		EnvVariables:     req.EnvVariables,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := s.store.CreateProject(ctx, project); err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	// Enqueue build job
	if s.buildService != nil {
		s.buildService.Enqueue(project.ID, userID)
	}

	return project.ToPublic(), nil
}

// GetProject retrieves a project by ID and verifies ownership
func (s *ProjectService) GetProject(ctx context.Context, userID, projectID string) (*models.Project, error) {
	project, err := s.store.GetProjectByID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	// Verify ownership
	if !project.IsOwnedBy(userID) {
		return nil, errors.New("access denied: project does not belong to user")
	}

	return project.ToPublic(), nil
}

// GetProjectsByUser retrieves all projects for a user
func (s *ProjectService) GetProjectsByUser(ctx context.Context, userID string) ([]*models.Project, error) {
	_, err := s.store.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	projects, err := s.store.GetProjectsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get projects: %w", err)
	}

	publicProjects := make([]*models.Project, len(projects))
	for i, project := range projects {
		publicProjects[i] = project.ToPublic()
	}

	return publicProjects, nil
}

// GetActiveProjectsByUser retrieves only active projects for a user
func (s *ProjectService) GetActiveProjectsByUser(ctx context.Context, userID string) ([]*models.Project, error) {
	_, err := s.store.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	projects, err := s.store.GetActiveProjectsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active projects: %w", err)
	}

	publicProjects := make([]*models.Project, len(projects))
	for i, project := range projects {
		publicProjects[i] = project.ToPublic()
	}

	return publicProjects, nil
}

// UpdateProject updates a project
func (s *ProjectService) UpdateProject(ctx context.Context, userID, projectID string, req *models.UpdateProjectRequest) (*models.Project, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	project, err := s.store.GetProjectByID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	if !project.IsOwnedBy(userID) {
		return nil, errors.New("access denied: project does not belong to user")
	}

	if req.Name != "" && req.Name != project.Name {
		exists, err := s.store.ProjectExistsByUserIDAndName(ctx, userID, req.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to check project name existence: %w", err)
		}
		if exists {
			return nil, errors.New("project with this name already exists for this user")
		}
	}

	updatedProject, err := s.store.UpdateProject(ctx, projectID, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}

	return updatedProject.ToPublic(), nil
}

// UpdateProjectStatus updates only the project status
func (s *ProjectService) UpdateProjectStatus(ctx context.Context, userID, projectID, status string) (*models.Project, error) {
	validStatuses := []string{"active", "inactive", "archived"}
	isValid := false
	for _, validStatus := range validStatuses {
		if status == validStatus {
			isValid = true
			break
		}
	}
	if !isValid {
		return nil, errors.New("invalid status: must be one of active, inactive, archived")
	}

	project, err := s.store.GetProjectByID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	if !project.IsOwnedBy(userID) {
		return nil, errors.New("access denied: project does not belong to user")
	}

	updatedProject, err := s.store.UpdateProjectStatus(ctx, projectID, status)
	if err != nil {
		return nil, fmt.Errorf("failed to update project status: %w", err)
	}

	return updatedProject.ToPublic(), nil
}

// ArchiveProject archives a project
func (s *ProjectService) ArchiveProject(ctx context.Context, userID, projectID string) (*models.Project, error) {
	project, err := s.store.GetProjectByID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	if !project.IsOwnedBy(userID) {
		return nil, errors.New("access denied: project does not belong to user")
	}

	archivedProject, err := s.store.ArchiveProject(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to archive project: %w", err)
	}

	return archivedProject.ToPublic(), nil
}

// ActivateProject activates a project
func (s *ProjectService) ActivateProject(ctx context.Context, userID, projectID string) (*models.Project, error) {
	project, err := s.store.GetProjectByID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	if !project.IsOwnedBy(userID) {
		return nil, errors.New("access denied: project does not belong to user")
	}

	activatedProject, err := s.store.ActivateProject(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to activate project: %w", err)
	}

	return activatedProject.ToPublic(), nil
}

// DeleteProject deletes a project
func (s *ProjectService) DeleteProject(ctx context.Context, userID, projectID string) error {
	project, err := s.store.GetProjectByID(ctx, projectID)
	if err != nil {
		return fmt.Errorf("project not found: %w", err)
	}

	if !project.IsOwnedBy(userID) {
		return errors.New("access denied: project does not belong to user")
	}

	if err := s.store.DeleteProject(ctx, projectID); err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	return nil
}

// SearchProjectsByUser searches projects for a specific user
func (s *ProjectService) SearchProjectsByUser(ctx context.Context, userID, query string, limit, offset int) ([]*models.Project, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	_, err := s.store.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	projects, err := s.store.SearchProjectsByUserID(ctx, userID, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search projects: %w", err)
	}

	publicProjects := make([]*models.Project, len(projects))
	for i, project := range projects {
		publicProjects[i] = project.ToPublic()
	}

	return publicProjects, nil
}

// ListProjectsByUser retrieves a paginated list of projects for a user
func (s *ProjectService) ListProjectsByUser(ctx context.Context, userID string, limit, offset int) ([]*models.Project, int64, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	_, err := s.store.GetUserByID(ctx, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("user not found: %w", err)
	}

	projects, err := s.store.GetProjectsByUserID(ctx, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get projects: %w", err)
	}

	total, err := s.store.CountProjectsByUserID(ctx, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count projects: %w", err)
	}

	start := offset
	end := offset + limit
	if start >= len(projects) {
		return []*models.Project{}, total, nil
	}
	if end > len(projects) {
		end = len(projects)
	}

	paginatedProjects := projects[start:end]

	publicProjects := make([]*models.Project, len(paginatedProjects))
	for i, project := range paginatedProjects {
		publicProjects[i] = project.ToPublic()
	}

	return publicProjects, total, nil
}

