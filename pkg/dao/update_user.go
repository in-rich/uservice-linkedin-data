package dao

import (
	"context"
	"github.com/in-rich/uservice-linkedin-data/pkg/entities"
	"github.com/samber/lo"
	"github.com/uptrace/bun"
	"time"
)

type UpdateUserData struct {
	FirstName string
	LastName  string
}

type UpdateUserRepository interface {
	UpdateUser(ctx context.Context, publicIdentifier string, data *UpdateUserData) (*entities.User, error)
}

type updateUserRepositoryImpl struct {
	db bun.IDB
}

func (r *updateUserRepositoryImpl) UpdateUser(ctx context.Context, publicIdentifier string, data *UpdateUserData) (*entities.User, error) {
	user := &entities.User{
		FirstName: data.FirstName,
		LastName:  data.LastName,
		UpdatedAt: lo.ToPtr(time.Now()),
	}

	res, err := r.db.NewUpdate().
		Model(user).
		Column("first_name", "last_name", "updated_at").
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
		return nil, ErrUserNotFound
	}

	return user, nil
}

func NewUpdateUserRepository(db bun.IDB) UpdateUserRepository {
	return &updateUserRepositoryImpl{
		db: db,
	}
}
