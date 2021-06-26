package integration

import (
	"bytes"
	"context"
	minio2 "github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/mock"
	mocks2 "goboilerplate/pkg/api/minio/mocks"
	"goboilerplate/pkg/httpapp/integration_test/request_builder"
	"goboilerplate/pkg/service/postservice/mocks"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func (suite *testSuite) Test_UploadFile() {
	path := "testdata/file.docx"
	file, err := os.Open(path)
	if err != nil {
		suite.Error(err)
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file.docx", filepath.Base(path))
	_, _ = io.Copy(part, file)
	_ = writer.Close()

	us := &mocks.UuidService{}
	us.On("Generate").Return("123e4567-e89b-12d3-a456-426655440000")

	mcl := &mocks2.Client{}
	mcl.On(
		"PutObject",
		context.Background(),
		"bucket",
		"123e4567-e89b-12d3-a456-426655440000",
		mock.AnythingOfType("sectionReadCloser"),
		int64(21635),
		minio2.PutObjectOptions{ContentType: "application/vnd.openxmlformats-officedocument.wordprocessingml.document"},
	).Return(minio2.UploadInfo{}, nil).Once()

	res := request_builder.NewRequestBuilder(http.MethodPost, "/upload").
		SetHeader("Content-Type", writer.FormDataContentType()).
		SetBody(body).
		SetUUIDService(us).
		SetMinioClient(mcl).
		DoIntegrationRequest()

	mcl.AssertExpectations(suite.T())
	expBody := `[
			{
				"file_uuid": "123e4567-e89b-12d3-a456-426655440000",
				"size": 21635,
			 	"name": "file.docx",
				"type": "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
			}
		]`
	suite.Equal(http.StatusOK, res.Code)
	suite.JSONEq(expBody, res.Body.String())
}
