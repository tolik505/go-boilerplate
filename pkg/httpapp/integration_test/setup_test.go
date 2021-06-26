package integration

import (
	"github.com/stretchr/testify/suite"
	"goboilerplate/pkg/httpapp/integration_test/request_builder"
	"goboilerplate/pkg/testhelper"
	"testing"
)

type testSuite struct {
	suite.Suite
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(testSuite))
}

func (suite *testSuite) SetupSuite() {
	request_builder.TestDB = testhelper.InitSpecificTestDB(request_builder.DBName)
}

func (suite *testSuite) SetupTest() {
	if err := testhelper.CleanDB(request_builder.TestDB); err != nil {
		suite.FailNow("Couldn't clean db", err)
	}
}
