package postgres

import (
	"dengovie/internal/utils/env"
	"embed"
	"log"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func (r *Repo) MigrationsUp() {
	env.InitEnvs(nil)

	// setup database
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(r.db, "migrations"); err != nil {
		panic(err)
	}

	log.Println("migrations successfully up")
}
