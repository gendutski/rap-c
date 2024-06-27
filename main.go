package main

import (
	"flag"
	"log"
	"rap-c/app/entity"
	"rap-c/config"
	"time"

	"gorm.io/gorm"
)

var installMode = flag.Bool("install", false, "run for installation only")

func main() {
	flag.Parse()

	if *installMode {
		config.GenerateDotEnv()
	} else {
		serve()
	}
}

func serve() {
	db := config.Connect()

	// auto migrate db
	// if cfg.AutoMigrateDB {
	migrateDB(db)
	// }
}

func migrateDB(db *gorm.DB) {
	start := time.Now()
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
}
