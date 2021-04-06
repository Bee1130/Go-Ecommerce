package main

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/smiletrl/micro_ecommerce/pkg/config"
	"github.com/smiletrl/micro_ecommerce/pkg/constants"
	"github.com/smiletrl/micro_ecommerce/pkg/dbcontext"
	rpcserver "github.com/smiletrl/micro_ecommerce/service.product/external/server"
	"github.com/smiletrl/micro_ecommerce/service.product/internal/product"
	"os"
)

func main() {
	// provide the .env
	if err := godotenv.Load(); err != nil {
		panic(err.Error())
	}

	// Echo instance
	e := echo.New()
	echoGroup := e.Group("api/v1")

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// initialize service
	stage := os.Getenv(constants.Stage)
	if stage == "" {
		stage = constants.StageLocal
	}
	config, err := config.Load(stage)
	if err != nil {
		panic(err)
	}
	db, err := dbcontext.InitDB(config)
	if err != nil {
		panic(err)
	}

	// product
	productRepo := product.NewRepository(db)
	productService := product.NewService(productRepo)
	product.RegisterHandlers(echoGroup, productService)

	// start grpc server
	go func() {
		rpcserver.Register()
	}()

	// Start rest server
	e.Logger.Fatal(e.Start(":1324"))
}