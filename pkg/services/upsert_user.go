package services

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/in-rich/uservice-linkedin-data/pkg/dao"
	"github.com/in-rich/uservice-linkedin-data/pkg/models"
)

var (
	ErrInvalidUpsertUser = errors.New("invalid upsert user data")
)

type UpsertUserService interface {
	Exec(ctx context.Context, publicIdentifier string, data *models.UpsertUser) (*models.User, error)
}

type upsertUserServiceImpl struct {
	createUserRepository           dao.CreateUserRepository
	updateUserRepository           dao.UpdateUserRepository
	upsertProfilePictureRepository dao.UpsertProfilePictureRepository
}

func (s *upsertUserServiceImpl) upsertProfilePicture(ctx context.Context, publicIdentifier string, data *models.UpsertUser) (string, error) {
	if data.ProfilePicture == "" {
		return "", nil
	}

	return s.upsertProfilePictureRepository.UpsertProfilePicture(ctx, publicIdentifier, data.ProfilePicture)
}

func (s *upsertUserServiceImpl) Exec(ctx context.Context, publicIdentifier string, data *models.UpsertUser) (*models.User, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(data); err != nil {
		return nil, errors.Join(ErrInvalidUpsertUser, err)
	}

	user, err := s.createUserRepository.CreateUser(ctx, publicIdentifier, &dao.CreateUserData{
		FirstName: data.FirstName,
		LastName:  data.LastName,
	})
	// User was successfully created.
	if err == nil {
		profilePicture, err := s.upsertProfilePicture(ctx, publicIdentifier, data)
		if err != nil {
			return nil, err
		}

		return &models.User{
			PublicIdentifier: user.PublicIdentifier,
			FirstName:        user.FirstName,
			LastName:         user.LastName,
			ProfilePicture:   profilePicture,
		}, nil
	}

	if !errors.Is(err, dao.ErrUserAlreadyExists) {
		return nil, err
	}

	// User already existed. Update it.
	user, err = s.updateUserRepository.UpdateUser(ctx, publicIdentifier, &dao.UpdateUserData{
		FirstName: data.FirstName,
		LastName:  data.LastName,
	})
	if err != nil {
		return nil, err
	}

	profilePicture, err := s.upsertProfilePicture(ctx, publicIdentifier, data)
	if err != nil {
		return nil, err
	}

	return &models.User{
		PublicIdentifier: user.PublicIdentifier,
		FirstName:        user.FirstName,
		LastName:         user.LastName,
		ProfilePicture:   profilePicture,
	}, nil
}

func NewUpsertUserService(
	createUserRepository dao.CreateUserRepository,
	updateUserRepository dao.UpdateUserRepository,
	upsertProfilePictureRepository dao.UpsertProfilePictureRepository,
) UpsertUserService {
	return &upsertUserServiceImpl{
		createUserRepository:           createUserRepository,
		updateUserRepository:           updateUserRepository,
		upsertProfilePictureRepository: upsertProfilePictureRepository,
	}
}
