package main

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/smiletrl/micro_ecommerce/pkg/config"
	"github.com/smiletrl/micro_ecommerce/pkg/constants"
	"github.com/smiletrl/micro_ecommerce/pkg/dbcontext"
	"github.com/smiletrl/micro_ecommerce/pkg/entity"
	"github.com/smiletrl/micro_ecommerce/service.cart/internal/cart"
	productClient "github.com/smiletrl/micro_ecommerce/service.product/external/client"
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

	// cart
	cartRepo := cart.NewRepository(db)
	productProxy := product{}
	cartService := cart.NewService(cartRepo, productProxy)
	cart.RegisterHandlers(echoGroup, cartService)

	// Start server
	e.Logger.Fatal(e.Start(":1325"))
}

// product proxy
type product struct{}

func (p product) GetDetail(c echo.Context, id int64) (entity.Product, error) {
	return productClient.GetProduct(id)
}
