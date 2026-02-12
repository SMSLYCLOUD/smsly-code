package handlers

import (
	"strconv"

	"github.com/SMSLYCLOUD/smsly-code/smsly-code-api/internal/models"
	"github.com/SMSLYCLOUD/smsly-code/smsly-code-api/internal/services"
	"github.com/SMSLYCLOUD/smsly-code/smsly-code-api/pkg/response"
	"github.com/gofiber/fiber/v2"
)

type ProjectHandler struct {
	projectService *services.ProjectService
}

func NewProjectHandler(projectService *services.ProjectService) *ProjectHandler {
	return &ProjectHandler{projectService: projectService}
}

type CreateProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (h *ProjectHandler) Create(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	var req CreateProjectRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	if req.Name == "" {
		return response.BadRequest(c, "Name is required")
	}

	project, err := h.projectService.CreateProject(c.Context(), user.ID, req.Name, req.Description)
	if err != nil {
		return response.InternalError(c, err)
	}

	return response.Success(c, project)
}

type CreateColumnRequest struct {
	Name string `json:"name"`
}

func (h *ProjectHandler) CreateColumn(c *fiber.Ctx) error {
	projectID, err := strconv.ParseInt(c.Params("project_id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid project ID")
	}

	var req CreateColumnRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	if req.Name == "" {
		return response.BadRequest(c, "Name is required")
	}

	column, err := h.projectService.CreateColumn(c.Context(), projectID, req.Name)
	if err != nil {
		return response.InternalError(c, err)
	}

	return response.Success(c, column)
}

type CreateCardRequest struct {
	ContentURL string `json:"content_url"`
	Note       string `json:"note"`
}

func (h *ProjectHandler) CreateCard(c *fiber.Ctx) error {
	columnID, err := strconv.ParseInt(c.Params("column_id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid column ID")
	}

	var req CreateCardRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	card, err := h.projectService.CreateCard(c.Context(), columnID, req.ContentURL, req.Note)
	if err != nil {
		return response.InternalError(c, err)
	}

	return response.Success(c, card)
}

func (h *ProjectHandler) Get(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid project ID")
	}

	project, err := h.projectService.GetProject(c.Context(), id)
	if err != nil {
		return response.InternalError(c, err)
	}
	if project == nil {
		return response.NotFound(c, "Project")
	}

	return response.Success(c, project)
}
