package dao

import (
	"context"
	"github.com/in-rich/uservice-linkedin-data/pkg/entities"
	"github.com/samber/lo"
	"github.com/uptrace/bun"
	"time"
)

type UpdateCompanyData struct {
	Name string
}

type UpdateCompanyRepository interface {
	UpdateCompany(ctx context.Context, publicIdentifier string, data *UpdateCompanyData) (*entities.Company, error)
}

type updateCompanyRepositoryImpl struct {
	db bun.IDB
}

func (r *updateCompanyRepositoryImpl) UpdateCompany(ctx context.Context, publicIdentifier string, data *UpdateCompanyData) (*entities.Company, error) {
	user := &entities.Company{
		Name:      data.Name,
		UpdatedAt: lo.ToPtr(time.Now()),
	}

	res, err := r.db.NewUpdate().
		Model(user).
		Column("name", "updated_at").
		Where("public_identifier = ?", publicIdentifier).
		Returning("*").
		Exec(ctx)

	if err != nil {
		return nil, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, ErrCompanyNotFound
	}

	return user, nil
}

func NewUpdateCompanyRepository(db bun.IDB) UpdateCompanyRepository {
	return &updateCompanyRepositoryImpl{
		db: db,
	}
}
