package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/in-rich/uservice-linkedin-data/pkg/entities"
	"github.com/uptrace/bun"
	"time"
)

type GetUserLastUpdateRepository interface {
	GetUserLastUpdate(ctx context.Context, publicIdentifier string) (*time.Time, error)
}

type getUserLastUpdateRepositoryImpl struct {
	db bun.IDB
}

func (r *getUserLastUpdateRepositoryImpl) GetUserLastUpdate(ctx context.Context, publicIdentifier string) (*time.Time, error) {
	user := new(entities.User)

	err := r.db.NewSelect().Model(user).Where("public_identifier = ?", publicIdentifier).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return user.UpdatedAt, nil
}

func NewGetUserLastUpdateRepository(db bun.IDB) GetUserLastUpdateRepository {
	return &getUserLastUpdateRepositoryImpl{
		db: db,
	}
}
