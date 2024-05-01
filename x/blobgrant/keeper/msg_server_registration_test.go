package keeper_test

import (
	"strconv"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	keepertest "github.com/sunriselayer/sunrise-app/testutil/keeper"
	"github.com/sunriselayer/sunrise-app/x/blobgrant/keeper"
	"github.com/sunriselayer/sunrise-app/x/blobgrant/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestRegistrationMsgServerCreate(t *testing.T) {
	k, ctx := keepertest.GrantKeeper(t)
	srv := keeper.NewMsgServerImpl(k)
	creator := "A"
	for i := 0; i < 5; i++ {
		expected := &types.MsgCreateRegistration{Creator: creator,
			Address: strconv.Itoa(i),
		}
		_, err := srv.CreateRegistration(ctx, expected)
		require.NoError(t, err)
		rst, found := k.GetRegistration(ctx,
			expected.Address,
		)
		require.True(t, found)
		require.Equal(t, expected.Creator, rst.Creator)
	}
}

func TestRegistrationMsgServerUpdate(t *testing.T) {
	creator := "A"

	tests := []struct {
		desc    string
		request *types.MsgUpdateRegistration
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgUpdateRegistration{Creator: creator,
				Address: strconv.Itoa(0),
			},
		},
		{
			desc: "Unauthorized",
			request: &types.MsgUpdateRegistration{Creator: "B",
				Address: strconv.Itoa(0),
			},
			err: sdkerrors.ErrUnauthorized,
		},
		{
			desc: "KeyNotFound",
			request: &types.MsgUpdateRegistration{Creator: creator,
				Address: strconv.Itoa(100000),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := keepertest.GrantKeeper(t)
			srv := keeper.NewMsgServerImpl(k)
			expected := &types.MsgCreateRegistration{Creator: creator,
				Address: strconv.Itoa(0),
			}
			_, err := srv.CreateRegistration(ctx, expected)
			require.NoError(t, err)

			_, err = srv.UpdateRegistration(ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				rst, found := k.GetRegistration(ctx,
					expected.Address,
				)
				require.True(t, found)
				require.Equal(t, expected.Creator, rst.Creator)
			}
		})
	}
}

func TestRegistrationMsgServerDelete(t *testing.T) {
	creator := "A"

	tests := []struct {
		desc    string
		request *types.MsgDeleteRegistration
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgDeleteRegistration{Creator: creator,
				Address: strconv.Itoa(0),
			},
		},
		{
			desc: "Unauthorized",
			request: &types.MsgDeleteRegistration{Creator: "B",
				Address: strconv.Itoa(0),
			},
			err: sdkerrors.ErrUnauthorized,
		},
		{
			desc: "KeyNotFound",
			request: &types.MsgDeleteRegistration{Creator: creator,
				Address: strconv.Itoa(100000),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := keepertest.GrantKeeper(t)
			srv := keeper.NewMsgServerImpl(k)

			_, err := srv.CreateRegistration(ctx, &types.MsgCreateRegistration{Creator: creator,
				Address: strconv.Itoa(0),
			})
			require.NoError(t, err)
			_, err = srv.DeleteRegistration(ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				_, found := k.GetRegistration(ctx,
					tc.request.Address,
				)
				require.False(t, found)
			}
		})
	}
}
