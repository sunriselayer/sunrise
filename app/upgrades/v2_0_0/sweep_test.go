package v2_0_0

import (
	"context"
	"os"
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v10/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v10/modules/core/04-channel/types"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("sunrise", "sunrisepub")
	config.SetBech32PrefixForValidator("sunrisevaloper", "sunrisevaloperpub")
	config.SetBech32PrefixForConsensusNode("sunrisevalcons", "sunrisevalconspub")
	os.Exit(m.Run())
}

type memoryBank struct {
	rows []balanceRow
	sent []sentCoins
}

type balanceRow struct {
	address sdk.AccAddress
	coin    sdk.Coin
}

type sentCoins struct {
	from  string
	to    string
	coins sdk.Coins
}

func (m *memoryBank) IterateAllBalances(_ context.Context, cb func(address sdk.AccAddress, coin sdk.Coin) bool) {
	for _, row := range m.rows {
		if cb(row.address, row.coin) {
			return
		}
	}
}

func (m *memoryBank) SendCoins(_ context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error {
	m.sent = append(m.sent, sentCoins{from: fromAddr.String(), to: toAddr.String(), coins: amt})
	return nil
}

type channelList struct {
	channels []channeltypes.IdentifiedChannel
}

func (c channelList) GetAllChannelsWithPortPrefix(sdk.Context, string) []channeltypes.IdentifiedChannel {
	return c.channels
}

func TestSweepIBCLeavesNativeAndEscrowBalances(t *testing.T) {
	user := sdk.AccAddress("user-address-00000001")
	other := sdk.AccAddress("other-address-0000001")
	recovery, err := sdk.AccAddressFromBech32(RecoveryAddress)
	require.NoError(t, err)
	escrow := transfertypes.GetEscrowAddress(transfertypes.PortID, "channel-0")

	bank := &memoryBank{rows: []balanceRow{
		{user, sdk.NewInt64Coin("urise", 20)},
		{user, sdk.NewInt64Coin("ibc/AAA", 5)},
		{user, sdk.NewInt64Coin("uusdrise", 7)},
		{user, sdk.NewInt64Coin("ibc/BBB", 9)},
		{escrow, sdk.NewInt64Coin("ibc/AAA", 100)},
		{recovery, sdk.NewInt64Coin("ibc/CCC", 3)},
		{other, sdk.NewInt64Coin("factory/sunrise1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq/test", 4)},
	}}
	ctx := sdk.Context{}.WithLogger(log.NewNopLogger())
	err = SweepIBC(ctx, bank, channelList{channels: []channeltypes.IdentifiedChannel{{
		PortId:    transfertypes.PortID,
		ChannelId: "channel-0",
	}}})
	require.NoError(t, err)
	require.Len(t, bank.sent, 1)
	require.Equal(t, user.String(), bank.sent[0].from)
	require.Equal(t, RecoveryAddress, bank.sent[0].to)
	require.Equal(t, sdk.NewCoins(sdk.NewInt64Coin("ibc/AAA", 5), sdk.NewInt64Coin("ibc/BBB", 9)), bank.sent[0].coins)
	require.True(t, math.NewInt(5).Equal(bank.sent[0].coins.AmountOf("ibc/AAA")))
}
