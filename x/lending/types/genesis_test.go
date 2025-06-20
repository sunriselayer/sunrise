package types_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/sunriselayer/sunrise/x/lending/types"
)

func TestGenesisState_Validate(t *testing.T) {
	tests := []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
		errMsg   string
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc: "valid genesis state with data",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				Markets: []types.Market{
					{
						Denom:             "usdc",
						TotalSupplied:     math.NewInt(1000000),
						TotalBorrowed:     math.NewInt(500000),
						GlobalRewardIndex: math.LegacyOneDec(),
						RiseDenom:         "riseusdc",
					},
				},
				UserPositions: []types.UserPosition{
					{
						UserAddress:     sdk.AccAddress("user1").String(),
						Denom:           "usdc",
						Amount:          math.NewInt(100000),
						LastRewardIndex: math.LegacyOneDec(),
					},
				},
				Borrows: []types.Borrow{
					{
						Id:                   1,
						Borrower:             sdk.AccAddress("borrower1").String(),
						Amount:               sdk.NewCoin("usdc", math.NewInt(50000)),
						CollateralPoolId:     1,
						CollateralPositionId: 100,
						BlockHeight:          1000,
					},
				},
				BorrowCount: 2,
			},
			valid: true,
		},
		{
			desc: "invalid params",
			genState: &types.GenesisState{
				Params: types.NewParams(
					math.LegacyNewDec(2), // LTV > 1
					math.LegacyNewDecWithPrec(85, 2),
					math.LegacyNewDecWithPrec(5, 2),
				),
			},
			valid:  false,
			errMsg: "invalid ltv ratio",
		},
		{
			desc: "duplicate market denom",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				Markets: []types.Market{
					{
						Denom:             "usdc",
						TotalSupplied:     math.NewInt(1000000),
						TotalBorrowed:     math.NewInt(500000),
						GlobalRewardIndex: math.LegacyOneDec(),
						RiseDenom:         "riseusdc",
					},
					{
						Denom:             "usdc",
						TotalSupplied:     math.NewInt(2000000),
						TotalBorrowed:     math.NewInt(1000000),
						GlobalRewardIndex: math.LegacyOneDec(),
						RiseDenom:         "riseusdc",
					},
				},
			},
			valid:  false,
			errMsg: "duplicate market denom: usdc",
		},
		{
			desc: "empty market denom",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				Markets: []types.Market{
					{
						Denom:             "",
						TotalSupplied:     math.NewInt(1000000),
						TotalBorrowed:     math.NewInt(500000),
						GlobalRewardIndex: math.LegacyOneDec(),
						RiseDenom:         "riseusdc",
					},
				},
			},
			valid:  false,
			errMsg: "market denom cannot be empty",
		},
		{
			desc: "negative total supplied",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				Markets: []types.Market{
					{
						Denom:             "usdc",
						TotalSupplied:     math.NewInt(-1000000),
						TotalBorrowed:     math.NewInt(500000),
						GlobalRewardIndex: math.LegacyOneDec(),
						RiseDenom:         "riseusdc",
					},
				},
			},
			valid:  false,
			errMsg: "total supplied cannot be negative",
		},
		{
			desc: "duplicate user position",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				Markets: []types.Market{
					{
						Denom:             "usdc",
						TotalSupplied:     math.NewInt(1000000),
						TotalBorrowed:     math.NewInt(500000),
						GlobalRewardIndex: math.LegacyOneDec(),
						RiseDenom:         "riseusdc",
					},
				},
				UserPositions: []types.UserPosition{
					{
						UserAddress:     sdk.AccAddress("user1").String(),
						Denom:           "usdc",
						Amount:          math.NewInt(100000),
						LastRewardIndex: math.LegacyOneDec(),
					},
					{
						UserAddress:     sdk.AccAddress("user1").String(),
						Denom:           "usdc",
						Amount:          math.NewInt(200000),
						LastRewardIndex: math.LegacyOneDec(),
					},
				},
			},
			valid:  false,
			errMsg: "duplicate user position",
		},
		{
			desc: "position references non-existent market",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				UserPositions: []types.UserPosition{
					{
						UserAddress:     sdk.AccAddress("user1").String(),
						Denom:           "usdc",
						Amount:          math.NewInt(100000),
						LastRewardIndex: math.LegacyOneDec(),
					},
				},
			},
			valid:  false,
			errMsg: "position references non-existent market: usdc",
		},
		{
			desc: "duplicate borrow id",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				Markets: []types.Market{
					{
						Denom:             "usdc",
						TotalSupplied:     math.NewInt(1000000),
						TotalBorrowed:     math.NewInt(500000),
						GlobalRewardIndex: math.LegacyOneDec(),
						RiseDenom:         "riseusdc",
					},
				},
				Borrows: []types.Borrow{
					{
						Id:                   1,
						Borrower:             sdk.AccAddress("borrower1").String(),
						Amount:               sdk.NewCoin("usdc", math.NewInt(50000)),
						CollateralPoolId:     1,
						CollateralPositionId: 100,
						BlockHeight:          1000,
					},
					{
						Id:                   1,
						Borrower:             sdk.AccAddress("borrower2").String(),
						Amount:               sdk.NewCoin("usdc", math.NewInt(60000)),
						CollateralPoolId:     2,
						CollateralPositionId: 200,
						BlockHeight:          1001,
					},
				},
				BorrowCount: 2,
			},
			valid:  false,
			errMsg: "duplicate borrow id: 1",
		},
		{
			desc: "borrow count less than max id",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				Markets: []types.Market{
					{
						Denom:             "usdc",
						TotalSupplied:     math.NewInt(1000000),
						TotalBorrowed:     math.NewInt(500000),
						GlobalRewardIndex: math.LegacyOneDec(),
						RiseDenom:         "riseusdc",
					},
				},
				Borrows: []types.Borrow{
					{
						Id:                   5,
						Borrower:             sdk.AccAddress("borrower1").String(),
						Amount:               sdk.NewCoin("usdc", math.NewInt(50000)),
						CollateralPoolId:     1,
						CollateralPositionId: 100,
						BlockHeight:          1000,
					},
				},
				BorrowCount: 3,
			},
			valid:  false,
			errMsg: "borrow count 3 is less than max borrow id 5",
		},
		{
			desc: "invalid borrow - zero collateral pool id",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				Markets: []types.Market{
					{
						Denom:             "usdc",
						TotalSupplied:     math.NewInt(1000000),
						TotalBorrowed:     math.NewInt(500000),
						GlobalRewardIndex: math.LegacyOneDec(),
						RiseDenom:         "riseusdc",
					},
				},
				Borrows: []types.Borrow{
					{
						Id:                   1,
						Borrower:             sdk.AccAddress("borrower1").String(),
						Amount:               sdk.NewCoin("usdc", math.NewInt(50000)),
						CollateralPoolId:     0,
						CollateralPositionId: 100,
						BlockHeight:          1000,
					},
				},
				BorrowCount: 2,
			},
			valid:  false,
			errMsg: "collateral pool id cannot be zero",
		},
		{
			desc: "borrow references non-existent market",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				Borrows: []types.Borrow{
					{
						Id:                   1,
						Borrower:             sdk.AccAddress("borrower1").String(),
						Amount:               sdk.NewCoin("atom", math.NewInt(50000)),
						CollateralPoolId:     1,
						CollateralPositionId: 100,
						BlockHeight:          1000,
					},
				},
				BorrowCount: 2,
			},
			valid:  false,
			errMsg: "borrow references non-existent market: atom",
		},
	}
	
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				if tc.errMsg != "" {
					require.Contains(t, err.Error(), tc.errMsg)
				}
			}
		})
	}
}
