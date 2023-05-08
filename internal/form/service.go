package form

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"code-challenge/internal/entity"
	"code-challenge/pkg/log"

	"github.com/h2non/filetype"
)

type repo interface {
	Create(ctx context.Context, metadata *entity.Metadata) (string, error)
}

type Service struct {
	imageDir string

	repo repo

	logger log.Logger
}

func NewService(imageDir string, repo repo, logger log.Logger) *Service {
	return &Service{
		imageDir: imageDir,
		repo:     repo,
		logger:   logger,
	}
}

func (svc *Service) Upload(ctx context.Context, fileHeader *multipart.FileHeader) (string, error) {
	svc.logger.Info("Uploading file: " + fileHeader.Filename)

	src, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	_ = src.Close()

	buf, err := io.ReadAll(src)
	if err != nil {
		return "", err
	}

	// Determine the mime type
	if !filetype.IsImage(buf) {
		return "", fmt.Errorf("not an image")
	}
	kind, _ := filetype.Match(buf)
	if kind == filetype.Unknown {
		return "", fmt.Errorf("unknown file type detected, file type: %s", kind.Extension)
	}

	// Save image to file system
	path := svc.imageDir + "/" + fileHeader.Filename
	if err = os.MkdirAll(filepath.Dir(path), 0750); err != nil {
		return "", err
	}
	out, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = out.Close()
	}()
	_, err = io.Copy(out, src)

	// Save metadata to database
	_, err = svc.repo.Create(ctx, &entity.Metadata{
		Name:     fileHeader.Filename,
		Type:     kind.MIME.Value,
		Size:     fileHeader.Size,
		Location: path,
	})
	if err != nil {
		return "", err
	}
	return path, nil
}
