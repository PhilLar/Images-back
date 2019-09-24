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

	"github.com/PhilLar/Images-back/handlers"
	"github.com/PhilLar/Images-back/mocks"
	"github.com/PhilLar/Images-back/models"
	gomock "github.com/golang/mock/gomock"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

//WORKS
func TestListImagesHandler(t *testing.T) {
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

	if assert.NoError(t, env.ListImagesHandler()(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		err := json.Unmarshal(rec.Body.Bytes(), &template)
		if err != nil {
			t.Fatal("Opps")
		}
		assert.Equal(t, outImgs, template)
	}
}

//WORKS
func TestDeleteImageHandler(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockImagesStore := mocks.NewMockImagesStore(mockCtrl)

	e := echo.New()
	env := &handlers.Env{Store: mockImagesStore}
	req := httptest.NewRequest(http.MethodPost, "/images/:id", nil)

	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	mockImagesStore.EXPECT().DeleteImage(1).Return(nil).Times(1)

	if assert.NoError(t, env.DeleteImageHandler()(c)) {
		assert.Equal(t, http.StatusNoContent, rec.Code)
	}
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
	values := map[string]io.Reader{
		"file":  mustOpen("cat.jpg"), // lets assume its this file
		"title": strings.NewReader("that's the cat!"),
	}
	b, w := Upload(values)

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

	mockImagesStore.EXPECT().InsertImage("that's the cat!", "cat.jpg").Return(1, nil).Times(1)
	mockFilesStore.EXPECT().SaveImage(gomock.Any(), 1).Return("1.jpg", nil).Times(1)

	var gotJSON handlers.ImageFile
	expectedJSON := handlers.ImageFile{
		ImgTitle: "that's the cat!",
		ImgURL:   "example.com/files/1.jpg",
		ImgID:    1,
	}

	if assert.NoError(t, env.UploadHandler()(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		log.Print(rec.Body.String())
		err := json.Unmarshal(rec.Body.Bytes(), &gotJSON)
		if err != nil {
			t.Fatal("Opps")
		}
		assert.Equal(t, expectedJSON, gotJSON)
	}
}

func CreateFormImagefile(fieldname, filename string, w *multipart.Writer) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			fieldname, filename))
	h.Set("Content-Type", "image/png")
	return w.CreatePart(h)
}

func Upload(values map[string]io.Reader) (bytes.Buffer, *multipart.Writer) {
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
			if fw, err = CreateFormImagefile(key, x.Name(), w); err != nil {
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
