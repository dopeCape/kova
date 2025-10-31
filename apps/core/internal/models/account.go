package models

import (
	"time"
)

type Account struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	GithubUsername string    `json:"github_username"`
	GithubID       int64     `json:"github_id"`
	AvatarURL      string    `json:"avatar_url"`
	AccessToken    string    `json:"-"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type CreateAccountRequest struct {
	GithubUsername string `json:"github_username" validate:"required,min=1,max=255"`
	GithubID       int64  `json:"github_id" validate:"required,min=1"`
	AvatarURL      string `json:"avatar_url" validate:"omitempty,url"`
	AccessToken    string `json:"access_token" validate:"required,min=40"`
}

type UpdateAccountRequest struct {
	GithubUsername string `json:"github_username" validate:"omitempty,min=1,max=255"`
	AvatarURL      string `json:"avatar_url" validate:"omitempty,url"`
}

type UpdateAccountByGithubIDRequest struct {
	GithubUsername string `json:"github_username" validate:"required,min=1,max=255"`
	AvatarURL      string `json:"avatar_url" validate:"omitempty,url"`
	AccessToken    string `json:"access_token" validate:"required,min=40"`
}

type UpdateAccountTokenRequest struct {
	AccessToken string `json:"access_token" validate:"required,min=40"`
}

type AccountWithUser struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	GithubUsername string    `json:"github_username"`
	GithubID       int64     `json:"github_id"`
	AvatarURL      string    `json:"avatar_url"`
	AccessToken    string    `json:"-"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
}

// ToPublic returns a copy of the Account without sensitive information (access token)
func (a *Account) ToPublic() *Account {
	return &Account{
		ID:             a.ID,
		UserID:         a.UserID,
		GithubUsername: a.GithubUsername,
		GithubID:       a.GithubID,
		AvatarURL:      a.AvatarURL,
		CreatedAt:      a.CreatedAt,
		UpdatedAt:      a.UpdatedAt,
	}
}

// ToPublic returns a copy of the AccountWithUser without sensitive information (access token)
func (awu *AccountWithUser) ToPublic() *AccountWithUser {
	return &AccountWithUser{
		ID:             awu.ID,
		UserID:         awu.UserID,
		GithubUsername: awu.GithubUsername,
		GithubID:       awu.GithubID,
		AvatarURL:      awu.AvatarURL,
		CreatedAt:      awu.CreatedAt,
		UpdatedAt:      awu.UpdatedAt,
		Username:       awu.Username,
		Email:          awu.Email,
	}
}

// HasValidToken checks if the account has a non-empty access token
func (a *Account) HasValidToken() bool {
	return len(a.AccessToken) >= 40
}

// IsOwnedBy checks if the account belongs to the specified user
func (a *Account) IsOwnedBy(userID string) bool {
	return a.UserID == userID
}

