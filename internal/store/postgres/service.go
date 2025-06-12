package postgres

import (
	"database/sql"
	"dengovie/internal/store/types"
	"fmt"

	_ "github.com/lib/pq"
)

var _ types.Storage = (*Repo)(nil)

type Repo struct {
	db *sql.DB
}

func New(connString string) (*Repo, error) {

	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	r := &Repo{
		db: db,
	}

	r.MigrationsUp()

	return r, nil
}
