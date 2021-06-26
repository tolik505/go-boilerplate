package request_builder

import (
	"goboilerplate/pkg/api/minio"
	mocks2 "goboilerplate/pkg/api/minio/mocks"
	mocks3 "goboilerplate/pkg/httpapp/gql/mocks"
	mocks4 "goboilerplate/pkg/httpapp/mocks"
	"goboilerplate/pkg/service/postservice/mocks"
	"gorm.io/gorm"
	"io"
	"net/http/httptest"
)

//TestDB contains test DB instance
var TestDB *gorm.DB

const DBName = "http_integration_test"

type RequestBuilder struct {
	method  string
	route   string
	us      *mocks.UuidService
	mcl     *mocks2.Client
	ps      *mocks3.PostService
	fs      *mocks4.FileService
	body    io.Reader
	headers map[string][]string
}

func NewRequestBuilder(method, route string) *RequestBuilder {
	rb := &RequestBuilder{
		method:  method,
		route:   route,
		headers: make(map[string][]string),
	}

	return rb
}

func (b *RequestBuilder) DoIntegrationRequest() *httptest.ResponseRecorder {
	req := httptest.NewRequest(b.method, b.route, b.body)
	req.Header = b.headers

	mc := minio.Config{
		Host:   "host",
		Key:    "key",
		Secret: "secret",
	}
	bucket := minio.Bucket("bucket")
	httpApp, _ := InitializeHTTPAppIntegration(TestDB, b.us, b.mcl, mc, bucket, "")

	res := httptest.NewRecorder()
	httpApp.Server.ServeHTTP(res, req)

	return res
}

func (b *RequestBuilder) DoFunctionalRequest() *httptest.ResponseRecorder {
	req := httptest.NewRequest(b.method, b.route, b.body)
	req.Header = b.headers

	httpApp, _ := InitializeHTTPAppFunctional(b.ps, b.fs, "")

	res := httptest.NewRecorder()
	httpApp.Server.ServeHTTP(res, req)

	return res
}

func (b *RequestBuilder) SetBody(body io.Reader) *RequestBuilder {
	b.body = body

	return b
}

func (b *RequestBuilder) SetHeader(name, value string) *RequestBuilder {
	b.headers[name] = []string{value}

	return b
}

func (b *RequestBuilder) RewriteHeaders(h map[string][]string) *RequestBuilder {
	b.headers = h

	return b
}

func (b *RequestBuilder) SetMinioClient(mcl *mocks2.Client) *RequestBuilder {
	b.mcl = mcl

	return b
}

func (b *RequestBuilder) SetUUIDService(us *mocks.UuidService) *RequestBuilder {
	b.us = us

	return b
}

func (b *RequestBuilder) SetPostService(ps *mocks3.PostService) *RequestBuilder {
	b.ps = ps

	return b
}

func (b *RequestBuilder) SetFileService(fs *mocks4.FileService) *RequestBuilder {
	b.fs = fs

	return b
}
