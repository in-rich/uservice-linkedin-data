package dao_test

import (
	"context"
	"github.com/google/uuid"
	"github.com/in-rich/uservice-linkedin-data/pkg/dao"
	"github.com/in-rich/uservice-linkedin-data/pkg/entities"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"testing"
)

var createUserFixtures = []*entities.User{
	{
		ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
		PublicIdentifier: "public-identifier-1",
		FirstName:        "first-name-1",
		LastName:         "last-name-1",
	},
}

func TestCreateUser(t *testing.T) {
	db := OpenDB()
	defer CloseDB(db)

	testData := []struct {
		name             string
		publicIdentifier string
		data             *dao.CreateUserData
		expect           *entities.User
		expectErr        error
	}{
		{
			name:             "CreateUser",
			publicIdentifier: "public-identifier-2",
			data: &dao.CreateUserData{
				FirstName: "first-name-2",
				LastName:  "last-name-2",
			},
			expect: &entities.User{
				PublicIdentifier: "public-identifier-2",
				FirstName:        "first-name-2",
				LastName:         "last-name-2",
			},
		},
		{
			name:             "UserAlreadyExists",
			publicIdentifier: "public-identifier-1",
			data: &dao.CreateUserData{
				FirstName: "first-name-1",
				LastName:  "last-name-1",
			},
			expectErr: dao.ErrUserAlreadyExists,
		},
	}

	stx := BeginTX(db, createUserFixtures)
	defer RollbackTX(stx)

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			tx := BeginTX[interface{}](stx, nil)
			defer RollbackTX(tx)

			repo := dao.NewCreateUserRepository(tx)
			user, err := repo.CreateUser(context.TODO(), data.publicIdentifier, data.data)

			if user != nil {
				// Since ID is random, nullify it for comparison.
				user.ID = nil
			}

			require.ErrorIs(t, err, data.expectErr)
			require.Empty(t, UsersCompare(user, data.expect))
		})
	}
}
