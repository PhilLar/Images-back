package models

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
)

func NewDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	if err = db.Ping(); err != nil {
		log.Print(err)
		return nil, err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Print(err)
		return nil, err
	}
	//defer driver.Close()
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Print(err)
		return nil, err
	}
	return db, nil
}
