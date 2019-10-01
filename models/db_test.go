package models_test

import (
	"strconv"
	"strings"
	"testing"

	"errors"

	"github.com/PhilLar/Images-back/mocks"
	"github.com/PhilLar/Images-back/models"
	"github.com/golang/mock/gomock"

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
}

func (s *StoreSuite) SetupSuite() {
	/*
		The database connection is opened in the setup, and
		stored as an instance variable,
		as is the higher level `store`, that wraps the `db`
	*/
	connString := "postgres://images:images_go@localhost/images_test?sslmode=disable"
	dbPsql, err := models.NewDB(connString, "file://../migrations")
	s.Require().NoError(err)

	s.store = &models.Store{DB: dbPsql}

}

func (s *StoreSuite) SetupTest() {
	/*
		We delete all entries from the table before each test runs, to ensure a
		consistent state before our tests run. In more complex applications, this
		is sometimes achieved in the form of migrations
	*/
	_, err := s.store.DB.Query("DELETE FROM images")
	s.Require().NoError(err)
}

func (s *StoreSuite) TearDownSuite() {
	// Close the connection after all tests in the suite
	s.store.DB.Close()
}

func TestStoreSuite(t *testing.T) {
	s := new(StoreSuite)
	suite.Run(t, s)
}

func (s *StoreSuite) TestInsertImage() {
	// Create a bird through the store `CreateBird` method
	s.Run("Correct insertion", func() {
		testTitle := "mytitle"
		testFileName := "cat.jpg"
		ID, err := s.store.InsertImage(testTitle, testFileName)
		s.Require().NoError(err)

		s.NotEqualf("-1", ID,
			"incorrect ID, expected != -1, actual is '%d'", ID)

		///////// CHECK IF IMAGES WERE INSERTED ///////////////////////////
		rows, err := s.store.DB.Query("SELECT * FROM images")
		s.Require().NoError(err)
		defer rows.Close()

		imgs := make([]*models.Image, 0)
		for rows.Next() {
			img := &models.Image{}
			err := rows.Scan(&img.ID, &img.SourceName, &img.StoredName)
			s.Require().NoError(err)
			imgs = append(imgs, img)
		}

		s.Require().NoError(rows.Err())

		//////////////////////////////////////////////////////////////////

		s.Equalf(testTitle, imgs[0].SourceName,
			"incorrect sourceName, expected %s, actual is '%s'", testTitle, imgs[0].SourceName)


		imgExt := strings.LastIndex(testFileName, ".")
		expectedStoredName := strconv.Itoa(imgs[0].ID) + testFileName[imgExt:]

		s.Equalf(expectedStoredName, imgs[0].StoredName,
			"incorrect storedName, expected %s, actual is '%s'", expectedStoredName, imgs[0].StoredName)
	})

	s.Run("incorrect insertion due to passing filename without extension", func() {
		testTitle := "mytitle"
		testFileName := "catjpg"
		expectedError := errors.New("filename must contain extension")

		ID, err := s.store.InsertImage(testTitle, testFileName)

		s.Equalf(expectedError, err,
			"incorrect error, expected %s, actual is '%s'", expectedError.Error(), err.Error())

		s.NotEqualf(ID, "-1",
			"incorrect ID, expected -1, actual is '%d'", ID)
	})
}

func (s *StoreSuite) TestDeleteImage() {
	s.Run("Correct deletion", func() {

		testTitle := "mytitle"
		testFileName := "cat.jpg"

		///////// CHECK IF IMAGES WERE INSERTED ///////////////////////////
		rows, err := s.store.DB.Query("SELECT * FROM images")
		s.Require().NoError(err)
		defer rows.Close()

		imgs := make([]*models.Image, 0)
		for rows.Next() {
			img := &models.Image{}
			err := rows.Scan(&img.ID, &img.SourceName, &img.StoredName)
			s.Require().NoError(err)
			imgs = append(imgs, img)
		}
		s.Require().NoError(rows.Err())
		//////////////////////////////////////////////////////////////////
		rowsCountBefore := len(imgs)

		ID, err := s.store.InsertImage(testTitle, testFileName)
		s.Require().NoError(err)

		mockCtrl := gomock.NewController(s.T())
		defer mockCtrl.Finish()
		mockSystem := mocks.NewMockSystem(mockCtrl)
		s.store.OS = mockSystem

		mockSystem.
			EXPECT().
			Remove(gomock.Any()).
			Return(nil)

		err = s.store.DeleteImage(ID)
		s.Require().NoError(err)

		///////// CHECK IF IMAGES WERE INSERTED ///////////////////////////
		rows, err = s.store.DB.Query("SELECT * FROM images")
		s.Require().NoError(err)
		defer rows.Close()

		imgs = make([]*models.Image, 0)
		for rows.Next() {
			img := &models.Image{}
			err := rows.Scan(&img.ID, &img.SourceName, &img.StoredName)
			s.Require().NoError(err)
			imgs = append(imgs, img)
		}
		if err = rows.Err(); err != nil {
			s.Require().NoError(err)
		}
		//////////////////////////////////////////////////////////////////
		rowsCountAfter := len(imgs)

		for _, img := range imgs {
			s.NotEqualf(img.ID, ID, "there must not be a row with ID = %d after its deletion, but one were found", ID)
		}

		s.Equalf(rowsCountBefore, rowsCountAfter,
			"rows quantity before deletion and after one must be equal, but they aren't\n"+
				"rowsCountBefore: %d\n"+
				"rowsCountAfter: %d\n", rowsCountBefore, rowsCountAfter)
	})

	s.Run("Incorrect deletion", func() {
		falseID := -1

		err := s.store.DeleteImage(falseID)
		expectedError := errors.New("image with such ID not found in database")
		s.Equalf(expectedError, err,
			"incorrect error, expected %s, actual is '%s'", expectedError.Error(), err.Error())
	})
}
