package services

import (
	"context"
	"errors"
	"github.com/in-rich/uservice-linkedin-data/pkg/dao"
	"github.com/in-rich/uservice-linkedin-data/pkg/models"
)

type GetCompanyService interface {
	Exec(ctx context.Context, publicIdentifier string) (*models.Company, error)
}

type getCompanyServiceImpl struct {
	getCompanyRepository        dao.GetCompanyRepository
	getProfilePictureRepository dao.GetProfilePictureRepository
}

func (s *getCompanyServiceImpl) Exec(ctx context.Context, publicIdentifier string) (*models.Company, error) {
	company, err := s.getCompanyRepository.GetCompany(ctx, publicIdentifier)
	if err != nil {
		return nil, err
	}

	profilePicture, err := s.getProfilePictureRepository.GetProfilePicture(ctx, publicIdentifier)
	if err != nil && !errors.Is(err, dao.ErrProfilePictureNotFound) {
		return nil, err
	}

	return &models.Company{
		PublicIdentifier: company.PublicIdentifier,
		Name:             company.Name,
		ProfilePicture:   profilePicture,
	}, nil
}

func NewGetCompanyService(getCompanyRepository dao.GetCompanyRepository, getProfilePictureRepository dao.GetProfilePictureRepository) GetCompanyService {
	return &getCompanyServiceImpl{
		getCompanyRepository:        getCompanyRepository,
		getProfilePictureRepository: getProfilePictureRepository,
	}
}
