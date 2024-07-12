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

var listUserFixtures = []*entities.User{
	{
		ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
		PublicIdentifier: "public-identifier-1",
		FirstName:        "first-name-1",
		LastName:         "last-name-1",
	},
	{
		ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000002")),
		PublicIdentifier: "public-identifier-2",
		FirstName:        "first-name-2",
		LastName:         "last-name-2",
	},
	{
		ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000003")),
		PublicIdentifier: "public-identifier-3",
		FirstName:        "first-name-3",
		LastName:         "last-name-3",
	},
}

func TestListUsers(t *testing.T) {
	db := OpenDB()
	defer CloseDB(db)

	testData := []struct {
		name              string
		publicIdentifiers []string
		expect            []*entities.User
	}{
		{
			name:              "ListUsers",
			publicIdentifiers: []string{"public-identifier-1", "public-identifier-3", "public-identifier-4"},
			expect: []*entities.User{
				{
					ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
					PublicIdentifier: "public-identifier-1",
					FirstName:        "first-name-1",
					LastName:         "last-name-1",
				},
				{
					ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000003")),
					PublicIdentifier: "public-identifier-3",
					FirstName:        "first-name-3",
					LastName:         "last-name-3",
				},
			},
		},
		{
			name:              "ListUsersEmpty",
			publicIdentifiers: []string{"public-identifier-4"},
			expect:            []*entities.User{},
		},
	}

	stx := BeginTX(db, listUserFixtures)
	defer RollbackTX(stx)

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			tx := BeginTX[interface{}](stx, nil)
			defer RollbackTX(tx)

			repo := dao.NewListUsersRepository(tx)
			users, err := repo.ListUsers(context.TODO(), data.publicIdentifiers)

			require.NoError(t, err)
			require.Empty(t, UsersCompareAll(users, data.expect))
		})
	}
}
