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

var createCompanyFixtures = []*entities.Company{
	{
		ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
		PublicIdentifier: "public-identifier-1",
		Name:             "company-1",
	},
}

func TestCreateCompany(t *testing.T) {
	db := OpenDB()
	defer CloseDB(db)

	testData := []struct {
		name             string
		publicIdentifier string
		data             *dao.CreateCompanyData
		expect           *entities.Company
		expectErr        error
	}{
		{
			name:             "CreateCompany",
			publicIdentifier: "public-identifier-2",
			data: &dao.CreateCompanyData{
				Name: "company-2",
			},
			expect: &entities.Company{
				PublicIdentifier: "public-identifier-2",
				Name:             "company-2",
			},
		},
		{
			name:             "CompanyAlreadyExists",
			publicIdentifier: "public-identifier-1",
			data: &dao.CreateCompanyData{
				Name: "company-1",
			},
			expectErr: dao.ErrCompanyAlreadyExists,
		},
	}

	stx := BeginTX(db, createCompanyFixtures)
	defer RollbackTX(stx)

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			tx := BeginTX[interface{}](stx, nil)
			defer RollbackTX(tx)

			repo := dao.NewCreateCompanyRepository(tx)
			company, err := repo.CreateCompany(context.TODO(), data.publicIdentifier, data.data)

			if company != nil {
				// Since ID is random, nullify it for comparison.
				company.ID = nil
			}

			require.ErrorIs(t, err, data.expectErr)
			require.Empty(t, CompaniesCompare(company, data.expect))
		})
	}
}
