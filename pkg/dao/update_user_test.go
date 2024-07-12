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

var updateUserFixtures = []*entities.User{
	{
		ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
		PublicIdentifier: "public-identifier-1",
		FirstName:        "first-name-1",
		LastName:         "last-name-1",
	},
}

func TestUpdateUser(t *testing.T) {
	db := OpenDB()
	defer CloseDB(db)

	testData := []struct {
		name             string
		publicIdentifier string
		data             *dao.UpdateUserData
		expect           *entities.User
		expectErr        error
	}{
		{
			name:             "UpdateUser",
			publicIdentifier: "public-identifier-1",
			data: &dao.UpdateUserData{
				FirstName: "first-name-2",
				LastName:  "last-name-2",
			},
			expect: &entities.User{
				ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
				PublicIdentifier: "public-identifier-1",
				FirstName:        "first-name-2",
				LastName:         "last-name-2",
			},
		},
		{
			name:             "UserNotFound",
			publicIdentifier: "public-identifier-2",
			data: &dao.UpdateUserData{
				FirstName: "first-name-2",
				LastName:  "last-name-2",
			},
			expectErr: dao.ErrUserNotFound,
		},
	}

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			tx := BeginTX(db, updateUserFixtures)
			defer RollbackTX(tx)

			repo := dao.NewUpdateUserRepository(tx)
			user, err := repo.UpdateUser(context.TODO(), data.publicIdentifier, data.data)

			require.ErrorIs(t, err, data.expectErr)
			require.Empty(t, UsersCompare(user, data.expect))
		})
	}
}
