package services

import (
	"context"
	"errors"
	"github.com/in-rich/uservice-linkedin-data/pkg/dao"
	"github.com/in-rich/uservice-linkedin-data/pkg/models"
)

type ListUsersService interface {
	Exec(ctx context.Context, publicIdentifiers []string) ([]*models.User, error)
}

type listUsersServiceImpl struct {
	listUsersRepository         dao.ListUsersRepository
	getProfilePictureRepository dao.GetProfilePictureRepository
}

func (s *listUsersServiceImpl) Exec(ctx context.Context, publicIdentifiers []string) ([]*models.User, error) {
	users, err := s.listUsersRepository.ListUsers(ctx, publicIdentifiers)
	if err != nil {
		return nil, err
	}

	result := make([]*models.User, len(users))
	for i, user := range users {
		profilePicture, err := s.getProfilePictureRepository.GetProfilePicture(ctx, user.PublicIdentifier)
		if err != nil && !errors.Is(err, dao.ErrProfilePictureNotFound) {
			return nil, err
		}

		result[i] = &models.User{
			PublicIdentifier: user.PublicIdentifier,
			FirstName:        user.FirstName,
			LastName:         user.LastName,
			ProfilePicture:   profilePicture,
		}
	}

	return result, nil
}

func NewListUsersService(listUsersRepository dao.ListUsersRepository, getProfilePictureRepository dao.GetProfilePictureRepository) ListUsersService {
	return &listUsersServiceImpl{
		listUsersRepository:         listUsersRepository,
		getProfilePictureRepository: getProfilePictureRepository,
	}
}
