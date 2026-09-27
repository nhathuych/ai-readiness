package routes

import (
	"backend/internal/http/handlers"
	"backend/internal/http/middleware"

	"github.com/gofiber/fiber/v3"
)

type RouteDependencies struct {
	AuthHandler            *handlers.AuthHandler
	AuthMiddleware         *middleware.AuthMiddleware
	SubmissionHandler      *handlers.SubmissionHandler
	AdminSubmissionHandler *handlers.AdminSubmissionHandler
}

func Register(app *fiber.App, deps RouteDependencies) {

	auth := app.Group("/auth")
	auth.Get("/google", deps.AuthHandler.StartGoogleAuth)
	auth.Get("/google/callback", deps.AuthHandler.GoogleCallback)
	auth.Post("/logout", deps.AuthHandler.Logout)
	auth.Get("/me", deps.AuthMiddleware.RequireAuth(), deps.AuthHandler.GetUserInfo)

	submissions := app.Group("/submissions", deps.AuthMiddleware.RequireAuth())
	submissions.Post("/", deps.SubmissionHandler.Create)
	submissions.Get("/", deps.SubmissionHandler.ListMine)
	submissions.Get("/:id", deps.SubmissionHandler.GetMineByID)

	admin := app.Group("/admin", deps.AuthMiddleware.RequireAuth(), deps.AuthMiddleware.RequireAdmin())
	admin.Get("/me", deps.AuthHandler.GetUserInfo)
	admin.Get("/submissions", deps.AdminSubmissionHandler.ListAllSubmissions)
	admin.Patch("/submissions/:id/status", deps.AdminSubmissionHandler.UpdateStatus)
}
