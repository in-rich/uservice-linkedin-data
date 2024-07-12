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

var getCompanyFixtures = []*entities.Company{
	{
		ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
		PublicIdentifier: "public-identifier-1",
		Name:             "company-1",
	},
}

func TestGetCompany(t *testing.T) {
	db := OpenDB()
	defer CloseDB(db)

	testData := []struct {
		name             string
		publicIdentifier string
		expect           *entities.Company
		expectErr        error
	}{
		{
			name:             "GetCompany",
			publicIdentifier: "public-identifier-1",
			expect: &entities.Company{
				ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
				PublicIdentifier: "public-identifier-1",
				Name:             "company-1",
			},
		},
		{
			name:             "CompanyNotFound",
			publicIdentifier: "public-identifier-2",
			expectErr:        dao.ErrCompanyNotFound,
		},
	}

	stx := BeginTX(db, getCompanyFixtures)
	defer RollbackTX(stx)

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			tx := BeginTX[interface{}](stx, nil)
			defer RollbackTX(tx)

			repo := dao.NewGetCompanyRepository(tx)
			company, err := repo.GetCompany(context.TODO(), data.publicIdentifier)

			require.ErrorIs(t, err, data.expectErr)
			require.Empty(t, CompaniesCompare(company, data.expect))
		})
	}
}
