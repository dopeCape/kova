package api

import (
	"strconv"
	"strings"

	"github.com/dopeCape/kova/internal/models"
	"github.com/dopeCape/kova/internal/services"
	"github.com/gofiber/fiber/v3"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

type CreateUserResponse struct {
	User    *models.User `json:"user"`
	Message string       `json:"message"`
}

type GetUserResponse struct {
	User *models.User `json:"user"`
}

type UpdateUserResponse struct {
	User    *models.User `json:"user"`
	Message string       `json:"message"`
}

type ListUsersResponse struct {
	Users   []*models.User `json:"users"`
	Total   int64          `json:"total"`
	Limit   int            `json:"limit"`
	Offset  int            `json:"offset"`
	HasMore bool           `json:"has_more"`
}

type SearchUsersResponse struct {
	Users  []*models.User `json:"users"`
	Query  string         `json:"query"`
	Limit  int            `json:"limit"`
	Offset int            `json:"offset"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}

// RegisterRoutes registers all user routes
func (h *UserHandler) RegisterRoutes(router fiber.Router) {
	router.Post("/", h.CreateUser)                        // POST /api/v1/users
	router.Get("/", h.ListUsers)                          // GET /api/v1/users
	router.Get("/search", h.SearchUsers)                  // GET /api/v1/users/search
	router.Get("/:id", h.GetUser)                         // GET /api/v1/users/:id
	router.Put("/:id", h.UpdateUser)                      // PUT /api/v1/users/:id
	router.Delete("/:id", h.DeleteUser)                   // DELETE /api/v1/users/:id
	router.Post("/:id/change-password", h.ChangePassword) // POST /api/v1/users/:id/change-password
}

func (h *UserHandler) CreateUser(c fiber.Ctx) error {
	var req models.CreateUserRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(400).JSON(ErrorResponse{
			Error: "Invalid request body",
			Code:  "INVALID_BODY",
		})
	}

	user, err := h.userService.CreateUser(c.RequestCtx(), &req)
	if err != nil {
		if strings.Contains(err.Error(), "validation failed") {
			return c.Status(400).JSON(ErrorResponse{
				Error:   "Validation failed",
				Code:    "VALIDATION_ERROR",
				Details: err.Error(),
			})
		}
		if strings.Contains(err.Error(), "already exists") {
			return c.Status(409).JSON(ErrorResponse{
				Error: err.Error(),
				Code:  "USER_EXISTS",
			})
		}
		return c.Status(500).JSON(ErrorResponse{
			Error: "Failed to create user",
			Code:  "INTERNAL_ERROR",
		})
	}

	return c.Status(201).JSON(CreateUserResponse{
		User:    user,
		Message: "User created successfully",
	})
}

// GetUser retrieves a user by ID
func (h *UserHandler) GetUser(c fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return c.Status(400).JSON(ErrorResponse{
			Error: "User ID is required",
			Code:  "MISSING_ID",
		})
	}

	user, err := h.userService.GetUser(c.RequestCtx(), userID)
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

	return c.JSON(GetUserResponse{
		User: user,
	})
}

// UpdateUser updates user information
func (h *UserHandler) UpdateUser(c fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return c.Status(400).JSON(ErrorResponse{
			Error: "User ID is required",
			Code:  "MISSING_ID",
		})
	}

	var req models.UpdateUserRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(400).JSON(ErrorResponse{
			Error: "Invalid request body",
			Code:  "INVALID_BODY",
		})
	}

	user, err := h.userService.UpdateUser(c.RequestCtx(), userID, &req)
	if err != nil {
		if strings.Contains(err.Error(), "validation failed") {
			return c.Status(400).JSON(ErrorResponse{
				Error:   "Validation failed",
				Code:    "VALIDATION_ERROR",
				Details: err.Error(),
			})
		}
		if strings.Contains(err.Error(), "not found") {
			return c.Status(404).JSON(ErrorResponse{
				Error: "User not found",
				Code:  "USER_NOT_FOUND",
			})
		}
		if strings.Contains(err.Error(), "already in use") || strings.Contains(err.Error(), "already exists") {
			return c.Status(409).JSON(ErrorResponse{
				Error: err.Error(),
				Code:  "CONFLICT",
			})
		}
		return c.Status(500).JSON(ErrorResponse{
			Error: "Failed to update user",
			Code:  "INTERNAL_ERROR",
		})
	}

	return c.JSON(UpdateUserResponse{
		User:    user,
		Message: "User updated successfully",
	})
}

// ChangePassword changes user password
func (h *UserHandler) ChangePassword(c fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return c.Status(400).JSON(ErrorResponse{
			Error: "User ID is required",
			Code:  "MISSING_ID",
		})
	}

	var req models.ChangePasswordRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(400).JSON(ErrorResponse{
			Error: "Invalid request body",
			Code:  "INVALID_BODY",
		})
	}

	err := h.userService.ChangePassword(c.RequestCtx(), userID, &req)
	if err != nil {
		if strings.Contains(err.Error(), "validation failed") {
			return c.Status(400).JSON(ErrorResponse{
				Error:   "Validation failed",
				Code:    "VALIDATION_ERROR",
				Details: err.Error(),
			})
		}
		if strings.Contains(err.Error(), "not found") {
			return c.Status(404).JSON(ErrorResponse{
				Error: "User not found",
				Code:  "USER_NOT_FOUND",
			})
		}
		if strings.Contains(err.Error(), "incorrect") {
			return c.Status(400).JSON(ErrorResponse{
				Error: "Current password is incorrect",
				Code:  "INVALID_PASSWORD",
			})
		}
		return c.Status(500).JSON(ErrorResponse{
			Error: "Failed to change password",
			Code:  "INTERNAL_ERROR",
		})
	}

	return c.JSON(MessageResponse{
		Message: "Password changed successfully",
	})
}

// DeleteUser deletes a user
func (h *UserHandler) DeleteUser(c fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return c.Status(400).JSON(ErrorResponse{
			Error: "User ID is required",
			Code:  "MISSING_ID",
		})
	}

	err := h.userService.DeleteUser(c.RequestCtx(), userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.Status(404).JSON(ErrorResponse{
				Error: "User not found",
				Code:  "USER_NOT_FOUND",
			})
		}
		return c.Status(500).JSON(ErrorResponse{
			Error: "Failed to delete user",
			Code:  "INTERNAL_ERROR",
		})
	}

	return c.JSON(MessageResponse{
		Message: "User deleted successfully",
	})
}

// ListUsers retrieves a paginated list of users
func (h *UserHandler) ListUsers(c fiber.Ctx) error {
	// Parse query parameters
	limit := 20 // default
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	offset := 0 // default
	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	users, total, err := h.userService.ListUsers(c.RequestCtx(), limit, offset)
	if err != nil {
		return c.Status(500).JSON(ErrorResponse{
			Error: "Failed to list users",
			Code:  "INTERNAL_ERROR",
		})
	}

	hasMore := int64(offset+limit) < total

	return c.JSON(ListUsersResponse{
		Users:   users,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
		HasMore: hasMore,
	})
}

// SearchUsers searches for users by username or email
func (h *UserHandler) SearchUsers(c fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return c.Status(400).JSON(ErrorResponse{
			Error: "Search query is required",
			Code:  "MISSING_QUERY",
		})
	}

	// Parse pagination parameters
	limit := 20 // default
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	offset := 0 // default
	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	users, err := h.userService.SearchUsers(c.RequestCtx(), query, limit, offset)
	if err != nil {
		return c.Status(500).JSON(ErrorResponse{
			Error: "Failed to search users",
			Code:  "INTERNAL_ERROR",
		})
	}

	return c.JSON(SearchUsersResponse{
		Users:  users,
		Query:  query,
		Limit:  limit,
		Offset: offset,
	})

}
