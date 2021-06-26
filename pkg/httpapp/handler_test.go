package httpapp_test

import (
	"bytes"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"goboilerplate/pkg/apperr"
	"goboilerplate/pkg/httpapp/integration_test/request_builder"
	"goboilerplate/pkg/httpapp/mocks"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

func Test_handler_UploadFile(t *testing.T) {
	t.Run("It Handles Known Error", func(t *testing.T) {
		fs := &mocks.FileService{}
		err := errors.Wrap(apperr.NewKnown("CODE", "known"), "error")
		ft := "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
		body, writer := prepareBody()

		fs.On("Upload", mock.AnythingOfType("multipart.sectionReadCloser"), ft).Return("", err)

		res := request_builder.NewRequestBuilder(http.MethodPost, "/upload").
			SetHeader("Content-Type", writer.FormDataContentType()).
			SetBody(body).
			SetFileService(fs).
			DoFunctionalRequest()

		assert.Equal(t, http.StatusBadRequest, res.Code)

		expBody := `{"code":"CODE","message":"known"}`
		actualBody := res.Body.String()

		assert.JSONEq(t, expBody, actualBody)
	})

	t.Run("It Logs Unexpected Error", func(t *testing.T) {
		hook := new(test.Hook)
		log.AddHook(hook)

		fs := &mocks.FileService{}
		err := errors.Wrap(errors.New("cause"), "error")
		ft := "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
		body, writer := prepareBody()

		fs.On("Upload", mock.AnythingOfType("multipart.sectionReadCloser"), ft).Return("", err)

		res := request_builder.NewRequestBuilder(http.MethodPost, "/upload").
			SetHeader("Content-Type", writer.FormDataContentType()).
			SetBody(body).
			SetFileService(fs).
			DoFunctionalRequest()

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		expBody := `{"code":"INTERNAL_SERVER_ERROR","message":"Internal server error"}`
		actualBody := res.Body.String()

		assert.JSONEq(t, expBody, actualBody)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, log.ErrorLevel, hook.LastEntry().Level)
		assert.Equal(t, "Couldn't upload file [file.docx]: error: cause", hook.LastEntry().Message)
	})
}

func prepareBody() (io.Reader, *multipart.Writer) {
	path := "integration_test/testdata/file.docx"
	file, _ := os.Open(path)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file.docx", filepath.Base(path))

	_, _ = io.Copy(part, file)
	_ = writer.Close()

	return body, writer
}
