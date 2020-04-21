package api

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func VNCProxy() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/api/:uuid/:vmname/vnc", WebSocketProxy())

	e.Start("0.0.0.0:8081")
}
