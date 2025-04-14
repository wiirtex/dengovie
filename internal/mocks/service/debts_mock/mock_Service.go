// Code generated by mockery. DO NOT EDIT.

package types

import (
	context "context"
	types "dengovie/internal/service/debts/types"

	mock "github.com/stretchr/testify/mock"
)

// MockService is an autogenerated mock type for the Service type
type MockService struct {
	mock.Mock
}

type MockService_Expecter struct {
	mock *mock.Mock
}

func (_m *MockService) EXPECT() *MockService_Expecter {
	return &MockService_Expecter{mock: &_m.Mock}
}

// PayDebt provides a mock function with given fields: ctx, input
func (_m *MockService) PayDebt(ctx context.Context, input types.PayDebtInput) error {
	ret := _m.Called(ctx, input)

	if len(ret) == 0 {
		panic("no return value specified for PayDebt")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, types.PayDebtInput) error); ok {
		r0 = rf(ctx, input)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockService_PayDebt_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PayDebt'
type MockService_PayDebt_Call struct {
	*mock.Call
}

// PayDebt is a helper method to define mock.On call
//   - ctx context.Context
//   - input types.PayDebtInput
func (_e *MockService_Expecter) PayDebt(ctx interface{}, input interface{}) *MockService_PayDebt_Call {
	return &MockService_PayDebt_Call{Call: _e.mock.On("PayDebt", ctx, input)}
}

func (_c *MockService_PayDebt_Call) Run(run func(ctx context.Context, input types.PayDebtInput)) *MockService_PayDebt_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(types.PayDebtInput))
	})
	return _c
}

func (_c *MockService_PayDebt_Call) Return(_a0 error) *MockService_PayDebt_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockService_PayDebt_Call) RunAndReturn(run func(context.Context, types.PayDebtInput) error) *MockService_PayDebt_Call {
	_c.Call.Return(run)
	return _c
}

// ShareDebt provides a mock function with given fields: ctx, input
func (_m *MockService) ShareDebt(ctx context.Context, input types.ShareDebtInput) error {
	ret := _m.Called(ctx, input)

	if len(ret) == 0 {
		panic("no return value specified for ShareDebt")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, types.ShareDebtInput) error); ok {
		r0 = rf(ctx, input)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockService_ShareDebt_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ShareDebt'
type MockService_ShareDebt_Call struct {
	*mock.Call
}

// ShareDebt is a helper method to define mock.On call
//   - ctx context.Context
//   - input types.ShareDebtInput
func (_e *MockService_Expecter) ShareDebt(ctx interface{}, input interface{}) *MockService_ShareDebt_Call {
	return &MockService_ShareDebt_Call{Call: _e.mock.On("ShareDebt", ctx, input)}
}

func (_c *MockService_ShareDebt_Call) Run(run func(ctx context.Context, input types.ShareDebtInput)) *MockService_ShareDebt_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(types.ShareDebtInput))
	})
	return _c
}

func (_c *MockService_ShareDebt_Call) Return(_a0 error) *MockService_ShareDebt_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockService_ShareDebt_Call) RunAndReturn(run func(context.Context, types.ShareDebtInput) error) *MockService_ShareDebt_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockService creates a new instance of MockService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockService {
	mock := &MockService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
