// Code generated by mockery v2.43.2. DO NOT EDIT.

package daomocks

import (
	context "context"

	dao "github.com/in-rich/uservice-linkedin-data/pkg/dao"
	entities "github.com/in-rich/uservice-linkedin-data/pkg/entities"

	mock "github.com/stretchr/testify/mock"
)

// MockCreateCompanyRepository is an autogenerated mock type for the CreateCompanyRepository type
type MockCreateCompanyRepository struct {
	mock.Mock
}

type MockCreateCompanyRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockCreateCompanyRepository) EXPECT() *MockCreateCompanyRepository_Expecter {
	return &MockCreateCompanyRepository_Expecter{mock: &_m.Mock}
}

// CreateCompany provides a mock function with given fields: ctx, publicIdentifier, data
func (_m *MockCreateCompanyRepository) CreateCompany(ctx context.Context, publicIdentifier string, data *dao.CreateCompanyData) (*entities.Company, error) {
	ret := _m.Called(ctx, publicIdentifier, data)

	if len(ret) == 0 {
		panic("no return value specified for CreateCompany")
	}

	var r0 *entities.Company
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *dao.CreateCompanyData) (*entities.Company, error)); ok {
		return rf(ctx, publicIdentifier, data)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, *dao.CreateCompanyData) *entities.Company); ok {
		r0 = rf(ctx, publicIdentifier, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entities.Company)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, *dao.CreateCompanyData) error); ok {
		r1 = rf(ctx, publicIdentifier, data)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCreateCompanyRepository_CreateCompany_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateCompany'
type MockCreateCompanyRepository_CreateCompany_Call struct {
	*mock.Call
}

// CreateCompany is a helper method to define mock.On call
//   - ctx context.Context
//   - publicIdentifier string
//   - data *dao.CreateCompanyData
func (_e *MockCreateCompanyRepository_Expecter) CreateCompany(ctx interface{}, publicIdentifier interface{}, data interface{}) *MockCreateCompanyRepository_CreateCompany_Call {
	return &MockCreateCompanyRepository_CreateCompany_Call{Call: _e.mock.On("CreateCompany", ctx, publicIdentifier, data)}
}

func (_c *MockCreateCompanyRepository_CreateCompany_Call) Run(run func(ctx context.Context, publicIdentifier string, data *dao.CreateCompanyData)) *MockCreateCompanyRepository_CreateCompany_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(*dao.CreateCompanyData))
	})
	return _c
}

func (_c *MockCreateCompanyRepository_CreateCompany_Call) Return(_a0 *entities.Company, _a1 error) *MockCreateCompanyRepository_CreateCompany_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCreateCompanyRepository_CreateCompany_Call) RunAndReturn(run func(context.Context, string, *dao.CreateCompanyData) (*entities.Company, error)) *MockCreateCompanyRepository_CreateCompany_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockCreateCompanyRepository creates a new instance of MockCreateCompanyRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockCreateCompanyRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockCreateCompanyRepository {
	mock := &MockCreateCompanyRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
