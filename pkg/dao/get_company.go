package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/in-rich/uservice-linkedin-data/pkg/entities"
	"github.com/uptrace/bun"
)

type GetCompanyRepository interface {
	GetCompany(ctx context.Context, publicIdentifier string) (*entities.Company, error)
}

type getCompanyRepositoryImpl struct {
	db bun.IDB
}

func (r *getCompanyRepositoryImpl) GetCompany(ctx context.Context, publicIdentifier string) (*entities.Company, error) {
	company := new(entities.Company)

	err := r.db.NewSelect().Model(company).Where("public_identifier = ?", publicIdentifier).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCompanyNotFound
		}

		return nil, err
	}

	return company, nil
}

func NewGetCompanyRepository(db bun.IDB) GetCompanyRepository {
	return &getCompanyRepositoryImpl{
		db: db,
	}
}
