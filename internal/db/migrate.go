package db

import (
	"context"
	"fmt"
	"io/fs"
	"sort"
	"strings"
)

// Migrate applies pending *.up.sql files from fsys, in name order.
//
// This is a tiny embedded migrator (not golang-migrate). Each file is one or
// more SQL statements. Statements are split on semicolons; do not put a
// semicolon inside a string literal. Down files (*.down.sql) are ignored.
func (p *Pool) Migrate(ctx context.Context, fsys fs.FS) error {
	if p == nil || p.inner == nil {
		return fmt.Errorf("db pool is nil")
	}

	if _, err := p.inner.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
		)
	`); err != nil {
		return fmt.Errorf("create schema_migrations: %w", err)
	}

	entries, err := fs.Glob(fsys, "*.up.sql")
	if err != nil {
		return fmt.Errorf("list migrations: %w", err)
	}
	sort.Strings(entries)

	for _, name := range entries {
		version := strings.TrimSuffix(name, ".up.sql")

		var applied bool
		if err := p.inner.QueryRow(ctx,
			`SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)`,
			version,
		).Scan(&applied); err != nil {
			return fmt.Errorf("check migration %s: %w", version, err)
		}
		if applied {
			continue
		}

		body, err := fs.ReadFile(fsys, name)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", name, err)
		}

		if err := p.applyFile(ctx, version, string(body)); err != nil {
			return err
		}
	}

	return nil
}

func (p *Pool) applyFile(ctx context.Context, version, body string) error {
	tx, err := p.inner.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin migration %s: %w", version, err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	for _, stmt := range splitSQL(body) {
		if _, err := tx.Exec(ctx, stmt); err != nil {
			return fmt.Errorf("apply migration %s: %w", version, err)
		}
	}

	if _, err := tx.Exec(ctx,
		`INSERT INTO schema_migrations (version) VALUES ($1)`,
		version,
	); err != nil {
		return fmt.Errorf("record migration %s: %w", version, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit migration %s: %w", version, err)
	}
	return nil
}

func splitSQL(body string) []string {
	parts := strings.Split(body, ";")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		stmt := strings.TrimSpace(part)
		if stmt == "" {
			continue
		}
		if stripSQLComments(stmt) == "" {
			continue
		}
		out = append(out, stmt)
	}
	return out
}

func stripSQLComments(stmt string) string {
	var b strings.Builder
	for _, line := range strings.Split(stmt, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "--") {
			continue
		}
		b.WriteString(trimmed)
		b.WriteByte(' ')
	}
	return strings.TrimSpace(b.String())
}
