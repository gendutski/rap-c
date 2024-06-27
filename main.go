package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"rap-c/app/entity"
	"rap-c/app/helper"
	"rap-c/config"
	"time"

	"gorm.io/gorm"
)

var installMode = flag.Bool("install", false, "run for installation .env only")

func main() {
	flag.Parse()

	if *installMode {
		config.GenerateDotEnv()
	} else {
		serve()
	}
}

func serve() {
	cfg := config.GetConfig()
	db := config.Connect()

	// auto migrate db
	migrateDB(db, cfg.EnableGuestLogin)
}

func migrateDB(db *gorm.DB, enableGuestLogin bool) {
	start := time.Now()

	// migrate tables
	log.Println("Start migrate db")
	db.AutoMigrate(&entity.User{})
	db.AutoMigrate(&entity.Unit{})
	db.AutoMigrate(&entity.Ingredient{})
	db.AutoMigrate(&entity.IngredientConvertionUnit{})
	db.AutoMigrate(&entity.Recipe{})
	db.AutoMigrate(&entity.RecipeIngredient{})
	db.AutoMigrate(&entity.StockMovement{})
	db.AutoMigrate(&entity.Product{})
	db.AutoMigrate(&entity.Account{})
	db.AutoMigrate(&entity.Transaction{})
	log.Printf("Migrate done in %v", time.Since(start))

	// check non guest user
	log.Println("Check non guest user")
	var totalNonGuestUser int64
	err := db.Model(entity.User{}).Where("is_guest = ? and disabled = ?", false, false).Count(&totalNonGuestUser).Error
	if err != nil {
		fmt.Println("Error get total non guest user:", err)
		os.Exit(1)
	}
	// create non guest user
	if totalNonGuestUser < 1 {
		log.Println("Create first non guest user")
		firstUser := config.GetFirstUser()
		pass, err := helper.EncryptPassword(firstUser.Password)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		user := entity.User{
			Username:           firstUser.Username,
			FullName:           firstUser.FullName,
			Email:              firstUser.Email,
			Password:           pass,
			PasswordMustChange: true,
			CreatedBy:          config.SystemUsername,
			UpdatedBy:          config.SystemUsername,
		}
		err = db.Save(&user).Error
		if err != nil {
			fmt.Println("Error create first user:", err)
			os.Exit(1)
		}
	}

	// auto create guest user
	if enableGuestLogin {
		log.Println("Check guest user")
		var totalGuestUser int64
		err = db.Model(entity.User{}).Where("is_guest = ?", true).Count(&totalGuestUser).Error
		if err != nil {
			fmt.Println("Error get total guest user:", err)
			os.Exit(1)
		}

		// create guest user
		if totalGuestUser < 1 {
			log.Println("Create guest user")
			pass, err := helper.EncryptPassword(config.GuestPassword)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			user := entity.User{
				Username:  config.GuestUsername,
				FullName:  config.GuestUsername,
				Email:     config.GuestEmail,
				Password:  pass,
				IsGuest:   true,
				CreatedBy: config.SystemUsername,
				UpdatedBy: config.SystemUsername,
			}
			err = db.Save(&user).Error
			if err != nil {
				fmt.Println("Error create guest user:", err)
				os.Exit(1)
			}
		}
	}
}
