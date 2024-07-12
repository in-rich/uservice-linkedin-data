package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/in-rich/uservice-linkedin-data/pkg/entities"
	"github.com/uptrace/bun"
	"time"
)

type GetCompanyLastUpdateRepository interface {
	GetCompanyLastUpdate(ctx context.Context, publicIdentifier string) (*time.Time, error)
}

type getCompanyLastUpdateRepositoryImpl struct {
	db bun.IDB
}

func (r *getCompanyLastUpdateRepositoryImpl) GetCompanyLastUpdate(ctx context.Context, publicIdentifier string) (*time.Time, error) {
	company := new(entities.Company)

	err := r.db.NewSelect().Model(company).Where("public_identifier = ?", publicIdentifier).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCompanyNotFound
		}

		return nil, err
	}

	return company.UpdatedAt, nil
}

func NewGetCompanyLastUpdateRepository(db bun.IDB) GetCompanyLastUpdateRepository {
	return &getCompanyLastUpdateRepositoryImpl{
		db: db,
	}
}
