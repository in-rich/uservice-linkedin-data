package dao

import (
	"context"
	"github.com/in-rich/uservice-linkedin-data/pkg/entities"
	"github.com/uptrace/bun"
)

type ListCompaniesRepository interface {
	ListCompanies(ctx context.Context, publicIdentifiers []string) ([]*entities.Company, error)
}

type listCompaniesRepositoryImpl struct {
	db bun.IDB
}

func (r *listCompaniesRepositoryImpl) ListCompanies(ctx context.Context, publicIdentifiers []string) ([]*entities.Company, error) {
	users := make([]*entities.Company, 0)

	err := r.db.NewSelect().Model(&users).Where("public_identifier IN (?)", bun.In(publicIdentifiers)).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func NewListCompaniesRepository(db bun.IDB) ListCompaniesRepository {
	return &listCompaniesRepositoryImpl{
		db: db,
	}
}
