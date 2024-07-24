package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	databaseentity "rap-c/app/entity/database-entity"
	"rap-c/app/handler/api"
	"rap-c/app/handler/middleware"
	"rap-c/app/handler/web"
	"rap-c/app/helper"
	authrepository "rap-c/app/repository/mysql/auth-repository"
	userrepository "rap-c/app/repository/mysql/user-repository"
	authusecase "rap-c/app/usecase/auth-usecase"
	mailusecase "rap-c/app/usecase/mail-usecase"
	sessionusecase "rap-c/app/usecase/session-usecase"
	userusecase "rap-c/app/usecase/user-usecase"
	"rap-c/config"
	"rap-c/route"
	"regexp"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
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
	// load config & connect db
	cfg := config.InitConfig()
	db := cfg.ConnectDB()

	// load route config
	router := config.InitRoute()

	// auto migrate db
	migrateDB(cfg, db)

	// load mysql repositories
	authRepo := authrepository.New(db)
	userRepo := userrepository.New(db)

	// load session store
	sessionStore := sessions.NewCookieStore([]byte(cfg.SessionKey()))

	// load modules
	authUsecase := authusecase.NewUsecase(cfg, authRepo)
	userUsecase := userusecase.NewUsecase(cfg, userRepo)
	mailUsecase := mailusecase.NewUsecase(cfg, router)
	sessionUsecase := sessionusecase.NewUsecase(cfg, router, sessionStore, authUsecase)

	// load api handler
	authAPI := api.NewAuthHandler(cfg, router, authUsecase, mailUsecase)
	userAPI := api.NewUserHandler(cfg, router, userUsecase, mailUsecase)

	// load web handler
	authWeb := web.NewAuthPage(cfg, router, authUsecase, sessionUsecase, mailUsecase)
	userWeb := web.NewUserPage(cfg, router, sessionUsecase, userUsecase, mailUsecase)
	dashboardWeb := web.NewDashboardPage(cfg, router, sessionUsecase)

	// init echo
	e := echo.New()

	e.Debug = cfg.EnableDebug()
	// custom http error handler
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		reg := regexp.MustCompile("^/api")
		if reg.MatchString(c.Request().RequestURI) {
			route.APIErrorHandler(e, err, c)
		} else {
			route.WebErrorHandler(e, err, c)
		}
	}
	// set general middleware
	e.Use(middleware.SetLog(cfg.LogMode(), cfg.EnableWarnFileLog()))
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.TimeoutWithConfig(echomiddleware.TimeoutConfig{
		ErrorMessage: "request timeout",
		Timeout:      time.Second * 30,
	}))

	// set API route
	route.SetAPIRoute(e, &route.APIHandler{
		Config:      cfg,
		Route:       router,
		AuthUsecase: authUsecase,
		AuthAPI:     authAPI,
		UserAPI:     userAPI,
	})
	// set web page route
	route.SetWebRoute(e, &route.WebHandler{
		Config:         cfg,
		Route:          router,
		AuthUsecase:    authUsecase,
		SessionUsecase: sessionUsecase,
		AuthPage:       authWeb,
		UserPage:       userWeb,
		DashboardPage:  dashboardWeb,
	})

	// set template renderer
	var err error
	e.Renderer, err = config.NewRenderer(cfg.AutoReloadTemplate())
	if err != nil {
		log.Fatal(err)
	}

	// run server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", cfg.Port())))
}

func migrateDB(cfg *config.Config, db *gorm.DB) {
	start := time.Now()

	// migrate tables
	log.Println("Start migrate db")
	db.AutoMigrate(&databaseentity.User{})
	db.AutoMigrate(&databaseentity.PasswordResetToken{})
	db.AutoMigrate(&databaseentity.Unit{})
	db.AutoMigrate(&databaseentity.Ingredient{})
	db.AutoMigrate(&databaseentity.IngredientConvertionUnit{})
	db.AutoMigrate(&databaseentity.Recipe{})
	db.AutoMigrate(&databaseentity.RecipeIngredient{})
	db.AutoMigrate(&databaseentity.StockMovement{})
	db.AutoMigrate(&databaseentity.Product{})
	db.AutoMigrate(&databaseentity.Account{})
	db.AutoMigrate(&databaseentity.Transaction{})
	log.Printf("Migrate done in %v", time.Since(start))

	// check non guest user
	log.Println("Check non guest user")
	var totalNonGuestUser int64
	err := db.Model(databaseentity.User{}).Where("is_guest = ? and disabled = ?", false, false).Count(&totalNonGuestUser).Error
	if err != nil {
		fmt.Println("Error get total non guest user:", err)
		os.Exit(1)
	}
	// create non guest user
	if totalNonGuestUser < 1 {
		log.Println("Create first non guest user")
		pass, err := helper.EncryptPassword(cfg.FirstUserPassword())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		token, err := helper.GenerateToken(64)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		user := databaseentity.User{
			Username:           cfg.FirstUserUsername(),
			FullName:           cfg.FirstUserFullName(),
			Email:              cfg.FirstUserEmail(),
			Password:           pass,
			PasswordMustChange: true,
			Token:              token,
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
	if cfg.EnableGuestLogin() {
		log.Println("Check guest user")
		var totalGuestUser int64
		err = db.Model(databaseentity.User{}).Where("is_guest = ?", true).Count(&totalGuestUser).Error
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
			token, err := helper.GenerateToken(64)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			user := databaseentity.User{
				Username:  config.GuestUsername,
				FullName:  config.GuestUsername,
				Email:     config.GuestEmail,
				Password:  pass,
				IsGuest:   true,
				Token:     token,
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
