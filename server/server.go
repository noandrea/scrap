package server

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"github.com/noandrea/scrap/pkg/scrap"
)

func _m(err error) interface{} {
	return map[string]string{"message": fmt.Sprint(err)}
}

// Start starts the rest api
func Start(settings ConfigSchema) (err error) {
	// configure scrap
	scrap.Configure(settings.ChromeAddress)
	// cache management
	if err := initCache(settings); err != nil {
		log.Warnf("cache initialization failed, will run without cache: %v", err)
	}
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
		var m scrap.Movie
		// check the cache
		if found := get(id, &m); found {
			return c.JSON(http.StatusOK, m)
		}
		// cache miss, scrape
		err = scrap.Run(scrap.AmazonPrime, id, settings.ScrapRegion, &m)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, _m(err))
		}
		// update cache
		if err := set(id, m); err != nil {
			log.Warnf("error updating cache: %v", err)
		}
		return c.JSON(http.StatusOK, m)
	})
	err = e.Start(settings.ListenAddress)
	return
}
