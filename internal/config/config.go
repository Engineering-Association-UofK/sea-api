package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

var App Config

type Config struct {
	Port         string `env:"PORT" envDefault:"8000"`
	LoggingLevel string `env:"LOGGING_LEVEL" envDefault:"INFO"`

	MainUrl   string `env:"MAIN_URL" envDefault:"http://localhost:5173"`
	ApiUrl    string `env:"API_URL" envDefault:"http://localhost:8000"`
	StoreUrl  string `env:"STORE_URL" envDefault:"http://localhost:8333"`
	HelperUrl string `env:"HELPER_URL" envDefault:"http://localhost:8888"`

	ResourcesDir  string `env:"RESOURCES_DIR" envDefault:"./resources"`
	MigrationsDir string `env:"MIGRATIONS_DIR" envDefault:"file://db/migrations"`

	DbName       string `env:"DB_NAME" envDefault:"mysql"`
	DbDSN        string `env:"DB_DSN" envDefault:"localhost"`
	DbHost       string `env:"DB_HOST" envDefault:"localhost"`
	DbPort       string `env:"DB_PORT" envDefault:"3306"`
	DbUsername   string `env:"DB_USERNAME" envDefault:"root"`
	DbPassword   string `env:"DB_PASSWORD" envDefault:"root"`
	DbDatabase   string `env:"DB_DATABASE" envDefault:"sea_db"`
	DbSkipVerify bool   `env:"DB_SKIP_VERIFY" envDefault:"true"`

	MailHost string `env:"MAIL_HOST" envDefault:"smtp.gmail.com"`
	MailPort string `env:"MAIL_PORT" envDefault:"587"`
	MailUser string `env:"MAIL_USER" required:"true"`
	MailPass string `env:"MAIL_PASS" required:"true"`

	CloudinaryUrl       string `env:"CLOUDINARY_URL" required:"true"`
	CloudinaryApiKey    string `env:"CLOUDINARY_API_KEY" required:"true"`
	CloudinaryApiSecret string `env:"CLOUDINARY_API_SECRET" required:"true"`

	JwtSecret        string `env:"JWT_SECRET" required:"true"`
	SecretSalt       string `env:"SECRET_SALT" required:"true"`
	KeystorePassword string `env:"KEYSTORE_PASSWORD" required:"true"`
	keystorePath     string `env:"KEYSTORE_PATH" envDefault:"/app/certs/sea_key.p12"`

	StoreMasterURL string `env:"STORE_MASTER_URL" envDefault:"http://localhost:9333"`
	StorePublicURL string `env:"STORE_PUBLIC_URL" envDefault:"http://localhost:8080"`
	StoreFilerUrl  string `env:"STORE_FILER_URL" envDefault:"http://localhost:8888"`
	StoreS3ApiUrl  string `env:"STORE_S3API_URL" envDefault:"http://localhost:8333"`

	S3AccessKey string `env:"S3_ACCESS_KEY" required:"true"`
	S3SecretKey string `env:"S3_SECRET_KEY" required:"true"`
}

func Load() error {
	godotenv.Load()
	return env.Parse(&App)
}
