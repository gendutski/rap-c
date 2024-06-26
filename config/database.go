package config

import (
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// mysql config
type database struct {
	Host                  string `envconfig:"MYSQL_HOST" default:"localhost"`
	Port                  string `envconfig:"MYSQL_PORT" default:"3306"`
	DBName                string `envconfig:"MYSQL_DB_NAME" default:"fin_track"`
	Username              string `envconfig:"MYSQL_USERNAME" default:""`
	Password              string `envconfig:"MYSQL_PASSWORD" default:""`
	LogMode               int    `envconfig:"MYSQL_LOG_MODE" default:"0"`
	ParseTime             bool   `envconfig:"MYSQL_PARSE_TIME" default:"true"`
	Charset               string `envconfig:"MYSQL_CHARSET" default:"utf8mb4"`
	Loc                   string `envconfig:"MYSQL_LOC" default:"Local"`
	MaxLifetimeConnection int    `envconfig:"MYSQL_MAX_LIFETIME_CONNECTION" default:"10"`
	MaxOpenConnection     int    `envconfig:"MYSQL_MAX_OPEN_CONNECTION" default:"50"`
	MaxIdleConnection     int    `envconfig:"MYSQL_MAX_IDLE_CONNECTION" default:"10"`
}

func Connect() *gorm.DB {
	var dbConfig database
	// load .env file if exists
	err := godotenv.Overload()
	if err != nil {
		log.Println(err)
	}
	envconfig.MustProcess("", &dbConfig)

	// construct connection string
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%+v&loc=%s",
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DBName,
		dbConfig.Charset,
		dbConfig.ParseTime,
		dbConfig.Loc)

	// open mysql connection
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.LogLevel(dbConfig.LogMode)),
	})
	if err != nil {
		panic(err)
	}

	// set configuration pooling connection
	mysqlDb, _ := db.DB()
	mysqlDb.SetMaxOpenConns(dbConfig.MaxOpenConnection)
	mysqlDb.SetConnMaxLifetime(time.Duration(dbConfig.MaxLifetimeConnection) * time.Minute)
	mysqlDb.SetMaxIdleConns(dbConfig.MaxIdleConnection)

	return db
}
