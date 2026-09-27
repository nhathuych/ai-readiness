package handlers

import (
	"backend/internal/models"
	"backend/internal/repositories"
	"context"
	"errors"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
)

type AdminSubmissionHandler struct {
	submissionRepo *repositories.SubmissionRepository
}

type UpdateAdminStatusRequest struct {
	AdminStatus string `json:"admin_status"`
}

func NewAdminSubmissionHandler(submissionRepo *repositories.SubmissionRepository) *AdminSubmissionHandler {
	return &AdminSubmissionHandler{
		submissionRepo: submissionRepo,
	}
}

func (h *AdminSubmissionHandler) ListAllSubmissions(c fiber.Ctx) error {
	submissions, err := h.submissionRepo.ListAll(context.Background())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to load submissions",
		})
	}

	return c.JSON(fiber.Map{"submissions": submissions})
}

func (h *AdminSubmissionHandler) UpdateStatus(c fiber.Ctx) error {
	var body UpdateAdminStatusRequest
	if err := c.Bind().Body(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalod request body",
		})
	}

	adminStatus := strings.TrimSpace(body.AdminStatus)
	if !isAllowedAdminStatus(adminStatus) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalod AdmiN status",
		})
	}

	submission, err := h.submissionRepo.UpdateAdminStatus(context.Background(), c.Params("id"), adminStatus)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Submission not found",
			})

		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to update submission status",
		})

	}

	return c.JSON(fiber.Map{"submission": submission})

}

func isAllowedAdminStatus(status string) bool {
	return status == models.AdminStatusPending ||
		status == models.AdminStatusApproved ||
		status == models.AdminStatusRejected
}
