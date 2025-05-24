package keeper_test

import (
	"testing"

	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

func TestBribeAllocationStore(t *testing.T) {
	// Create test accounts
	_, _, addr1 := testdata.KeyTestPubAddr()
	addr1Str := addr1.String()
	_, _, addr2 := testdata.KeyTestPubAddr()
	addr2Str := addr2.String()

	bech32Codec := addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())

	tests := []struct {
		name      string
		setup     func(fx *fixture, ctx sdk.Context)
		run       func(fx *fixture, ctx sdk.Context) error
		check     func(fx *fixture, ctx sdk.Context)
		expectErr bool
	}{
		{
			name: "set and get bribe allocation",
			setup: func(fx *fixture, ctx sdk.Context) {
				fx.mocks.AcctKeeper.EXPECT().AddressCodec().Return(bech32Codec).AnyTimes()
			},
			run: func(fx *fixture, ctx sdk.Context) error {
				allocation := types.BribeAllocation{
					Address: addr1Str,
					EpochId: 1,
					PoolId:  1,
					Weight:  "1.0",
				}
				return fx.keeper.SetBribeAllocation(ctx, allocation)
			},
			check: func(fx *fixture, ctx sdk.Context) {
				allocation, err := fx.keeper.GetBribeAllocation(ctx, addr1, 1, 1)
				require.NoError(t, err)
				require.Equal(t, addr1Str, allocation.Address)
				require.Equal(t, uint64(1), allocation.EpochId)
				require.Equal(t, uint64(1), allocation.PoolId)
				require.Equal(t, "1.0", allocation.Weight)
			},
			expectErr: false,
		},
		{
			name: "remove bribe allocation",
			setup: func(fx *fixture, ctx sdk.Context) {
				fx.mocks.AcctKeeper.EXPECT().AddressCodec().Return(bech32Codec).AnyTimes()
				allocation := types.BribeAllocation{
					Address: addr1Str,
					EpochId: 1,
					PoolId:  1,
					Weight:  "1.0",
				}
				err := fx.keeper.SetBribeAllocation(ctx, allocation)
				require.NoError(t, err)
			},
			run: func(fx *fixture, ctx sdk.Context) error {
				allocation := types.BribeAllocation{
					Address: addr1Str,
					EpochId: 1,
					PoolId:  1,
					Weight:  "1.0",
				}
				return fx.keeper.RemoveBribeAllocation(ctx, allocation)
			},
			check: func(fx *fixture, ctx sdk.Context) {
				_, err := fx.keeper.GetBribeAllocation(ctx, addr1, 1, 1)
				require.Error(t, err)
			},
			expectErr: false,
		},
		{
			name: "get all bribe allocations",
			setup: func(fx *fixture, ctx sdk.Context) {
				fx.mocks.AcctKeeper.EXPECT().AddressCodec().Return(bech32Codec).AnyTimes()
				allocations := []types.BribeAllocation{
					{
						Address: addr1Str,
						EpochId: 1,
						PoolId:  1,
						Weight:  "1.0",
					},
					{
						Address: addr2Str,
						EpochId: 1,
						PoolId:  2,
						Weight:  "0.5",
					},
				}
				for _, allocation := range allocations {
					err := fx.keeper.SetBribeAllocation(ctx, allocation)
					require.NoError(t, err)
				}
			},
			run: func(fx *fixture, ctx sdk.Context) error {
				return nil
			},
			check: func(fx *fixture, ctx sdk.Context) {
				allocations, err := fx.keeper.GetAllBribeAllocations(ctx)
				require.NoError(t, err)
				require.Len(t, allocations, 2)
				require.Contains(t, allocations, types.BribeAllocation{
					Address: addr1Str,
					EpochId: 1,
					PoolId:  1,
					Weight:  "1.0",
				})
				require.Contains(t, allocations, types.BribeAllocation{
					Address: addr2Str,
					EpochId: 1,
					PoolId:  2,
					Weight:  "0.5",
				})
			},
			expectErr: false,
		},
		{
			name: "get bribe allocations by address",
			setup: func(fx *fixture, ctx sdk.Context) {
				fx.mocks.AcctKeeper.EXPECT().AddressCodec().Return(bech32Codec).AnyTimes()
				allocations := []types.BribeAllocation{
					{
						Address: addr1Str,
						EpochId: 1,
						PoolId:  1,
						Weight:  "1.0",
					},
					{
						Address: addr1Str,
						EpochId: 2,
						PoolId:  2,
						Weight:  "0.5",
					},
				}
				for _, allocation := range allocations {
					err := fx.keeper.SetBribeAllocation(ctx, allocation)
					require.NoError(t, err)
				}
			},
			run: func(fx *fixture, ctx sdk.Context) error {
				return nil
			},
			check: func(fx *fixture, ctx sdk.Context) {
				allocations, err := fx.keeper.GetBribeAllocationsByAddress(ctx, addr1)
				require.NoError(t, err)
				require.Len(t, allocations, 2)
				require.Contains(t, allocations, types.BribeAllocation{
					Address: addr1Str,
					EpochId: 1,
					PoolId:  1,
					Weight:  "1.0",
				})
				require.Contains(t, allocations, types.BribeAllocation{
					Address: addr1Str,
					EpochId: 2,
					PoolId:  2,
					Weight:  "0.5",
				})
			},
			expectErr: false,
		},
		{
			name: "get bribe allocations by address and epoch id",
			setup: func(fx *fixture, ctx sdk.Context) {
				fx.mocks.AcctKeeper.EXPECT().AddressCodec().Return(bech32Codec).AnyTimes()
				allocations := []types.BribeAllocation{
					{
						Address: addr1Str,
						EpochId: 1,
						PoolId:  1,
						Weight:  "1.0",
					},
					{
						Address: addr1Str,
						EpochId: 2,
						PoolId:  2,
						Weight:  "0.5",
					},
				}
				for _, allocation := range allocations {
					err := fx.keeper.SetBribeAllocation(ctx, allocation)
					require.NoError(t, err)
				}
			},
			run: func(fx *fixture, ctx sdk.Context) error {
				return nil
			},
			check: func(fx *fixture, ctx sdk.Context) {
				allocations, err := fx.keeper.GetBribeAllocationsByAddressAndEpochId(ctx, addr1, 1)
				require.NoError(t, err)
				require.Len(t, allocations, 1)
				require.Contains(t, allocations, types.BribeAllocation{
					Address: addr1Str,
					EpochId: 1,
					PoolId:  1,
					Weight:  "1.0",
				})
			},
			expectErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fx := initFixture(t)
			ctx := sdk.UnwrapSDKContext(fx.ctx)

			if tc.setup != nil {
				tc.setup(fx, ctx)
			}

			err := tc.run(fx, ctx)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			if tc.check != nil {
				tc.check(fx, ctx)
			}
		})
	}
}
