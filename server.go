package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type (
	message struct {
		ID   int    `json:"id"`
		Body string `json:"body"`
	}
)

var (
	seq = 1
)

//----------
// Handlers
//----------

func sendMessage(c echo.Context) error {
	u := &message{
		ID: seq,
	}
	if err := c.Bind(u); err != nil {
		return err
	}
	seq++
	return c.JSON(http.StatusCreated, u)
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.POST("/kafka", sendMessage)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
