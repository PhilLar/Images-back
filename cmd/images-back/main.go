package main

import (
	"context"
	"database/sql"
	"flag"
	//"path/filepath"

	//"fmt"
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
	"time"
)

type imageFile struct {
	imgTitle string `json:"title"`
	imgURL   string `json:"url"`
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

func main() {
	flag.Parse()
	e := echo.New()
	e.Logger.SetLevel(log.INFO)

	dbPsql, err := sql.Open("postgres", db)
	if err != nil {
		log.Fatal(err)
	}
	defer dbPsql.Close()

	//dot, err := dotsql.LoadFromFile("migrations/20190909154444_image_files_table.up.sql")
	//_, err = dot.Exec(dbPsql, "create-users-table")

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

	e.POST("files", upload)
	//e.GET("test", testJSON)

	e.Static("/", "static")
	//e.Static("/files", "files")
	// the line below loads only one file - bad
	//e.File("/index", "static/index.html")

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

func upload(c echo.Context) error {
	// Read form fields
	imgPath := c.FormValue("file") //name
	imgTitle := c.FormValue("title") // email


	//-----------
	// Read file
	//-----------

	// Source
	file, err := c.FormFile("file")
	if err != nil {
		log.Print(err)
	}
	src, err := file.Open()
	if err != nil {
		log.Print(err)
	}
	defer src.Close()

	// Destination
	dst, err := os.Create("files/"+file.Filename)
	if err != nil {
		log.Print(err)
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		log.Print(err)
	}


	//filedirectory := filepath.Dir(file.Filename)
	//thepath, err := filepath.Abs(filedirectory)

	if err != nil {
		log.Fatal(err)
	}
	//log.Print(thepath)
	//log.Print(filedirectory)

	log.Print(imgTitle)
	log.Print(file.Filename)
	outJSON := &imageFile{
		imgTitle: imgTitle,
		imgURL:   imgPath,
	}
	return c.JSON(http.StatusOK, outJSON)
}

//func testJSON(c echo.Context) error {
//	// Read form fields
//	outJSON := &imageFile{
//		imgTitle: "imgTitle",
//		imgURL:   "imgPath",
//	}
//	return c.JSON(http.StatusOK, outJSON)
//}


