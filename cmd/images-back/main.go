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
	_ "github.com/lib/pq"
	"database/sql"
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

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("dbTest", func(c echo.Context) error {
		//connStr := "user=images password=images_go dbname=imagesapp sslmode=disable"
		url := db
		dbPsql, err := sql.Open("postgres", url)
		if err != nil {
			log.Fatal(err)
		}
		defer dbPsql.Close()

		err = dbPsql.Ping()
		if err != nil {
			panic(err)
		}
		log.Print("DB OK!")
		_, err = dbPsql.Exec("CREATE TABLE IF NOT EXISTS images_table3 (" +
			"id serial PRIMARY KEY," +
			"size integer NOT NULL," +
			"log VARCHAR (50) NOT NULL" +
			")")
		if err != nil {
			panic(err)
		}
		return c.String(http.StatusOK, "DB_TABLE images CREATED!")
	})

	e.Static("/", "static")
	e.Static("/files", "files")
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
