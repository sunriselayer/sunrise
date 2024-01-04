package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	testutil "github.com/sunrise-zone/sunrise-app/test"
	"github.com/sunrise-zone/sunrise-app/x/blobstream/types"

	"github.com/sunrise-zone/sunrise-app/x/blobstream/keeper"
)

func TestRegisterEVMAddress(t *testing.T) {
	input, sdkCtx := testutil.SetupFiveValChain(t)
	k := input.BlobstreamKeeper
	vals, err := input.StakingKeeper.GetValidators(sdkCtx, 100)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(vals), 1)
	val := vals[0]
	valoper, err := sdk.ValAddressFromBech32(val.GetOperator())
	require.NoError(t, err)
	evmAddr, exists := k.GetEVMAddress(sdkCtx, valoper)
	require.True(t, exists)

	// test again with an address that is not the validator
	valAddr, err := sdk.ValAddressFromBech32("celestiavaloper1xcy3els9ua75kdm783c3qu0rfa2eplestc6sqc")
	require.NoError(t, err)
	msg := types.NewMsgRegisterEvmAddress(valAddr, evmAddr)

	msgServer := keeper.NewMsgServerImpl(k)

	_, err = msgServer.RegisterEvmAddress(sdkCtx, msg)
	require.Error(t, err)

	// override the previous EVM address with a new one
	evmAddr = common.BytesToAddress([]byte("evm_address"))
	msg = types.NewMsgRegisterEvmAddress(valoper, evmAddr)
	_, err = msgServer.RegisterEvmAddress(sdkCtx, msg)
	require.NoError(t, err)

	addr, _ := k.GetEVMAddress(sdkCtx, valoper)
	require.Equal(t, evmAddr, addr)
}
