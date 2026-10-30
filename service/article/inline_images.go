package articleservice

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"regexp"
	"strings"

	mediaparam "github.com/abolfazlnorzad/steve/param/media"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

// inlineImageRegex finds an <img> src holding a base64 data URI, the way rich
// text editors embed pasted images. Group 1 is the whole data URI, group 2 the
// subtype, group 3 the payload.
var inlineImageRegex = regexp.MustCompile(`(?i)<img\b[^>]*?\bsrc\s*=\s*["'](data:image/([a-z0-9.+-]+);base64,([^"']+))["']`)

// externalizeInlineImages uploads every base64 image embedded in body to
// object storage and swaps it for the public URL, so article bodies stay
// small and images get cached by CDNs and indexed by search engines. The
// same data URI appearing twice is uploaded once.
func (s BackofficeService) externalizeInlineImages(ctx context.Context, body string, uploaderID string) (string, int, error) {
	const op = "articleservice.externalizeInlineImages"

	matches := inlineImageRegex.FindAllStringSubmatch(body, -1)
	if len(matches) == 0 {
		return body, 0, nil
	}

	replaced := make(map[string]string, len(matches))
	for i, m := range matches {
		dataURI, subtype, payload := m[1], strings.ToLower(m[2]), m[3]
		if _, done := replaced[dataURI]; done {
			continue
		}

		raw, err := base64.StdEncoding.DecodeString(strings.Join(strings.Fields(payload), ""))
		if err != nil {
			return "", 0, richerror.New(op).
				WithKind(richerror.KindInvalid).
				WithMessageKey("article.invalid_inline_image").
				WithWarpError(err)
		}

		ext := subtype
		if ext == "jpeg" {
			ext = "jpg"
		}
		name := fmt.Sprintf("article-inline-%d.%s", i+1, ext)

		res, err := s.mediaSvc.UploadFile(ctx, mediaparam.UploadMediaRequest{
			Domain:        "article",
			UploaderID:    uploaderID,
			UploaderType:  "backoffice",
			IsPublic:      true,
			File:          memFile{bytes.NewReader(raw)},
			Header:        &multipart.FileHeader{Filename: name, Size: int64(len(raw))},
			ConvertToWebp: true,
		})
		if err != nil {
			return "", 0, richerror.New(op).WithWarpError(err)
		}
		replaced[dataURI] = res.URL
	}

	pairs := make([]string, 0, len(replaced)*2)
	for dataURI, url := range replaced {
		pairs = append(pairs, dataURI, url)
	}
	return strings.NewReplacer(pairs...).Replace(body), len(replaced), nil
}

// memFile adapts an in-memory image to multipart.File for the media service.
type memFile struct{ *bytes.Reader }

func (memFile) Close() error { return nil }
