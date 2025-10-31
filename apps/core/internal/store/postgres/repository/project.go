package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dopeCape/kova/internal/models"
	"github.com/dopeCape/kova/internal/store/postgres/generated"
	"github.com/jackc/pgx/v5/pgtype"
)

// CreateProject creates a new project
func (s *Store) CreateProject(ctx context.Context, project *models.Project) error {
	// Marshal env variables to JSON
	envJSON, err := json.Marshal(project.EnvVariables)
	if err != nil {
		return fmt.Errorf("failed to marshal env variables: %w", err)
	}

	domain := pgtype.Text{String: project.Domain, Valid: project.Domain != ""}
	port := pgtype.Int4{Int32: int32(project.Port), Valid: project.Port > 0}

	params := generated.CreateProjectParams{
		Name:         project.Name,
		UserID:       project.UserID,
		RepoID:       project.RepoID,
		RepoName:     project.RepoName,
		RepoFullName: project.RepoFullName,
		RepoUrl:      project.RepoURL,
		RepoBranch:   project.RepoBranch,
		EnvVariables: envJSON,
		Domain:       domain,
		Port:         port,
	}

	dbProject, err := s.queries.CreateProject(ctx, params)
	if err != nil {
		return err
	}

	*project = s.toDomainProject(dbProject)
	return nil
}

// GetProjectByID retrieves a project by ID
func (s *Store) GetProjectByID(ctx context.Context, id string) (*models.Project, error) {
	dbProject, err := s.queries.GetProjectByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrProjectNotFound
		}
		return nil, err
	}

	project := s.toDomainProject(dbProject)
	return &project, nil
}

// GetProjectByUserIDAndName retrieves a project by user ID and name
func (s *Store) GetProjectByUserIDAndName(ctx context.Context, userID, name string) (*models.Project, error) {
	params := generated.GetProjectByUserIDAndNameParams{
		UserID: userID,
		Name:   name,
	}

	dbProject, err := s.queries.GetProjectByUserIDAndName(ctx, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrProjectNotFound
		}
		return nil, err
	}

	project := s.toDomainProject(generated.Project(dbProject))
	return &project, nil
}

// GetProjectsByUserID retrieves all projects for a user
func (s *Store) GetProjectsByUserID(ctx context.Context, userID string) ([]*models.Project, error) {
	dbProjects, err := s.queries.GetProjectsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	projects := make([]*models.Project, len(dbProjects))
	for i, dbProject := range dbProjects {
		project := s.toDomainProject(dbProject)
		projects[i] = &project
	}

	return projects, nil
}

// GetProjectsByUserIDAndStatus retrieves projects for a user filtered by status
func (s *Store) GetProjectsByUserIDAndStatus(ctx context.Context, userID, status string) ([]*models.Project, error) {
	params := generated.GetProjectsByUserIDAndStatusParams{
		UserID: userID,
		Status: status,
	}

	dbProjects, err := s.queries.GetProjectsByUserIDAndStatus(ctx, params)
	if err != nil {
		return nil, err
	}

	projects := make([]*models.Project, len(dbProjects))
	for i, dbProject := range dbProjects {
		project := s.toDomainProject(dbProject)
		projects[i] = &project
	}

	return projects, nil
}

// GetProjectsByRepoID retrieves all projects associated with a repository
func (s *Store) GetProjectsByRepoID(ctx context.Context, repoID int64) ([]*models.Project, error) {
	dbProjects, err := s.queries.GetProjectsByRepoID(ctx, repoID)
	if err != nil {
		return nil, err
	}

	projects := make([]*models.Project, len(dbProjects))
	for i, dbProject := range dbProjects {
		project := s.toDomainProject(dbProject)
		projects[i] = &project
	}

	return projects, nil
}

// GetActiveProjectsByUserID retrieves only active projects for a user
func (s *Store) GetActiveProjectsByUserID(ctx context.Context, userID string) ([]*models.Project, error) {
	dbProjects, err := s.queries.GetActiveProjectsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	projects := make([]*models.Project, len(dbProjects))
	for i, dbProject := range dbProjects {
		project := s.toDomainProject(dbProject)
		projects[i] = &project
	}

	return projects, nil
}

// UpdateProject updates a project with partial data
func (s *Store) UpdateProject(ctx context.Context, projectID string, req *models.UpdateProjectRequest) (*models.Project, error) {
	// Get current project to use existing values for unchanged fields
	currentProject, err := s.queries.GetProjectByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	// Use existing values if not provided in request
	name := currentProject.Name
	if req.Name != "" {
		name = req.Name
	}

	branch := currentProject.RepoBranch
	if req.RepoBranch != "" {
		branch = req.RepoBranch
	}

	status := currentProject.Status
	if req.Status != "" {
		status = req.Status
	}

	domain := currentProject.Domain
	if req.Domain != "" {
		domain = pgtype.Text{String: req.Domain, Valid: true}
	}

	params := generated.UpdateProjectParams{
		ID:         projectID,
		Name:       name,
		RepoBranch: branch,
		Status:     status,
		Domain:     domain,
	}

	dbProject, err := s.queries.UpdateProject(ctx, params)
	if err != nil {
		return nil, err
	}

	project := s.toDomainProject(dbProject)
	return &project, nil
}

// UpdateProjectStatus updates only the project status
func (s *Store) UpdateProjectStatus(ctx context.Context, projectID, status string) (*models.Project, error) {
	params := generated.UpdateProjectStatusParams{
		ID:     projectID,
		Status: status,
	}

	dbProject, err := s.queries.UpdateProjectStatus(ctx, params)
	if err != nil {
		return nil, err
	}

	project := s.toDomainProject(dbProject)
	return &project, nil
}

// UpdateProjectDeploymentStatus updates only the deployment status
func (s *Store) UpdateProjectDeploymentStatus(ctx context.Context, projectID, status string) (*models.Project, error) {
	params := generated.UpdateProjectDeploymentStatusParams{
		ID:               projectID,
		DeploymentStatus: status,
	}

	dbProject, err := s.queries.UpdateProjectDeploymentStatus(ctx, params)
	if err != nil {
		return nil, err
	}

	project := s.toDomainProject(dbProject)
	return &project, nil
}

// UpdateProjectBranch updates only the repository branch
func (s *Store) UpdateProjectBranch(ctx context.Context, projectID, branch string) (*models.Project, error) {
	params := generated.UpdateProjectBranchParams{
		ID:         projectID,
		RepoBranch: branch,
	}

	dbProject, err := s.queries.UpdateProjectBranch(ctx, params)
	if err != nil {
		return nil, err
	}

	project := s.toDomainProject(dbProject)
	return &project, nil
}

// UpdateProjectPort updates the port assigned to a project
func (s *Store) UpdateProjectPort(ctx context.Context, projectID string, port int) error {
	params := generated.UpdateProjectPortParams{
		ID:   projectID,
		Port: pgtype.Int4{Int32: int32(port), Valid: true},
	}

	return s.queries.UpdateProjectPort(ctx, params)
}

// ArchiveProject sets project status to archived
func (s *Store) ArchiveProject(ctx context.Context, projectID string) (*models.Project, error) {
	dbProject, err := s.queries.ArchiveProject(ctx, projectID)
	if err != nil {
		return nil, err
	}

	project := s.toDomainProject(dbProject)
	return &project, nil
}

// ActivateProject sets project status to active
func (s *Store) ActivateProject(ctx context.Context, projectID string) (*models.Project, error) {
	dbProject, err := s.queries.ActivateProject(ctx, projectID)
	if err != nil {
		return nil, err
	}

	project := s.toDomainProject(dbProject)
	return &project, nil
}

// DeleteProject deletes a project by ID
func (s *Store) DeleteProject(ctx context.Context, id string) error {
	return s.queries.DeleteProject(ctx, id)
}

// DeleteProjectsByUserID deletes all projects for a user
func (s *Store) DeleteProjectsByUserID(ctx context.Context, userID string) error {
	return s.queries.DeleteProjectsByUserID(ctx, userID)
}

// ListProjects retrieves a paginated list of all projects
func (s *Store) ListProjects(ctx context.Context, limit, offset int) ([]*models.Project, error) {
	params := generated.ListProjectsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	dbProjects, err := s.queries.ListProjects(ctx, params)
	if err != nil {
		return nil, err
	}

	projects := make([]*models.Project, len(dbProjects))
	for i, dbProject := range dbProjects {
		project := s.toDomainProject(dbProject)
		projects[i] = &project
	}

	return projects, nil
}

// ListProjectsByStatus retrieves a paginated list of projects by status
func (s *Store) ListProjectsByStatus(ctx context.Context, status string, limit, offset int) ([]*models.Project, error) {
	params := generated.ListProjectsByStatusParams{
		Status: status,
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	dbProjects, err := s.queries.ListProjectsByStatus(ctx, params)
	if err != nil {
		return nil, err
	}

	projects := make([]*models.Project, len(dbProjects))
	for i, dbProject := range dbProjects {
		project := s.toDomainProject(dbProject)
		projects[i] = &project
	}

	return projects, nil
}

// CountProjects returns the total number of projects
func (s *Store) CountProjects(ctx context.Context) (int64, error) {
	return s.queries.CountProjects(ctx)
}

// CountProjectsByUserID returns the number of projects for a user
func (s *Store) CountProjectsByUserID(ctx context.Context, userID string) (int64, error) {
	return s.queries.CountProjectsByUserID(ctx, userID)
}

// CountProjectsByStatus returns the number of projects with a specific status
func (s *Store) CountProjectsByStatus(ctx context.Context, status string) (int64, error) {
	return s.queries.CountProjectsByStatus(ctx, status)
}

// ProjectExistsByUserIDAndName checks if a project with the given name exists for a user
func (s *Store) ProjectExistsByUserIDAndName(ctx context.Context, userID, name string) (bool, error) {
	params := generated.ProjectExistsByUserIDAndNameParams{
		UserID: userID,
		Name:   name,
	}
	return s.queries.ProjectExistsByUserIDAndName(ctx, params)
}

// ProjectExistsByID checks if a project exists by ID
func (s *Store) ProjectExistsByID(ctx context.Context, id string) (bool, error) {
	return s.queries.ProjectExistsByID(ctx, id)
}

// SearchProjects searches all projects by query string
func (s *Store) SearchProjects(ctx context.Context, query string, limit, offset int) ([]*models.Project, error) {
	params := generated.SearchProjectsParams{
		Column1: pgtype.Text{String: query, Valid: true},
		Limit:   int32(limit),
		Offset:  int32(offset),
	}

	dbProjects, err := s.queries.SearchProjects(ctx, params)
	if err != nil {
		return nil, err
	}

	projects := make([]*models.Project, len(dbProjects))
	for i, dbProject := range dbProjects {
		project := s.toDomainProject(dbProject)
		projects[i] = &project
	}

	return projects, nil
}

// SearchProjectsByUserID searches projects for a specific user
func (s *Store) SearchProjectsByUserID(ctx context.Context, userID, query string, limit, offset int) ([]*models.Project, error) {
	params := generated.SearchProjectsByUserIDParams{
		UserID:  userID,
		Column2: pgtype.Text{String: query, Valid: true},
		Limit:   int32(limit),
		Offset:  int32(offset),
	}

	dbProjects, err := s.queries.SearchProjectsByUserID(ctx, params)
	if err != nil {
		return nil, err
	}

	projects := make([]*models.Project, len(dbProjects))
	for i, dbProject := range dbProjects {
		project := s.toDomainProject(dbProject)
		projects[i] = &project
	}

	return projects, nil
}

// GetUsedPorts retrieves all ports currently assigned to projects
func (s *Store) GetUsedPorts(ctx context.Context) ([]int, error) {
	ports, err := s.queries.GetUsedPorts(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]int, 0, len(ports))
	for _, p := range ports {
		if p.Valid {
			result = append(result, int(p.Int32))
		}
	}

	return result, nil
}

// toDomainProject converts a database project to a domain model
func (s *Store) toDomainProject(dbProject generated.Project) models.Project {
	// Unmarshal env variables from JSON
	var envVars []models.EnvironmentVariable
	if err := json.Unmarshal(dbProject.EnvVariables, &envVars); err != nil {
		envVars = []models.EnvironmentVariable{}
	}

	domain := ""
	if dbProject.Domain.Valid {
		domain = dbProject.Domain.String
	}

	port := 0
	if dbProject.Port.Valid {
		port = int(dbProject.Port.Int32)
	}

	return models.Project{
		ID:               dbProject.ID,
		Name:             dbProject.Name,
		UserID:           dbProject.UserID,
		RepoID:           dbProject.RepoID,
		RepoName:         dbProject.RepoName,
		RepoFullName:     dbProject.RepoFullName,
		RepoURL:          dbProject.RepoUrl,
		RepoBranch:       dbProject.RepoBranch,
		Status:           dbProject.Status,
		DeploymentStatus: dbProject.DeploymentStatus,
		Domain:           domain,
		Port:             port,
		EnvVariables:     envVars,
		CreatedAt:        dbProject.CreatedAt,
		UpdatedAt:        dbProject.UpdatedAt,
	}
}

// Error definitions
var (
	ErrProjectNotFound = errors.New("project not found")
)
