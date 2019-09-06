package main

import (
	"context"
	"flag"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

var port int

func init() {
	defPort := 5000
	if portVar, ok := os.LookupEnv("PORT"); ok {
		if portValue, err := strconv.Atoi(portVar); err == nil {
			defPort = portValue
		}
	}
	flag.IntVar(&port, "port", defPort, "port to listen on")
}

func main() {
	flag.Parse()
	e := echo.New()
	e.Static("/static", "static")
	e.Logger.SetLevel(log.INFO)
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.File("/index", "static/index.html")

	go func() {
		if err := e.Start(":" + strconv.Itoa(port)); err != nil {
			e.Logger.Info("shutting down the server", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
