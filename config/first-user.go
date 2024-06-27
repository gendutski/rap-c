package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// will create user with this value if users table is empty
type FirstUser struct {
	FullName string `envconfig:"FIRST_USER_FULL_NAME" default:"Firman Darmawan" prompt:"Enter user full name"`
	Email    string `envconfig:"FIRST_USER_EMAIL" default:"mvp.firman.darmawan@gmail.com" prompt:"Enter user email"`
	Password string `envconfig:"FIRST_USER_PASSWORD" default:"password" prompt:"Enter user password" secret:"true"`
}

func GetFirstUser() FirstUser {
	var cfg FirstUser
	err := godotenv.Overload()
	if err != nil {
		log.Println(err)
	}
	envconfig.MustProcess("", &cfg)
	return cfg
}
