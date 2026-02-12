package handlers

import (
	"strconv"

	"github.com/SMSLYCLOUD/smsly-code/smsly-code-api/internal/models"
	"github.com/SMSLYCLOUD/smsly-code/smsly-code-api/internal/services"
	"github.com/SMSLYCLOUD/smsly-code/smsly-code-api/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

type IssueHandler struct {
	issueService *services.IssueService
}

func NewIssueHandler(issueService *services.IssueService) *IssueHandler {
	return &IssueHandler{issueService: issueService}
}

type CreateIssueRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func (h *IssueHandler) Create(c *fiber.Ctx) error {
	repo := c.Locals("repo").(*models.Repository)
	user := c.Locals("user").(*models.User)

	var req CreateIssueRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	if req.Title == "" {
		return response.BadRequest(c, "Title is required")
	}

	issue, err := h.issueService.Create(c.Context(), repo.ID, req.Title, req.Body, user.ID)
	if err != nil {
		return response.InternalError(c, err)
	}

	return response.Success(c, issue)
}

type UpdateIssueRequest struct {
	State *string `json:"state"` // "open" or "closed"
}

func (h *IssueHandler) Update(c *fiber.Ctx) error {
	repo := c.Locals("repo").(*models.Repository)
	number, _ := strconv.ParseInt(c.Params("number"), 10, 64)

	var req UpdateIssueRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	if req.State != nil {
		if *req.State == "closed" {
			if err := h.issueService.Close(c.Context(), repo.ID, number); err != nil {
				if err == pgx.ErrNoRows {
					return response.NotFound(c, "Issue")
				}
				return response.InternalError(c, err)
			}
		} else if *req.State == "open" {
			if err := h.issueService.Reopen(c.Context(), repo.ID, number); err != nil {
				if err == pgx.ErrNoRows {
					return response.NotFound(c, "Issue")
				}
				return response.InternalError(c, err)
			}
		} else {
			return response.BadRequest(c, "Invalid state (must be 'open' or 'closed')")
		}
	}

	issue, err := h.issueService.Get(c.Context(), repo.ID, number)
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.Success(c, issue)
}

func (h *IssueHandler) List(c *fiber.Ctx) error {
	repo := c.Locals("repo").(*models.Repository)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "20"))
	if page < 1 { page = 1 }
	if perPage < 1 { perPage = 20 }
	if perPage > 100 { perPage = 100 }

	issues, total, err := h.issueService.List(c.Context(), repo.ID, page, perPage)
	if err != nil {
		return response.InternalError(c, err)
	}

	return response.SuccessWithMeta(c, issues, &response.Meta{
		Page:    page,
		PerPage: perPage,
		Total:   total,
	})
}

func (h *IssueHandler) Get(c *fiber.Ctx) error {
	repo := c.Locals("repo").(*models.Repository)
	number, err := strconv.ParseInt(c.Params("number"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid issue number")
	}

	issue, err := h.issueService.Get(c.Context(), repo.ID, number)
	if err != nil {
		return response.InternalError(c, err)
	}
	if issue == nil {
		return response.NotFound(c, "Issue")
	}

	return response.Success(c, issue)
}
