package handlers

import (
	"github.com/SMSLYCLOUD/smsly-code/smsly-code-api/internal/config"
	"github.com/SMSLYCLOUD/smsly-code/smsly-code-api/internal/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type IssueHandler struct {
	DB     *gorm.DB
	Config *config.Config
}

func NewIssueHandler(db *gorm.DB, cfg *config.Config) *IssueHandler {
	return &IssueHandler{
		DB:     db,
		Config: cfg,
	}
}

type CreateIssueRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func (h *IssueHandler) Create(c *fiber.Ctx) error {
	repoName := c.Params("name")
	var req CreateIssueRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Title is required"})
	}

	user := c.Locals("user").(models.User)

	var repo models.Repository
	if result := h.DB.Where("name = ?", repoName).First(&repo); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Repository not found"})
	}

	issue := models.Issue{
		RepoID:    repo.ID,
		Title:     req.Title,
		Body:      req.Body,
		CreatorID: user.ID,
		// Creator:   user, // GORM handles relation automatically via ID, no need to set object for create
		State:     "open",
	}

	if result := h.DB.Create(&issue); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create issue"})
	}

	// Reload to get Creator loaded
	h.DB.Preload("Creator").First(&issue, issue.ID)

	return c.Status(fiber.StatusCreated).JSON(issue)
}

func (h *IssueHandler) List(c *fiber.Ctx) error {
	repoName := c.Params("name")

	var repo models.Repository
	if result := h.DB.Where("name = ?", repoName).First(&repo); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Repository not found"})
	}

	var issues []models.Issue
	if result := h.DB.Preload("Creator").Preload("Assignee").Where("repo_id = ?", repo.ID).Order("created_at desc").Find(&issues); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch issues"})
	}

	return c.JSON(issues)
}

func (h *IssueHandler) Get(c *fiber.Ctx) error {
	repoName := c.Params("name")
	issueID := c.Params("id")

	var repo models.Repository
	if result := h.DB.Where("name = ?", repoName).First(&repo); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Repository not found"})
	}

	var issue models.Issue
	if result := h.DB.Preload("Creator").Preload("Assignee").Where("repo_id = ? AND id = ?", repo.ID, issueID).First(&issue); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Issue not found"})
	}

	return c.JSON(issue)
}

type UpdateIssueRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	State string `json:"state"`
}

func (h *IssueHandler) Update(c *fiber.Ctx) error {
	repoName := c.Params("name")
	issueID := c.Params("id")
	var req UpdateIssueRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	var repo models.Repository
	if result := h.DB.Where("name = ?", repoName).First(&repo); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Repository not found"})
	}

	var issue models.Issue
	if result := h.DB.Where("repo_id = ? AND id = ?", repo.ID, issueID).First(&issue); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Issue not found"})
	}

	if req.Title != "" {
		issue.Title = req.Title
	}
	if req.Body != "" {
		issue.Body = req.Body
	}
	if req.State != "" {
		issue.State = req.State
	}

	if result := h.DB.Save(&issue); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update issue"})
	}

	// Reload to get associations
	h.DB.Preload("Creator").Preload("Assignee").First(&issue, issue.ID)

	return c.JSON(issue)
}

type CreateCommentRequest struct {
	Body string `json:"body"`
}

func (h *IssueHandler) CreateComment(c *fiber.Ctx) error {
	repoName := c.Params("name")
	issueID := c.Params("id")
	var req CreateCommentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.Body == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Body is required"})
	}

	user := c.Locals("user").(models.User)

	var repo models.Repository
	if result := h.DB.Where("name = ?", repoName).First(&repo); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Repository not found"})
	}

	var issue models.Issue
	if result := h.DB.Where("repo_id = ? AND id = ?", repo.ID, issueID).First(&issue); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Issue not found"})
	}

	comment := models.Comment{
		IssueID: issue.ID,
		UserID:  user.ID,
		Body:    req.Body,
	}

	if result := h.DB.Create(&comment); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create comment"})
	}

	h.DB.Preload("User").First(&comment, comment.ID)

	return c.Status(fiber.StatusCreated).JSON(comment)
}

func (h *IssueHandler) ListComments(c *fiber.Ctx) error {
	repoName := c.Params("name")
	issueID := c.Params("id")

	var repo models.Repository
	if result := h.DB.Where("name = ?", repoName).First(&repo); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Repository not found"})
	}

	var issue models.Issue
	if result := h.DB.Where("repo_id = ? AND id = ?", repo.ID, issueID).First(&issue); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Issue not found"})
	}

	var comments []models.Comment
	if result := h.DB.Preload("User").Where("issue_id = ?", issue.ID).Order("created_at asc").Find(&comments); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch comments"})
	}

	return c.JSON(comments)
}
