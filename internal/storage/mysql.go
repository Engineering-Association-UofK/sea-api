package storage

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"sea-api/internal/config"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	mysqlMigrate "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

func NewMySQLConnection() *sqlx.DB {
	slog.Info("Starting MySQL connection...")
	app := &config.App

	mysql.RegisterTLSConfig(app.DbName, &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: app.DbSkipVerify,
		ServerName:         app.DbHost,
	})

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?tls=%s&parseTime=true&multiStatements=true",
		app.DbUsername,
		app.DbPassword,
		app.DbHost,
		app.DbPort,
		app.DbDatabase,
		app.DbName,
	)

	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		slog.Error("Failed to open database connection: ", "error", err)
		panic(err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(time.Minute * 5)

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	if app.DbName == "tidb" {
		_, err := db.Exec("SET SESSION tidb_skip_isolation_level_check=1")
		if err != nil {
			panic(err)
		}
	}

	runMigrations(db)

	slog.Info("Database connection and migrations ready.")
	return db
}

func runMigrations(db *sqlx.DB) {
	slog.Info("Running database migrations...")

	driver, err := mysqlMigrate.WithInstance(db.DB, &mysqlMigrate.Config{
		MigrationsTable: "schema_migrations",
	})
	if err != nil {
		slog.Error("Could not create migration driver", "error", err)
		panic(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:///app/db/migrations",
		"mysql",
		driver,
	)
	if err != nil {
		slog.Error("Could not initialize migrate instance", "error", err)
		panic(err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		slog.Error("Migration Up failed", "error", err)
		panic(err)
	}

	slog.Info("Migrations finished successfully.")
}
