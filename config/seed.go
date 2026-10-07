package config

import (
	"embed"
	"errors"
	"io/fs"
	"os"
	"path"
)

// The seed scripts ship inside the binary so a database can be created from
// wherever the binary runs. They used to be read from ./config/ only, which
// fails for anything not started from the repository root - a Docker image
// that carries just the binary and its settings file, for one.
//
//go:embed db.sql pg.sql db-begin-mysql.sql db-end-mysql.sql
var seedFS embed.FS

// ReadSeed returns the seed script at name (e.g. "config/db.sql").
// A file on disk wins, so a deployment that edits its own copy keeps doing
// so; the embedded copy is the fallback when there is none.
func ReadSeed(name string) ([]byte, error) {
	b, err := os.ReadFile(name)
	if err == nil || !errors.Is(err, fs.ErrNotExist) {
		return b, err
	}
	if emb, embErr := seedFS.ReadFile(path.Base(name)); embErr == nil {
		return emb, nil
	}
	return nil, err
}
