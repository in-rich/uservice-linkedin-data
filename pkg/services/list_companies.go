package services

import (
	"context"
	"errors"
	"github.com/in-rich/uservice-linkedin-data/pkg/dao"
	"github.com/in-rich/uservice-linkedin-data/pkg/models"
)

type ListCompaniesService interface {
	Exec(ctx context.Context, publicIdentifiers []string) ([]*models.Company, error)
}

type listCompaniesServiceImpl struct {
	listCompaniesRepository     dao.ListCompaniesRepository
	getProfilePictureRepository dao.GetProfilePictureRepository
}

func (s *listCompaniesServiceImpl) Exec(ctx context.Context, publicIdentifiers []string) ([]*models.Company, error) {
	companies, err := s.listCompaniesRepository.ListCompanies(ctx, publicIdentifiers)
	if err != nil {
		return nil, err
	}

	result := make([]*models.Company, len(companies))
	for i, company := range companies {
		profilePicture, err := s.getProfilePictureRepository.GetProfilePicture(ctx, company.PublicIdentifier)
		if err != nil && !errors.Is(err, dao.ErrProfilePictureNotFound) {
			return nil, err
		}

		result[i] = &models.Company{
			PublicIdentifier: company.PublicIdentifier,
			Name:             company.Name,
			ProfilePicture:   profilePicture,
		}
	}

	return result, nil
}

func NewListCompaniesService(listCompaniesRepository dao.ListCompaniesRepository, getProfilePictureRepository dao.GetProfilePictureRepository) ListCompaniesService {
	return &listCompaniesServiceImpl{
		listCompaniesRepository:     listCompaniesRepository,
		getProfilePictureRepository: getProfilePictureRepository,
	}
}
