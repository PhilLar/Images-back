package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.com/PhilLar/Images-back/models"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

type imageFile struct {
	ImgTitle string `json:"title"`
	ImgURL   string `json:"url"`
}

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
type Env struct {
	db *sql.DB
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

	env := &Env{db: dbPsql}

	e.POST("files", env.uploadHandler())

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
		e.Logger.Fatal(err)
	}
}
func (env *Env) uploadHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		imgTitle := c.FormValue("title") //name
		ID, err := models.InsertImage(env.db, imgTitle)
		if err != nil {
			log.Fatal(err)
		}

		file, err := c.FormFile("file")
		if err != nil {
			log.Print(err)
		}
		imgNewTitle, err := models.SaveImage(file, ID)
		if err != nil {
			log.Fatal(err)
		}

		imgURL := c.Request().Host + c.Request().URL.String() + "/" + imgNewTitle
		outJSON := &imageFile{
			ImgTitle: imgTitle,
			ImgURL:   imgURL,
		}
		respHeader := c.Response().Header()
		for i, j := range respHeader {
			fmt.Println(i, j)
		}
		return c.JSON(http.StatusOK, outJSON)
	}
}
