package testutil

import (
	"context"
	"reflect"

	"cosmossdk.io/core/address"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"go.uber.org/mock/gomock"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	liquiditypooltypes "github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

// MockAccountKeeper is a mock of AccountKeeper interface
type MockAccountKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockAccountKeeperMockRecorder
}

// MockAccountKeeperMockRecorder is the mock recorder for MockAccountKeeper
type MockAccountKeeperMockRecorder struct {
	mock *MockAccountKeeper
}

// NewMockAccountKeeper creates a new mock instance
func NewMockAccountKeeper(ctrl *gomock.Controller) *MockAccountKeeper {
	mock := &MockAccountKeeper{ctrl: ctrl}
	mock.recorder = &MockAccountKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAccountKeeper) EXPECT() *MockAccountKeeperMockRecorder {
	return m.recorder
}

// GetModuleAddress mocks base method
func (m *MockAccountKeeper) GetModuleAddress(moduleName string) sdk.AccAddress {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetModuleAddress", moduleName)
	ret0, _ := ret[0].(sdk.AccAddress)
	return ret0
}

// GetModuleAddress indicates an expected call of GetModuleAddress
func (mr *MockAccountKeeperMockRecorder) GetModuleAddress(moduleName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetModuleAddress", reflect.TypeOf((*MockAccountKeeper)(nil).GetModuleAddress), moduleName)
}

// AddressCodec mocks base method
func (m *MockAccountKeeper) AddressCodec() address.Codec {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddressCodec")
	ret0, _ := ret[0].(address.Codec)
	return ret0
}

// AddressCodec indicates an expected call of AddressCodec
func (mr *MockAccountKeeperMockRecorder) AddressCodec() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddressCodec", reflect.TypeOf((*MockAccountKeeper)(nil).AddressCodec))
}

// GetAccount mocks base method
func (m *MockAccountKeeper) GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccount", ctx, addr)
	ret0, _ := ret[0].(sdk.AccountI)
	return ret0
}

// GetAccount indicates an expected call of GetAccount
func (mr *MockAccountKeeperMockRecorder) GetAccount(ctx, addr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccount", reflect.TypeOf((*MockAccountKeeper)(nil).GetAccount), ctx, addr)
}

// GetModuleAccount mocks base method
func (m *MockAccountKeeper) GetModuleAccount(ctx context.Context, name string) sdk.ModuleAccountI {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetModuleAccount", ctx, name)
	ret0, _ := ret[0].(sdk.ModuleAccountI)
	return ret0
}

// GetModuleAccount indicates an expected call of GetModuleAccount
func (mr *MockAccountKeeperMockRecorder) GetModuleAccount(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetModuleAccount", reflect.TypeOf((*MockAccountKeeper)(nil).GetModuleAccount), ctx, name)
}

// SetModuleAccount mocks base method
func (m *MockAccountKeeper) SetModuleAccount(ctx context.Context, acc sdk.ModuleAccountI) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetModuleAccount", ctx, acc)
}

// SetModuleAccount indicates an expected call of SetModuleAccount
func (mr *MockAccountKeeperMockRecorder) SetModuleAccount(ctx, acc interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetModuleAccount", reflect.TypeOf((*MockAccountKeeper)(nil).SetModuleAccount), ctx, acc)
}

// MockBankKeeper is a mock of BankKeeper interface
type MockBankKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockBankKeeperMockRecorder
}

// MockBankKeeperMockRecorder is the mock recorder for MockBankKeeper
type MockBankKeeperMockRecorder struct {
	mock *MockBankKeeper
}

// NewMockBankKeeper creates a new mock instance
func NewMockBankKeeper(ctrl *gomock.Controller) *MockBankKeeper {
	mock := &MockBankKeeper{ctrl: ctrl}
	mock.recorder = &MockBankKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockBankKeeper) EXPECT() *MockBankKeeperMockRecorder {
	return m.recorder
}

// SendCoinsFromModuleToAccount mocks base method
func (m *MockBankKeeper) SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendCoinsFromModuleToAccount", ctx, senderModule, recipientAddr, amt)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendCoinsFromModuleToAccount indicates an expected call of SendCoinsFromModuleToAccount
func (mr *MockBankKeeperMockRecorder) SendCoinsFromModuleToAccount(ctx, senderModule, recipientAddr, amt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendCoinsFromModuleToAccount", reflect.TypeOf((*MockBankKeeper)(nil).SendCoinsFromModuleToAccount), ctx, senderModule, recipientAddr, amt)
}

// GetAllBalances mocks base method
func (m *MockBankKeeper) GetAllBalances(ctx context.Context, addr sdk.AccAddress) sdk.Coins {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllBalances", ctx, addr)
	ret0, _ := ret[0].(sdk.Coins)
	return ret0
}

// GetAllBalances indicates an expected call of GetAllBalances
func (mr *MockBankKeeperMockRecorder) GetAllBalances(ctx, addr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllBalances", reflect.TypeOf((*MockBankKeeper)(nil).GetAllBalances), ctx, addr)
}

// GetBalance mocks base method
func (m *MockBankKeeper) GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBalance", ctx, addr, denom)
	ret0, _ := ret[0].(sdk.Coin)
	return ret0
}

// GetBalance indicates an expected call of GetBalance
func (mr *MockBankKeeperMockRecorder) GetBalance(ctx, addr, denom interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBalance", reflect.TypeOf((*MockBankKeeper)(nil).GetBalance), ctx, addr, denom)
}

// IsSendEnabledCoins mocks base method (variadic)
func (m *MockBankKeeper) IsSendEnabledCoins(ctx context.Context, coins ...sdk.Coin) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range coins {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "IsSendEnabledCoins", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// IsSendEnabledCoins indicates an expected call of IsSendEnabledCoins
func (mr *MockBankKeeperMockRecorder) IsSendEnabledCoins(ctx interface{}, coins ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, coins...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsSendEnabledCoins", reflect.TypeOf((*MockBankKeeper)(nil).IsSendEnabledCoins), varargs...)
}

// SendCoinsFromAccountToModule mocks base method
func (m *MockBankKeeper) SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendCoinsFromAccountToModule", ctx, senderAddr, recipientModule, amt)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendCoinsFromAccountToModule indicates an expected call of SendCoinsFromAccountToModule
func (mr *MockBankKeeperMockRecorder) SendCoinsFromAccountToModule(ctx, senderAddr, recipientModule, amt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendCoinsFromAccountToModule", reflect.TypeOf((*MockBankKeeper)(nil).SendCoinsFromAccountToModule), ctx, senderAddr, recipientModule, amt)
}

// MintCoins mocks base method (for test usage)
func (m *MockBankKeeper) MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MintCoins", ctx, moduleName, amt)
	ret0, _ := ret[0].(error)
	return ret0
}

// MintCoins indicates an expected call of MintCoins
func (mr *MockBankKeeperMockRecorder) MintCoins(ctx, moduleName, amt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MintCoins", reflect.TypeOf((*MockBankKeeper)(nil).MintCoins), ctx, moduleName, amt)
}

// SendCoinsFromModuleToModule mocks base method
func (m *MockBankKeeper) SendCoinsFromModuleToModule(ctx context.Context, senderModule, recipientModule string, amt sdk.Coins) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendCoinsFromModuleToModule", ctx, senderModule, recipientModule, amt)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendCoinsFromModuleToModule indicates an expected call of SendCoinsFromModuleToModule
func (mr *MockBankKeeperMockRecorder) SendCoinsFromModuleToModule(ctx, senderModule, recipientModule, amt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendCoinsFromModuleToModule", reflect.TypeOf((*MockBankKeeper)(nil).SendCoinsFromModuleToModule), ctx, senderModule, recipientModule, amt)
}

// MockStakingKeeper is a mock of StakingKeeper interface
type MockStakingKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockStakingKeeperMockRecorder
}

// MockStakingKeeperMockRecorder is the mock recorder for MockStakingKeeper
type MockStakingKeeperMockRecorder struct {
	mock *MockStakingKeeper
}

// NewMockStakingKeeper creates a new mock instance
func NewMockStakingKeeper(ctrl *gomock.Controller) *MockStakingKeeper {
	mock := &MockStakingKeeper{ctrl: ctrl}
	mock.recorder = &MockStakingKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockStakingKeeper) EXPECT() *MockStakingKeeperMockRecorder {
	return m.recorder
}

// ValidatorAddressCodec mocks base method
func (m *MockStakingKeeper) ValidatorAddressCodec() address.Codec {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidatorAddressCodec")
	ret0, _ := ret[0].(address.Codec)
	return ret0
}

// ValidatorAddressCodec indicates an expected call of ValidatorAddressCodec
func (mr *MockStakingKeeperMockRecorder) ValidatorAddressCodec() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidatorAddressCodec", reflect.TypeOf((*MockStakingKeeper)(nil).ValidatorAddressCodec))
}

// BondDenom mocks base method
func (m *MockStakingKeeper) BondDenom(ctx context.Context) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BondDenom", ctx)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BondDenom indicates an expected call of BondDenom
func (mr *MockStakingKeeperMockRecorder) BondDenom(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BondDenom", reflect.TypeOf((*MockStakingKeeper)(nil).BondDenom), ctx)
}

// IterateBondedValidatorsByPower mocks base method (with error return)
func (m *MockStakingKeeper) IterateBondedValidatorsByPower(ctx context.Context, fn func(index int64, validator stakingtypes.ValidatorI) (stop bool)) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IterateBondedValidatorsByPower", ctx, fn)
	ret0, _ := ret[0].(error)
	return ret0
}

// IterateBondedValidatorsByPower indicates an expected call of IterateBondedValidatorsByPower
func (mr *MockStakingKeeperMockRecorder) IterateBondedValidatorsByPower(ctx, fn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IterateBondedValidatorsByPower", reflect.TypeOf((*MockStakingKeeper)(nil).IterateBondedValidatorsByPower), ctx, fn)
}

// IterateDelegations mocks base method (with error return)
func (m *MockStakingKeeper) IterateDelegations(ctx context.Context, delegator sdk.AccAddress, fn func(index int64, delegation stakingtypes.DelegationI) (stop bool)) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IterateDelegations", ctx, delegator, fn)
	ret0, _ := ret[0].(error)
	return ret0
}

// IterateDelegations indicates an expected call of IterateDelegations
func (mr *MockStakingKeeperMockRecorder) IterateDelegations(ctx, delegator, fn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IterateDelegations", reflect.TypeOf((*MockStakingKeeper)(nil).IterateDelegations), ctx, delegator, fn)
}

// TotalBondedTokens mocks base method
func (m *MockStakingKeeper) TotalBondedTokens(ctx context.Context) (math.Int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TotalBondedTokens", ctx)
	ret0, _ := ret[0].(math.Int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TotalBondedTokens indicates an expected call of TotalBondedTokens
func (mr *MockStakingKeeperMockRecorder) TotalBondedTokens(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TotalBondedTokens", reflect.TypeOf((*MockStakingKeeper)(nil).TotalBondedTokens), ctx)
}

// MockFeeKeeper is a mock of FeeKeeper interface
type MockFeeKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockFeeKeeperMockRecorder
}

// MockFeeKeeperMockRecorder is the mock recorder for MockFeeKeeper
type MockFeeKeeperMockRecorder struct {
	mock *MockFeeKeeper
}

// NewMockFeeKeeper creates a new mock instance
func NewMockFeeKeeper(ctrl *gomock.Controller) *MockFeeKeeper {
	mock := &MockFeeKeeper{ctrl: ctrl}
	mock.recorder = &MockFeeKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockFeeKeeper) EXPECT() *MockFeeKeeperMockRecorder {
	return m.recorder
}

// FeeDenom mocks base method
func (m *MockFeeKeeper) FeeDenom(ctx context.Context) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FeeDenom", ctx)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FeeDenom indicates an expected call of FeeDenom
func (mr *MockFeeKeeperMockRecorder) FeeDenom(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FeeDenom", reflect.TypeOf((*MockFeeKeeper)(nil).FeeDenom), ctx)
}

// MockTokenConverterKeeper is a mock of TokenConverterKeeper interface
type MockTokenConverterKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockTokenConverterKeeperMockRecorder
}

// MockTokenConverterKeeperMockRecorder is the mock recorder for MockTokenConverterKeeper
type MockTokenConverterKeeperMockRecorder struct {
	mock *MockTokenConverterKeeper
}

// NewMockTokenConverterKeeper creates a new mock instance
func NewMockTokenConverterKeeper(ctrl *gomock.Controller) *MockTokenConverterKeeper {
	mock := &MockTokenConverterKeeper{ctrl: ctrl}
	mock.recorder = &MockTokenConverterKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTokenConverterKeeper) EXPECT() *MockTokenConverterKeeperMockRecorder {
	return m.recorder
}

// ConvertReverse mocks base method
func (m *MockTokenConverterKeeper) ConvertReverse(ctx context.Context, amount math.Int, address sdk.AccAddress) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConvertReverse", ctx, amount, address)
	ret0, _ := ret[0].(error)
	return ret0
}

// ConvertReverse indicates an expected call of ConvertReverse
func (mr *MockTokenConverterKeeperMockRecorder) ConvertReverse(ctx, amount, address interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConvertReverse", reflect.TypeOf((*MockTokenConverterKeeper)(nil).ConvertReverse), ctx, amount, address)
}

// MockLiquidityPoolKeeper is a mock of LiquidityPoolKeeper interface
type MockLiquidityPoolKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockLiquidityPoolKeeperMockRecorder
}

// MockLiquidityPoolKeeperMockRecorder is the mock recorder for MockLiquidityPoolKeeper
type MockLiquidityPoolKeeperMockRecorder struct {
	mock *MockLiquidityPoolKeeper
}

// NewMockLiquidityPoolKeeper creates a new mock instance
func NewMockLiquidityPoolKeeper(ctrl *gomock.Controller) *MockLiquidityPoolKeeper {
	mock := &MockLiquidityPoolKeeper{ctrl: ctrl}
	mock.recorder = &MockLiquidityPoolKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLiquidityPoolKeeper) EXPECT() *MockLiquidityPoolKeeperMockRecorder {
	return m.recorder
}

// GetPool mocks base method
func (m *MockLiquidityPoolKeeper) GetPool(ctx context.Context, id uint64) (liquiditypooltypes.Pool, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPool", ctx, id)
	ret0, _ := ret[0].(liquiditypooltypes.Pool)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetPool indicates an expected call of GetPool
func (mr *MockLiquidityPoolKeeperMockRecorder) GetPool(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPool", reflect.TypeOf((*MockLiquidityPoolKeeper)(nil).GetPool), ctx, id)
}

// AllocateIncentive mocks base method
func (m *MockLiquidityPoolKeeper) AllocateIncentive(ctx sdk.Context, poolId uint64, sender sdk.AccAddress, incentiveCoins sdk.Coins) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllocateIncentive", ctx, poolId, sender, incentiveCoins)
	ret0, _ := ret[0].(error)
	return ret0
}

// AllocateIncentive indicates an expected call of AllocateIncentive
func (mr *MockLiquidityPoolKeeperMockRecorder) AllocateIncentive(ctx, poolId, sender, incentiveCoins interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllocateIncentive", reflect.TypeOf((*MockLiquidityPoolKeeper)(nil).AllocateIncentive), ctx, poolId, sender, incentiveCoins)
}
