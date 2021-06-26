package integration

import (
	"goboilerplate/pkg/httpapp/integration_test/request_builder"
	"goboilerplate/pkg/service/postservice/mocks"
	"goboilerplate/pkg/testhelper"
	"net/http"
)

func (suite *testSuite) Test_CreatePost() {
	uuidService := &mocks.UuidService{}
	uuidService.On("Generate").Return("123e4567-e89b-12d3-a456-426655440000")
	vars := `{
		"p": {
			"content": "Test post",
			"metadata": "{\"name\": \"John\"}"
		}
	}`
	query := `
			mutation cp($p: CreatePostInput!) {
				createPost (data: $p)
			}
		`
	res := request_builder.NewGraphQLRequestBuilder(query).
		SetGqlVars(vars).
		SetUUIDService(uuidService).
		DoIntegrationRequest()
	expected := `{
			  "data": {
				"createPost": "123e4567-e89b-12d3-a456-426655440000"
			  }
			}`
	suite.JSONEq(expected, res.Body.String())
	suite.Equal(http.StatusOK, res.Code)
}

func (suite *testSuite) Test_QueryPosts() {
	testhelper.SeedFixtures(suite.T(), request_builder.TestDB, "testdata/fixtures/posts")
	query := `
		query p {
			posts {
				uuid
				content
				metadata
			}
		}
	`
	res := request_builder.NewGraphQLRequestBuilder(query).DoIntegrationRequest()
	expected := `{
			"data": {
				"posts": [
					{
						"uuid": "uuid-1",
						"content": "Post 1",
						"metadata": null
					},
					{
						"uuid": "uuid-2",
						"content": "Post 2",
						"metadata": "{\"title\": \"test\"}"
					}
				]
			}
		}
	`
	suite.JSONEq(expected, res.Body.String())
	suite.Equal(http.StatusOK, res.Code)
}
