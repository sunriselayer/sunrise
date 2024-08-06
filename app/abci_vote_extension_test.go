package app_test

import (
	"fmt"
	"testing"

	"github.com/cometbft/cometbft/crypto/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/sunriselayer/sunrise/app"
	"github.com/sunriselayer/sunrise/x/da/types"
)

func TestConfirmVotedData(t *testing.T) {
	data1 := types.PublishedData{
		MetadataUri:       "metadata/data1",
		ShardDoubleHashes: [][]byte{[]byte("data1")},
	}
	data2 := types.PublishedData{
		MetadataUri:       "metadata/data2",
		ShardDoubleHashes: [][]byte{[]byte("data2")},
	}
	data2v := types.PublishedData{
		MetadataUri:       "metadata/data2v",
		ShardDoubleHashes: [][]byte{[]byte("data2v")},
	}
	data3 := types.PublishedData{
		MetadataUri:       "metadata/data3",
		ShardDoubleHashes: [][]byte{[]byte("data3")},
	}
	testCases := []struct {
		expected map[string]*types.PublishedData
		actual   map[string]*types.PublishedData
		expEqual bool
	}{
		{
			expected: map[string]*types.PublishedData{
				data1.MetadataUri: &data1,
				data2.MetadataUri: &data2,
			},
			actual: map[string]*types.PublishedData{
				data1.MetadataUri: &data1,
				data2.MetadataUri: &data2,
			},
			expEqual: true,
		},
		{
			expected: map[string]*types.PublishedData{
				data1.MetadataUri: &data1,
				data2.MetadataUri: &data2,
			},
			actual: map[string]*types.PublishedData{
				data1.MetadataUri: &data1,
				data2.MetadataUri: &data2v,
			},
			expEqual: false,
		},
		{
			expected: map[string]*types.PublishedData{
				data1.MetadataUri: &data1,
				data2.MetadataUri: &data2,
			},
			actual: map[string]*types.PublishedData{
				data1.MetadataUri: &data1,
				data2.MetadataUri: &data2,
				data3.MetadataUri: &data3,
			},
			expEqual: false,
		},
		{
			expected: map[string]*types.PublishedData{
				data1.MetadataUri: &data1,
				data2.MetadataUri: &data2,
				data3.MetadataUri: &data3,
			},
			actual: map[string]*types.PublishedData{
				data1.MetadataUri: &data1,
				data2.MetadataUri: &data2,
			},
			expEqual: false,
		},
		{
			expected: nil,
			actual: map[string]*types.PublishedData{
				data1.MetadataUri: &data1,
				data2.MetadataUri: &data2,
			},
			expEqual: false,
		},
		{
			expected: map[string]*types.PublishedData{
				data1.MetadataUri: &data1,
				data2.MetadataUri: &data2,
			},
			actual:   nil,
			expEqual: false,
		},
		{
			expected: nil,
			actual:   nil,
			expEqual: true,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			err := app.ConfirmVotedData(tc.actual, tc.expected)
			if tc.expEqual {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestConfirmFaultValidators(t *testing.T) {
	var valAddrs = []sdk.ValAddress{
		sdk.ValAddress(secp256k1.GenPrivKey().PubKey().Address()),
		sdk.ValAddress(secp256k1.GenPrivKey().PubKey().Address()),
		sdk.ValAddress(secp256k1.GenPrivKey().PubKey().Address()),
	}

	testCases := []struct {
		expected map[string]sdk.ValAddress
		actual   map[string]sdk.ValAddress
		expEqual bool
	}{
		{
			expected: map[string]sdk.ValAddress{
				valAddrs[0].String(): valAddrs[0],
				valAddrs[1].String(): valAddrs[1],
			},
			actual: map[string]sdk.ValAddress{
				valAddrs[0].String(): valAddrs[0],
				valAddrs[1].String(): valAddrs[1],
			},
			expEqual: true,
		},
		{
			expected: map[string]sdk.ValAddress{
				valAddrs[0].String(): valAddrs[0],
				valAddrs[1].String(): valAddrs[1],
			},
			actual: map[string]sdk.ValAddress{
				valAddrs[0].String(): valAddrs[0],
				valAddrs[2].String(): valAddrs[2],
			},
			expEqual: false,
		},
		{
			expected: map[string]sdk.ValAddress{
				valAddrs[0].String(): valAddrs[0],
				valAddrs[1].String(): valAddrs[1],
			},
			actual: map[string]sdk.ValAddress{
				valAddrs[0].String(): valAddrs[0],
				valAddrs[1].String(): valAddrs[2],
			},
			expEqual: false,
		},
		{
			expected: map[string]sdk.ValAddress{
				valAddrs[0].String(): valAddrs[0],
				valAddrs[1].String(): valAddrs[1],
			},
			actual: map[string]sdk.ValAddress{
				valAddrs[0].String(): valAddrs[0],
				valAddrs[2].String(): valAddrs[1],
			},
			expEqual: false,
		},
		{
			expected: map[string]sdk.ValAddress{
				valAddrs[0].String(): valAddrs[0],
				valAddrs[1].String(): valAddrs[1],
			},
			actual: map[string]sdk.ValAddress{
				valAddrs[0].String(): valAddrs[0],
			},
			expEqual: false,
		},
		{
			expected: map[string]sdk.ValAddress{
				valAddrs[0].String(): valAddrs[0],
			},
			actual: map[string]sdk.ValAddress{
				valAddrs[0].String(): valAddrs[0],
				valAddrs[1].String(): valAddrs[1],
			},
			expEqual: false,
		},
		{
			expected: map[string]sdk.ValAddress{},
			actual:   map[string]sdk.ValAddress{},
			expEqual: true,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			err := app.ConfirmFaultValidators(tc.actual, tc.expected)
			if tc.expEqual {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
