package services

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/in-rich/uservice-linkedin-data/pkg/dao"
	"github.com/in-rich/uservice-linkedin-data/pkg/models"
)

var (
	ErrInvalidUpsertCompany = errors.New("invalid upsert company data")
)

type UpsertCompanyService interface {
	Exec(ctx context.Context, publicIdentifier string, data *models.UpsertCompany) (*models.Company, error)
}

type upsertCompanyServiceImpl struct {
	createCompanyRepository        dao.CreateCompanyRepository
	updateCompanyRepository        dao.UpdateCompanyRepository
	upsertProfilePictureRepository dao.UpsertProfilePictureRepository
}

func (s *upsertCompanyServiceImpl) upsertProfilePicture(ctx context.Context, publicIdentifier string, data *models.UpsertCompany) (string, error) {
	if data.ProfilePicture == "" {
		return "", nil
	}

	return s.upsertProfilePictureRepository.UpsertProfilePicture(ctx, publicIdentifier, data.ProfilePicture)
}

func (s *upsertCompanyServiceImpl) Exec(ctx context.Context, publicIdentifier string, data *models.UpsertCompany) (*models.Company, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(data); err != nil {
		return nil, errors.Join(ErrInvalidUpsertCompany, err)
	}

	company, err := s.createCompanyRepository.CreateCompany(ctx, publicIdentifier, &dao.CreateCompanyData{
		Name: data.Name,
	})
	// Company was successfully created.
	if err == nil {
		profilePicture, err := s.upsertProfilePicture(ctx, publicIdentifier, data)
		if err != nil {
			return nil, err
		}

		return &models.Company{
			PublicIdentifier: company.PublicIdentifier,
			Name:             company.Name,
			ProfilePicture:   profilePicture,
		}, nil
	}

	if !errors.Is(err, dao.ErrCompanyAlreadyExists) {
		return nil, err
	}

	// Company already existed. Update it.
	company, err = s.updateCompanyRepository.UpdateCompany(ctx, publicIdentifier, &dao.UpdateCompanyData{
		Name: data.Name,
	})
	if err != nil {
		return nil, err
	}

	profilePicture, err := s.upsertProfilePicture(ctx, publicIdentifier, data)
	if err != nil {
		return nil, err
	}

	return &models.Company{
		PublicIdentifier: company.PublicIdentifier,
		Name:             company.Name,
		ProfilePicture:   profilePicture,
	}, nil
}

func NewUpsertCompanyService(
	createCompanyRepository dao.CreateCompanyRepository,
	updateCompanyRepository dao.UpdateCompanyRepository,
	upsertProfilePictureRepository dao.UpsertProfilePictureRepository,
) UpsertCompanyService {
	return &upsertCompanyServiceImpl{
		createCompanyRepository:        createCompanyRepository,
		updateCompanyRepository:        updateCompanyRepository,
		upsertProfilePictureRepository: upsertProfilePictureRepository,
	}
}
