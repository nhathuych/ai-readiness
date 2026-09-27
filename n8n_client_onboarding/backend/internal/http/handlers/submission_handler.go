package handlers

import (
	"backend/internal/http/middleware"
	"backend/internal/repositories"
	"backend/internal/services"
	"context"
	"errors"
	"net/mail"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
)

type SubmissionHandler struct {
	submissionRepo *repositories.SubmissionRepository
	n8nService     *services.N8NService
}

type CreateSubmissionRequest struct {
	ClientName      string `json:"client_name"`
	ClientEmail     string `json:"client_email"`
	ServicePackage  string `json:"service_package"`
	ProjectGoal     string `json:"project_goal"`
	DesiredTimeline string `json:"desired_timeline"`
	AssetsProvided  string `json:"assets_provided"`
}

func NewSubmissionHandler(submissionRepo *repositories.SubmissionRepository, n8nService *services.N8NService) *SubmissionHandler {
	return &SubmissionHandler{
		submissionRepo: submissionRepo,
		n8nService:     n8nService,
	}
}

func (h *SubmissionHandler) Create(c fiber.Ctx) error {
	currentUser, ok := middleware.CurrentUserFromContext(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	var body CreateSubmissionRequest

	if err := c.Bind().Body(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	form := normalizeCreateSubmission(body)
	if message := validateCreateSubmissionRequest(form); message != "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": message,
		})
	}

	readiness, err := h.n8nService.CheckReadiness(context.Background(), form)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"message": "n8n CheckReadiness failed", "error": err.Error(),
		})
	}

	submission, err := h.submissionRepo.Create(context.Background(),
		services.ReadinessToCreateSubmissionInput(currentUser.ID, form, readiness))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to create submission", "error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"submission": submission,
	})
}

func (h *SubmissionHandler) ListMine(c fiber.Ctx) error {
	currentUser, ok := middleware.CurrentUserFromContext(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	submissions, err := h.submissionRepo.ListByUserID(context.Background(), currentUser.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to load submissions",
		})
	}

	return c.JSON(fiber.Map{"submissions": submissions})
}

func (h *SubmissionHandler) GetMineByID(c fiber.Ctx) error {
	currentUser, ok := middleware.CurrentUserFromContext(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	submission, err := h.submissionRepo.FindByIDForUser(context.Background(), c.Params("id"), currentUser.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Submission not found",
			})

		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to load submission",
		})

	}

	return c.JSON(fiber.Map{"submission": submission})

}

func normalizeCreateSubmission(body CreateSubmissionRequest) services.N8NReadinessInput {
	return services.N8NReadinessInput{
		ClientName:      strings.TrimSpace(body.ClientName),
		ClientEmail:     strings.TrimSpace(strings.ToLower(body.ClientEmail)),
		ServicePackage:  strings.TrimSpace(body.ServicePackage),
		ProjectGoal:     strings.TrimSpace(body.ProjectGoal),
		DesiredTimeline: strings.TrimSpace(body.DesiredTimeline),
		AssetsProvided:  strings.TrimSpace(body.AssetsProvided),
	}

}

func validateCreateSubmissionRequest(body services.N8NReadinessInput) string {

	switch {
	case body.ClientName == "":
		return "client name is required"
	case body.ClientEmail == "":
		return "client email is required"
	case body.ServicePackage == "":
		return "service package is required"
	case body.ProjectGoal == "":
		return "project goal is required"
	case body.DesiredTimeline == "":
		return "desired timeline is required"
	case body.AssetsProvided == "":
		return "asset readiness is required"
	}

	if _, err := mail.ParseAddress(body.ClientEmail); err != nil {
		return "Client email is invalid"
	}

	return ""

}
