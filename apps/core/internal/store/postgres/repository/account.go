package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dopeCape/kova/internal/models"
	"github.com/dopeCape/kova/internal/store/postgres/generated"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Store) CreateAccount(ctx context.Context, account *models.Account) error {
	params := generated.CreateAccountParams{
		UserID:         account.UserID,
		GithubUsername: account.GithubUsername,
		GithubID:       account.GithubID,
		AvatarUrl:      pgtype.Text{String: account.AvatarURL, Valid: account.AvatarURL != ""},
		AccessToken:    account.AccessToken,
	}

	dbAccount, err := s.queries.CreateAccount(ctx, params)
	if err != nil {
		return err
	}

	*account = models.Account{
		ID:             dbAccount.ID,
		UserID:         dbAccount.UserID,
		GithubUsername: dbAccount.GithubUsername,
		GithubID:       dbAccount.GithubID,
		AvatarURL:      dbAccount.AvatarUrl.String,
		AccessToken:    account.AccessToken,
		CreatedAt:      dbAccount.CreatedAt,
		UpdatedAt:      dbAccount.UpdatedAt,
	}

	return nil
}

func (s *Store) GetAccountByID(ctx context.Context, id string) (*models.Account, error) {
	dbAccount, err := s.queries.GetAccountByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrAccountNotFound
		}
		return nil, err
	}

	account := s.toDomainAccount(dbAccount)
	return &account, nil
}

func (s *Store) GetAccountByGithubID(ctx context.Context, githubID int64) (*models.Account, error) {
	dbAccount, err := s.queries.GetAccountByGithubID(ctx, githubID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrAccountNotFound
		}
		return nil, err
	}

	account := s.toDomainAccount(dbAccount)
	return &account, nil
}

func (s *Store) GetAccountByGithubUsername(ctx context.Context, githubUsername string) (*models.Account, error) {
	dbAccount, err := s.queries.GetAccountByGithubUsername(ctx, githubUsername)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrAccountNotFound
		}
		return nil, err
	}

	account := s.toDomainAccount(dbAccount)
	return &account, nil
}

func (s *Store) GetAccountsByUserID(ctx context.Context, userID string) ([]*models.Account, error) {
	dbAccounts, err := s.queries.GetAccountsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	accounts := make([]*models.Account, len(dbAccounts))
	for i, dbAccount := range dbAccounts {
		account := models.Account{
			ID:             dbAccount.ID,
			UserID:         dbAccount.UserID,
			GithubUsername: dbAccount.GithubUsername,
			GithubID:       dbAccount.GithubID,
			AvatarURL:      dbAccount.AvatarUrl.String,
			CreatedAt:      dbAccount.CreatedAt,
			UpdatedAt:      dbAccount.UpdatedAt,
			// AccessToken intentionally omitted for security
		}
		accounts[i] = &account
	}

	return accounts, nil
}

func (s *Store) GetAccountsByUserIDWithTokens(ctx context.Context, userID string) ([]*models.Account, error) {
	dbAccounts, err := s.queries.GetAccountsByUserIDWithTokens(ctx, userID)
	if err != nil {
		return nil, err
	}

	accounts := make([]*models.Account, len(dbAccounts))
	for i, dbAccount := range dbAccounts {
		account := s.toDomainAccount(dbAccount)
		accounts[i] = &account
	}

	return accounts, nil
}

func (s *Store) UpdateAccount(ctx context.Context, accountID string, req *models.UpdateAccountRequest) (*models.Account, error) {
	params := generated.UpdateAccountParams{
		ID:             accountID,
		GithubUsername: req.GithubUsername,
		AvatarUrl:      pgtype.Text{String: req.AvatarURL, Valid: req.AvatarURL != ""},
	}

	dbAccount, err := s.queries.UpdateAccount(ctx, params)
	if err != nil {
		return nil, err
	}

	account := &models.Account{
		ID:             dbAccount.ID,
		UserID:         dbAccount.UserID,
		GithubUsername: dbAccount.GithubUsername,
		GithubID:       dbAccount.GithubID,
		AvatarURL:      dbAccount.AvatarUrl.String,
		CreatedAt:      dbAccount.CreatedAt,
		UpdatedAt:      dbAccount.UpdatedAt,
	}

	return account, nil
}

func (s *Store) UpdateAccountToken(ctx context.Context, accountID, accessToken string) (*models.Account, error) {
	params := generated.UpdateAccountTokenParams{
		ID:          accountID,
		AccessToken: accessToken,
	}

	dbAccount, err := s.queries.UpdateAccountToken(ctx, params)
	if err != nil {
		return nil, err
	}

	account := &models.Account{
		ID:             dbAccount.ID,
		UserID:         dbAccount.UserID,
		GithubUsername: dbAccount.GithubUsername,
		GithubID:       dbAccount.GithubID,
		AvatarURL:      dbAccount.AvatarUrl.String,
		AccessToken:    accessToken,
		CreatedAt:      dbAccount.CreatedAt,
		UpdatedAt:      dbAccount.UpdatedAt,
	}

	return account, nil
}

func (s *Store) UpdateAccountByGithubID(ctx context.Context, githubID int64, req *models.UpdateAccountByGithubIDRequest) (*models.Account, error) {
	params := generated.UpdateAccountByGithubIDParams{
		GithubID:       githubID,
		GithubUsername: req.GithubUsername,
		AvatarUrl:      pgtype.Text{String: req.AvatarURL, Valid: req.AvatarURL != ""},
		AccessToken:    req.AccessToken,
	}

	dbAccount, err := s.queries.UpdateAccountByGithubID(ctx, params)
	if err != nil {
		return nil, err
	}

	account := &models.Account{
		ID:             dbAccount.ID,
		UserID:         dbAccount.UserID,
		GithubUsername: dbAccount.GithubUsername,
		GithubID:       dbAccount.GithubID,
		AvatarURL:      dbAccount.AvatarUrl.String,
		AccessToken:    req.AccessToken,
		CreatedAt:      dbAccount.CreatedAt,
		UpdatedAt:      dbAccount.UpdatedAt,
	}

	return account, nil
}

func (s *Store) DeleteAccount(ctx context.Context, id string) error {
	return s.queries.DeleteAccount(ctx, id)
}

func (s *Store) DeleteAccountByGithubID(ctx context.Context, githubID int64) error {
	return s.queries.DeleteAccountByGithubID(ctx, githubID)
}

func (s *Store) DeleteAccountsByUserID(ctx context.Context, userID string) error {
	return s.queries.DeleteAccountsByUserID(ctx, userID)
}

func (s *Store) ListAccounts(ctx context.Context, limit, offset int) ([]*models.Account, error) {
	params := generated.ListAccountsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	dbAccounts, err := s.queries.ListAccounts(ctx, params)
	if err != nil {
		return nil, err
	}

	accounts := make([]*models.Account, len(dbAccounts))
	for i, dbAccount := range dbAccounts {
		account := models.Account{
			ID:             dbAccount.ID,
			UserID:         dbAccount.UserID,
			GithubUsername: dbAccount.GithubUsername,
			GithubID:       dbAccount.GithubID,
			AvatarURL:      dbAccount.AvatarUrl.String,
			CreatedAt:      dbAccount.CreatedAt,
			UpdatedAt:      dbAccount.UpdatedAt,
		}
		accounts[i] = &account
	}

	return accounts, nil
}

func (s *Store) CountAccounts(ctx context.Context) (int64, error) {
	return s.queries.CountAccounts(ctx)
}

func (s *Store) CountAccountsByUserID(ctx context.Context, userID string) (int64, error) {
	return s.queries.CountAccountsByUserID(ctx, userID)
}

func (s *Store) AccountExistsByGithubID(ctx context.Context, githubID int64) (bool, error) {
	return s.queries.AccountExistsByGithubID(ctx, githubID)
}

func (s *Store) AccountExistsByGithubUsername(ctx context.Context, githubUsername string) (bool, error) {
	return s.queries.AccountExistsByGithubUsername(ctx, githubUsername)
}
func (s *Store) AccountExistsByUserIDAndGithubID(ctx context.Context, userId string, githubUsername int64) (bool, error) {
	return s.queries.AccountExistsByUserIDAndGithubID(ctx, generated.AccountExistsByUserIDAndGithubIDParams{UserID: userId, GithubID: githubUsername})
}

func (s *Store) AccountExistsForUser(ctx context.Context, userID string, githubID int64) (bool, error) {
	params := generated.AccountExistsForUserParams{
		UserID:   userID,
		GithubID: githubID,
	}
	return s.queries.AccountExistsForUser(ctx, params)
}

func (s *Store) SearchAccounts(ctx context.Context, query string, limit, offset int) ([]*models.Account, error) {
	params := generated.SearchAccountsParams{
		Column1: pgtype.Text{String: query},
		Limit:   int32(limit),
		Offset:  int32(offset),
	}

	dbAccounts, err := s.queries.SearchAccounts(ctx, params)
	if err != nil {
		return nil, err
	}

	accounts := make([]*models.Account, len(dbAccounts))
	for i, dbAccount := range dbAccounts {
		account := models.Account{
			ID:             dbAccount.ID,
			UserID:         dbAccount.UserID,
			GithubUsername: dbAccount.GithubUsername,
			GithubID:       dbAccount.GithubID,
			AvatarURL:      dbAccount.AvatarUrl.String,
			CreatedAt:      dbAccount.CreatedAt,
			UpdatedAt:      dbAccount.UpdatedAt,
		}
		accounts[i] = &account
	}

	return accounts, nil
}

func (s *Store) SearchAccountsByUserID(ctx context.Context, userID, query string, limit, offset int) ([]*models.Account, error) {
	params := generated.SearchAccountsByUserIDParams{
		UserID:  userID,
		Column2: pgtype.Text{String: query},
		Limit:   int32(limit),
		Offset:  int32(offset),
	}

	dbAccounts, err := s.queries.SearchAccountsByUserID(ctx, params)
	if err != nil {
		return nil, err
	}

	accounts := make([]*models.Account, len(dbAccounts))
	for i, dbAccount := range dbAccounts {
		account := models.Account{
			ID:             dbAccount.ID,
			UserID:         dbAccount.UserID,
			GithubUsername: dbAccount.GithubUsername,
			GithubID:       dbAccount.GithubID,
			AvatarURL:      dbAccount.AvatarUrl.String,
			CreatedAt:      dbAccount.CreatedAt,
			UpdatedAt:      dbAccount.UpdatedAt,
		}
		accounts[i] = &account
	}

	return accounts, nil
}

func (s *Store) GetAccountWithUser(ctx context.Context, accountID string) (*models.AccountWithUser, error) {
	dbAccountWithUser, err := s.queries.GetAccountWithUser(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrAccountNotFound
		}
		return nil, err
	}

	accountWithUser := &models.AccountWithUser{
		ID:             dbAccountWithUser.ID,
		UserID:         dbAccountWithUser.UserID,
		GithubUsername: dbAccountWithUser.GithubUsername,
		GithubID:       dbAccountWithUser.GithubID,
		AvatarURL:      dbAccountWithUser.AvatarUrl.String,
		AccessToken:    dbAccountWithUser.AccessToken,
		CreatedAt:      dbAccountWithUser.CreatedAt,
		UpdatedAt:      dbAccountWithUser.UpdatedAt,
		Username:       dbAccountWithUser.Username,
		Email:          dbAccountWithUser.Email,
	}

	return accountWithUser, nil
}

func (s *Store) toDomainAccount(dbAccount generated.Account) models.Account {
	return models.Account{
		ID:             dbAccount.ID,
		UserID:         dbAccount.UserID,
		GithubUsername: dbAccount.GithubUsername,
		GithubID:       dbAccount.GithubID,
		AvatarURL:      dbAccount.AvatarUrl.String,
		AccessToken:    dbAccount.AccessToken,
		CreatedAt:      dbAccount.CreatedAt,
		UpdatedAt:      dbAccount.UpdatedAt,
	}
}

var (
	ErrAccountNotFound = errors.New("account not found")
)
