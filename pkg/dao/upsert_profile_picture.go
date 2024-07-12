package dao

import (
	"cloud.google.com/go/storage"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/samber/lo"
	"strings"
)

var supported_formats = []string{"data:image/png", "data:image/jpeg", "data:image/jpg"}

type UpsertProfilePictureRepository interface {
	UpsertProfilePicture(ctx context.Context, publicIdentifier string, base64 string) (string, error)
}

type upsertProfilePictureRepositoryImpl struct {
	bucket *storage.BucketHandle
}

// Extract raw data from base64 string.
func (r *upsertProfilePictureRepositoryImpl) decodeBase64(raw string) ([]byte, error) {
	parts := strings.Split(raw, ";base64,")

	if len(parts) != 2 {
		return nil, errors.New("invalid base64 string")
	}

	// Check if the format is supported.
	if !lo.Contains(supported_formats, parts[0]) {
		return nil, fmt.Errorf("unsupported format: %s", parts[0])
	}

	decoded, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 string: %w", err)
	}

	return decoded, nil
}

func (r *upsertProfilePictureRepositoryImpl) UpsertProfilePicture(ctx context.Context, publicIdentifier string, base64 string) (string, error) {
	decoded, err := r.decodeBase64(base64)
	if err != nil {
		return "", errors.Join(err, ErrInvalidProfilePicture)
	}

	obj := r.bucket.Object(publicIdentifier)
	writer := obj.NewWriter(ctx)

	// Delete object if exists.
	if err := obj.Delete(ctx); err != nil && !errors.Is(err, storage.ErrObjectNotExist) {
		return "", fmt.Errorf("failed to delete object %s: %w", publicIdentifier, err)
	}

	if _, err := writer.Write(decoded); err != nil {
		_ = writer.Close()
		return "", fmt.Errorf("failed to write object %s: %w", publicIdentifier, err)
	}

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer for object %s: %w", publicIdentifier, err)
	}

	attrs, err := obj.Attrs(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get attributes for object %s: %w", publicIdentifier, err)
	}

	return attrs.MediaLink, nil
}

func NewUpsertProfilePictureRepository(bucket *storage.BucketHandle) UpsertProfilePictureRepository {
	return &upsertProfilePictureRepositoryImpl{
		bucket: bucket,
	}
}
