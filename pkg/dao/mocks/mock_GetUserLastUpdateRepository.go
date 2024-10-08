// Code generated by mockery v2.43.2. DO NOT EDIT.

package daomocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// MockGetUserLastUpdateRepository is an autogenerated mock type for the GetUserLastUpdateRepository type
type MockGetUserLastUpdateRepository struct {
	mock.Mock
}

type MockGetUserLastUpdateRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockGetUserLastUpdateRepository) EXPECT() *MockGetUserLastUpdateRepository_Expecter {
	return &MockGetUserLastUpdateRepository_Expecter{mock: &_m.Mock}
}

// GetUserLastUpdate provides a mock function with given fields: ctx, publicIdentifier
func (_m *MockGetUserLastUpdateRepository) GetUserLastUpdate(ctx context.Context, publicIdentifier string) (*time.Time, error) {
	ret := _m.Called(ctx, publicIdentifier)

	if len(ret) == 0 {
		panic("no return value specified for GetUserLastUpdate")
	}

	var r0 *time.Time
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*time.Time, error)); ok {
		return rf(ctx, publicIdentifier)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *time.Time); ok {
		r0 = rf(ctx, publicIdentifier)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*time.Time)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, publicIdentifier)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockGetUserLastUpdateRepository_GetUserLastUpdate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUserLastUpdate'
type MockGetUserLastUpdateRepository_GetUserLastUpdate_Call struct {
	*mock.Call
}

// GetUserLastUpdate is a helper method to define mock.On call
//   - ctx context.Context
//   - publicIdentifier string
func (_e *MockGetUserLastUpdateRepository_Expecter) GetUserLastUpdate(ctx interface{}, publicIdentifier interface{}) *MockGetUserLastUpdateRepository_GetUserLastUpdate_Call {
	return &MockGetUserLastUpdateRepository_GetUserLastUpdate_Call{Call: _e.mock.On("GetUserLastUpdate", ctx, publicIdentifier)}
}

func (_c *MockGetUserLastUpdateRepository_GetUserLastUpdate_Call) Run(run func(ctx context.Context, publicIdentifier string)) *MockGetUserLastUpdateRepository_GetUserLastUpdate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockGetUserLastUpdateRepository_GetUserLastUpdate_Call) Return(_a0 *time.Time, _a1 error) *MockGetUserLastUpdateRepository_GetUserLastUpdate_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockGetUserLastUpdateRepository_GetUserLastUpdate_Call) RunAndReturn(run func(context.Context, string) (*time.Time, error)) *MockGetUserLastUpdateRepository_GetUserLastUpdate_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockGetUserLastUpdateRepository creates a new instance of MockGetUserLastUpdateRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockGetUserLastUpdateRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockGetUserLastUpdateRepository {
	mock := &MockGetUserLastUpdateRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
