package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

const (
	SystemUsername string = "SYSTEM"
	GuestUsername  string = "Guest"
	GuestEmail     string = "guest@example.com"
	GuestPassword  string = "$GuestP@ssw0rD!"
	AppName        string = "Rap-C"
)

type Config struct {
	Port               int    `envconfig:"HTTP_PORT" default:"8080" prompt:"Enter port to serve http"`
	EnableDebug        bool   `envconfig:"ENABLE_DEBUG" default:"false" prompt:"Enable debug to show error received"`
	EnableWarnFileLog  bool   `envconfig:"ENABLE_WARN_FILE_LOG" default:"false" prompt:"Enable log for warning type error (eg: http bad request error)"`
	EnableGuestLogin   bool   `envconfig:"ENABLE_GUEST_LOGIN" default:"false" prompt:"Enable guest login"`
	AutoReloadTemplate bool   `envconfig:"AUTO_RELOAD_TEMPLATE" default:"false" prompt:"Auto reload template (not recommended for production)"`
	SessionKey         string `envconfig:"SESSION_KEY" default:"session secret" prompt:"Enter http session key secret"`

	AppURL string `envconfig:"APP_URL" default:"http://localhost:8080" prompt:"Enter website location"`

	// jwt
	JwtSecret              string `envconfig:"JWT_SECRET" default:"secret" prompt:"Enter secret to generate JWT token"`
	JwtExpirationInMinutes int    `envconfig:"JWT_EXPIRATION_IN_MINUTES" default:"60" prompt:"Enter token expired in minute"`
	JwtRememberInDays      int    `envconfig:"JWT_REMEMBER_IN_DAYS" default:"30" prompt:"Enter token remember(for remember login) in days"`
	JwtUserContextKey      string

	// smtp
	MailHost          string `envconfig:"MAIL_HOST" default:"smtp.gmail.com" prompt:"Enter smtp server host"`
	MailPort          int    `envconfig:"MAIL_PORT" default:"465" prompt:"Enter smtp server port"`
	MailUser          string `envconfig:"MAIL_USER" prompt:"Enter smtp server user"`
	MailPassword      string `envconfig:"MAIL_PASSWORD" prompt:"Enter smtp server password" secret:"true"`
	MailSenderName    string `envconfig:"MAIL_SENDER_NAME" default:"Rap-C" prompt:"Enter email sender name"`
	MailSenderAddress string `envconfig:"MAIL_SENDER_ADDRESS" prompt:"Enter email sender address"`
}

func (cfg Config) URL(path string) string {
	if cfg.AppURL == "" {
		return path
	}
	var result = cfg.AppURL
	if result[len(result)-1] != '/' {
		result += "/"
	}
	if path != "" {
		if path[0] == '/' {
			result += path[1:]
		} else {
			result += path
		}
	}
	return result
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
