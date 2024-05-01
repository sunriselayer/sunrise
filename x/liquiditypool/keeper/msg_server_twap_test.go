package keeper_test

import (
    "strconv"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

    keepertest "github.com/sunriselayer/sunrise-app/testutil/keeper"
    "github.com/sunriselayer/sunrise-app/x/liquiditypool/keeper"
    "github.com/sunriselayer/sunrise-app/x/liquiditypool/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestTwapMsgServerCreate(t *testing.T) {
	k, ctx := keepertest.LiquiditypoolKeeper(t)
	srv := keeper.NewMsgServerImpl(k)
	creator := "A"
	for i := 0; i < 5; i++ {
		expected := &types.MsgCreateTwap{Creator: creator,
		    Index: strconv.Itoa(i),
            
		}
		_, err := srv.CreateTwap(ctx, expected)
		require.NoError(t, err)
		rst, found := k.GetTwap(ctx,
		    expected.Index,
            
		)
		require.True(t, found)
		require.Equal(t, expected.Creator, rst.Creator)
	}
}

func TestTwapMsgServerUpdate(t *testing.T) {
	creator := "A"

	tests := []struct {
		desc    string
		request *types.MsgUpdateTwap
		err     error
	}{
		{
			desc:    "Completed",
			request: &types.MsgUpdateTwap{Creator: creator,
			    Index: strconv.Itoa(0),
                
			},
		},
		{
			desc:    "Unauthorized",
			request: &types.MsgUpdateTwap{Creator: "B",
			    Index: strconv.Itoa(0),
                
			},
			err:     sdkerrors.ErrUnauthorized,
		},
		{
			desc:    "KeyNotFound",
			request: &types.MsgUpdateTwap{Creator: creator,
			    Index: strconv.Itoa(100000),
                
			},
			err:     sdkerrors.ErrKeyNotFound,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := keepertest.LiquiditypoolKeeper(t)
			srv := keeper.NewMsgServerImpl(k)
			expected := &types.MsgCreateTwap{Creator: creator,
			    Index: strconv.Itoa(0),
                
			}
			_, err := srv.CreateTwap(ctx, expected)
			require.NoError(t, err)

			_, err = srv.UpdateTwap(ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				rst, found := k.GetTwap(ctx,
				    expected.Index,
                    
				)
				require.True(t, found)
				require.Equal(t, expected.Creator, rst.Creator)
			}
		})
	}
}

func TestTwapMsgServerDelete(t *testing.T) {
	creator := "A"

	tests := []struct {
		desc    string
		request *types.MsgDeleteTwap
		err     error
	}{
		{
			desc:    "Completed",
			request: &types.MsgDeleteTwap{Creator: creator,
			    Index: strconv.Itoa(0),
                
			},
		},
		{
			desc:    "Unauthorized",
			request: &types.MsgDeleteTwap{Creator: "B",
			    Index: strconv.Itoa(0),
                
			},
			err:     sdkerrors.ErrUnauthorized,
		},
		{
			desc:    "KeyNotFound",
			request: &types.MsgDeleteTwap{Creator: creator,
			    Index: strconv.Itoa(100000),
                
			},
			err:     sdkerrors.ErrKeyNotFound,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := keepertest.LiquiditypoolKeeper(t)
			srv := keeper.NewMsgServerImpl(k)

			_, err := srv.CreateTwap(ctx, &types.MsgCreateTwap{Creator: creator,
			    Index: strconv.Itoa(0),
                
			})
			require.NoError(t, err)
			_, err = srv.DeleteTwap(ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				_, found := k.GetTwap(ctx,
				    tc.request.Index,
                    
				)
				require.False(t, found)
			}
		})
	}
}
