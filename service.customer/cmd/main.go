package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/smiletrl/micro_ecommerce/pkg/config"
	"github.com/smiletrl/micro_ecommerce/pkg/constants"
	"github.com/smiletrl/micro_ecommerce/pkg/dbcontext"
	"github.com/smiletrl/micro_ecommerce/service.customer/internal/balance"
	"github.com/smiletrl/micro_ecommerce/service.customer/internal/customer"
	"os"
)

func main() {
	// Echo instance
	e := echo.New()
	echoGroup := e.Group("api/v1")

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// initialize service
	stage := os.Getenv(constants.Stage)
	if stage == "" {
		//stage = constants.StageLocal
		stage = "/Users/smiletrl/go/src/github.com/smiletrl/micro_ecommerce/config/local.yml"
	}
	config, err := config.Load(stage)
	if err != nil {
		panic(err)
	}
	db, err := dbcontext.InitDB(config)
	if err != nil {
		panic(err)
	}

	// balance
	balanceRepo := balance.NewRepository(db)
	balanceService := balance.NewService(balanceRepo)
	balance.RegisterHandlers(echoGroup, balanceService)

	// customer
	customerRepo := customer.NewRepository(db)
	customerService := customer.NewService(customerRepo)
	customer.RegisterHandlers(echoGroup, customerService)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
