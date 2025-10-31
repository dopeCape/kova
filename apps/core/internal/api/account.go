package api

import (
	"strings"

	"github.com/dopeCape/kova/internal/models"
	"github.com/dopeCape/kova/internal/services"
	"github.com/gofiber/fiber/v3"
)

type AccountHandler struct {
	accountService *services.AccountService
}

func NewAccountHandler(accountService *services.AccountService) *AccountHandler {
	return &AccountHandler{
		accountService: accountService,
	}
}

type CreateAccountRequest struct {
	AccessToken string `json:"access_token" validate:"required,min=40"`
}

type CreateAccountResponse struct {
	Account *models.Account `json:"account"`
	Message string          `json:"message"`
}

type GetAccountsResponse struct {
	Accounts []*models.Account `json:"accounts"`
	Total    int               `json:"total"`
	UserID   string            `json:"user_id"`
}

// RegisterRoutes registers all account routes
func (h *AccountHandler) RegisterRoutes(router fiber.Router) {
	router.Get("/:id/accounts", h.GetAccountsByUserID)   // GET /api/v1/users/:id/accounts
	router.Post("/:id/accounts", h.CreateAccountForUser) // POST /api/v1/users/:id/accounts
}

// GetAccountsByUserID retrieves all GitHub accounts for a user
func (h *AccountHandler) GetAccountsByUserID(c fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return c.Status(400).JSON(ErrorResponse{
			Error: "User ID is required",
			Code:  "MISSING_ID",
		})
	}

	accounts, err := h.accountService.GetAccountsByUserID(c.RequestCtx(), userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.Status(404).JSON(ErrorResponse{
				Error: "User not found",
				Code:  "USER_NOT_FOUND",
			})
		}
		return c.Status(500).JSON(ErrorResponse{
			Error: "Failed to get accounts",
			Code:  "INTERNAL_ERROR",
		})
	}

	return c.JSON(GetAccountsResponse{
		Accounts: accounts,
		Total:    len(accounts),
		UserID:   userID,
	})
}

// CreateAccountForUser creates a new GitHub account for a user using access token
func (h *AccountHandler) CreateAccountForUser(c fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return c.Status(400).JSON(ErrorResponse{
			Error: "User ID is required",
			Code:  "MISSING_ID",
		})
	}

	var req CreateAccountRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(400).JSON(ErrorResponse{
			Error: "Invalid request body",
			Code:  "INVALID_BODY",
		})
	}

	account, err := h.accountService.CreateAccount(c.RequestCtx(), userID, req.AccessToken)
	if err != nil {
		if strings.Contains(err.Error(), "user not found") {
			return c.Status(404).JSON(ErrorResponse{
				Error: "User not found",
				Code:  "USER_NOT_FOUND",
			})
		}
		if strings.Contains(err.Error(), "access token is required") || strings.Contains(err.Error(), "invalid access token") {
			return c.Status(400).JSON(ErrorResponse{
				Error: "Invalid access token",
				Code:  "INVALID_TOKEN",
			})
		}
		if strings.Contains(err.Error(), "invalid or expired") {
			return c.Status(401).JSON(ErrorResponse{
				Error: "Invalid or expired GitHub access token",
				Code:  "INVALID_GITHUB_TOKEN",
			})
		}
		if strings.Contains(err.Error(), "already linked") {
			return c.Status(409).JSON(ErrorResponse{
				Error: err.Error(),
				Code:  "ACCOUNT_ALREADY_LINKED",
			})
		}
		if strings.Contains(err.Error(), "GitHub user details") {
			return c.Status(502).JSON(ErrorResponse{
				Error: "Failed to fetch GitHub user details",
				Code:  "GITHUB_API_ERROR",
			})
		}
		return c.Status(500).JSON(ErrorResponse{
			Error: "Failed to create account",
			Code:  "INTERNAL_ERROR",
		})
	}

	return c.Status(201).JSON(CreateAccountResponse{
		Account: account,
		Message: "GitHub account linked successfully",
	})
}
