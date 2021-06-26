package httpapp

import (
	"context"
	"github.com/gabriel-vasile/mimetype"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"goboilerplate/pkg/apperr"
	"io"
	"mime/multipart"
	"net/http"
	"syscall"
)

type FileService interface {
	Upload(src io.Reader, contentType string) (string, error)
}

type FileResponse struct {
	UUID     string `json:"file_uuid"`
	MimeType string `json:"type"`
	Size     int    `json:"size"`
	Name     string `json:"name"`
}

type Handler interface {
	UploadFile(c echo.Context) error
}

type handler struct {
	fileService FileService
}

func NewHandler(
	fileService FileService,
) Handler {
	return &handler{
		fileService,
	}
}

func (h *handler) UploadFile(c echo.Context) error {
	var files []*FileResponse

	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	formFiles := form.File
	for _, f := range formFiles {
		fHeader := f[0]
		src, err := fHeader.Open()
		if err != nil {
			return err
		}

		UUID, fType, err := h.uploadFile(src, fHeader, c.Request().Context())
		if err != nil {
			return err
		}

		fr := &FileResponse{
			UUID:     UUID,
			MimeType: fType,
			Size:     int(fHeader.Size),
			Name:     fHeader.Filename,
		}
		files = append(files, fr)
	}

	return c.JSON(http.StatusOK, files)
}

func (h *handler) uploadFile(
	src multipart.File,
	header *multipart.FileHeader,
	ctx context.Context,
) (string, string, error) {

	defer func(src multipart.File) {
		err := src.Close()
		if err != nil {
			log.Error(errors.Wrap(err, "Couldn't close file"), ctx)
		}
	}(src)

	mime, err := mimetype.DetectReader(src)
	if err != nil {
		return "", "", errors.Wrap(err, "Couldn't detect mime type")
	}
	fType := mime.String()
	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return "", "", errors.Wrap(err, "Couldn't set pointer to the file start")
	}
	UUID, err := h.fileService.Upload(src, fType)
	if err != nil {
		return "", "", errors.Wrapf(err, "Couldn't upload file [%s]", header.Filename)
	}

	return UUID, fType, nil
}

func customHTTPErrorHandler(err error, c echo.Context) {
	var code int
	var msg interface{}

	causeErr := errors.Cause(err)
	switch e := causeErr.(type) {
	case apperr.Known:
		code = http.StatusBadRequest
		msg = echo.Map{"message": e.Error(), "code": e.Code()}
	case *echo.HTTPError:
		code = e.Code
		msg = echo.Map{"message": e.Error()}
	default:
		if errors.Is(err, syscall.EPIPE) || errors.Is(err, syscall.ECONNRESET) {
			log.Warn(err)
			return
		}
		code = http.StatusInternalServerError
		msg = echo.Map{"message": "Internal server error", "code": InternalServerErrorCode}

		log.Error(err)
	}
	if err := c.JSON(code, msg); err != nil {
		log.Error(err)
	}
}
