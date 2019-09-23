package main

import (
	"context"
	"flag"
	"github.com/PhilLar/Images-back/handlers"
	"github.com/PhilLar/Images-back/models"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
	"net/http"

	//"mime/multipart"
	//"net/http"
	"os"
	"os/signal"
	"strconv"
	//"strings"
	"time"
)

var port int
var db string

func init() {
	defPort := 8080
	var defDB string
	if portVar, ok := os.LookupEnv("PORT"); ok {
		if portValue, err := strconv.Atoi(portVar); err == nil {
			defPort = portValue
		}
	}
	if dbVar, ok := os.LookupEnv("DATABASE_URL"); ok {
		defDB = dbVar
	}
	flag.IntVar(&port, "port", defPort, "port to listen on")
	flag.StringVar(&db, "db", defDB, "database to connect to")
}



func main() {
	flag.Parse()
	e := echo.New()
	e.Logger.SetLevel(log.INFO)

	dbPsql, err := models.NewDB(db)
	if err != nil {
		log.Panic(err)
	}
	defer dbPsql.Close()
	fs := &models.FilesSystem{"/files"}
	env := &handlers.Env{Store: dbPsql, FilesSystem:fs}


	e.GET("images", env.ListImagesHandler())
	e.POST("files", env.UploadHandler())
	e.DELETE("images/:id", env.DeleteImageHandler())

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodDelete},
	}))

	e.Static("/", "static")
	e.Static("/files", "files")

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
		e.Logger.Print(err)
	}
}
