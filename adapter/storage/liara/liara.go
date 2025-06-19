package liarastorage

import (
	"bytes"
	"context"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"mime/multipart"
	"net/http"
	"time"
)

type Config struct {
	AwsAccessKeyID     string `koanf:"aws_access_key_id"`
	AwsSecretAccessKey string `koanf:"aws_secret_access_key"`
	Endpoint           string `koanf:"endpoint"`
	BucketName         string `koanf:"bucket_name"`
}

type Disk struct {
	config Config
	client *s3.Client
}

func New(credential Config) *Disk {
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
		config: credential,
		client: client,
	}
}

func (d *Disk) Upload(ctx context.Context, f multipart.FileHeader, filePath string) error {
	size := f.Size
	buffer := make([]byte, size)
	file, oErr := f.Open()
	if oErr != nil {
		return oErr
	}
	if _, rErr := file.Read(buffer); rErr != nil {
		return rErr
	}
	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)
	uploader := manager.NewUploader(d.client, func(u *manager.Uploader) {
		u.PartSize = entity.MaxFileUploadSize
	})
	param := &s3.PutObjectInput{
		Bucket:        aws.String(d.config.BucketName),
		Key:           aws.String(filePath),
		Body:          fileBytes,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(fileType),
	}
	_, uErr := uploader.Upload(ctx, param)

	if uErr != nil {
		return uErr
	}

	return nil
}

func (d *Disk) GetURL(ctx context.Context, filePath string) (string, error) {
	pClient := s3.NewPresignClient(d.client)
	request, err := pClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(d.config.BucketName),
		Key:    aws.String(filePath),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(entity.FileLinkExpirationDuration)
	})
	if err != nil {
		return "", err
	}

	return request.URL, nil
}

func (d *Disk) Delete(ctx context.Context, filePath string) error {
	_, err := d.client.DeleteObject(ctx, &s3.DeleteObjectInput{Bucket: aws.String(d.config.BucketName), Key: aws.String(filePath)})

	return err
}
