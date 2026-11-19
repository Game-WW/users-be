package avatar

import (
	"bytes"
	"fmt"
	"image"

	"github.com/google/uuid"

	_ "image/jpeg"
	_ "image/png"

	_ "golang.org/x/image/webp"
)

const (
	MaxBytes     = 5 << 20
	minDimension = 64
	maxDimension = 4096
)

var allowedFormats = map[string]string{
	"jpeg": "image/jpeg",
	"png":  "image/png",
	"webp": "image/webp",
}

func Key(userID string) string {
	return fmt.Sprintf("avatars/%s", userID)
}

func NormalizeID(value string) (string, error) {
	id, err := uuid.Parse(value)
	if err != nil || id == uuid.Nil {
		return "", ErrInvalidUserID
	}
	return id.String(), nil
}

func Validate(data []byte) (contentType string, err error) {
	if len(data) == 0 {
		return "", ErrInvalidImage
	}
	if len(data) > MaxBytes {
		return "", ErrTooLarge
	}

	cfg, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return "", ErrInvalidImage
	}

	contentType, ok := allowedFormats[format]
	if !ok {
		return "", ErrUnsupportedFormat
	}

	if cfg.Width < minDimension || cfg.Height < minDimension ||
		cfg.Width > maxDimension || cfg.Height > maxDimension {
		return "", ErrInvalidDimensions
	}

	if _, _, err := image.Decode(bytes.NewReader(data)); err != nil {
		return "", ErrInvalidImage
	}

	return contentType, nil
}
