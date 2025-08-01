package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dopeCape/kova/internal/models"
	"github.com/dopeCape/kova/internal/store/postgres/generated"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Store) CreateUser(ctx context.Context, user *models.User) error {
	params := generated.CreateUserParams{
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
	}

	dbUser, err := s.queries.CreateUser(ctx, params)
	if err != nil {
		return err
	}

	// Convert CreateUserRow to User model
	*user = models.User{
		ID:           dbUser.ID,
		Username:     dbUser.Username,
		Email:        dbUser.Email,
		PasswordHash: user.PasswordHash,
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
	}

	return nil
}

func (s *Store) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	dbUser, err := s.queries.GetUserByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	user := s.toDomainUser(dbUser)
	return &user, nil
}

func (s *Store) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	dbUser, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	user := s.toDomainUser(dbUser)
	return &user, nil
}

func (s *Store) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	dbUser, err := s.queries.GetUserByUsername(ctx, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	user := s.toDomainUser(dbUser)
	return &user, nil
}

func (s *Store) GetUserByEmailOrUsername(ctx context.Context, login string) (*models.User, error) {
	dbUser, err := s.queries.GetUserByEmailOrUsername(ctx, login)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	user := s.toDomainUser(dbUser)
	return &user, nil
}

func (s *Store) UpdateUser(ctx context.Context, userID string, req *models.UpdateUserRequest) (*models.User, error) {
	params := generated.UpdateUserParams{
		ID:       userID,
		Username: req.Username,
		Email:    req.Email,
	}

	dbUser, err := s.queries.UpdateUser(ctx, params)
	if err != nil {
		return nil, err
	}

	// Convert UpdateUserRow to User model
	user := &models.User{
		ID:        dbUser.ID,
		Username:  dbUser.Username,
		Email:     dbUser.Email,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
	}

	return user, nil
}

func (s *Store) UpdateUserPassword(ctx context.Context, userID, passwordHash string) (*models.User, error) {
	params := generated.UpdateUserPasswordParams{
		ID:           userID,
		PasswordHash: passwordHash,
	}

	dbUser, err := s.queries.UpdateUserPassword(ctx, params)
	if err != nil {
		return nil, err
	}

	// Convert UpdateUserPasswordRow to User model
	user := &models.User{
		ID:           dbUser.ID,
		Username:     dbUser.Username,
		Email:        dbUser.Email,
		PasswordHash: passwordHash,
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
	}

	return user, nil
}

func (s *Store) DeleteUser(ctx context.Context, id string) error {
	return s.queries.DeleteUser(ctx, id)
}

func (s *Store) ListUsers(ctx context.Context, limit, offset int) ([]*models.User, error) {
	params := generated.ListUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	dbUsers, err := s.queries.ListUsers(ctx, params)
	if err != nil {
		return nil, err
	}

	users := make([]*models.User, len(dbUsers))
	for i, dbUser := range dbUsers {
		user := models.User{
			ID:        dbUser.ID,
			Username:  dbUser.Username,
			Email:     dbUser.Email,
			CreatedAt: dbUser.CreatedAt,
			UpdatedAt: dbUser.UpdatedAt,
		}
		users[i] = &user
	}

	return users, nil
}

func (s *Store) CountUsers(ctx context.Context) (int64, error) {
	return s.queries.CountUsers(ctx)
}

func (s *Store) UserExistsByEmail(ctx context.Context, email string) (bool, error) {
	return s.queries.UserExistsByEmail(ctx, email)
}

func (s *Store) UserExistsByUsername(ctx context.Context, username string) (bool, error) {
	return s.queries.UserExistsByUsername(ctx, username)
}

func (s *Store) SearchUsers(ctx context.Context, query string, limit, offset int) ([]*models.User, error) {
	params := generated.SearchUsersParams{
		Column1: pgtype.Text{String: query},
		Limit:   int32(limit),
		Offset:  int32(offset),
	}

	dbUsers, err := s.queries.SearchUsers(ctx, params)
	if err != nil {
		return nil, err
	}

	users := make([]*models.User, len(dbUsers))
	for i, dbUser := range dbUsers {
		user := models.User{
			ID:        dbUser.ID,
			Username:  dbUser.Username,
			Email:     dbUser.Email,
			CreatedAt: dbUser.CreatedAt,
			UpdatedAt: dbUser.UpdatedAt,
		}
		users[i] = &user
	}

	return users, nil
}

func (s *Store) toDomainUser(dbUser generated.User) models.User {
	return models.User{
		ID:           dbUser.ID,
		Username:     dbUser.Username,
		Email:        dbUser.Email,
		PasswordHash: dbUser.PasswordHash,
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
	}
}

var (
	ErrUserNotFound = errors.New("user not found")
)
