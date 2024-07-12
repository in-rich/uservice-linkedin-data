package dao

import (
	"context"
	"github.com/in-rich/uservice-linkedin-data/pkg/entities"
	"github.com/uptrace/bun"
)

type ListUsersRepository interface {
	ListUsers(ctx context.Context, publicIdentifiers []string) ([]*entities.User, error)
}

type listUsersRepositoryImpl struct {
	db bun.IDB
}

func (r *listUsersRepositoryImpl) ListUsers(ctx context.Context, publicIdentifiers []string) ([]*entities.User, error) {
	users := make([]*entities.User, 0)

	err := r.db.NewSelect().Model(&users).Where("public_identifier IN (?)", bun.In(publicIdentifiers)).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func NewListUsersRepository(db bun.IDB) ListUsersRepository {
	return &listUsersRepositoryImpl{
		db: db,
	}
}
