package mediaservice

import (
	"context"
	"errors"

	"github.com/abolfazlnorzad/steve/entity"
	mediaparam "github.com/abolfazlnorzad/steve/param/media"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (s Service) UploadFile(ctx context.Context, req mediaparam.UploadMediaRequest) (mediaparam.UploadMediaResponse, error) {
	const op = "mediaservice.UploadFile"

	res, uErr := s.provider.Upload(ctx, req.Domain, req.IsPublic, req.File, req.Header)
	if uErr != nil {
		var rErr *richerror.RichError
		if errors.As(uErr, &rErr) && rErr.Kind() == richerror.KindInvalid {
			return mediaparam.UploadMediaResponse{}, uErr
		}
		return mediaparam.UploadMediaResponse{}, richerror.New(op).WithWarpError(uErr).WithKind(richerror.KindUnexpected)
	}
	media := &entity.Media{
		OriginalFileName: req.Header.Filename,
		StorageProvider:  "liara-s3", // todo - add to cfg
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
