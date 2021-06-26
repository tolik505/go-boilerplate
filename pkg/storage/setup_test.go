package storage

import (
	"github.com/stretchr/testify/suite"
	"goboilerplate/pkg/testhelper"
	"testing"
)

type storageTestSuite struct {
	suite.Suite
}

func TestStorageTestSuite(t *testing.T) {
	suite.Run(t, new(storageTestSuite))
}

func (suite *storageTestSuite) SetupSuite() {
	testhelper.InitTestDB()
}

func (suite *storageTestSuite) SetupTest() {
	if err := testhelper.CleanDB(testhelper.TestDB); err != nil {
		suite.FailNow("Couldn't clean db", err)
	}
}
