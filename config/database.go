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

// mysql config
type database struct {
	Host                  string `envconfig:"MYSQL_HOST" default:"localhost" prompt:"Enter mysql host"`
	Port                  string `envconfig:"MYSQL_PORT" default:"3306" prompt:"Enter mysql port"`
	DBName                string `envconfig:"MYSQL_DB_NAME" default:"rap_c" prompt:"Enter database name"`
	Username              string `envconfig:"MYSQL_USERNAME" default:"" prompt:"Enter mysql username"`
	Password              string `envconfig:"MYSQL_PASSWORD" default:"" prompt:"Enter mysql password" secret:"true"`
	LogMode               int    `envconfig:"MYSQL_LOG_MODE" default:"1" prompt:"Enter gorm log mode 1-4"`
	ParseTime             bool   `envconfig:"MYSQL_PARSE_TIME" default:"true" prompt:"Parse mysql time to local"`
	Charset               string `envconfig:"MYSQL_CHARSET" default:"utf8mb4" prompt:"Enter mysql database charset"`
	Loc                   string `envconfig:"MYSQL_LOC" default:"Local" prompt:"Enter mysql local time"`
	MaxLifetimeConnection int    `envconfig:"MYSQL_MAX_LIFETIME_CONNECTION" default:"10" prompt:"Enter mysql maximum amount of time a connection may be reused, in minute"`
	MaxOpenConnection     int    `envconfig:"MYSQL_MAX_OPEN_CONNECTION" default:"50" prompt:"Enter mysql maximum number of open connections to the database"`
	MaxIdleConnection     int    `envconfig:"MYSQL_MAX_IDLE_CONNECTION" default:"10" prompt:"Enter mysql maximum number of connections in the idle connection pool"`
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
		fmt.Println(err)
		os.Exit(1)
	}

	// set configuration pooling connection
	mysqlDb, _ := db.DB()
	mysqlDb.SetMaxOpenConns(dbConfig.MaxOpenConnection)
	mysqlDb.SetConnMaxLifetime(time.Duration(dbConfig.MaxLifetimeConnection) * time.Minute)
	mysqlDb.SetMaxIdleConns(dbConfig.MaxIdleConnection)

	return db
}
