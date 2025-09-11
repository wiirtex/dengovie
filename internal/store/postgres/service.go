package postgres

import (
	"database/sql"
	"dengovie/internal/store/types"
	"dengovie/internal/utils/env"
	"fmt"

	_ "github.com/lib/pq"
)

var _ types.Storage = (*Repo)(nil)

type Repo struct {
	db *sql.DB
}

func New(connStrings ...string) (*Repo, error) {

	connString := ""
	if len(connStrings) != 0 {
		connString = connStrings[0]
	} else {
		var err error
		connString, err = env.GetEnv(env.KeyPostgresConnString)
		if err != nil {
			panic(fmt.Sprintf("can not get connection string: %v", err))
		}
	}

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
