package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/SMSLYCLOUD/smsly-code/smsly-code-api/internal/config"
	"github.com/SMSLYCLOUD/smsly-code/smsly-code-api/internal/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type RepoHandler struct {
	DB     *gorm.DB
	Config *config.Config
}

func NewRepoHandler(db *gorm.DB, cfg *config.Config) *RepoHandler {
	return &RepoHandler{
		DB:     db,
		Config: cfg,
	}
}

type CreateRepoRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsPrivate   bool   `json:"is_private"`
}

func (h *RepoHandler) Create(c *fiber.Ctx) error {
	var req CreateRepoRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	user := c.Locals("user").(models.User)

	// Check if repo already exists
	var existingRepo models.Repository
	if result := h.DB.Where("name = ?", req.Name).First(&existingRepo); result.Error == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Repository name already taken"})
	}

	// 1. Call Rust service to create bare repo
	gitHost := os.Getenv("GIT_HOST")
	if gitHost == "" {
		gitHost = "localhost"
	}
	gitPort := os.Getenv("GIT_PORT")
	if gitPort == "" {
		gitPort = "8081"
	}

	gitServiceURL := fmt.Sprintf("http://%s:%s/repo", gitHost, gitPort)

	gitPayload := map[string]string{"name": req.Name}
	jsonPayload, _ := json.Marshal(gitPayload)

	resp, err := http.Post(gitServiceURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to contact Git engine", "details": err.Error()})
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Git engine failed to create repository"})
	}

	// 2. Create DB entry
	repo := models.Repository{
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     user.ID,
		IsPrivate:   req.IsPrivate,
		Owner:       user,
	}

	if result := h.DB.Create(&repo); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save repository metadata"})
	}

	return c.Status(fiber.StatusCreated).JSON(repo)
}

func (h *RepoHandler) List(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	var repos []models.Repository
	if result := h.DB.Where("owner_id = ?", user.ID).Find(&repos); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch repositories"})
	}

	return c.JSON(repos)
}

type Commit struct {
	ID          string `json:"id"`
	Message     string `json:"message"`
	Author      string `json:"author"`
	Date        string `json:"date"`
	MIPVerified bool   `json:"mip_verified"`
}

func (h *RepoHandler) GetCommits(c *fiber.Ctx) error {
	repoName := c.Params("name")

	gitHost := os.Getenv("GIT_HOST")
	if gitHost == "" { gitHost = "localhost" }
	gitPort := os.Getenv("GIT_PORT")
	if gitPort == "" { gitPort = "8081" }

	url := fmt.Sprintf("http://%s:%s/repo/%s/commits", gitHost, gitPort, repoName)

	resp, err := http.Get(url)
	if err != nil {
		// Log error for debugging
		fmt.Printf("Error contacting git engine: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to contact Git engine"})
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Forward error from git engine
		var errResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errResp)
		return c.Status(resp.StatusCode).JSON(errResp)
	}

	var commits []Commit
	// Need to check if body is empty or handle errors gracefully
	if err := json.NewDecoder(resp.Body).Decode(&commits); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to parse commits"})
	}

	return c.JSON(commits)
}

type TreeEntry struct {
	Name string `json:"name"`
	Kind string `json:"kind"`
	ID   string `json:"id"`
}

func (h *RepoHandler) GetTree(c *fiber.Ctx) error {
	repoName := c.Params("name")
	ref := c.Params("ref", "HEAD")
	path := c.Params("*")

	// If path is empty (root), we might not get it in params depending on route definition
	// Fiber's wildcard params include slashes.

	gitHost := os.Getenv("GIT_HOST")
	if gitHost == "" { gitHost = "localhost" }
	gitPort := os.Getenv("GIT_PORT")
	if gitPort == "" { gitPort = "8081" }

	// Construct URL carefully.
	// The Rust route is /repo/{name}/tree/{ref_name}/{path:.*}
	url := fmt.Sprintf("http://%s:%s/repo/%s/tree/%s/%s", gitHost, gitPort, repoName, ref, path)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error contacting git engine: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to contact Git engine"})
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errResp)
		return c.Status(resp.StatusCode).JSON(errResp)
	}

	var entries []TreeEntry
	if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to parse tree"})
	}

	return c.JSON(entries)
}
