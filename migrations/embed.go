// Package migrations holds the SQL schema changes, embedded into the binary.
package migrations

import "embed"

// FS contains every *.sql file in this folder, built into the program.
//
//go:embed *.sql
var FS embed.FS
