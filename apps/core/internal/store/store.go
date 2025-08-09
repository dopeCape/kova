package store

import (
	"context"

	"github.com/dopeCape/kova/internal/models"
)

type Store interface {
	UserStore
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
