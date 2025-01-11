package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sunriselayer/sunrise/x/da/keeper"
	"github.com/sunriselayer/sunrise/x/da/types"
)

func TestMsgUpdateParams(t *testing.T) {
	f := initFixture(t)
	ms := keeper.NewMsgServerImpl(f.keeper)

	params := types.DefaultParams()
	require.NoError(t, f.keeper.Params.Set(f.ctx, params))

	authorityStr, err := f.addressCodec.BytesToString(f.keeper.GetAuthority())
	require.NoError(t, err)

	// default params
	testCases := []struct {
		name      string
		input     *types.MsgUpdateParams
		expErr    bool
		expErrMsg string
	}{
		{
			name: "invalid authority",
			input: &types.MsgUpdateParams{
				Authority: "invalid",
				Params:    params,
			},
			expErr:    true,
			expErrMsg: "invalid authority",
		},
		{
			name: "send enabled param",
			input: &types.MsgUpdateParams{
				Authority: authorityStr,
				Params: types.Params{
					VoteThreshold:     "0.67",
					SlashEpoch:        1,
					EpochMaxFault:     1,
					SlashFraction:     "0.001",
					ReplicationFactor: "3",
					MinShardCount:     2,
					MaxShardCount:     3,
					MaxShardSize:      1000000000000000000,
					ChallengePeriod:   1,
					ProofPeriod:       1,
				},
			},
			expErr: false,
		},
		{
			name: "all good",
			input: &types.MsgUpdateParams{
				Authority: authorityStr,
				Params:    params,
			},
			expErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ms.UpdateParams(f.ctx, tc.input)

			if tc.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expErrMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
