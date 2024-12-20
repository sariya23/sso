package main

import (
	"errors"
	"flag"
	"fmt"
	"sso/interanal/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/ilyakaznacheev/cleanenv"
)

func main() {
	var migrationsPath, migrationsTable string
	var cfg config.ConfigDataBase
	err := cleanenv.ReadConfig("config/db.yaml", &cfg)
	if err != nil {
		panic(err)
	}
	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "name of migrations table")
	flag.Parse()

	if migrationsPath == "" {
		panic("storage-path is required")
	}
	m, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf("pgx5://%s:%s@%s:%s/%s?x-migrations-table=%s&sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, migrationsTable),
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")
			return
		}
		panic(err)
	}
	fmt.Println("migrations applied")
}
