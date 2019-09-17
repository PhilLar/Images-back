package models

import (
	"database/sql"
	"errors"
	//"github.com/golang-migrate/migrate/source/file"
	"io"
	"log"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
)

type Image struct {
	ID         int
	SourceName string
	StoredName string
}

func InsertImage(db *sql.DB, imgTitle, fileName string) (int, error) {
	var ID int
	err := db.QueryRow("INSERT INTO images(source_name) VALUES($1) RETURNING id", imgTitle).Scan(&ID)
	if err != nil {
		return -1, err
	}
	imgExt := strings.LastIndex(fileName, ".")
	imgNewTitle := strconv.Itoa(ID) + fileName[imgExt:]
	_, err = db.Exec("UPDATE images SET stored_name=$1 WHERE id=$2", imgNewTitle, ID)
	if err != nil {
		return -1, err
	}
	return ID, nil
}

func AllImages(db *sql.DB) ([]*Image, error) {
	rows, err := db.Query("SELECT * FROM images")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	imgs := make([]*Image, 0)
	for rows.Next() {
		img := &Image{}
		err := rows.Scan(&img.ID, &img.SourceName, &img.StoredName)
		if err != nil {
			return nil, err
		}
		imgs = append(imgs, img)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return imgs, nil
}

func SaveImage(file *multipart.FileHeader, ID int) (string, error) {
	src, err := file.Open()
	if err != nil {
		log.Print(err)
		return "", err
	}
	defer src.Close()
	//a := src.(os.File)
	//contentType, err := getFileContentType(*src.(os.File))
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

func DeleteImage(db *sql.DB, ID int) error {
	var storedName string
	err := db.QueryRow("DELETE FROM images WHERE id=$1 RETURNING stored_name", ID).Scan(&storedName)
	if err != nil {
		return errors.New("image with such ID not found in database")
	}
	err = os.Remove("files/" + storedName)
	if err != nil {
		return errors.New("image with such ID not found in '/files' directory")
	}
	return nil
}
