package repositories

import (
	"backend/internal/models"
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SubmissionRepository struct {
	db *pgxpool.Pool
}

type CreateSubmissionInput struct {
	UserID                string
	ClientName            string
	ClientEmail           string
	ServicePackage        string
	ProjectGoal           string
	DesiredTimeline       string
	AssetsProvided        string
	ReadinessStatus       string
	MissingItems          []string
	AISummary             string
	RecommendedNextAction string
}

func NewSubmissionRepository(db *pgxpool.Pool) *SubmissionRepository {
	return &SubmissionRepository{db: db}
}

func (r *SubmissionRepository) Create(ctx context.Context, input CreateSubmissionInput) (*models.Submission, error) {
	var submission models.Submission

	err := r.db.QueryRow(ctx, `
			INSERT INTO submissions (
			user_id,
			client_name,
			client_email,
			service_package,
			project_goal,
			desired_timeline,
			assets_provided,
			readiness_status,
			missing_items,
			ai_summary,
			recommended_next_action
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING
			id,
			user_id,
			client_name,
			client_email,
			service_package,
			project_goal,
			desired_timeline,
			assets_provided,
			readiness_status,
			missing_items,
			ai_summary,
			recommended_next_action,
			admin_status,
			created_at,
			updated_at
	`,
		strings.TrimSpace(input.UserID),
		strings.TrimSpace(input.ClientName),
		strings.TrimSpace(strings.ToLower(input.ClientEmail)),
		strings.TrimSpace(input.ServicePackage),
		strings.TrimSpace(input.ProjectGoal),
		strings.TrimSpace(input.DesiredTimeline),
		strings.TrimSpace(input.AssetsProvided),
		strings.TrimSpace(input.ReadinessStatus),
		input.MissingItems,
		strings.TrimSpace(input.AISummary),
		strings.TrimSpace(input.RecommendedNextAction),
	).Scan(scanSubmissionFields(&submission)...)

	if err != nil {
		return nil, fmt.Errorf("create submission failed: %w", err)
	}

	return &submission, nil
}

func (r *SubmissionRepository) ListByUserID(ctx context.Context, userID string) ([]models.Submission, error) {
	rows, err := r.db.Query(ctx, selectSubmissionSQL()+`
	WHERE user_id = $1
	ORDER BY created_at DESC
	`, strings.TrimSpace(userID))
	if err != nil {
		return nil, fmt.Errorf("list submission by user failed")
	}

	defer rows.Close()

	return scanSubmissions(rows)

}

func (r *SubmissionRepository) ListAll(ctx context.Context) ([]models.Submission, error) {
	rows, err := r.db.Query(ctx, selectSubmissionSQL()+`
	ORDER BY created_at DESC
	`)

	if err != nil {
		return nil, fmt.Errorf("list all submissions failed")
	}

	defer rows.Close()

	return scanSubmissions(rows)
}

func (r *SubmissionRepository) FindByIDForUser(ctx context.Context, submissionID string, userID string) (*models.Submission, error) {
	var submission models.Submission

	err := r.db.QueryRow(ctx, selectSubmissionSQL()+`
	WHERE id = $1 AND user_id = $2
	`,
		strings.TrimSpace(submissionID),
		strings.TrimSpace(userID),
	).Scan(scanSubmissionFields(&submission)...)

	if err != nil {
		return nil, fmt.Errorf("find submission for user by id failed")
	}

	return &submission, nil
}

func (r *SubmissionRepository) UpdateAdminStatus(ctx context.Context, submissionID string, adminStatus string) (*models.Submission, error) {
	var submission models.Submission

	err := r.db.QueryRow(ctx, `
		UPDATE submissions
		SET admin_status = $2, updated_at = NOW()
		WHERE id = $1
		RETURNING
			id,
			user_id,
			client_name,
			client_email,
			service_package,
			project_goal,
			desired_timeline,
			assets_provided,
			readiness_status,
			missing_items,
			ai_summary,
			recommended_next_action,
			admin_status,
			created_at,
			updated_at
	`, strings.TrimSpace(submissionID), strings.TrimSpace(adminStatus)).Scan(
		scanSubmissionFields(&submission)...,
	)

	if err != nil {
		return nil, fmt.Errorf("Update status failed")
	}

	return &submission, nil
}

func scanSubmissions(rows pgx.Rows) ([]models.Submission, error) {
	submissions := make([]models.Submission, 0)

	for rows.Next() {
		var submission models.Submission

		if err := rows.Scan(scanSubmissionFields(&submission)...); err != nil {
			return nil, fmt.Errorf("scan submission failed")
		}

		// add the scanned submission into the submissions struct
		submissions = append(submissions, submission)
	}

	return submissions, rows.Err()
}

func scanSubmissionFields(submission *models.Submission) []any {
	return []any{
		&submission.ID,
		&submission.UserID,
		&submission.ClientName,
		&submission.ClientEmail,
		&submission.ServicePackage,
		&submission.ProjectGoal,
		&submission.DesiredTimeline,
		&submission.AssetsProvided,
		&submission.ReadinessStatus,
		&submission.MissingItems,
		&submission.AISummary,
		&submission.RecommendedNextAction,
		&submission.AdminStatus,
		&submission.CreatedAt,
		&submission.UpdatedAt,
	}
}

func selectSubmissionSQL() string {
	return `
		SELECT
			id,
			user_id,
			client_name,
			client_email,
			service_package,
			project_goal,
			desired_timeline,
			assets_provided,
			readiness_status,
			missing_items,
			ai_summary,
			recommended_next_action,
			admin_status,
			created_at,
			updated_at
		FROM submissions
	`
}
