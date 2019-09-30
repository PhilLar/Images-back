package models_test

import (
	"database/sql"
	"github.com/PhilLar/Images-back/mocks"
	"github.com/PhilLar/Images-back/models"
	"github.com/golang/mock/gomock"
	"log"
	"strconv"
	"strings"
	"testing"

	// The "testify/suite" package is used to make the test suite
	"github.com/stretchr/testify/suite"
)

type StoreSuite struct {
	suite.Suite
	/*
		The suite is defined as a struct, with the store and db as its
		attributes. Any variables that are to be shared between tests in a
		suite should be stored as attributes of the suite instance
	*/
	store *models.Store
	db    *sql.DB
}

func (s *StoreSuite) SetupSuite() {
	/*
		The database connection is opened in the setup, and
		stored as an instance variable,
		as is the higher level `store`, that wraps the `db`
	*/
	connString := "postgres://images:images_go@localhost/images_test?sslmode=disable"
	db, err := sql.Open("postgres", connString)
	if err != nil {
		s.T().Fatal(err)
	}
	s.db = db
	dbPsql, err := models.NewDB(connString, "file://../migrations")
	if err != nil {
		log.Printf(err.Error())
	}

	s.store = &models.Store{DB: dbPsql}
	//s.store, err = models.NewDB(connString, "file://../migrations")
	//if err != nil {
	//	log.Printf(err.Error())
	//}
}

func (s *StoreSuite) SetupTest() {
	/*
		We delete all entries from the table before each test runs, to ensure a
		consistent state before our tests run. In more complex applications, this
		is sometimes achieved in the form of migrations
	*/
	_, err := s.db.Query("DELETE FROM images")
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *StoreSuite) TearDownSuite() {
	// Close the connection after all tests in the suite
	s.db.Close()
}

func TestStoreSuite(t *testing.T) {
	s := new(StoreSuite)
	suite.Run(t, s)
}

func (s *StoreSuite) TestInsertImage() {
	// Create a bird through the store `CreateBird` method
	s.T().Run("Correct insertion", func(t *testing.T) {
		testTitle := "mytitle"
		testFileName := "cat.jpg"
		ID, err := s.store.InsertImage(testTitle, testFileName)
		if err != nil {
			s.T().Error(err)
		}

		if ID == -1 {
			s.T().Errorf("incorrect ID, expected != -1, actual is '%d'", ID)
		}

		imgs, err := s.store.AllImages()
		if err != nil {
			s.T().Error(err)
		}

		if imgs[0].SourceName != testTitle {
			s.T().Errorf("incorrect sourceName, expected %s, actual is '%s'", testTitle, imgs[0].SourceName)
		}

		imgExt := strings.LastIndex(testFileName, ".")
		expectedStoredName := strconv.Itoa(imgs[0].ID) + testFileName[imgExt:]
		if imgs[0].StoredName != expectedStoredName {
			s.T().Errorf("incorrect storedName, expected %s, actual is '%s'", expectedStoredName, imgs[0].StoredName)
		}
	})

	s.T().Run("incorrect insertion due to passing filename without extension", func(t *testing.T) {
		testTitle := "mytitle"
		testFileName := "catjpg"
		expectedError := "filename must contain extension"
		imgs, err := s.store.AllImages()
		if err != nil {
			s.T().Error(err)
		}
		rowsCountBefore := len(imgs)
		ID, err := s.store.InsertImage(testTitle, testFileName)

		if err != nil && err.Error() != "filename must contain extension" {
			s.T().Errorf("incorrect error, expected %s, actual is '%s'", expectedError, err.Error())
		}

		if ID != -1 {
			s.T().Errorf("incorrect ID, expected -1, actual is '%d'", ID)
		}

		imgs, err = s.store.AllImages()
		if err != nil {
			s.T().Error(err)
		}
		rowsCountAfter := len(imgs)

		if rowsCountBefore != rowsCountAfter {
			s.T().Errorf("there must be no row added, but the query returned %d new were", rowsCountBefore - rowsCountAfter)
		}

	})


}

func (s *StoreSuite) TestDeleteImage() {
	s.T().Run("Correct deletion", func(t *testing.T) {

		testTitle := "mytitle"
		testFileName := "cat.jpg"

		imgs, err := s.store.AllImages()
		if err != nil {
			s.T().Error(err)
		}
		rowsCountBefore := len(imgs)

		ID, err := s.store.InsertImage(testTitle, testFileName)
		if err != nil {
			s.T().Error(err)
		}


		mockCtrl := gomock.NewController(&testing.T{})
		defer mockCtrl.Finish()
		mockSystem := mocks.NewMockSystem(mockCtrl)
		s.store.OS = mockSystem


		mockSystem.
			EXPECT().
			Remove(gomock.Any()).
			Return(nil).
			AnyTimes()

		err = s.store.DeleteImage(ID)
		if err != nil {
			s.T().Error(err)
		}

		imgs, err = s.store.AllImages()
		if err != nil {
			s.T().Error(err)
		}
		rowsCountAfter := len(imgs)

		for _, img := range imgs {
			if img.ID == ID {
				s.T().Errorf("there must not be a row with ID = %d after its deletion, but one were found", ID)
			}
		}

		if rowsCountBefore != rowsCountAfter {
			s.T().Errorf("rows quantity before deletion and after one must be equal, but they aren't\n" +
				"rowsCountBefore: %d\n" +
				"rowsCountAfter: %d\n", rowsCountBefore, rowsCountAfter)
		}
	})

	s.T().Run("Incorrect deletion", func(t *testing.T) {
		falseID := -1

		imgs, err := s.store.AllImages()
		if err != nil {
			s.T().Error(err)
		}
		rowsCountBefore := len(imgs)


		err = s.store.DeleteImage(falseID)
		expectedError := "image with such ID not found in database"
		if err != nil && err.Error() != expectedError {
			s.T().Errorf("incorrect error, expected %s, actual is '%s'", expectedError, err.Error())
		}

		imgs, err = s.store.AllImages()
		if err != nil {
			s.T().Error(err)
		}
		rowsCountAfter := len(imgs)

		for _, img := range imgs {
			if img.ID == falseID {
				s.T().Errorf("there must not be a row with ID = %d after its deletion, but one were found", falseID)
			}
		}

		if rowsCountBefore != rowsCountAfter {
			s.T().Errorf("rows quantity before deletion and after one must be equal, but they aren't\n" +
				"rowsCountBefore: %d\n" +
				"rowsCountAfter: %d\n", rowsCountBefore, rowsCountAfter)
		}
	})
}


