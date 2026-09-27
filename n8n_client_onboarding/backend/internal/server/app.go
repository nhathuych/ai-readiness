package server

import (
	"backend/internal/config"
	"backend/internal/http/handlers"
	"backend/internal/http/middleware"
	"backend/internal/http/routes"
	"backend/internal/repositories"
	"backend/internal/services"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	recovermw "github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(cfg config.Config, db *pgxpool.Pool) *fiber.App {

	app := fiber.New()

	app.Use(logger.New())

	app.Use(recovermw.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.FrontendURL},
		AllowCredentials: true,
		AllowMethods:     []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))

	userRepo := repositories.NewUserRepository(db)
	submissionRepo := repositories.NewSubmissionRepository(db)

	authService := services.NewAuthService(cfg)
	n8nService := services.NewN8NService(cfg)

	authHandler := handlers.NewAuthHandler(cfg, authService, userRepo)
	submissionHandler := handlers.NewSubmissionHandler(submissionRepo, n8nService)
	adminSubmissionHandler := handlers.NewAdminSubmissionHandler(submissionRepo)

	authMiddleware := middleware.NewAuthMiddleware(cfg, authService, userRepo)

	routes.Register(app, routes.RouteDependencies{
		AuthHandler:            authHandler,
		AuthMiddleware:         authMiddleware,
		SubmissionHandler:      submissionHandler,
		AdminSubmissionHandler: adminSubmissionHandler,
	})

	return app
}
