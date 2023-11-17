package main

import (
	"gocd/internal/docker"
	"gocd/web"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

func main() {
	go func() {
		// Always listen for new events
		for {
			docker.StartListener()
		}
	}()

	e := echo.New()

	// Use Logrus
	log := logrus.New()
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			log.WithFields(logrus.Fields{
				"URI":    values.URI,
				"status": values.Status,
			}).Info("HTTP Request")

			return nil
		},
	}))

	web.Routes(e)

	e.Logger.Fatal(e.Start(":8080"))
}
