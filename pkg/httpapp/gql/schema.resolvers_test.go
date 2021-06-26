package gql_test

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"goboilerplate/pkg/apperr"
	"goboilerplate/pkg/httpapp/gql/mocks"
	"goboilerplate/pkg/httpapp/integration_test/request_builder"
	"net/http"
	"testing"
)

func Test_queryResolver_Posts(t *testing.T) {
	t.Run("It Handles Known Error", func(t *testing.T) {
		query := `
		query p {
			posts {
				content
			}
		}`
		s := &mocks.PostService{}
		s.On("GetAll", mock.Anything).Return(nil, errors.Wrap(apperr.NewKnown("CODE", "Known"), "error"))
		res := request_builder.NewGraphQLRequestBuilder(query).
			SetPostService(s).
			DoFunctionalRequest()
		expected := `{
			"errors": [
				{
				  "message": "Known",
				  "extensions": { "code": "CODE" }
				}
			],
			"data": null
		}`
		assert.JSONEq(t, expected, res.Body.String())
		assert.Equal(t, http.StatusOK, res.Code)
	})

	t.Run("It Logs Unexpected Error", func(t *testing.T) {
		hook := new(test.Hook)
		log.AddHook(hook)

		query := `
		query p {
			posts {
				content
			}
		}`
		s := &mocks.PostService{}
		s.On("GetAll", mock.Anything).Return(nil, errors.Wrap(errors.New("cause"), "error"))
		res := request_builder.NewGraphQLRequestBuilder(query).
			SetPostService(s).
			DoFunctionalRequest()
		expected := `{
			"errors": [
				{
					"message": "Internal server error",
					"extensions": { "code": "INTERNAL_SERVER_ERROR" }
				}
			],
			"data": null
		}`
		assert.JSONEq(t, expected, res.Body.String())
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, log.ErrorLevel, hook.LastEntry().Level)
		assert.Contains(t, hook.LastEntry().Message, "error: cause, gqlQuery: \n\t\tquery p {\n\t\t\tposts {\n\t\t\t\tcontent\n\t\t\t}\n\t\t} | gqlVariables: map[] | stacktrace: goroutine")
	})
}
