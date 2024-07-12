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

var updateCompanyFixtures = []*entities.Company{
	{
		ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
		PublicIdentifier: "public-identifier-1",
		Name:             "company-1",
	},
}

func TestUpdateCompany(t *testing.T) {
	db := OpenDB()
	defer CloseDB(db)

	testData := []struct {
		name             string
		publicIdentifier string
		data             *dao.UpdateCompanyData
		expect           *entities.Company
		expectErr        error
	}{
		{
			name:             "UpdateCompany",
			publicIdentifier: "public-identifier-1",
			data: &dao.UpdateCompanyData{
				Name: "company-2",
			},
			expect: &entities.Company{
				ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
				PublicIdentifier: "public-identifier-1",
				Name:             "company-2",
			},
		},
		{
			name:             "CompanyNotFound",
			publicIdentifier: "public-identifier-2",
			data: &dao.UpdateCompanyData{
				Name: "company-2",
			},
			expectErr: dao.ErrCompanyNotFound,
		},
	}

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			tx := BeginTX(db, updateCompanyFixtures)
			defer RollbackTX(tx)

			repo := dao.NewUpdateCompanyRepository(tx)
			company, err := repo.UpdateCompany(context.TODO(), data.publicIdentifier, data.data)

			require.ErrorIs(t, err, data.expectErr)
			require.Empty(t, CompaniesCompare(data.expect, company))
		})
	}
}
