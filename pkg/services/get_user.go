package services

import (
	"context"
	"errors"
	"github.com/in-rich/uservice-linkedin-data/pkg/dao"
	"github.com/in-rich/uservice-linkedin-data/pkg/models"
)

type GetUserService interface {
	Exec(ctx context.Context, publicIdentifier string) (*models.User, error)
}

type getUserServiceImpl struct {
	getUserRepository           dao.GetUserRepository
	getProfilePictureRepository dao.GetProfilePictureRepository
}

func (s *getUserServiceImpl) Exec(ctx context.Context, publicIdentifier string) (*models.User, error) {
	user, err := s.getUserRepository.GetUser(ctx, publicIdentifier)
	if err != nil {
		return nil, err
	}

	profilePicture, err := s.getProfilePictureRepository.GetProfilePicture(ctx, publicIdentifier)
	if err != nil && !errors.Is(err, dao.ErrProfilePictureNotFound) {
		return nil, err
	}

	return &models.User{
		PublicIdentifier: user.PublicIdentifier,
		FirstName:        user.FirstName,
		LastName:         user.LastName,
		ProfilePicture:   profilePicture,
	}, nil
}

func NewGetUserService(getUserRepository dao.GetUserRepository, getProfilePictureRepository dao.GetProfilePictureRepository) GetUserService {
	return &getUserServiceImpl{
		getUserRepository:           getUserRepository,
		getProfilePictureRepository: getProfilePictureRepository,
	}
}
