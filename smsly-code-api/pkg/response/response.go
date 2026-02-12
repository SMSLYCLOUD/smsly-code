package response

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// Standard success response
type APIResponse struct {
    Data  interface{} `json:"data,omitempty"`
    Meta  *Meta       `json:"meta,omitempty"`
    Error *APIError   `json:"error,omitempty"`
}

type Meta struct {
    Page    int `json:"page"`
    PerPage int `json:"per_page"`
    Total   int `json:"total"`
}

type APIError struct {
    Code    string      `json:"code"`
    Message string      `json:"message"`
    Details interface{} `json:"details,omitempty"`
}

func Success(c *fiber.Ctx, data interface{}) error {
    return c.JSON(APIResponse{Data: data})
}

func SuccessWithMeta(c *fiber.Ctx, data interface{}, meta *Meta) error {
    return c.JSON(APIResponse{Data: data, Meta: meta})
}

func Error(c *fiber.Ctx, status int, code, message string) error {
    return c.Status(status).JSON(APIResponse{
        Error: &APIError{Code: code, Message: message},
    })
}

func NotFound(c *fiber.Ctx, resource string) error {
    return Error(c, 404, "NOT_FOUND", resource+" not found")
}

func Unauthorized(c *fiber.Ctx) error {
    return Error(c, 401, "UNAUTHORIZED", "Authentication required")
}

func Forbidden(c *fiber.Ctx) error {
    return Error(c, 403, "FORBIDDEN", "Insufficient permissions")
}

func BadRequest(c *fiber.Ctx, message string) error {
    return Error(c, 400, "BAD_REQUEST", message)
}

func InternalError(c *fiber.Ctx, err error) error {
    // Log the actual error, return generic message
    log.Error().Err(err).Msg("Internal server error")
    return Error(c, 500, "INTERNAL_ERROR", "An internal error occurred")
}
