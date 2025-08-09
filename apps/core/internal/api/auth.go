package api

import (
	"strings"

	"github.com/dopeCape/kova/internal/models"
	"github.com/dopeCape/kova/internal/services"
	"github.com/gofiber/fiber/v3"
)

type AuthHandler struct {
	authService *services.AuthService
	userService *services.UserService
}

func NewAuthHandler(authService *services.AuthService, userService *services.UserService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userService: userService,
	}
}

// Response structures
type LoginResponse struct {
	User    *models.User        `json:"user"`
	Tokens  *services.TokenPair `json:"tokens"`
	Message string              `json:"message"`
}

type RefreshResponse struct {
	Tokens  *services.TokenPair `json:"tokens"`
	Message string              `json:"message"`
}

type MeResponse struct {
	User *models.User `json:"user"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// RegisterRoutes registers all auth routes
func (h *AuthHandler) RegisterRoutes(router fiber.Router) {
	router.Post("/login", h.Login)
	router.Post("/logout", h.Logout)
	router.Post("/refresh", h.RefreshToken)
	router.Get("/me", h.AuthMiddleware(), h.GetCurrentUser)
}

// Login authenticates a user
func (h *AuthHandler) Login(c fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(400).JSON(ErrorResponse{
			Error: "Invalid request body",
			Code:  "INVALID_BODY",
		})
	}

	user, tokens, err := h.authService.Login(c.RequestCtx(), &req)
	if err != nil {
		if strings.Contains(err.Error(), "validation failed") {
			return c.Status(400).JSON(ErrorResponse{
				Error:   "Validation failed",
				Code:    "VALIDATION_ERROR",
				Details: err.Error(),
			})
		}
		if strings.Contains(err.Error(), "invalid credentials") {
			return c.Status(401).JSON(ErrorResponse{
				Error: "Invalid credentials",
				Code:  "INVALID_CREDENTIALS",
			})
		}
		return c.Status(500).JSON(ErrorResponse{
			Error: "Authentication failed",
			Code:  "AUTH_ERROR",
		})
	}

	return c.JSON(LoginResponse{
		User:    user,
		Tokens:  tokens,
		Message: "Login successful",
	})
}

// Logout invalidates the current session
func (h *AuthHandler) Logout(c fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	tokenString := h.authService.ExtractTokenFromHeader(authHeader)

	if tokenString == "" {
		return c.Status(400).JSON(ErrorResponse{
			Error: "Authorization header required",
			Code:  "MISSING_TOKEN",
		})
	}

	err := h.authService.Logout(c.RequestCtx(), tokenString)
	if err != nil {
		return c.Status(400).JSON(ErrorResponse{
			Error: "Invalid token",
			Code:  "INVALID_TOKEN",
		})
	}

	return c.JSON(MessageResponse{
		Message: "Logout successful",
	})
}

// RefreshToken generates new access token using refresh token
func (h *AuthHandler) RefreshToken(c fiber.Ctx) error {
	var req RefreshRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(400).JSON(ErrorResponse{
			Error: "Invalid request body",
			Code:  "INVALID_BODY",
		})
	}

	tokens, err := h.authService.RefreshToken(c.RequestCtx(), req.RefreshToken)
	if err != nil {
		if strings.Contains(err.Error(), "invalid") {
			return c.Status(401).JSON(ErrorResponse{
				Error: "Invalid refresh token",
				Code:  "INVALID_REFRESH_TOKEN",
			})
		}
		return c.Status(500).JSON(ErrorResponse{
			Error: "Token refresh failed",
			Code:  "REFRESH_ERROR",
		})
	}

	return c.JSON(RefreshResponse{
		Tokens:  tokens,
		Message: "Token refreshed successfully",
	})
}

// GetCurrentUser returns the current authenticated user
func (h *AuthHandler) GetCurrentUser(c fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	user, err := h.authService.GetCurrentUser(c.RequestCtx(), userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.Status(404).JSON(ErrorResponse{
				Error: "User not found",
				Code:  "USER_NOT_FOUND",
			})
		}
		return c.Status(500).JSON(ErrorResponse{
			Error: "Failed to get user",
			Code:  "INTERNAL_ERROR",
		})
	}

	return c.JSON(MeResponse{
		User: user,
	})
}

// AuthMiddleware validates JWT tokens
func (h *AuthHandler) AuthMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		tokenString := h.authService.ExtractTokenFromHeader(authHeader)

		if tokenString == "" {
			return c.Status(401).JSON(ErrorResponse{

				Code: "MISSING_TOKEN",
			})
		}

		claims, err := h.authService.ValidateToken(tokenString)
		if err != nil {
			return c.Status(401).JSON(ErrorResponse{
				Error: "Invalid or expired token",
				Code:  "INVALID_TOKEN",
			})
		}

		// Set user information in context
		c.Locals("user_id", claims.UserID)
		c.Locals("username", claims.Username)
		c.Locals("email", claims.Email)

		return c.Next()
	}
}

// OptionalAuthMiddleware validates JWT tokens but doesn't require them
func (h *AuthHandler) OptionalAuthMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		tokenString := h.authService.ExtractTokenFromHeader(authHeader)

		if tokenString != "" {
			claims, err := h.authService.ValidateToken(tokenString)
			if err == nil {
				// Set user information in context if token is valid
				c.Locals("user_id", claims.UserID)
				c.Locals("username", claims.Username)
				c.Locals("email", claims.Email)
				c.Locals("authenticated", true)
			}
		}

		return c.Next()
	}
}

// RequireAuthMiddleware is a helper to get the auth middleware
func (h *AuthHandler) RequireAuthMiddleware() fiber.Handler {
	return h.AuthMiddleware()
}

// GetUserFromContext extracts user information from Fiber context
func GetUserFromContext(c fiber.Ctx) (string, string, string, bool) {
	userID, _ := c.Locals("user_id").(string)
	username, _ := c.Locals("username").(string)
	email, _ := c.Locals("email").(string)
	authenticated := userID != ""

	return userID, username, email, authenticated
}

// RequireAuth is a utility function to check if user is authenticated
func RequireAuth(c fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(401).JSON(ErrorResponse{
			Error: "Authentication required",
			Code:  "AUTH_REQUIRED",
		})
	}
	return nil
}
