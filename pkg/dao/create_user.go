package dao

import (
	"context"
	"errors"
	"github.com/in-rich/uservice-linkedin-data/pkg/entities"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/driver/pgdriver"
)

type CreateUserData struct {
	FirstName string
	LastName  string
}

type CreateUserRepository interface {
	CreateUser(ctx context.Context, publicIdentifier string, data *CreateUserData) (*entities.User, error)
}

type createUserRepositoryImpl struct {
	db bun.IDB
}

func (r *createUserRepositoryImpl) CreateUser(ctx context.Context, publicIdentifier string, data *CreateUserData) (*entities.User, error) {
	user := &entities.User{
		PublicIdentifier: publicIdentifier,
		FirstName:        data.FirstName,
		LastName:         data.LastName,
	}

	if _, err := r.db.NewInsert().Model(user).Returning("*").Exec(ctx); err != nil {
		var pgErr pgdriver.Error
		if errors.As(err, &pgErr) && pgErr.IntegrityViolation() {
			return nil, errors.Join(ErrUserAlreadyExists, err)
		}

		return nil, err
	}

	return user, nil
}

func NewCreateUserRepository(db bun.IDB) CreateUserRepository {
	return &createUserRepositoryImpl{
		db: db,
	}
}
