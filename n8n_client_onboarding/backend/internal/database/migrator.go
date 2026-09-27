package database

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func RunMigrations(ctx context.Context, db *pgxpool.Pool) error {

	if _, err := db.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS schema_migrations (
	version VARCHAR PRIMARY KEY,
	applied_at TIMESTAMP NOT NULL DEFAULT NOW()	
	)
	`); err != nil {
		return fmt.Errorf("create schema migration: %w", err)
	}

	entries, err := fs.ReadDir(migrationFiles, "migrations")
	if err != nil {
		return fmt.Errorf("read migrations dir: %w", err)
	}

	var filenames []string

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if strings.HasSuffix(entry.Name(), ".sql") {
			filenames = append(filenames, entry.Name())
		}
	}

	sort.Strings(filenames)

	for _, filename := range filenames {
		var alreadyApplied bool

		if err := db.QueryRow(ctx, `
		SELECT EXISTS (
		SELECT 1
		FROM schema_migrations
		WHERE version = $1
		)
		`, filename).Scan(&alreadyApplied); err != nil {
			return fmt.Errorf("check migration %s: %w", filename, err)
		}

		if alreadyApplied {
			continue
		}

		sqlBytes, err := migrationFiles.ReadFile("migrations/" + filename)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", filename, err)
		}

		tx, err := db.Begin(ctx)
		if err != nil {
			return fmt.Errorf("begin tx")
		}

		if _, err := tx.Exec(ctx, string(sqlBytes)); err != nil {

			_ = tx.Rollback(ctx)

			return fmt.Errorf("execute migration %s: %w", filename, err)
		}

		if _, err := tx.Exec(ctx, `
		INSERT INTO schema_migrations (version)
		VALUES ($1)
		`, filename); err != nil {
			_ = tx.Rollback(ctx)

			return fmt.Errorf("record migration %s: %w", filename, err)
		}

		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("commit migration %s: %w", filename, err)
		}
	}

	return nil

}
