package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/sunriselayer/sunrise/x/da/types"
)

func TestMsgRegisterCommitmentKey(t *testing.T) {
	k, _, ms, ctx := setupMsgServer(t)

	sender := sdk.AccAddress("sender")
	deputy := sdk.AccAddress("deputy")
	invalidAddr := "invalid_address"

	testCases := []struct {
		name      string
		input     *types.MsgRegisterProofDeputy
		expErr    bool
		expErrMsg string
	}{
		{
			name: "invalid sender address",
			input: &types.MsgRegisterProofDeputy{
				Sender:        invalidAddr,
				DeputyAddress: deputy.String(),
			},
			expErr:    true,
			expErrMsg: "invalid sender address",
		},
		{
			name: "invalid deputy address",
			input: &types.MsgRegisterProofDeputy{
				Sender:        sender.String(),
				DeputyAddress: invalidAddr,
			},
			expErr:    true,
			expErrMsg: "invalid deputy address",
		},
		{
			name: "normal case",
			input: &types.MsgRegisterProofDeputy{
				Sender:        sender.String(),
				DeputyAddress: deputy.String(),
			},
			expErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ms.RegisterProofDeputy(ctx, tc.input)

			if tc.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expErrMsg)
			} else {
				require.NoError(t, err)

				// in normal case, proof deputy is saved in storage
				deputyAddr, found, err := k.GetProofDeputy(ctx, sender)
				require.NoError(t, err)
				require.True(t, found)

				require.Equal(t, deputy.Bytes(), deputyAddr)
			}
		})
	}
}

func TestMsgUnregisterProofDeputy(t *testing.T) {
	k, _, ms, ctx := setupMsgServer(t)

	sender := sdk.AccAddress("sender")
	deputy := sdk.AccAddress("deputy")
	invalidAddr := "invalid_address"

	// set proof deputy for test
	err := k.SetProofDeputy(ctx, sender, deputy)
	require.NoError(t, err)

	testCases := []struct {
		name      string
		input     *types.MsgUnregisterProofDeputy
		expErr    bool
		expErrMsg string
	}{
		{
			name: "invalid sender address",
			input: &types.MsgUnregisterProofDeputy{
				Sender: invalidAddr,
			},
			expErr:    true,
			expErrMsg: "invalid sender address",
		},
		{
			name: "deputy not found",
			input: &types.MsgUnregisterProofDeputy{
				Sender: deputy.String(),
			},
			expErr:    true,
			expErrMsg: "deputy not found",
		},
		{
			name: "normal case",
			input: &types.MsgUnregisterProofDeputy{
				Sender: sender.String(),
			},
			expErr: false,
		},
		{
			name: "deputy already deleted",
			input: &types.MsgUnregisterProofDeputy{
				Sender: sender.String(),
			},
			expErr:    true,
			expErrMsg: "deputy not found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ms.UnregisterProofDeputy(ctx, tc.input)

			if tc.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expErrMsg)
			} else {
				require.NoError(t, err)

				// in normal case, proof deputy is deleted from storage
				_, found, err := k.GetProofDeputy(ctx, sender)
				require.NoError(t, err)
				require.False(t, found)
			}
		})
	}
}
