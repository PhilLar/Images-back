package handlers

import (
	"database/sql"
	"fmt"
	"github.com/PhilLar/Images-back/models"
	"github.com/labstack/echo"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
)

type Env struct {
	DB *sql.DB
}

type imageFile struct {
	ImgID    int    `json:"id"`
	ImgTitle string `json:"title"`
	ImgURL   string `json:"url"`
}

func (env *Env) UploadHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		file, err := c.FormFile("file")
		if err != nil {
			log.Print(err)
			return echo.NewHTTPError(http.StatusBadRequest, "Please provide valid type of file (image)")
		}
		if getFileContentType(file) != "image" {
			return echo.NewHTTPError(http.StatusBadRequest, "Please provide valid type of file (image)")
		}

		imgTitle := c.FormValue("title") //name
		ID, err := models.InsertImage(env.DB, imgTitle, file.Filename)
		if err != nil {
			log.Print(err)
			return echo.NewHTTPError(http.StatusBadRequest, "Please provide valid type of file (image)")
		}

		imgNewTitle, err := models.SaveImage(file, ID)
		if err != nil {
			log.Print(err)
			return echo.NewHTTPError(http.StatusBadRequest, "Please provide valid type of file (image)")
		}

		imgURL := c.Request().Host + c.Request().URL.String() + "/" + imgNewTitle
		outJSON := &imageFile{
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
		err = models.DeleteImage(env.DB, ID)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return c.NoContent(http.StatusNoContent)
	}
}

func (env *Env) ListImagesHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		imgs, err := models.AllImages(env.DB)
		if err != nil {
			log.Print(err)
			return echo.NewHTTPError(http.StatusBadRequest, "Bad Request")
		}
		outImgs := make([]*imageFile, 0)
		for _, i := range imgs {
			imgURL := c.Request().Host + "/files" + "/" + i.StoredName
			outImgs = append(outImgs, &imageFile{
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
