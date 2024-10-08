// Code generated by mockery v2.43.2. DO NOT EDIT.

package daomocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockGetProfilePictureRepository is an autogenerated mock type for the GetProfilePictureRepository type
type MockGetProfilePictureRepository struct {
	mock.Mock
}

type MockGetProfilePictureRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockGetProfilePictureRepository) EXPECT() *MockGetProfilePictureRepository_Expecter {
	return &MockGetProfilePictureRepository_Expecter{mock: &_m.Mock}
}

// GetProfilePicture provides a mock function with given fields: ctx, publicIdentifier
func (_m *MockGetProfilePictureRepository) GetProfilePicture(ctx context.Context, publicIdentifier string) (string, error) {
	ret := _m.Called(ctx, publicIdentifier)

	if len(ret) == 0 {
		panic("no return value specified for GetProfilePicture")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (string, error)); ok {
		return rf(ctx, publicIdentifier)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, publicIdentifier)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, publicIdentifier)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockGetProfilePictureRepository_GetProfilePicture_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetProfilePicture'
type MockGetProfilePictureRepository_GetProfilePicture_Call struct {
	*mock.Call
}

// GetProfilePicture is a helper method to define mock.On call
//   - ctx context.Context
//   - publicIdentifier string
func (_e *MockGetProfilePictureRepository_Expecter) GetProfilePicture(ctx interface{}, publicIdentifier interface{}) *MockGetProfilePictureRepository_GetProfilePicture_Call {
	return &MockGetProfilePictureRepository_GetProfilePicture_Call{Call: _e.mock.On("GetProfilePicture", ctx, publicIdentifier)}
}

func (_c *MockGetProfilePictureRepository_GetProfilePicture_Call) Run(run func(ctx context.Context, publicIdentifier string)) *MockGetProfilePictureRepository_GetProfilePicture_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockGetProfilePictureRepository_GetProfilePicture_Call) Return(_a0 string, _a1 error) *MockGetProfilePictureRepository_GetProfilePicture_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockGetProfilePictureRepository_GetProfilePicture_Call) RunAndReturn(run func(context.Context, string) (string, error)) *MockGetProfilePictureRepository_GetProfilePicture_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockGetProfilePictureRepository creates a new instance of MockGetProfilePictureRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockGetProfilePictureRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockGetProfilePictureRepository {
	mock := &MockGetProfilePictureRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
