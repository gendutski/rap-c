package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port              int  `envconfig:"HTTP_PORT" default:"8080" prompt:"Enter port to serve http"`
	AutoMigrateDB     bool `envconfig:"AUTO_MIGRATE_DB" default:"false" prompt:"Auto migrate database tables"`
	EnableDebug       bool `envconfig:"ENABLE_DEBUG" default:"false" prompt:"Enable debug to show error received"`
	EnableWarnFileLog bool `envconfig:"ENABLE_WARN_FILE_LOG" default:"false" prompt:"Enable log for warning type error (eg: http bad request error)"`

	// jwt
	JwtSecret              string `envconfig:"JWT_SECRET" default:"secret" prompt:"Enter secret to generate JWT token"`
	JwtExpirationInMinutes int    `envconfig:"JWT_EXPIRATION_IN_MINUTES" default:"60" prompt:"Enter token expired in minute"`
	JwtUserContextKey      string

	// smtp
	MailHost        string `envconfig:"MAIL_HOST" prompt:"Enter smtp mail host"`
	MailPort        int    `envconfig:"MAIL_PORT" default:"465" prompt:"Enter smtp mail port"`
	MailUser        string `envconfig:"MAIL_USER" prompt:"Enter smtp mail user"`
	MailPassword    string `envconfig:"MAIL_PASSWORD" prompt:"Enter smtp mail password" secret:"true"`
	MailFromName    string `envconfig:"MAIL_FROM_NAME" default:"Rap-C" prompt:"Enter name to send email"`
	MailFromAddress string `envconfig:"MAIL_FROM_ADDRESS" prompt:"Enter email address to send email"`
}

func GetConfig() Config {
	var cfg Config
	err := godotenv.Overload()
	if err != nil {
		log.Println(err)
	}
	envconfig.MustProcess("", &cfg)
	cfg.JwtUserContextKey = "user" // context key to get from echo
	return cfg
}
