package dao

import (
	"context"
	"errors"
	"github.com/in-rich/uservice-linkedin-data/pkg/entities"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/driver/pgdriver"
)

type CreateCompanyData struct {
	Name string
}

type CreateCompanyRepository interface {
	CreateCompany(ctx context.Context, publicIdentifier string, data *CreateCompanyData) (*entities.Company, error)
}

type createCompanyRepositoryImpl struct {
	db bun.IDB
}

func (r *createCompanyRepositoryImpl) CreateCompany(ctx context.Context, publicIdentifier string, data *CreateCompanyData) (*entities.Company, error) {
	user := &entities.Company{
		PublicIdentifier: publicIdentifier,
		Name:             data.Name,
	}

	if _, err := r.db.NewInsert().Model(user).Returning("*").Exec(ctx); err != nil {
		var pgErr pgdriver.Error
		if errors.As(err, &pgErr) && pgErr.IntegrityViolation() {
			return nil, ErrCompanyAlreadyExists
		}

		return nil, err
	}

	return user, nil
}

func NewCreateCompanyRepository(db bun.IDB) CreateCompanyRepository {
	return &createCompanyRepositoryImpl{
		db: db,
	}
}
