package storage

import (
	"crypto/tls"
	"fmt"
	"log/slog"
	"os"
	"sea-api/internal/config"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
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

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?tls=%s&parseTime=true",
		app.DbUsername,
		app.DbPassword,
		app.DbHost,
		app.DbPort,
		app.DbDatabase,
		app.DbName,
	)

	db, err := sqlx.Open(app.DbName, dsn)
	if err != nil {
		slog.Error("Failed to open database connection: ", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(time.Minute * 5)

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	slog.Info("Opening initial SQL file...")
	file, err := os.ReadFile("migrations/00001.init.up.sql")
	if err != nil {
		panic(err)
	}
	slog.Info("Running initial SQL script...")

	queries := strings.Split(string(file), ";")

	for _, query := range queries {
		query = strings.TrimSpace(query)
		if query == "" {
			continue
		}
		_, err = db.Exec(query)
		if err != nil {
			panic(err)
		}
	}

	slog.Info("MySQL connection ready.")
	return db
}
