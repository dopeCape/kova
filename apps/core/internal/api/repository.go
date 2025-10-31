package api

import (
	"strconv"
	"strings"
	"time"

	"github.com/dopeCape/kova/internal/services"
	"github.com/gofiber/fiber/v3"
)

type RepositoryHandler struct {
	accountService *services.AccountService
}

func NewRepositoryHandler(accountService *services.AccountService) *RepositoryHandler {
	return &RepositoryHandler{
		accountService: accountService,
	}
}

type Repository struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	FullName      string `json:"full_name"`
	Private       bool   `json:"private"`
	Description   string `json:"description"`
	Language      string `json:"language"`
	Stars         int    `json:"stars"`
	UpdatedAt     string `json:"updated_at"`
	DefaultBranch string `json:"default_branch"`
	URL           string `json:"url"`
}

type GetRepositoriesResponse struct {
	Repositories []*Repository `json:"repositories"`
	Total        int           `json:"total"`
	Page         int           `json:"page"`
	PerPage      int           `json:"per_page"`
	AccountID    string        `json:"account_id"`
}

// RegisterRoutes registers repository routes under account routes
func (h *RepositoryHandler) RegisterRoutes(router fiber.Router) {
	router.Get(":id/:accountId/repositories", h.GetRepositoriesByAccount) // GET /api/v1/users/:id/accounts/:accountId/repositories
}

// GetRepositoriesByAccount retrieves all repositories for a specific GitHub account
func (h *RepositoryHandler) GetRepositoriesByAccount(c fiber.Ctx) error {
	userID := c.Params("id")
	accountID := c.Params("accountId")

	if userID == "" {
		return c.Status(400).JSON(ErrorResponse{
			Error: "User ID is required",
			Code:  "MISSING_USER_ID",
		})
	}

	if accountID == "" {
		return c.Status(400).JSON(ErrorResponse{
			Error: "Account ID is required",
			Code:  "MISSING_ACCOUNT_ID",
		})
	}

	// Parse pagination parameters
	page := 1
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	perPage := 30
	if pp := c.Query("per_page"); pp != "" {
		if parsed, err := strconv.Atoi(pp); err == nil && parsed > 0 && parsed <= 100 {
			perPage = parsed
		}
	}

	repositories, err := h.accountService.GetAccountRepositories(c.RequestCtx(), userID, accountID, page, perPage)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.Status(404).JSON(ErrorResponse{
				Error: "Account not found",
				Code:  "ACCOUNT_NOT_FOUND",
			})
		}
		if strings.Contains(err.Error(), "access denied") {
			return c.Status(403).JSON(ErrorResponse{
				Error: "Access denied",
				Code:  "ACCESS_DENIED",
			})
		}
		if strings.Contains(err.Error(), "invalid or expired") {
			return c.Status(401).JSON(ErrorResponse{
				Error: "GitHub access token is invalid or expired",
				Code:  "INVALID_GITHUB_TOKEN",
			})
		}
		if strings.Contains(err.Error(), "GitHub") {
			return c.Status(502).JSON(ErrorResponse{
				Error: "Failed to fetch repositories from GitHub",
				Code:  "GITHUB_API_ERROR",
			})
		}
		return c.Status(500).JSON(ErrorResponse{
			Error: "Failed to get repositories",
			Code:  "INTERNAL_ERROR",
		})
	}

	// Convert GitHub repositories to API response format
	apiRepositories := make([]*Repository, len(repositories))
	for i, repo := range repositories {
		// Parse updated_at time and format it relative to now
		updatedAt := "unknown"
		if repo.UpdatedAt != "" {
			if t, err := time.Parse(time.RFC3339, repo.UpdatedAt); err == nil {
				updatedAt = formatRelativeTime(t)
			}
		}

		apiRepositories[i] = &Repository{
			ID:            strconv.FormatInt(repo.ID, 10),
			Name:          repo.Name,
			FullName:      repo.FullName,
			Private:       repo.Private,
			Description:   repo.Description,
			Language:      repo.Language,
			Stars:         repo.StarCount,
			UpdatedAt:     updatedAt,
			DefaultBranch: repo.DefaultBranch,
			URL:           repo.HTMLURL,
		}
	}

	return c.JSON(GetRepositoriesResponse{
		Repositories: apiRepositories,
		Total:        len(apiRepositories),
		Page:         page,
		PerPage:      perPage,
		AccountID:    accountID,
	})
}

// formatRelativeTime formats a time.Time to a relative string like "2h", "1d", "3w"
func formatRelativeTime(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	if diff < time.Hour {
		minutes := int(diff.Minutes())
		if minutes < 1 {
			return "now"
		}
		return strconv.Itoa(minutes) + "m"
	}

	if diff < 24*time.Hour {
		hours := int(diff.Hours())
		return strconv.Itoa(hours) + "h"
	}

	if diff < 7*24*time.Hour {
		days := int(diff.Hours() / 24)
		return strconv.Itoa(days) + "d"
	}

	if diff < 30*24*time.Hour {
		weeks := int(diff.Hours() / (24 * 7))
		return strconv.Itoa(weeks) + "w"
	}

	if diff < 365*24*time.Hour {
		months := int(diff.Hours() / (24 * 30))
		return strconv.Itoa(months) + "mo"
	}

	years := int(diff.Hours() / (24 * 365))
	return strconv.Itoa(years) + "y"
}
