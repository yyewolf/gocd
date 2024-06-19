package web

import (
	"gocd/internal/docker"
	"sync"

	"github.com/labstack/echo/v4"
)

var working = make(map[string]*sync.Mutex)

func Routes(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		return c.String(418, "I'm a teapot")
	})

	e.GET("/containers/:token", func(c echo.Context) error {
		token := c.Param("token")

		m, f := working[token]
		if !f {
			m = &sync.Mutex{}
			working[token] = m
		}

		m.Lock()
		defer m.Unlock()

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

		m, f := working[token]
		if !f {
			m = &sync.Mutex{}
			working[token] = m
		}

		m.Lock()
		defer m.Unlock()

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
