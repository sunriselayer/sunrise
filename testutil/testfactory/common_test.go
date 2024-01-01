package testfactory_test

import (
	"testing"

	"sunrise/app/encoding"
	"sunrise/testutil"
	"sunrise/testutil/testfactory"
	"sunrise/testutil/testnode"

	"github.com/stretchr/testify/require"
)

func TestTestAccount(t *testing.T) {
	enc := encoding.MakeConfig(testutil.ModuleBasics...)
	kr := testfactory.TestKeyring(enc.Codec)
	record, err := kr.Key(testfactory.TestAccName)
	require.NoError(t, err)
	addr, err := record.GetAddress()
	require.NoError(t, err)
	require.Equal(t, testfactory.TestAccAddr, addr.String())
	require.Equal(t, testnode.TestAddress(), addr)
}
