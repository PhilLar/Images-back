package models

import (
	"database/sql"
	"io"
	"log"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
)

type Image struct {
	ID   		int
	SourceNAME 	string
}

func InsertImage(db *sql.DB, imgTitle string) (int, error) {
	var ID int
	err := db.QueryRow("INSERT INTO images(source_name) VALUES($1) RETURNING id", imgTitle).Scan(&ID)
	if err != nil {
		return -1, err
	}
	return ID, nil
}

func SaveImage(file *multipart.FileHeader, ID int) (string, error ){
	src, err := file.Open()
	if err != nil {
		log.Print(err)
		return "", err
	}
	defer src.Close()

	imgExt := strings.LastIndex(file.Filename, ".")
	imgNewTitle := strconv.Itoa(ID) + file.Filename[imgExt:]
	dst, err := os.Create("files/" + imgNewTitle)
	if err != nil {
		log.Print(err)
		return "", err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		log.Print(err)
		return "", err
	}

	return imgNewTitle, nil
}

