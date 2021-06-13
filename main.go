package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/zercle/docker-snowflake/pkg/datamodels"
	"github.com/zercle/docker-snowflake/pkg/routers"
)

func main() {

	// Running flag
	perFork := flag.Bool("prefork", false, "A env file name without .env")
	flag.Parse()

	// Init app
	app := fiber.New(fiber.Config{
		ErrorHandler: customErrorHandler,
		ReadTimeout:  60 * time.Second,
		Prefork:      *perFork,
	})

	// setup routers
	routers.SetRouters(app)

	// Listen from a different goroutine
	go func() {
		if err := app.Listen(":8080"); err != nil {
			log.Panic(err)
		}
	}()

	// Create channel to signify a signal being sent
	ch := make(chan os.Signal, 1)
	// When an interrupt or termination signal is sent, notify the channel
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	// This blocks the main thread until an interrupt is received
	<-ch
	fmt.Println("Gracefully shutting down...")
	_ = app.Shutdown()

	fmt.Println("Running cleanup tasks...")
	// Your cleanup tasks go here
	// db.Close()
	// redisConn.Close()
	fmt.Println("Successful shutdown.")
}

var customErrorHandler = func(c *fiber.Ctx, err error) error {
	// Default 500 statuscode
	code := http.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		// Override status code if fiber.Error type
		code = e.Code
	}

	responseData := datamodels.ResponseForm{
		Success: false,
		Errors: []*datamodels.ResposeError{
			{
				Code:    code,
				Message: err.Error(),
			},
		},
	}

	// Return statuscode with error message
	err = c.Status(code).JSON(responseData)
	if err != nil {
		// In case the JSON fails
		log.Printf("customErrorHandler: %+v", err)
		return c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
	}

	// Return from handler
	return nil
}
