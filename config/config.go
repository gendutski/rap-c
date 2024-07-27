package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	GuestUsername         string = "Guest"
	GuestEmail            string = "guest@example.com"
	GuestPassword         string = "$GuestP@ssw0rD!"
	AppName               string = "Rap-C"
	EchoJwtUserContextKey string = "user"
	EchoTokenContextKey   string = "token"
)

type LogMode int

const (
	LogModeErrorOnly LogMode = iota + 1
	LogModeErrorAndWarnOnly
	LogModeAll
)

// immutable config
type config struct {
	Port               int     `envconfig:"HTTP_PORT" default:"8080" prompt:"Enter port to serve http"`
	EnableDebug        bool    `envconfig:"ENABLE_DEBUG" default:"false" prompt:"Enable debug to show error received"`
	LogMode            LogMode `envconfig:"LOG_MODE" default:"1" prompt:"Enter log mode (1:error, 2:error & warn, 3:all)"`
	EnableWarnFileLog  bool    `envconfig:"ENABLE_WARN_FILE_LOG" default:"false" prompt:"Enable log for warning type error (eg: http bad request error)"`
	EnableGuestLogin   bool    `envconfig:"ENABLE_GUEST_LOGIN" default:"false" prompt:"Enable guest login"`
	AutoReloadTemplate bool    `envconfig:"AUTO_RELOAD_TEMPLATE" default:"false" prompt:"Auto reload template (not recommended for production)"`
	SessionKey         string  `envconfig:"SESSION_KEY" default:"session secret" prompt:"Enter http session key secret"`
	AppURL             string  `envconfig:"APP_URL" default:"http://localhost:8080" prompt:"Enter website location"`

	// jwt
	JwtSecret              string `envconfig:"JWT_SECRET" default:"secret" prompt:"Enter secret to generate JWT token"`
	JwtExpirationInMinutes int    `envconfig:"JWT_EXPIRATION_IN_MINUTES" default:"60" prompt:"Enter token expired in minute"`
	JwtRememberInDays      int    `envconfig:"JWT_REMEMBER_IN_DAYS" default:"30" prompt:"Enter token remember(for remember login) in days"`

	// smtp
	MailHost          string `envconfig:"MAIL_HOST" default:"smtp.gmail.com" prompt:"Enter smtp server host"`
	MailPort          int    `envconfig:"MAIL_PORT" default:"465" prompt:"Enter smtp server port"`
	MailUser          string `envconfig:"MAIL_USER" prompt:"Enter smtp server user"`
	MailPassword      string `envconfig:"MAIL_PASSWORD" prompt:"Enter smtp server password" secret:"true"`
	MailSenderName    string `envconfig:"MAIL_SENDER_NAME" default:"Rap-C" prompt:"Enter email sender name"`
	MailSenderAddress string `envconfig:"MAIL_SENDER_ADDRESS" prompt:"Enter email sender address"`

	// first user installation
	FirstUserUsername string `envconfig:"FIRST_USER_USERNAME" default:"gendutski" prompt:"Enter first user username"`
	FirstUserFullName string `envconfig:"FIRST_USER_FULL_NAME" default:"Firman Darmawan" prompt:"Enter first user full name"`
	FirstUserEmail    string `envconfig:"FIRST_USER_EMAIL" default:"mvp.firman.darmawan@gmail.com" prompt:"Enter first user email"`
	FirstUserPassword string `envconfig:"FIRST_USER_PASSWORD" default:"password" prompt:"Enter first user password" secret:"true"`

	// mysql
	MysqlHost                  string `envconfig:"MYSQL_HOST" default:"localhost" prompt:"Enter mysql host"`
	MysqlPort                  int    `envconfig:"MYSQL_PORT" default:"3306" prompt:"Enter mysql port"`
	MysqlDBName                string `envconfig:"MYSQL_DB_NAME" default:"rap_c" prompt:"Enter database name"`
	MysqlUsername              string `envconfig:"MYSQL_USERNAME" default:"" prompt:"Enter mysql username"`
	MysqlPassword              string `envconfig:"MYSQL_PASSWORD" default:"" prompt:"Enter mysql password" secret:"true"`
	MysqlLogMode               int    `envconfig:"MYSQL_LOG_MODE" default:"1" prompt:"Enter gorm log mode 1-4"`
	MysqlParseTime             bool   `envconfig:"MYSQL_PARSE_TIME" default:"true" prompt:"Parse mysql time to local"`
	MysqlCharset               string `envconfig:"MYSQL_CHARSET" default:"utf8mb4" prompt:"Enter mysql database charset"`
	MysqlLoc                   string `envconfig:"MYSQL_LOC" default:"Local" prompt:"Enter mysql local time"`
	MysqlMaxLifetimeConnection int    `envconfig:"MYSQL_MAX_LIFETIME_CONNECTION" default:"10" prompt:"Enter mysql maximum amount of time a connection may be reused, in minute"`
	MysqlMaxOpenConnection     int    `envconfig:"MYSQL_MAX_OPEN_CONNECTION" default:"50" prompt:"Enter mysql maximum number of open connections to the database"`
	MysqlMaxIdleConnection     int    `envconfig:"MYSQL_MAX_IDLE_CONNECTION" default:"10" prompt:"Enter mysql maximum number of connections in the idle connection pool"`
}

type Config struct {
	config *config
}

func InitConfig() *Config {
	var cfg config
	err := godotenv.Overload()
	if err != nil {
		log.Println(err)
	}
	envconfig.MustProcess("", &cfg)
	return &Config{config: &cfg}
}

// for testing purposes
func InitTestConfig(payload map[string]string) *Config {
	for key, val := range payload {
		os.Setenv(key, val)
	}
	var cfg config
	envconfig.MustProcess("", &cfg)
	return &Config{config: &cfg}
}

// render url from app url and path
func (cfg *Config) URL(path string) string {
	if cfg.config.AppURL == "" {
		return path
	}
	var result = cfg.config.AppURL
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

// connect to db
func (cfg *Config) ConnectDB() *gorm.DB {
	// construct connection string
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%+v&loc=%s",
		cfg.config.MysqlUsername,
		cfg.config.MysqlPassword,
		cfg.config.MysqlHost,
		cfg.config.MysqlPort,
		cfg.config.MysqlDBName,
		cfg.config.MysqlCharset,
		cfg.config.MysqlParseTime,
		cfg.config.MysqlLoc)

	// open mysql connection
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.LogLevel(cfg.config.MysqlLogMode)),
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// set configuration pooling connection
	mysqlDb, _ := db.DB()
	mysqlDb.SetMaxOpenConns(cfg.config.MysqlMaxOpenConnection)
	mysqlDb.SetConnMaxLifetime(time.Duration(cfg.config.MysqlMaxLifetimeConnection) * time.Minute)
	mysqlDb.SetMaxIdleConns(cfg.config.MysqlMaxIdleConnection)

	return db
}

// ----------------------private to public field-----------------------------\\
func (cfg *Config) Port() int                { return cfg.config.Port }
func (cfg *Config) EnableDebug() bool        { return cfg.config.EnableDebug }
func (cfg *Config) LogMode() LogMode         { return cfg.config.LogMode }
func (cfg *Config) EnableWarnFileLog() bool  { return cfg.config.EnableWarnFileLog }
func (cfg *Config) EnableGuestLogin() bool   { return cfg.config.EnableGuestLogin }
func (cfg *Config) AutoReloadTemplate() bool { return cfg.config.AutoReloadTemplate }
func (cfg *Config) SessionKey() string       { return cfg.config.SessionKey }
func (cfg *Config) AppURL() string           { return cfg.config.AppURL }

// jwt
func (cfg *Config) JwtSecret() string           { return cfg.config.JwtSecret }
func (cfg *Config) JwtExpirationInMinutes() int { return cfg.config.JwtExpirationInMinutes }
func (cfg *Config) JwtRememberInDays() int      { return cfg.config.JwtRememberInDays }

// smtp
func (cfg *Config) MailHost() string          { return cfg.config.MailHost }
func (cfg *Config) MailPort() int             { return cfg.config.MailPort }
func (cfg *Config) MailUser() string          { return cfg.config.MailUser }
func (cfg *Config) MailPassword() string      { return cfg.config.MailPassword }
func (cfg *Config) MailSenderName() string    { return cfg.config.MailSenderName }
func (cfg *Config) MailSenderAddress() string { return cfg.config.MailSenderAddress }

// first user installation
func (cfg *Config) FirstUserUsername() string { return cfg.config.FirstUserUsername }
func (cfg *Config) FirstUserFullName() string { return cfg.config.FirstUserFullName }
func (cfg *Config) FirstUserEmail() string    { return cfg.config.FirstUserEmail }
func (cfg *Config) FirstUserPassword() string { return cfg.config.FirstUserPassword }
