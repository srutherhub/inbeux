package db

import (
	"context"
	"embed"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5"
)

//go:embed migrations/*
var MIGRATIONS_FOLDER embed.FS

func InitializeDb(ctx context.Context) (*pgx.Conn, func(), error) {
	postgresConnection := os.Getenv("DATABASE_URL")

	conn, err := pgx.Connect(ctx, postgresConnection)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	cleanup := func() {
		if err := conn.Close(context.Background()); err != nil {
			fmt.Fprintf(os.Stderr, "failed to close database connection: %v\n", err)
		}
	}

	return conn, cleanup, nil
}

func RunMigrations() {
	sourceDriver, err := iofs.New(MIGRATIONS_FOLDER, "migrations")

	if err != nil {
		panic(err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", sourceDriver, os.Getenv("DATABASE_URL"))

	if err != nil {
		panic(err)
	}

	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		panic(err)
	}

}
