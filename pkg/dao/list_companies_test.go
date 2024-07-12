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

var listCompanyFixtures = []*entities.Company{
	{
		ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
		PublicIdentifier: "public-identifier-1",
		Name:             "company-1",
	},
	{
		ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000002")),
		PublicIdentifier: "public-identifier-2",
		Name:             "company-2",
	},
	{
		ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000003")),
		PublicIdentifier: "public-identifier-3",
		Name:             "company-3",
	},
}

func TestListCompanies(t *testing.T) {
	db := OpenDB()
	defer CloseDB(db)

	testData := []struct {
		name              string
		publicIdentifiers []string
		expect            []*entities.Company
	}{
		{
			name:              "ListCompanies",
			publicIdentifiers: []string{"public-identifier-1", "public-identifier-3", "public-identifier-4"},
			expect: []*entities.Company{
				{
					ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
					PublicIdentifier: "public-identifier-1",
					Name:             "company-1",
				},
				{
					ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000003")),
					PublicIdentifier: "public-identifier-3",
					Name:             "company-3",
				},
			},
		},
		{
			name:              "ListCompaniesEmpty",
			publicIdentifiers: []string{"public-identifier-4"},
			expect:            []*entities.Company{},
		},
	}

	stx := BeginTX(db, listCompanyFixtures)
	defer RollbackTX(stx)

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			tx := BeginTX[interface{}](stx, nil)
			defer RollbackTX(tx)

			repo := dao.NewListCompaniesRepository(tx)
			users, err := repo.ListCompanies(context.TODO(), data.publicIdentifiers)

			require.NoError(t, err)
			require.Empty(t, CompaniesCompareAll(users, data.expect))
		})
	}
}
