package integration

import (
	"goboilerplate/pkg/httpapp/integration_test/request_builder"
	"net/http"
)

func (suite *testSuite) Test_Health() {
	res := request_builder.NewRequestBuilder(http.MethodGet, "/health").
		DoIntegrationRequest()

	suite.Equal(http.StatusOK, res.Code)
}
