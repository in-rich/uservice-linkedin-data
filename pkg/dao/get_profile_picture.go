package dao

import (
	"cloud.google.com/go/storage"
	"context"
	"errors"
)

type GetProfilePictureRepository interface {
	GetProfilePicture(ctx context.Context, publicIdentifier string) (string, error)
}

type getProfilePictureRepositoryImpl struct {
	bucket *storage.BucketHandle
}

func (r *getProfilePictureRepositoryImpl) GetProfilePicture(ctx context.Context, publicIdentifier string) (string, error) {
	obj := r.bucket.Object(publicIdentifier)

	attrs, err := obj.Attrs(ctx)
	if err != nil {
		if errors.Is(err, storage.ErrObjectNotExist) {
			return "", ErrProfilePictureNotFound
		}

		return "", err
	}

	return attrs.MediaLink, nil
}

func NewGetProfilePictureRepository(bucket *storage.BucketHandle) GetProfilePictureRepository {
	return &getProfilePictureRepositoryImpl{
		bucket: bucket,
	}
}
