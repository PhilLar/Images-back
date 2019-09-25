package handlers

import (
	//"database/sql"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/PhilLar/Images-back/models"
	"github.com/labstack/echo"
)

type ImagesStore interface {
	InsertImage(imgTitle, fileName string) (int, error)
	DeleteImage(ID int) error
	AllImages() ([]*models.Image, error)
}

type FilesStore interface {
	SaveImage(file *multipart.FileHeader, ID int) (string, error)
}

type Env struct {
	Store       ImagesStore
	FilesSystem FilesStore
}

type ImageFile struct {
	ImgID    int    `json:"id"`
	ImgTitle string `json:"title"`
	ImgURL   string `json:"url"`
}

func (env *Env) UploadHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		file, err := c.FormFile("file")
		if err != nil {
			log.Print(err)
			return echo.NewHTTPError(http.StatusBadRequest, "Please provide valid type of file (image): "+err.Error())
		}
		if getFileContentType(file) != "image" {
			return echo.NewHTTPError(http.StatusBadRequest, "Please provide valid type of file (image), actual: "+getFileContentType(file))
		}

		imgTitle := c.FormValue("title") //name
		ID, err := env.Store.InsertImage(imgTitle, file.Filename)
		if err != nil {
			log.Print(err)
			return echo.NewHTTPError(http.StatusBadRequest, "Please provide valid type of file (image)"+err.Error())
		}

		imgNewTitle, err := env.FilesSystem.SaveImage(file, ID)
		if err != nil {
			log.Print(err)
			return echo.NewHTTPError(http.StatusBadRequest, "Please provide valid type of file (image)"+err.Error())
		}

		imgURL := c.Request().Host + c.Request().URL.String() + "/" + imgNewTitle
		outJSON := &ImageFile{
			ImgTitle: imgTitle,
			ImgURL:   imgURL,
			ImgID:    ID,
		}
		respHeader := c.Response().Header()
		for i, j := range respHeader {
			fmt.Println(i, j)
		}
		return c.JSON(http.StatusOK, outJSON)
	}
}
func (env *Env) DeleteImageHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		ID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Print(err)
			return echo.NewHTTPError(http.StatusBadRequest, "ID must be integer (BIGSERIAL)")
		}
		err = env.Store.DeleteImage(ID)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return c.NoContent(http.StatusNoContent)
	}
}

func (env *Env) ListImagesHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Method != "GET" {
			return echo.NewHTTPError(http.StatusBadRequest, "Bad Request [only method GET allowed]")
		}
		imgs, err := env.Store.AllImages()
		if err != nil {
			log.Print(err)
			return echo.NewHTTPError(http.StatusBadRequest, "Bad Request")
		}
		outImgs := make([]*ImageFile, 0)
		for _, i := range imgs {
			imgURL := c.Request().Host + "/files" + "/" + i.StoredName
			outImgs = append(outImgs, &ImageFile{
				ImgTitle: i.SourceName,
				ImgURL:   imgURL,
				ImgID:    i.ID,
			})
		}
		return c.JSONPretty(http.StatusOK, outImgs, "  ")
	}
}

func getFileContentType(file *multipart.FileHeader) string {

	contentType := file.Header["Content-Type"][0]
	imgExt := strings.Index(contentType, "/")

	return contentType[:imgExt]
}
