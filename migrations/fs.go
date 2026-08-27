// Package migrations embeds SQL files applied at process start.
package migrations

import "embed"

// FS is the embedded migration directory (*.sql).
//
//go:embed *.sql
var FS embed.FS
