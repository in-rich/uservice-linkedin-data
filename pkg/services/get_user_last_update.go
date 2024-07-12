package services

import (
	"context"
	"github.com/in-rich/uservice-linkedin-data/pkg/dao"
	"time"
)

type GetUserLastUpdateService interface {
	Exec(ctx context.Context, publicIdentifier string) (*time.Time, error)
}

type getUserLastUpdateServiceImpl struct {
	getUserLastUpdateRepository dao.GetUserLastUpdateRepository
}

func (s *getUserLastUpdateServiceImpl) Exec(ctx context.Context, publicIdentifier string) (*time.Time, error) {
	userLastUpdate, err := s.getUserLastUpdateRepository.GetUserLastUpdate(ctx, publicIdentifier)
	if err != nil {
		return nil, err
	}

	return userLastUpdate, nil
}

func NewGetUserLastUpdateService(getUserLastUpdateRepository dao.GetUserLastUpdateRepository) GetUserLastUpdateService {
	return &getUserLastUpdateServiceImpl{
		getUserLastUpdateRepository: getUserLastUpdateRepository,
	}
}
