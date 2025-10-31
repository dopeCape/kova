package api

import (
	"strings"

	"github.com/dopeCape/kova/internal/models"
	"github.com/dopeCape/kova/internal/services"
	"github.com/gofiber/fiber/v3"
)

type AnalyzerHandler struct {
	analyzerService *services.RepositoryAnalyzerService
	accountService  *services.AccountService
}

func NewAnalyzerHandler(analyzerService *services.RepositoryAnalyzerService, accountService *services.AccountService) *AnalyzerHandler {
	return &AnalyzerHandler{
		analyzerService: analyzerService,
		accountService:  accountService,
	}
}

type AnalyzeRepositoryResponse struct {
	Analysis *models.RepositoryAnalysis `json:"analysis"`
	Message  string                     `json:"message"`
}

// RegisterRoutes registers analyzer routes
func (h *AnalyzerHandler) RegisterRoutes(router fiber.Router) {
	router.Post(":id/:accountId/repositorie/analyze", h.AnalyzeRepository)
}

// AnalyzeRepository analyzes a repository and returns build commands
func (h *AnalyzerHandler) AnalyzeRepository(c fiber.Ctx) error {
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

	var req models.AnalyzeRepositoryRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(400).JSON(ErrorResponse{
			Error: "Invalid request body",
			Code:  "INVALID_BODY",
		})
	}

	// Get account to verify ownership and get access token
	_, err := h.accountService.GetAccountByID(c.RequestCtx(), userID, accountID)
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
		return c.Status(500).JSON(ErrorResponse{
			Error: "Failed to get account",
			Code:  "INTERNAL_ERROR",
		})
	}

	// We need the account with token for cloning
	accountWithToken, err := h.accountService.GetAccountWithToken(c.RequestCtx(), userID, accountID)
	if err != nil {
		return c.Status(500).JSON(ErrorResponse{
			Error: "Failed to get account credentials",
			Code:  "INTERNAL_ERROR",
		})
	}

	// Analyze repository
	analysis, err := h.analyzerService.AnalyzeRepository(c.RequestCtx(), accountWithToken.AccessToken, &req)
	if err != nil {
		if strings.Contains(err.Error(), "git clone failed") {
			return c.Status(400).JSON(ErrorResponse{
				Error: "Failed to clone repository. Check if the repository exists and access token has proper permissions",
				Code:  "CLONE_FAILED",
			})
		}
		if strings.Contains(err.Error(), "railpack") {
			return c.Status(500).JSON(ErrorResponse{
				Error: "Failed to analyze repository",
				Code:  "ANALYSIS_FAILED",
			})
		}
		return c.Status(500).JSON(ErrorResponse{
			Error: "Repository analysis failed",
			Code:  "INTERNAL_ERROR",
		})
	}

	// Check if analysis was successful
	if !analysis.Success {
		return c.Status(200).JSON(AnalyzeRepositoryResponse{
			Analysis: analysis,
			Message:  "Repository is not supported by the build system",
		})
	}

	return c.Status(200).JSON(AnalyzeRepositoryResponse{
		Analysis: analysis,
		Message:  "Repository analyzed successfully",
	})
}
