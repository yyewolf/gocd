package web

import (
	"gocd/internal/docker"

	"github.com/labstack/echo/v4"
)

func Routes(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		return c.String(418, "I'm a teapot")
	})

	e.GET("/containers/:token", func(c echo.Context) error {
		token := c.Param("token")

		err := docker.UpdateContainers(token)
		if err != nil {
			return c.JSON(403, map[string]interface{}{
				"message": "forbidden",
			})
		}

		return c.JSON(200, map[string]interface{}{
			"message": "ok",
		})
	})

	e.POST("/containers", func(c echo.Context) error {
		token := c.FormValue("token")

		err := docker.UpdateContainers(token)
		if err != nil {
			return c.JSON(403, map[string]interface{}{
				"message": "forbidden",
			})
		}

		return c.JSON(200, map[string]interface{}{
			"message": "ok",
		})
	})
}
