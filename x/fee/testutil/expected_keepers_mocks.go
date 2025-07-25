// Code generated by MockGen. DO NOT EDIT.
// Source: x/fee/types/expected_keepers.go
//
// Generated by this command:
//
//	mockgen -source=x/fee/types/expected_keepers.go -destination=x/fee/testutil/expected_keepers_mocks.go -package=testutil
//

// Package testutil is a generated GoMock package.
package testutil

import (
	context "context"
	reflect "reflect"

	math "cosmossdk.io/math"
	types "github.com/cosmos/cosmos-sdk/types"
	types0 "github.com/sunriselayer/sunrise/x/liquiditypool/types"
	types1 "github.com/sunriselayer/sunrise/x/swap/types"
	gomock "go.uber.org/mock/gomock"
)

// MockAccountKeeper is a mock of AccountKeeper interface.
type MockAccountKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockAccountKeeperMockRecorder
	isgomock struct{}
}

// MockAccountKeeperMockRecorder is the mock recorder for MockAccountKeeper.
type MockAccountKeeperMockRecorder struct {
	mock *MockAccountKeeper
}

// NewMockAccountKeeper creates a new mock instance.
func NewMockAccountKeeper(ctrl *gomock.Controller) *MockAccountKeeper {
	mock := &MockAccountKeeper{ctrl: ctrl}
	mock.recorder = &MockAccountKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAccountKeeper) EXPECT() *MockAccountKeeperMockRecorder {
	return m.recorder
}

// GetAccount mocks base method.
func (m *MockAccountKeeper) GetAccount(arg0 context.Context, arg1 types.AccAddress) types.AccountI {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccount", arg0, arg1)
	ret0, _ := ret[0].(types.AccountI)
	return ret0
}

// GetAccount indicates an expected call of GetAccount.
func (mr *MockAccountKeeperMockRecorder) GetAccount(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccount", reflect.TypeOf((*MockAccountKeeper)(nil).GetAccount), arg0, arg1)
}

// GetModuleAddress mocks base method.
func (m *MockAccountKeeper) GetModuleAddress(name string) types.AccAddress {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetModuleAddress", name)
	ret0, _ := ret[0].(types.AccAddress)
	return ret0
}

// GetModuleAddress indicates an expected call of GetModuleAddress.
func (mr *MockAccountKeeperMockRecorder) GetModuleAddress(name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetModuleAddress", reflect.TypeOf((*MockAccountKeeper)(nil).GetModuleAddress), name)
}

// MockBankKeeper is a mock of BankKeeper interface.
type MockBankKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockBankKeeperMockRecorder
	isgomock struct{}
}

// MockBankKeeperMockRecorder is the mock recorder for MockBankKeeper.
type MockBankKeeperMockRecorder struct {
	mock *MockBankKeeper
}

// NewMockBankKeeper creates a new mock instance.
func NewMockBankKeeper(ctrl *gomock.Controller) *MockBankKeeper {
	mock := &MockBankKeeper{ctrl: ctrl}
	mock.recorder = &MockBankKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBankKeeper) EXPECT() *MockBankKeeperMockRecorder {
	return m.recorder
}

// BurnCoins mocks base method.
func (m *MockBankKeeper) BurnCoins(ctx context.Context, moduleName string, amt types.Coins) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BurnCoins", ctx, moduleName, amt)
	ret0, _ := ret[0].(error)
	return ret0
}

// BurnCoins indicates an expected call of BurnCoins.
func (mr *MockBankKeeperMockRecorder) BurnCoins(ctx, moduleName, amt any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BurnCoins", reflect.TypeOf((*MockBankKeeper)(nil).BurnCoins), ctx, moduleName, amt)
}

// SendCoinsFromModuleToModule mocks base method.
func (m *MockBankKeeper) SendCoinsFromModuleToModule(ctx context.Context, senderModule, recipientModule string, amt types.Coins) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendCoinsFromModuleToModule", ctx, senderModule, recipientModule, amt)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendCoinsFromModuleToModule indicates an expected call of SendCoinsFromModuleToModule.
func (mr *MockBankKeeperMockRecorder) SendCoinsFromModuleToModule(ctx, senderModule, recipientModule, amt any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendCoinsFromModuleToModule", reflect.TypeOf((*MockBankKeeper)(nil).SendCoinsFromModuleToModule), ctx, senderModule, recipientModule, amt)
}

// SpendableCoins mocks base method.
func (m *MockBankKeeper) SpendableCoins(arg0 context.Context, arg1 types.AccAddress) types.Coins {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SpendableCoins", arg0, arg1)
	ret0, _ := ret[0].(types.Coins)
	return ret0
}

// SpendableCoins indicates an expected call of SpendableCoins.
func (mr *MockBankKeeperMockRecorder) SpendableCoins(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SpendableCoins", reflect.TypeOf((*MockBankKeeper)(nil).SpendableCoins), arg0, arg1)
}

// MockLiquidityPoolKeeper is a mock of LiquidityPoolKeeper interface.
type MockLiquidityPoolKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockLiquidityPoolKeeperMockRecorder
	isgomock struct{}
}

// MockLiquidityPoolKeeperMockRecorder is the mock recorder for MockLiquidityPoolKeeper.
type MockLiquidityPoolKeeperMockRecorder struct {
	mock *MockLiquidityPoolKeeper
}

// NewMockLiquidityPoolKeeper creates a new mock instance.
func NewMockLiquidityPoolKeeper(ctrl *gomock.Controller) *MockLiquidityPoolKeeper {
	mock := &MockLiquidityPoolKeeper{ctrl: ctrl}
	mock.recorder = &MockLiquidityPoolKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLiquidityPoolKeeper) EXPECT() *MockLiquidityPoolKeeperMockRecorder {
	return m.recorder
}

// GetPool mocks base method.
func (m *MockLiquidityPoolKeeper) GetPool(ctx context.Context, id uint64) (types0.Pool, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPool", ctx, id)
	ret0, _ := ret[0].(types0.Pool)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetPool indicates an expected call of GetPool.
func (mr *MockLiquidityPoolKeeperMockRecorder) GetPool(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPool", reflect.TypeOf((*MockLiquidityPoolKeeper)(nil).GetPool), ctx, id)
}

// MockSwapKeeper is a mock of SwapKeeper interface.
type MockSwapKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockSwapKeeperMockRecorder
	isgomock struct{}
}

// MockSwapKeeperMockRecorder is the mock recorder for MockSwapKeeper.
type MockSwapKeeperMockRecorder struct {
	mock *MockSwapKeeper
}

// NewMockSwapKeeper creates a new mock instance.
func NewMockSwapKeeper(ctrl *gomock.Controller) *MockSwapKeeper {
	mock := &MockSwapKeeper{ctrl: ctrl}
	mock.recorder = &MockSwapKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSwapKeeper) EXPECT() *MockSwapKeeperMockRecorder {
	return m.recorder
}

// SwapExactAmountIn mocks base method.
func (m *MockSwapKeeper) SwapExactAmountIn(ctx types.Context, sender types.AccAddress, interfaceProvider string, route types1.Route, amountIn, minAmountOut math.Int) (types1.RouteResult, math.Int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SwapExactAmountIn", ctx, sender, interfaceProvider, route, amountIn, minAmountOut)
	ret0, _ := ret[0].(types1.RouteResult)
	ret1, _ := ret[1].(math.Int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// SwapExactAmountIn indicates an expected call of SwapExactAmountIn.
func (mr *MockSwapKeeperMockRecorder) SwapExactAmountIn(ctx, sender, interfaceProvider, route, amountIn, minAmountOut any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SwapExactAmountIn", reflect.TypeOf((*MockSwapKeeper)(nil).SwapExactAmountIn), ctx, sender, interfaceProvider, route, amountIn, minAmountOut)
}
