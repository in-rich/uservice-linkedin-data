package services

import (
	"context"
	"github.com/in-rich/uservice-linkedin-data/pkg/dao"
	"time"
)

type GetCompanyLastUpdateService interface {
	Exec(ctx context.Context, publicIdentifier string) (*time.Time, error)
}

type getCompanyLastUpdateServiceImpl struct {
	getCompanyLastUpdateRepository dao.GetCompanyLastUpdateRepository
}

func (s *getCompanyLastUpdateServiceImpl) Exec(ctx context.Context, publicIdentifier string) (*time.Time, error) {
	companyLastUpdate, err := s.getCompanyLastUpdateRepository.GetCompanyLastUpdate(ctx, publicIdentifier)
	if err != nil {
		return nil, err
	}

	return companyLastUpdate, nil
}

func NewGetCompanyLastUpdateService(getCompanyLastUpdateRepository dao.GetCompanyLastUpdateRepository) GetCompanyLastUpdateService {
	return &getCompanyLastUpdateServiceImpl{
		getCompanyLastUpdateRepository: getCompanyLastUpdateRepository,
	}
}
