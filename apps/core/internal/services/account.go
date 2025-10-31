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

type AccountService struct {
	store         store.Store
	validator     *validator.Validate
	githubService *GitHubService
}

func NewAccountService(store store.Store, githubService *GitHubService) *AccountService {
	return &AccountService{
		store:         store,
		validator:     validator.New(),
		githubService: githubService,
	}
}

// GetAccountsByUserID retrieves all accounts for a user
func (s *AccountService) GetAccountsByUserID(ctx context.Context, userID string) ([]*models.Account, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	// Check if user exists
	_, err := s.store.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Get accounts (without tokens for security)
	accounts, err := s.store.GetAccountsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get accounts: %w", err)
	}

	// Convert to public accounts (tokens are already excluded from GetAccountsByUserID)
	publicAccounts := make([]*models.Account, len(accounts))
	for i, account := range accounts {
		publicAccounts[i] = account.ToPublic()
	}

	return publicAccounts, nil
}

// CreateAccount creates a new GitHub account for a user using an access token
func (s *AccountService) CreateAccount(ctx context.Context, userID, accessToken string) (*models.Account, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	if accessToken == "" {
		return nil, errors.New("access token is required")
	}

	// Validate access token length (GitHub personal access tokens are typically 40+ characters)
	if len(accessToken) < 40 {
		return nil, errors.New("invalid access token format")
	}

	// Check if user exists
	_, err := s.store.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Get GitHub user details using the access token
	githubUser, err := s.githubService.GetUser(ctx, accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get GitHub user details: %w", err)
	}

	// Check if this GitHub account is already linked to any user
	exists, err := s.store.AccountExistsByGithubID(ctx, githubUser.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to check if GitHub account exists: %w", err)
	}
	if exists {
		return nil, errors.New("this GitHub account is already linked to another user")
	}

	// Check if this GitHub account is already linked to this specific user
	accountExists, err := s.store.AccountExistsForUser(ctx, userID, githubUser.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to check if account exists for user: %w", err)
	}
	if accountExists {
		return nil, errors.New("this GitHub account is already linked to your account")
	}

	// Create account
	account := &models.Account{
		UserID:         userID,
		GithubUsername: githubUser.Login,
		GithubID:       githubUser.ID,
		AvatarURL:      githubUser.AvatarURL,
		AccessToken:    accessToken,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.store.CreateAccount(ctx, account); err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	return account.ToPublic(), nil
}

// GetAccountByID retrieves an account by ID and verifies ownership
func (s *AccountService) GetAccountByID(ctx context.Context, userID, accountID string) (*models.Account, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	if accountID == "" {
		return nil, errors.New("account ID is required")
	}

	account, err := s.store.GetAccountByID(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("account not found: %w", err)
	}

	// Verify ownership
	if !account.IsOwnedBy(userID) {
		return nil, errors.New("access denied: account does not belong to user")
	}

	return account.ToPublic(), nil
}

// DeleteAccount removes a GitHub account
func (s *AccountService) DeleteAccount(ctx context.Context, userID, accountID string) error {
	if userID == "" {
		return errors.New("user ID is required")
	}
	if accountID == "" {
		return errors.New("account ID is required")
	}

	// Get account to verify ownership
	account, err := s.store.GetAccountByID(ctx, accountID)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}

	// Verify ownership
	if !account.IsOwnedBy(userID) {
		return errors.New("access denied: account does not belong to user")
	}

	if err := s.store.DeleteAccount(ctx, accountID); err != nil {
		return fmt.Errorf("failed to delete account: %w", err)
	}

	return nil
}

// RefreshAccountToken updates the access token for an account
func (s *AccountService) RefreshAccountToken(ctx context.Context, userID, accountID, newAccessToken string) (*models.Account, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	if accountID == "" {
		return nil, errors.New("account ID is required")
	}
	if newAccessToken == "" {
		return nil, errors.New("new access token is required")
	}

	// Validate new access token length
	if len(newAccessToken) < 40 {
		return nil, errors.New("invalid access token format")
	}

	// Get account to verify ownership
	account, err := s.store.GetAccountByID(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("account not found: %w", err)
	}

	// Verify ownership
	if !account.IsOwnedBy(userID) {
		return nil, errors.New("access denied: account does not belong to user")
	}

	// Validate new token with GitHub
	githubUser, err := s.githubService.GetUser(ctx, newAccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to validate new access token: %w", err)
	}

	// Ensure the token belongs to the same GitHub account
	if githubUser.ID != account.GithubID {
		return nil, errors.New("access token does not belong to the linked GitHub account")
	}

	// Update token
	updatedAccount, err := s.store.UpdateAccountToken(ctx, accountID, newAccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to update access token: %w", err)
	}
	return updatedAccount, nil
}

// GetAccountRepositories retrieves repositories for a specific account
func (s *AccountService) GetAccountRepositories(ctx context.Context, userID, accountID string, page, perPage int) ([]*GitHubRepository, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	if accountID == "" {
		return nil, errors.New("account ID is required")
	}

	// Get account to verify ownership and get access token
	account, err := s.store.GetAccountByID(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("account not found: %w", err)
	}

	// Verify ownership
	if !account.IsOwnedBy(userID) {
		return nil, errors.New("access denied: account does not belong to user")
	}

	// Fetch repositories from GitHub using the stored access token
	repositories, err := s.githubService.GetRepositories(ctx, account.AccessToken, page, perPage)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repositories from GitHub: %w", err)
	}

	return repositories, nil
}

// GetAccountWithToken retrieves an account with its access token (for internal use)
func (s *AccountService) GetAccountWithToken(ctx context.Context, userID, accountID string) (*models.Account, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	if accountID == "" {
		return nil, errors.New("account ID is required")
	}

	// Get account to verify ownership
	account, err := s.store.GetAccountByID(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("account not found: %w", err)
	}

	// Verify ownership
	if !account.IsOwnedBy(userID) {
		return nil, errors.New("access denied: account does not belong to user")
	}

	// Return account with token (don't call ToPublic)
	return account, nil
}
