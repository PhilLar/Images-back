package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"os"
	"strings"
	"testing"

	"github.com/pkg/errors"

	"github.com/PhilLar/Images-back/handlers"
	"github.com/PhilLar/Images-back/mocks"
	"github.com/PhilLar/Images-back/models"
	gomock "github.com/golang/mock/gomock"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListImagesHandler(t *testing.T) {
	t.Run("returns StatusOK", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockImagesStore := mocks.NewMockImagesStore(mockCtrl)

		e := echo.New()
		env := &handlers.Env{Store: mockImagesStore}
		req := httptest.NewRequest(http.MethodGet, "/images", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		imgs := []*models.Image{
			&models.Image{
				1,
				"cat",
				"1.jpg",
			},
			&models.Image{
				2,
				"dog",
				"2.jpg",
			},
			&models.Image{
				3,
				"frog",
				"3.jpg",
			},
		}

		mockImagesStore.EXPECT().AllImages().Return(imgs, nil).Times(1)

		var template []handlers.ImageFile
		outImgs := make([]handlers.ImageFile, 0)
		outImgs = append(outImgs, handlers.ImageFile{
			ImgTitle: "cat",
			ImgURL:   "example.com/files/1.jpg",
			ImgID:    1,
		})
		outImgs = append(outImgs, handlers.ImageFile{
			ImgTitle: "dog",
			ImgURL:   "example.com/files/2.jpg",
			ImgID:    2,
		})
		outImgs = append(outImgs, handlers.ImageFile{
			ImgTitle: "frog",
			ImgURL:   "example.com/files/3.jpg",
			ImgID:    3,
		})

		require.NoError(t, env.ListImagesHandler()(c))
		assert.Equal(t, http.StatusOK, rec.Code)
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &template))
		assert.Equal(t, outImgs, template)
	})

	t.Run("returns BadRequest due to scanError in AllImages()", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockImagesStore := mocks.NewMockImagesStore(mockCtrl)

		e := echo.New()
		env := &handlers.Env{Store: mockImagesStore}
		req := httptest.NewRequest(http.MethodGet, "/images", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockImagesStore.EXPECT().AllImages().Return(nil, errors.New("Error while scanning db rows"))

		outImgs := make([]handlers.ImageFile, 0)
		outImgs = append(outImgs, handlers.ImageFile{
			ImgTitle: "cat",
			ImgURL:   "example.com/files/1.jpg",
			ImgID:    1,
		})
		outImgs = append(outImgs, handlers.ImageFile{
			ImgTitle: "dog",
			ImgURL:   "example.com/files/2.jpg",
			ImgID:    2,
		})
		outImgs = append(outImgs, handlers.ImageFile{
			ImgTitle: "frog",
			ImgURL:   "example.com/files/3.jpg",
			ImgID:    3,
		})

		err := env.ListImagesHandler()(c)

		assert.Equal(t, 200, rec.Code)
		assert.Equal(t, err.Error(), "code=400, message=Bad Request")

	})
}

func TestDeleteImageHandler(t *testing.T) {
	t.Run("returns NoContent", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockImagesStore := mocks.NewMockImagesStore(mockCtrl)

		e := echo.New()
		env := &handlers.Env{Store: mockImagesStore}
		req := httptest.NewRequest(http.MethodPost, "/images/:id", nil)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		mockImagesStore.EXPECT().DeleteImage(1).Return(nil)

		require.NoError(t, env.DeleteImageHandler()(c))
		assert.Equal(t, http.StatusNoContent, rec.Code)

	})

	t.Run("Returns BadRequest on invalid id", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockImagesStore := mocks.NewMockImagesStore(mockCtrl)
		e := echo.New()
		env := &handlers.Env{Store: mockImagesStore}
		req := httptest.NewRequest(http.MethodPost, "/images/:id", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("")

		err := env.DeleteImageHandler()(c)
		assert.Equal(t, 200, rec.Code)
		assert.Equal(t, "code=400, message=ID must be integer (BIGSERIAL)", err.Error())
	})

	t.Run("Returns BadRequest on not found id", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockImagesStore := mocks.NewMockImagesStore(mockCtrl)
		e := echo.New()
		env := &handlers.Env{Store: mockImagesStore}
		req := httptest.NewRequest(http.MethodPost, "/images/:id", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("2")

		mockImagesStore.
			EXPECT().
			DeleteImage(2).
			Return(errors.New("image with such ID not found in database")).
			AnyTimes()

		err := env.DeleteImageHandler()(c)
		assert.Equal(t, 200, rec.Code)
		assert.Equal(t, "code=400, message=image with such ID not found in database", err.Error())
	})

}

//
func createMultipartFormData(t *testing.T, fieldName, fileName string) (bytes.Buffer, *multipart.Writer) {
	var b bytes.Buffer
	var err error
	w := multipart.NewWriter(&b)
	var fw io.Writer
	file := mustOpen(fileName)
	if fw, err = w.CreateFormFile(fieldName, file.Name()); err != nil {
		t.Errorf("Error creating writer: %v", err)
	}
	if _, err = io.Copy(fw, file); err != nil {
		t.Errorf("Error with io.Copy: %v", err)
	}
	w.Close()
	return b, w
}

func TestUploadHandler(t *testing.T) {
	t.Run("returns StatusOK", func(t *testing.T) {
		values := map[string]io.Reader{
			"file":  mustOpen("cat.jpg"), // lets assume its this file
			"title": strings.NewReader("that's the cat!"),
		}

		b, w := Upload(values, "image/png")

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockImagesStore := mocks.NewMockImagesStore(mockCtrl)
		mockFilesStore := mocks.NewMockFilesStore(mockCtrl)

		e := echo.New()
		env := &handlers.Env{Store: mockImagesStore, FilesSystem: mockFilesStore}
		req := httptest.NewRequest(http.MethodPost, "/files", &b)
		req.Header.Set(echo.HeaderContentType, w.FormDataContentType())

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockImagesStore.EXPECT().InsertImage("that's the cat!", "cat.jpg").Return(1, nil).AnyTimes()
		mockFilesStore.EXPECT().SaveImage(gomock.Any(), 1).Return("1.jpg", nil).AnyTimes()

		var gotJSON handlers.ImageFile
		expectedJSON := handlers.ImageFile{
			ImgTitle: "that's the cat!",
			ImgURL:   "example.com/files/1.jpg",
			ImgID:    1,
		}

		require.NoError(t, env.UploadHandler()(c))
		assert.Equal(t, http.StatusOK, rec.Code)
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &gotJSON))
		assert.Equal(t, expectedJSON, gotJSON)

	})
	t.Run("returns BadRequest due to empty req body", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockImagesStore := mocks.NewMockImagesStore(mockCtrl)
		mockFilesStore := mocks.NewMockFilesStore(mockCtrl)

		e := echo.New()
		env := &handlers.Env{Store: mockImagesStore, FilesSystem: mockFilesStore}
		req := httptest.NewRequest(http.MethodPost, "/files", nil)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockImagesStore.EXPECT().InsertImage("that's the cat!", "cat.jpg").Return(1, nil).AnyTimes()
		mockFilesStore.EXPECT().SaveImage(gomock.Any(), 1).Return("1.jpg", nil).AnyTimes()

		err := env.UploadHandler()(c)
		assert.Equal(t, 200, rec.Code)
		assert.Equal(t, "code=400, message=Please provide valid type of file (image): request Content-Type isn't multipart/form-data", err.Error())
	})

	t.Run("returns BadRequest due to wrong Content-Type", func(t *testing.T) {
		values := map[string]io.Reader{
			"file":  mustOpen("cat.jpg"), // lets assume its this file
			"title": strings.NewReader("that's the cat!"),
		}

		b, w := Upload(values, "nonImage/png")

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockImagesStore := mocks.NewMockImagesStore(mockCtrl)
		mockFilesStore := mocks.NewMockFilesStore(mockCtrl)

		e := echo.New()
		env := &handlers.Env{Store: mockImagesStore, FilesSystem: mockFilesStore}
		req := httptest.NewRequest(http.MethodPost, "/files", &b)
		req.Header.Set(echo.HeaderContentType, w.FormDataContentType())

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockImagesStore.EXPECT().InsertImage("that's the cat!", "cat.jpg").Return(1, nil).AnyTimes()
		mockFilesStore.EXPECT().SaveImage(gomock.Any(), 1).Return("1.jpg", nil).AnyTimes()

		err := env.UploadHandler()(c)
		assert.Equal(t, 200, rec.Code)
		assert.Equal(t, "code=400, message=Please provide valid type of file (image), actual: nonImage", err.Error())

	})

	t.Run("returns BadRequest due error while scanning db rows", func(t *testing.T) {
		values := map[string]io.Reader{
			"file":  mustOpen("cat.jpg"), // lets assume its this file
			"title": strings.NewReader("that's the cat!"),
		}

		b, w := Upload(values, "image/png")

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockImagesStore := mocks.NewMockImagesStore(mockCtrl)
		mockFilesStore := mocks.NewMockFilesStore(mockCtrl)

		e := echo.New()
		env := &handlers.Env{Store: mockImagesStore, FilesSystem: mockFilesStore}
		req := httptest.NewRequest(http.MethodPost, "/files", &b)
		req.Header.Set(echo.HeaderContentType, w.FormDataContentType())

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockImagesStore.EXPECT().InsertImage("that's the cat!", "cat.jpg").Return(0, errors.New("Error while scanning db rows")).AnyTimes()
		mockFilesStore.EXPECT().SaveImage(gomock.Any(), 1).Return("1.jpg", nil).AnyTimes()

		err := env.UploadHandler()(c)
		assert.Equal(t, 200, rec.Code)
		assert.Equal(t, "code=400, message=Please provide valid type of file (image)Error while scanning db rows", err.Error())
	})

	t.Run("returns BadRequest due error while scanning db rows", func(t *testing.T) {
		values := map[string]io.Reader{
			"file":  mustOpen("cat.jpg"), // lets assume its this file
			"title": strings.NewReader("that's the cat!"),
		}

		b, w := Upload(values, "image/png")

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockImagesStore := mocks.NewMockImagesStore(mockCtrl)
		mockFilesStore := mocks.NewMockFilesStore(mockCtrl)

		e := echo.New()
		env := &handlers.Env{Store: mockImagesStore, FilesSystem: mockFilesStore}
		req := httptest.NewRequest(http.MethodPost, "/files", &b)
		req.Header.Set(echo.HeaderContentType, w.FormDataContentType())

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockImagesStore.EXPECT().InsertImage("that's the cat!", "cat.jpg").Return(1, nil).AnyTimes()
		mockFilesStore.EXPECT().SaveImage(gomock.Any(), 1).Return("", errors.New("Error while scanning db rows")).AnyTimes()

		err := env.UploadHandler()(c)
		assert.Equal(t, 200, rec.Code)
		assert.Equal(t, "code=400, message=Please provide valid type of file (image)Error while scanning db rows", err.Error())
	})

}

func CreateFormImagefile(fieldname, filename string, w *multipart.Writer, contType string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			fieldname, filename))
	h.Set("Content-Type", contType)
	return w.CreatePart(h)
}

func Upload(values map[string]io.Reader, contType string) (bytes.Buffer, *multipart.Writer) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range values {
		var fw io.Writer
		var err error
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		// Add an image file
		if x, ok := r.(*os.File); ok {
			if fw, err = CreateFormImagefile(key, x.Name(), w, contType); err != nil {
				log.Fatal(err)
			}
		} else {
			// Add other fields
			if fw, err = w.CreateFormField(key); err != nil {
				log.Fatal(err)
			}
		}
		if _, err := io.Copy(fw, r); err != nil {
			log.Fatal(err)
		}

	}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()
	return b, w
}

func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	return r
}
