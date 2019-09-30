package models

import (
	"database/sql"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
	"syscall"
)



type Image struct {
	ID         int
	SourceName string
	StoredName string
}
type System interface {
	Remove(name string) error
}

type Store struct {
	DB	*sql.DB
	OS 	System
}

type FilesSystem struct {
	Root string
}

type OS struct {}

func(*OS) Remove(name string) error {
	// System call interface forces us to know
	// whether name is a file or directory.
	// Try both: it is cheaper on average than
	// doing a Stat plus the right one.
	e := syscall.Unlink(name)
	if e == nil {
		return nil
	}
	e1 := syscall.Rmdir(name)
	if e1 == nil {
		return nil
	}

	// Both failed: figure out which error to return.
	// OS X and Linux differ on whether unlink(dir)
	// returns EISDIR, so can't use that. However,
	// both agree that rmdir(file) returns ENOTDIR,
	// so we can use that to decide which error is real.
	// Rmdir might also return ENOTDIR if given a bad
	// file path, like /etc/passwd/foo, but in that case,
	// both errors will be ENOTDIR, so it's okay to
	// use the error from unlink.
	if e1 != syscall.ENOTDIR {
		e = e1
	}
	return &os.PathError{"remove", name, e}
}

func (s *Store) InsertImage(imgTitle, fileName string) (int, error) {
	var ID int
	imgExt := strings.LastIndex(fileName, ".")
	if imgExt == -1 {
		return -1, errors.New("filename must contain extension")
	}
	err := s.DB.QueryRow("INSERT INTO images(source_name) VALUES($1) RETURNING id", imgTitle).Scan(&ID)
	if err != nil {
		return -1, err
	}
	imgNewTitle := strconv.Itoa(ID) + fileName[imgExt:]
	_, err = s.DB.Exec("UPDATE images SET stored_name=$1 WHERE id=$2", imgNewTitle, ID)
	if err != nil {
		return -1, err
	}
	return ID, nil
}

func (s *Store) AllImages() ([]*Image, error) {
	rows, err := s.DB.Query("SELECT * FROM images")
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

func (fs *FilesSystem) SaveImage(file *multipart.FileHeader, ID int) (string, error) {
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

func (s *Store) DeleteImage(ID int) error {
	var storedName string
	err := s.DB.QueryRow("DELETE FROM images WHERE id=$1 RETURNING stored_name", ID).Scan(&storedName)
	if err != nil {
		return errors.New("image with such ID not found in database")
	}
	err = s.OS.Remove("files/" + storedName)
	if err != nil {
		return errors.New("image with such ID not found in '/files' directory "+storedName)
	}
	return nil
}
