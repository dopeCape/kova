package api

import (
	"strconv"
	"strings"

	"github.com/dopeCape/kova/internal/models"
	"github.com/dopeCape/kova/internal/services"
	"github.com/gofiber/fiber/v3"
)

type ProjectHandler struct {
	projectService *services.ProjectService
}

func NewProjectHandler(projectService *services.ProjectService) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
	}
}

type CreateProjectResponse struct {
	Project *models.Project `json:"project"`
	Message string          `json:"message"`
}

type GetProjectResponse struct {
	Project *models.Project `json:"project"`
}

type UpdateProjectResponse struct {
	Project *models.Project `json:"project"`
	Message string          `json:"message"`
}

type ListProjectsResponse struct {
	Projects []*models.Project `json:"projects"`
	Total    int64             `json:"total"`
	Limit    int               `json:"limit"`
	Offset   int               `json:"offset"`
	HasMore  bool              `json:"has_more"`
}

type SearchProjectsResponse struct {
	Projects []*models.Project `json:"projects"`
	Query    string            `json:"query"`
	Limit    int               `json:"limit"`
	Offset   int               `json:"offset"`
}

// RegisterRoutes registers all project routes
func (h *ProjectHandler) RegisterRoutes(router fiber.Router) {
	router.Post("/:id/projects", h.CreateProject)                        // POST /api/v1/users/:id/projects
	router.Get("/:id/projects", h.GetProjectsByUser)                     // GET /api/v1/users/:id/projects
	router.Get("/:id/projects/search", h.SearchProjectsByUser)           // GET /api/v1/users/:id/projects/search
	router.Get("/:id/projects/active", h.GetActiveProjectsByUser)        // GET /api/v1/users/:id/projects/active
	router.Get("/:id/projects/:projectId", h.GetProject)                 // GET /api/v1/users/:id/projects/:projectId
	router.Put("/:id/projects/:projectId", h.UpdateProject)              // PUT /api/v1/users/:id/projects/:projectId
	router.Put("/:id/projects/:projectId/status", h.UpdateProjectStatus) // PUT /api/v1/users/:id/projects/:projectId/status
	router.Put("/:id/projects/:projectId/archive", h.ArchiveProject)     // PUT /api/v1/users/:id/projects/:projectId/archive
	router.Put("/:id/projects/:projectId/activate", h.ActivateProject)   // PUT /api/v1/users/:id/projects/:projectId/activate
	router.Delete("/:id/projects/:projectId", h.DeleteProject)           // DELETE /api/v1/users/:id/projects/:projectId
}

// CreateProject creates a new project for a user
func (h *ProjectHandler) CreateProject(c fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return c.Status(400).JSON(ErrorResponse{
			Error: "User ID is required",
			Code:  "MISSING_USER_ID",
		})
	}

	var req models.CreateProjectRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(400).JSON(ErrorResponse{
			Error: "Invalid request body",
			Code:  "INVALID_BODY",
		})
	}

	project, err := h.projectService.CreateProject(c.RequestCtx(), userID, &req)
	if err != nil {
		if strings.Contains(err.Error(), "validation failed") {
			return c.Status(400).JSON(ErrorResponse{
				Error:   "Validation failed",
				Code:    "VALIDATION_ERROR",
				Details: err.Error(),
			})
		}
		if strings.Contains(err.Error(), "user not found") {
			return c.Status(404).JSON(ErrorResponse{
				Error: "User not found",
				Code:  "USER_NOT_FOUND",
			})
		}
		if strings.Contains(err.Error(), "already exists") {
			return c.Status(409).JSON(ErrorResponse{
				Error: err.Error(),
				Code:  "PROJECT_EXISTS",
			})
		}
		return c.Status(500).JSON(ErrorResponse{
			Error: "Failed to create project",
			Code:  "INTERNAL_ERROR",
		})
	}

	return c.Status(201).JSON(CreateProjectResponse{
		Project: project,
		Message: "Project created successfully",
	})
}

// GetProject retrieves a project by ID
func (h *ProjectHandler) GetProject(c fiber.Ctx) error {
	userID := c.Params("id")
	projectID := c.Params("projectId")

	if userID == "" {
		return c.Status(400).JSON(ErrorResponse{
			Error: "User ID is required",
			Code:  "MISSING_USER_ID",
		})
	}

	if projectID == "" {
		return c.Status(400).JSON(ErrorResponse{
			Error: "Project ID is required",
			Code:  "MISSING_PROJECT_ID",
		})
	}

	project, err := h.projectService.GetProject(c.RequestCtx(), userID, projectID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.Status(404).JSON(ErrorResponse{
				Error: "Project not found",
				Code:  "PROJECT_NOT_FOUND",
			})
		}
		if strings.Contains(err.Error(), "access denied") {
			return c.Status(403).JSON(ErrorResponse{
				Error: "Access denied",
				Code:  "ACCESS_DENIED",
			})
		}
		return c.Status(500).JSON(ErrorResponse{
			Error: "Failed to get project",
			Code:  "INTERNAL_ERROR",
		})
	}

	return c.JSON(GetProjectResponse{
		Project: project,
	})
}

// GetProjectsByUser retrieves all projects for a user
func (h *ProjectHandler) GetProjectsByUser(c fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return c.Status(400).JSON(ErrorResponse{
			Error: "User ID is required",
			Code:  "MISSING_USER_ID",
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

	projects, total, err := h.projectService.ListProjectsByUser(c.RequestCtx(), userID, limit, offset)
	if err != nil {
		if strings.Contains(err.Error(), "user not found") {
			return c.Status(404).JSON(ErrorResponse{
				Error: "User not found",
				Code:  "USER_NOT_FOUND",
			})
		}
		return c.Status(500).JSON(ErrorResponse{
			Error: "Failed to get projects",
			Code:  "INTERNAL_ERROR",
		})
	}

	hasMore := int64(offset+limit) < total

	return c.JSON(ListProjectsResponse{
		Projects: projects,
		Total:    total,
		Limit:    limit,
		Offset:   offset,
		HasMore:  hasMore,
	})
}

// GetActiveProjectsByUser retrieves only active projects for a user
func (h *ProjectHandler) GetActiveProjectsByUser(c fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return c.Status(400).JSON(ErrorResponse{
			Error: "User ID is required",
			Code:  "MISSING_USER_ID",
		})
	}

	projects, err := h.projectService.GetActiveProjectsByUser(c.RequestCtx(), userID)
	if err != nil {
		if strings.Contains(err.Error(), "user not found") {
			return c.Status(404).JSON(ErrorResponse{
				Error: "User not found",
				Code:  "USER_NOT_FOUND",
			})
		}
		return c.Status(500).JSON(ErrorResponse{
			Error: "Failed to get active projects",
			Code:  "INTERNAL_ERROR",
		})
	}

	return c.JSON(fiber.Map{
		"projects": projects,
		"total":    len(projects),
	})
}

// UpdateProject updates a project
func (h *ProjectHandler) UpdateProject(c fiber.Ctx) error {
	userID := c.Params("id")
	projectID := c.Params("projectId")

	if userID == "" {
		return c.Status(400).JSON(ErrorResponse{
			Error: "User ID is required",
			Code:  "MISSING_USER_ID",
		})
	}

	if projectID == "" {
		return c.Status(400).JSON(ErrorResponse{
			Error: "Project ID is required",
			Code:  "MISSING_PROJECT_ID",
		})
	}

	var req models.UpdateProjectRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(400).JSON(ErrorResponse{
			Error: "Invalid request body",
			Code:  "INVALID_BODY",
		})
	}

	project, err := h.projectService.UpdateProject(c.RequestCtx(), userID, projectID, &req)
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
				Error: "Project not found",
				Code:  "PROJECT_NOT_FOUND",
			})
		}
		if strings.Contains(err.Error(), "access denied") {
			return c.Status(403).JSON(ErrorResponse{
				Error: "Access denied",
				Code:  "ACCESS_DENIED",
			})
		}
		if strings.Contains(err.Error(), "already exists") {
			return c.Status(409).JSON(ErrorResponse{
				Error: err.Error(),
				Code:  "PROJECT_EXISTS",
			})
		}
		return c.Status(500).JSON(ErrorResponse{
			Error: "Failed to update project",
			Code:  "INTERNAL_ERROR",
		})
	}

	return c.JSON(UpdateProjectResponse{
		Project: project,
		Message: "Project updated successfully",
	})
}

// UpdateProjectStatus updates only the project status
func (h *ProjectHandler) UpdateProjectStatus(c fiber.Ctx) error {
	userID := c.Params("id")
	projectID := c.Params("projectId")

	if userID == "" {
		return c.Status(400).JSON(ErrorResponse{
			Error: "User ID is required",
			Code:  "MISSING_USER_ID",
		})
	}

	if projectID == "" {
		return c.Status(400).JSON(ErrorResponse{
			Error: "Project ID is required",
			Code:  "MISSING_PROJECT_ID",
		})
	}

	var req struct {
		Status string `json:"status" validate:"required,oneof=active inactive archived"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(400).JSON(ErrorResponse{
			Error: "Invalid request body",
			Code:  "INVALID_BODY",
		})
	}

	project, err := h.projectService.UpdateProjectStatus(c.RequestCtx(), userID, projectID, req.Status)
	if err != nil {
		if strings.Contains(err.Error(), "invalid status") {
			return c.Status(400).JSON(ErrorResponse{
				Error: err.Error(),
				Code:  "INVALID_STATUS",
			})
		}
		if strings.Contains(err.Error(), "not found") {
			return c.Status(404).JSON(ErrorResponse{
				Error: "Project not found",
				Code:  "PROJECT_NOT_FOUND",
			})
		}
		if strings.Contains(err.Error(), "access denied") {
			return c.Status(403).JSON(ErrorResponse{
				Error: "Access denied",
				Code:  "ACCESS_DENIED",
			})
		}
		return c.Status(500).JSON(ErrorResponse{
			Error: "Failed to update project status",
			Code:  "INTERNAL_ERROR",
		})
	}

	return c.JSON(UpdateProjectResponse{
		Project: project,
		Message: "Project status updated successfully",
	})
}

// ArchiveProject archives a project
func (h *ProjectHandler) ArchiveProject(c fiber.Ctx) error {
	userID := c.Params("id")
	projectID := c.Params("projectId")

	if userID == "" {
		return c.Status(400).JSON(ErrorResponse{
			Error: "User ID is required",
			Code:  "MISSING_USER_ID",
		})
	}

	if projectID == "" {
		return c.Status(400).JSON(ErrorResponse{
			Error: "Project ID is required",
			Code:  "MISSING_PROJECT_ID",
		})
	}

	project, err := h.projectService.ArchiveProject(c.RequestCtx(), userID, projectID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.Status(404).JSON(ErrorResponse{
				Error: "Project not found",
				Code:  "PROJECT_NOT_FOUND",
			})
		}
		if strings.Contains(err.Error(), "access denied") {
			return c.Status(403).JSON(ErrorResponse{
				Error: "Access denied",
				Code:  "ACCESS_DENIED",
			})
		}
		return c.Status(500).JSON(ErrorResponse{
			Error: "Failed to archive project",
			Code:  "INTERNAL_ERROR",
		})
	}

	return c.JSON(UpdateProjectResponse{
		Project: project,
		Message: "Project archived successfully",
	})
}

// ActivateProject activates a project
func (h *ProjectHandler) ActivateProject(c fiber.Ctx) error {
	userID := c.Params("id")
	projectID := c.Params("projectId")

	if userID == "" {
		return c.Status(400).JSON(ErrorResponse{
			Error: "User ID is required",
			Code:  "MISSING_USER_ID",
		})
	}

	if projectID == "" {
		return c.Status(400).JSON(ErrorResponse{
			Error: "Project ID is required",
			Code:  "MISSING_PROJECT_ID",
		})
	}

	project, err := h.projectService.ActivateProject(c.RequestCtx(), userID, projectID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.Status(404).JSON(ErrorResponse{
				Error: "Project not found",
				Code:  "PROJECT_NOT_FOUND",
			})
		}
		if strings.Contains(err.Error(), "access denied") {
			return c.Status(403).JSON(ErrorResponse{
				Error: "Access denied",
				Code:  "ACCESS_DENIED",
			})
		}
		return c.Status(500).JSON(ErrorResponse{
			Error: "Failed to activate project",
			Code:  "INTERNAL_ERROR",
		})
	}

	return c.JSON(UpdateProjectResponse{
		Project: project,
		Message: "Project activated successfully",
	})
}

// DeleteProject deletes a project
func (h *ProjectHandler) DeleteProject(c fiber.Ctx) error {
	userID := c.Params("id")
	projectID := c.Params("projectId")

	if userID == "" {
		return c.Status(400).JSON(ErrorResponse{
			Error: "User ID is required",
			Code:  "MISSING_USER_ID",
		})
	}

	if projectID == "" {
		return c.Status(400).JSON(ErrorResponse{
			Error: "Project ID is required",
			Code:  "MISSING_PROJECT_ID",
		})
	}

	err := h.projectService.DeleteProject(c.RequestCtx(), userID, projectID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.Status(404).JSON(ErrorResponse{
				Error: "Project not found",
				Code:  "PROJECT_NOT_FOUND",
			})
		}
		if strings.Contains(err.Error(), "access denied") {
			return c.Status(403).JSON(ErrorResponse{
				Error: "Access denied",
				Code:  "ACCESS_DENIED",
			})
		}
		return c.Status(500).JSON(ErrorResponse{
			Error: "Failed to delete project",
			Code:  "INTERNAL_ERROR",
		})
	}

	return c.JSON(MessageResponse{
		Message: "Project deleted successfully",
	})
}

// SearchProjectsByUser searches projects for a user
func (h *ProjectHandler) SearchProjectsByUser(c fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return c.Status(400).JSON(ErrorResponse{
			Error: "User ID is required",
			Code:  "MISSING_USER_ID",
		})
	}

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

	projects, err := h.projectService.SearchProjectsByUser(c.RequestCtx(), userID, query, limit, offset)
	if err != nil {
		if strings.Contains(err.Error(), "user not found") {
			return c.Status(404).JSON(ErrorResponse{
				Error: "User not found",
				Code:  "USER_NOT_FOUND",
			})
		}
		return c.Status(500).JSON(ErrorResponse{
			Error: "Failed to search projects",
			Code:  "INTERNAL_ERROR",
		})
	}

	return c.JSON(SearchProjectsResponse{
		Projects: projects,
		Query:    query,
		Limit:    limit,
		Offset:   offset,
	})
}
