package request_builder

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type GqlRequestBuilder struct {
	gqlQuery string
	gqlVars  string
	*RequestBuilder
}

func NewGraphQLRequestBuilder(gqlQuery string) *GqlRequestBuilder {
	bb := NewRequestBuilder(http.MethodPost, "/graphql")

	b := &GqlRequestBuilder{
		gqlQuery:       gqlQuery,
		gqlVars:        "{}",
		RequestBuilder: bb,
	}
	b.SetHeader("Content-Type", "application/json")
	b.body = strings.NewReader(b.reqBody())

	return b
}

func (b *GqlRequestBuilder) reqBody() string {
	return fmt.Sprintf(`{
		"operationName": "",
		"variables": %v,
		"query": %v
	}`, b.gqlVars, strconv.Quote(b.gqlQuery))
}

func (b *GqlRequestBuilder) SetGqlVars(gqlVars string) *GqlRequestBuilder {
	b.gqlVars = gqlVars
	b.body = strings.NewReader(b.reqBody())

	return b
}
