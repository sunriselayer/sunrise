package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestGenerateValidAddresses(t *testing.T) {
	// Generate some valid test addresses using default cosmos prefix
	addr1 := sdk.AccAddress([]byte("test1")).String()
	addr2 := sdk.AccAddress([]byte("test2")).String()
	addr3 := sdk.AccAddress([]byte("test3")).String()
	
	t.Logf("Address 1: %s", addr1)
	t.Logf("Address 2: %s", addr2)
	t.Logf("Address 3: %s", addr3)
	
	// Test that they're valid
	_, err := sdk.AccAddressFromBech32(addr1)
	require.NoError(t, err)
}