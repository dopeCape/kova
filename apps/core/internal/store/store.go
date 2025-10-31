package store

import (
	"context"

	"github.com/dopeCape/kova/internal/models"
)

type Store interface {
	UserStore
	AccountStore
	ProjectStore
	Ping(ctx context.Context) error
}

type UserStore interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetUserByEmailOrUsername(ctx context.Context, login string) (*models.User, error)
	UpdateUser(ctx context.Context, userID string, req *models.UpdateUserRequest) (*models.User, error)
	UpdateUserPassword(ctx context.Context, userID, passwordHash string) (*models.User, error)
	DeleteUser(ctx context.Context, id string) error
	ListUsers(ctx context.Context, limit, offset int) ([]*models.User, error)
	CountUsers(ctx context.Context) (int64, error)
	UserExistsByEmail(ctx context.Context, email string) (bool, error)
	UserExistsByUsername(ctx context.Context, username string) (bool, error)
	SearchUsers(ctx context.Context, query string, limit, offset int) ([]*models.User, error)
}

type AccountStore interface {
	CreateAccount(ctx context.Context, account *models.Account) error
	GetAccountByID(ctx context.Context, id string) (*models.Account, error)
	GetAccountByGithubID(ctx context.Context, githubID int64) (*models.Account, error)
	GetAccountByGithubUsername(ctx context.Context, githubUsername string) (*models.Account, error)
	GetAccountsByUserID(ctx context.Context, userID string) ([]*models.Account, error)
	GetAccountsByUserIDWithTokens(ctx context.Context, userID string) ([]*models.Account, error)
	UpdateAccount(ctx context.Context, accountID string, req *models.UpdateAccountRequest) (*models.Account, error)
	UpdateAccountToken(ctx context.Context, accountID, accessToken string) (*models.Account, error)
	UpdateAccountByGithubID(ctx context.Context, githubID int64, req *models.UpdateAccountByGithubIDRequest) (*models.Account, error)
	DeleteAccount(ctx context.Context, id string) error
	DeleteAccountsByUserID(ctx context.Context, userID string) error
	AccountExistsByGithubID(ctx context.Context, githubID int64) (bool, error)
	AccountExistsByUserIDAndGithubID(ctx context.Context, userID string, githubID int64) (bool, error)
	AccountExistsForUser(ctx context.Context, userID string, githubID int64) (bool, error) // NEW METHOD
}

type ProjectStore interface {
	CreateProject(ctx context.Context, project *models.Project) error
	GetProjectByID(ctx context.Context, id string) (*models.Project, error)
	GetProjectByUserIDAndName(ctx context.Context, userID, name string) (*models.Project, error)
	GetProjectsByUserID(ctx context.Context, userID string) ([]*models.Project, error)
	GetProjectsByUserIDAndStatus(ctx context.Context, userID, status string) ([]*models.Project, error)
	GetProjectsByRepoID(ctx context.Context, repoID int64) ([]*models.Project, error)
	GetActiveProjectsByUserID(ctx context.Context, userID string) ([]*models.Project, error)
	UpdateProject(ctx context.Context, projectID string, req *models.UpdateProjectRequest) (*models.Project, error)
	UpdateProjectStatus(ctx context.Context, projectID, status string) (*models.Project, error)
	UpdateProjectDeploymentStatus(ctx context.Context, projectID, status string) (*models.Project, error)
	UpdateProjectBranch(ctx context.Context, projectID, branch string) (*models.Project, error)
	UpdateProjectPort(ctx context.Context, projectID string, port int) error
	ArchiveProject(ctx context.Context, projectID string) (*models.Project, error)
	ActivateProject(ctx context.Context, projectID string) (*models.Project, error)
	DeleteProject(ctx context.Context, id string) error
	DeleteProjectsByUserID(ctx context.Context, userID string) error
	ListProjects(ctx context.Context, limit, offset int) ([]*models.Project, error)
	ListProjectsByStatus(ctx context.Context, status string, limit, offset int) ([]*models.Project, error)
	CountProjects(ctx context.Context) (int64, error)
	CountProjectsByUserID(ctx context.Context, userID string) (int64, error)
	CountProjectsByStatus(ctx context.Context, status string) (int64, error)
	ProjectExistsByUserIDAndName(ctx context.Context, userID, name string) (bool, error)
	ProjectExistsByID(ctx context.Context, id string) (bool, error)
	SearchProjects(ctx context.Context, query string, limit, offset int) ([]*models.Project, error)
	SearchProjectsByUserID(ctx context.Context, userID, query string, limit, offset int) ([]*models.Project, error)
	GetUsedPorts(ctx context.Context) ([]int, error)
}
