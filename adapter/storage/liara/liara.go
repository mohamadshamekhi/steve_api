package liarastorage

import (
	"context"
	"fmt"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	mediaservice "github.com/abolfazlnorzad/steve/service/media"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/oklog/ulid/v2"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

type Config struct {
	AwsAccessKeyID     string `koanf:"aws_access_key_id"`
	AwsSecretAccessKey string `koanf:"aws_secret_access_key"`
	Endpoint           string `koanf:"endpoint"`
	BucketName         string `koanf:"bucket_name"`
}

type Disk struct {
	config   Config
	client   *s3.Client
	MediaCfg mediaservice.Config
}

func New(credential Config, mcfg mediaservice.Config) *Disk {
	cfg, lErr := config.LoadDefaultConfig(context.Background(), config.WithRegion("us-west-2"))
	if lErr != nil {
		logger.L().Error("error in load config in liara object storage :", lErr)
	}
	cfg.Credentials = aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
		return aws.Credentials{
			AccessKeyID:     credential.AwsAccessKeyID,
			SecretAccessKey: credential.AwsSecretAccessKey,
		}, nil
	})

	cfg.BaseEndpoint = aws.String(credential.Endpoint)

	client := s3.NewFromConfig(cfg)

	return &Disk{
		config:   credential,
		client:   client,
		MediaCfg: mcfg,
	}
}

func (d *Disk) Upload(ctx context.Context, domain string, file multipart.File, header *multipart.FileHeader) (string, error) {
	const op = "liarastorage.Upload"

	if header.Size > d.MediaCfg.MaxUploadSizeBytes {
		return "", richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("media.file_too_large").
			WithMeta(map[string]any{"max_size_mb": d.MediaCfg.MaxUploadSizeBytes, "actual_size_bytes": header.Size})
	}

	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	if _, seekErr := file.Seek(0, 0); seekErr != nil {
		return "", richerror.New(op).WithWarpError(seekErr).WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	mimeType := http.DetectContentType(buffer)
	isValidMimeType := false
	for _, allowedType := range d.MediaCfg.AllowedExtensions {
		if mimeType == allowedType {
			isValidMimeType = true
			break
		}
	}
	if !isValidMimeType {
		return "", richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("media.mime_type_not_allowed").
			WithMeta(map[string]any{"mime_type": mimeType})
	}

	fileType := getFileTypeFromMime(mimeType)

	datePath := time.Now().Format("2006/01/02")

	ext := filepath.Ext(header.Filename)
	uniqueFileName := fmt.Sprintf("%s%s", ulid.Make().String(), ext)

	filePath := fmt.Sprintf("%s/%s/%s/%s", domain, fileType, datePath, uniqueFileName)

	_, err = d.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(d.config.BucketName),
		Key:         aws.String(filePath),
		Body:        file,
		ContentType: aws.String(mimeType),
	})
	if err != nil {
		return "", richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	return filePath, nil
}

func (d *Disk) GetURL(ctx context.Context, path string, expiration time.Duration) (string, error) {
	pClient := s3.NewPresignClient(d.client)
	request, err := pClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(d.config.BucketName),
		Key:    aws.String(path),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expiration
	})
	if err != nil {
		return "", richerror.New("adapter.s3.GeneratePresignedURL").WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	return request.URL, nil
}

func (d *Disk) Delete(ctx context.Context, filePath string) error {
	_, err := d.client.DeleteObject(ctx, &s3.DeleteObjectInput{Bucket: aws.String(d.config.BucketName), Key: aws.String(filePath)})

	return err
}
func getFileTypeFromMime(mimeType string) string {
	if strings.HasPrefix(mimeType, "image/") {
		return "images"
	}
	if strings.HasPrefix(mimeType, "video/") {
		return "videos"
	}
	if strings.HasPrefix(mimeType, "audio/") {
		return "audios"
	}
	return "documents"
}
