package handlers

import (
	"backend/internal/config"
	"backend/internal/http/middleware"
	"backend/internal/models"
	"backend/internal/repositories"
	"backend/internal/services"
	"context"

	"github.com/gofiber/fiber/v3"
)

type AuthHandler struct {
	config config.Config

	// authservice
	authservice *services.AuthService

	userRepo *repositories.UserRepository
}

func NewAuthHandler(cfg config.Config, authService *services.AuthService,
	userRepo *repositories.UserRepository,
) *AuthHandler {
	return &AuthHandler{
		config:      cfg,
		authservice: authService,
		userRepo:    userRepo,
	}
}

func (h *AuthHandler) StartGoogleAuth(c fiber.Ctx) error {
	state, err := h.authservice.GenerateStateToken()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to generate the state token",
		})
	}

	h.authservice.SetOauthStateCookie(c, state)

	return c.Redirect().To(h.authservice.BuildGoogleAuthUrl(state))
}

func (h *AuthHandler) GoogleCallback(c fiber.Ctx) error {
	stateFromQuery := c.Query("state")

	stateFromCookie := h.authservice.ReadOauthStateCookie(c)

	if stateFromQuery == "" || stateFromCookie == "" || stateFromQuery != stateFromCookie {
		h.authservice.ClearOauthStateCookie(c)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid oauth state",
		})
	}

	h.authservice.ClearOauthStateCookie(c)

	// auth code from google's callback url
	// temp code and will be exchanged for a Google access token
	code := c.Query("code")
	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing auth code",
		})
	}

	googleUser, err := h.authservice.ExchangeCodeForGoogleUser(context.Background(), code)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"message": "Failed to fetch google user information",
			"error":   err.Error(),
		})
	}

	user, err := h.userRepo.UpsertByEmail(context.Background(), repositories.UpsertUserInput{
		Email:     googleUser.Email,
		Name:      googleUser.Name,
		AvatarURL: googleUser.Picture,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to upsert user in our DB",
		})
	}

	token, err := h.authservice.SignJWT(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to sign jwt",
		})
	}

	h.authservice.SetAuthCookie(c, token)

	return c.Redirect().To(h.config.FrontendURL + "/")
}

func (h *AuthHandler) GetUserInfo(c fiber.Ctx) error {
	currentUser, ok := c.Locals(middleware.ExtractCurrentUserLocalKey()).(*models.User)

	if !ok || currentUser == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	return c.JSON(fiber.Map{
		"user": currentUser,
	})
}

func (h *AuthHandler) Logout(c fiber.Ctx) error {
	h.authservice.ClearAuthCookie(c)

	return c.JSON(fiber.Map{
		"message": "Logged out successfully",
	})
}
