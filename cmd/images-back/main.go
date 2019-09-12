package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
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

//type store struct{
//	dbPsql	*sql.DB
//}
//
//func (s *store) insertFile

func main() {
	flag.Parse()
	e := echo.New()
	e.Logger.SetLevel(log.INFO)

	dbPsql, err := sql.Open("postgres", db)
	if err != nil {
		log.Fatal(err)
	}
	defer dbPsql.Close()

	driver, err := postgres.WithInstance(dbPsql, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}
	defer driver.Close()
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		log.Fatal(err)
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	e.GET("dbTest", func(c echo.Context) error {
		err = dbPsql.Ping()
		if err != nil {
			panic(err)
		}
		log.Print("DB OK!")
		return c.String(http.StatusOK, "DB_TABLE images CREATED!")
	})

	e.POST("files", uploadHandler(dbPsql))

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
func uploadHandler(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		imgTitle := c.FormValue("title") //name

		var id int
		err := db.QueryRow("INSERT INTO images(source_name) VALUES($1) RETURNING id", imgTitle).Scan(&id)
		if err != nil {
			log.Fatal(err)
		}

		file, err := c.FormFile("file")
		if err != nil {
			log.Print(err)
		}
		log.Print(file.Filename)
		src, err := file.Open()
		if err != nil {
			log.Print(err)
		}
		defer src.Close()

		imgNewTitle := strconv.Itoa(id) + ".jpg"
		dst, err := os.Create("files/" + imgNewTitle)
		if err != nil {
			log.Print(err)
		}
		defer dst.Close()

		// Copy
		if _, err = io.Copy(dst, src); err != nil {
			log.Print(err)
		}
		if err != nil {
			log.Fatal(err)
		}

		imgExt := strings.LastIndex(file.Filename, ".")
		imgURL := c.Request().Host + c.Request().URL.String() + "/" + strconv.Itoa(id) + file.Filename[imgExt:]
		outJSON := &imageFile{
			ImgTitle: imgTitle,
			ImgURL:   imgURL,
		}
		respHeadder := c.Response().Header()
		for i, j := range respHeadder {
			fmt.Println(i, j)
		}
		return c.JSON(http.StatusOK, outJSON)
	}
}
