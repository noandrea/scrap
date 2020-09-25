package server

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/noandrea/scrap/pkg/scrap"
)

// Start starts the rest api
func Start(settings ConfigSchema) (err error) {
	// echo start
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())
	// health check :)
	e.GET("/status", func(c echo.Context) (err error) {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  "ok",
			"version": settings.RuntimeVersion,
		})
	})
	//
	e.GET("/movie/amazon/:amazonID", func(c echo.Context) (err error) {
		id := c.Param("amazonID")
		m, err := scrap.Run(scrap.AmazonPrime, id, settings.ScrapRegion)
		if err != nil {
			return c.JSON(http.StatusExpectationFailed, map[string]string{"message": fmt.Sprint(err)})
		}
		return c.JSON(http.StatusOK, m)
	})
	err = e.Start(settings.ListenAddress)
	return
}
