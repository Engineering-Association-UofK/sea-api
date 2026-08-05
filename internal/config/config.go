package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

var App Config

type Config struct {
	Port         string `env:"PORT" envDefault:"8000"`
	LoggingLevel int    `env:"LOGGING_LEVEL" envDefault:"0"`

	// Resouses directories relative to where the backend server is launched in the file system
	ResourcesDir  string `env:"RESOURCES_DIR" envDefault:"./resources"`
	MigrationsDir string `env:"MIGRATIONS_DIR" envDefault:"file://db/migrations"`

	// Database details
	DbHost       string `env:"DB_HOST" envDefault:"localhost"`
	DbPort       string `env:"DB_PORT" envDefault:"3306"`
	DbUsername   string `env:"DB_USERNAME" envDefault:"root"`
	DbPassword   string `env:"DB_PASSWORD" envDefault:"root"`
	DbDatabase   string `env:"DB_DATABASE" envDefault:"sea_db"`
	DbSkipVerify bool   `env:"DB_SKIP_VERIFY" envDefault:"true"` // Whether to skip SSL verification

	// STMP details for sending emails
	MailHost string `env:"MAIL_HOST" envDefault:"smtp.gmail.com"`
	MailPort string `env:"MAIL_PORT" envDefault:"587"`
	MailUser string `env:"MAIL_USER" required:"true"`
	MailPass string `env:"MAIL_PASS" required:"true"`

	// Secret keys
	JwtSecret  string `env:"JWT_SECRET" required:"true"`
	SecretSalt string `env:"SECRET_SALT" required:"true"`

	///////////////////////////
	// ### S3 Properties ### //
	///////////////////////////

	// Access keys for the S3 service
	S3AccessKey string `env:"S3_ACCESS_KEY" required:"true"`
	S3SecretKey string `env:"S3_SECRET_KEY" required:"true"`

	StoreUrl      string `env:"STORE_S3_URL" envDefault:"http://localhost:8333"`     // The url for the S3 store, for image links generation
	StoreS3ApiUrl string `env:"STORE_S3_API_URL" envDefault:"http://localhost:8333"` // S3 url relative to the backend server
}

func Load() error {
	godotenv.Load()        // godotenv library to load .env files
	return env.Parse(&App) // Parse the .env file into the Config struct
}
