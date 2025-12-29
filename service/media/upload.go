package mediaservice

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/abolfazlnorzad/steve/entity"
	mediaparam "github.com/abolfazlnorzad/steve/param/media"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (s Service) UploadFile(ctx context.Context, req mediaparam.UploadMediaRequest) (mediaparam.UploadMediaResponse, error) {
	const op = "mediaservice.UploadFile"
	defer func() {
		if err := req.File.Close(); err != nil {
			logger.L().Error("failed to close file", "err", err)
		}
	}()
	isImg, mimeType, err := IsImage(req.File)
	if err != nil {
		return mediaparam.UploadMediaResponse{}, richerror.New(op).WithWarpError(err)
	}

	if seeker, ok := req.File.(io.Seeker); ok {
		_, _ = seeker.Seek(0, io.SeekStart)
	}
	meta := mediaparam.FileMetadata{
		Filename: req.Header.Filename,
		//nolint
		Size:        uint64(req.Header.Size),
		ContentType: mimeType,
	}
	var finalFile io.Reader = req.File
	fmt.Println("isImg", isImg)
	fmt.Println("req.ConvertToWebp", req.ConvertToWebp)
	if isImg && req.ConvertToWebp && mimeType != "image/webp" {
		convertedReader, newSize, err := s.provider.TransformToWebP(req.File)
		if err != nil {
			return mediaparam.UploadMediaResponse{}, richerror.New(op).WithWarpError(err)
		}
		finalFile = convertedReader
		//nolint
		meta.Size = uint64(newSize)
		meta.ContentType = "image/webp"
		meta.Filename = strings.TrimSuffix(meta.Filename, filepath.Ext(meta.Filename)) + ".webp"
	}

	res, uErr := s.provider.Upload(ctx, req.Domain, req.UploaderType, req.IsPublic, finalFile, meta)
	if uErr != nil {
		var rErr *richerror.RichError
		if errors.As(uErr, &rErr) && rErr.Kind() == richerror.KindInvalid {
			return mediaparam.UploadMediaResponse{}, uErr
		}
		return mediaparam.UploadMediaResponse{}, richerror.New(op).WithWarpError(uErr).WithKind(richerror.KindUnexpected)
	}
	media := &entity.Media{
		OriginalFileName: req.Header.Filename,
		StorageProvider:  "arvan", // todo - add to cfg
		Path:             res.Path,
		IsPublic:         req.IsPublic,
		MimeType:         res.MimeType,
		Size:             req.Header.Size,
		UploaderID:       req.UploaderID,
		UploaderType:     req.UploaderType,
		Bucket:           res.BucketName,
	}
	media.BeforeCreate()
	if _, err := s.repository.AddMedia(ctx, *media); err != nil {
		_ = s.provider.Delete(context.Background(), media)
		return mediaparam.UploadMediaResponse{}, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}
	r, _ := s.GetURLByMediaID(ctx, mediaparam.GetUrlByIDRequest{
		ID: media.ID,
	})
	return mediaparam.UploadMediaResponse{
		ID:  media.ID,
		URL: r.URL,
	}, nil
}

func (s Service) StorefrontUploadFile(ctx context.Context, req mediaparam.StorefrontUploadMediaRequest) (mediaparam.StorefrontUploadMediaResponse, error) {
	const op = "mediaservice.StorefrontUploadFile"

	defer func() {
		if err := req.File.Close(); err != nil {
			logger.L().Error("failed to close file", "err", err)
		}
	}()
	isImg, mimeType, err := IsImage(req.File)
	if err != nil {
		return mediaparam.StorefrontUploadMediaResponse{}, richerror.New(op).WithWarpError(err)
	}

	if seeker, ok := req.File.(io.Seeker); ok {
		_, _ = seeker.Seek(0, io.SeekStart)
	}
	meta := mediaparam.FileMetadata{
		Filename: req.Header.Filename,
		//nolint
		Size:        uint64(req.Header.Size),
		ContentType: mimeType,
	}
	var finalFile io.Reader = req.File
	if isImg && req.ConvertToWebp && mimeType != "image/webp" {
		convertedReader, newSize, err := s.provider.TransformToWebP(req.File)
		if err != nil {
			return mediaparam.StorefrontUploadMediaResponse{}, richerror.New(op).WithWarpError(err)
		}
		finalFile = convertedReader
		//nolint
		meta.Size = uint64(newSize)
		meta.ContentType = "image/webp"
		meta.Filename = strings.TrimSuffix(meta.Filename, filepath.Ext(meta.Filename)) + ".webp"
	}

	res, uErr := s.provider.Upload(ctx, req.Domain, "user", true, finalFile, meta)

	if uErr != nil {
		var rErr *richerror.RichError
		if errors.As(uErr, &rErr) && rErr.Kind() == richerror.KindInvalid {
			return mediaparam.StorefrontUploadMediaResponse{}, uErr
		}
		return mediaparam.StorefrontUploadMediaResponse{}, richerror.New(op).WithWarpError(uErr).WithKind(richerror.KindUnexpected)
	}
	media := &entity.Media{
		OriginalFileName: req.Header.Filename,
		StorageProvider:  "arvan", // todo - add to cfg
		Path:             res.Path,
		IsPublic:         true,
		MimeType:         res.MimeType,
		Size:             req.Header.Size,
		UploaderID:       req.UploaderID,
		UploaderType:     "user",
		Bucket:           res.BucketName,
	}
	media.BeforeCreate()
	if _, err := s.repository.AddMedia(ctx, *media); err != nil {
		_ = s.provider.Delete(context.Background(), media)
		return mediaparam.StorefrontUploadMediaResponse{}, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}
	r, _ := s.GetURLByMediaID(ctx, mediaparam.GetUrlByIDRequest{
		ID: media.ID,
	})
	return mediaparam.StorefrontUploadMediaResponse{
		ID:  media.ID,
		URL: r.URL,
	}, nil
}

func IsImage(file io.ReadSeeker) (bool, string, error) {
	buffer := make([]byte, 512)
	if _, err := file.Read(buffer); err != nil && err != io.EOF {
		return false, "", err
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return false, "", err
	}

	mimeType := http.DetectContentType(buffer)
	return strings.HasPrefix(mimeType, "image/"), mimeType, nil
}
