package liarastorage

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"log"
	"math/rand/v2"

	"github.com/abolfazlnorzad/steve/entity"
	mediaparam "github.com/abolfazlnorzad/steve/param/media"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	mediaservice "github.com/abolfazlnorzad/steve/service/media"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/nickalie/go-webpbin"
)

type Config struct {
	AwsAccessKeyID     string `koanf:"aws_access_key_id"`
	AwsSecretAccessKey string `koanf:"aws_secret_access_key"`
	Endpoint           string `koanf:"endpoint"`
	PublicBucketName   string `koanf:"public_bucket_name"`
	PrivateBucketName  string `koanf:"private_bucket_name"`
}

type Disk struct {
	config        Config
	client        *s3.Client
	presignClient *s3.PresignClient
	MediaCfg      mediaservice.Config
}

func New(credential Config, mcfg mediaservice.Config) *Disk {

	creds := credentials.NewStaticCredentialsProvider(credential.AwsAccessKeyID, credential.AwsSecretAccessKey, "")
	fmt.Println("credentials: ", creds)
	cfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithCredentialsProvider(creds),
		awsconfig.WithRegion("us-east-1"),
	)
	if err != nil {
		log.Fatalf("Failed to create sdk config: %v", err)
	}
	fmt.Println("eeeee", credential.Endpoint)

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(credential.Endpoint)
		o.UsePathStyle = true
	})

	presignClient := s3.NewPresignClient(client)

	return &Disk{
		config:        credential,
		client:        client,
		MediaCfg:      mcfg,
		presignClient: presignClient,
	}
}

func (d *Disk) GetURL(ctx context.Context, media *entity.Media) (string, error) {

	if media.IsPublic {
		return fmt.Sprintf("%s/%s", "https://steve-storage-dev.storage.c2.liara.space", media.Path), nil
	}

	targetBucket := d.config.PrivateBucketName
	if media.IsPublic {
		targetBucket = d.config.PublicBucketName
	}
	presignedReq, err := d.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(targetBucket),
		Key:    aws.String(media.Path),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Hour
	})
	if err != nil {
		return "", err
	}

	return presignedReq.URL, nil
}

func (d *Disk) Delete(ctx context.Context, media *entity.Media) error {
	targetBucket := d.config.PrivateBucketName
	if media.IsPublic {
		targetBucket = d.config.PublicBucketName
	}
	_, err := d.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(targetBucket),
		Key:    aws.String(media.Path),
	})

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

func (d *Disk) Upload(
	ctx context.Context,
	domain string,
	role string,
	isPublic bool,
	file io.Reader,
	info mediaparam.FileMetadata,
) (mediaservice.UploadResult, error) {
	const op = "liarastorage.Upload"
	policy := d.MediaCfg.GetPolicy(role)

	if int64(info.Size) > policy.MaxUploadSizeBytes {
		return mediaservice.UploadResult{}, richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("media.file_too_large").
			WithMeta(map[string]any{
				"max_size_bytes":    policy.MaxUploadSizeBytes,
				"actual_size_bytes": info.Size,
			})
	}

	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return mediaservice.UploadResult{}, richerror.New(op).
			WithWarpError(err).
			WithKind(richerror.KindUnexpected).
			WithMessageKey("global.something_wrong")
	}

	mimeType := http.DetectContentType(buffer)

	isValidMimeType := false
	for _, allowedType := range policy.AllowedMimeTypes {
		if mimeType == allowedType {
			isValidMimeType = true
			break
		}
	}
	if !isValidMimeType {
		return mediaservice.UploadResult{}, richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("media.mime_type_not_allowed").
			WithMeta(map[string]any{"mime_type": mimeType})
	}

	fileType := getFileTypeFromMime(mimeType)

	datePath := time.Now().Format("2006/01/02")
	uniqueFileName := fmt.Sprintf("%d%s", rand.Int(), filepath.Ext(info.Filename))
	filePath := fmt.Sprintf("%s/%s/%s/%s", domain, fileType, datePath, uniqueFileName)

	targetBucket := d.config.PrivateBucketName
	if isPublic {
		targetBucket = d.config.PublicBucketName
	}
	fullReader := io.MultiReader(bytes.NewReader(buffer[:512]), file)

	_, err = d.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(targetBucket),
		Key:         aws.String(filePath),
		Body:        fullReader,
		ContentType: aws.String(info.ContentType),
		//nolint
		ContentLength: aws.Int64(int64(info.Size)),
		ACL:           types.ObjectCannedACLPublicRead,
	})

	if err != nil {
		return mediaservice.UploadResult{}, richerror.New(op).
			WithWarpError(err).
			WithKind(richerror.KindUnexpected)
	}

	return mediaservice.UploadResult{
		Path:       filePath,
		BucketName: targetBucket,
		MimeType:   mimeType,
	}, nil
}

func (d *Disk) TransformToWebP(input io.Reader) (io.ReadSeeker, int64, error) {
	const OP = "mediaadapter.Disk.TransformToWebP"

	img, _, err := image.Decode(input)
	if err != nil {
		return nil, 0, richerror.New(OP).WithWarpError(err).WithMessageKey("فرمت تصویر قابل شناسایی نیست")
	}

	var b bytes.Buffer
	if err := webpbin.NewCWebP().
		Quality(75).
		InputImage(img).
		Output(&b).
		Run(); err != nil {
		return nil, 0, richerror.New(OP).WithWarpError(err)
	}

	return bytes.NewReader(b.Bytes()), int64(b.Len()), nil
}
