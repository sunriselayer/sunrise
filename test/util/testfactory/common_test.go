package testfactory_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/sunriselayer/sunrise-app/app"
	"github.com/sunriselayer/sunrise-app/app/encoding"
	"github.com/sunriselayer/sunrise-app/test/util/testfactory"
	"github.com/sunriselayer/sunrise-app/test/util/testnode"
)

func TestTestAccount(t *testing.T) {
	enc := encoding.MakeConfig(app.ModuleEncodingRegisters...)
	kr := testfactory.TestKeyring(enc.Codec)
	record, err := kr.Key(testfactory.TestAccName)
	require.NoError(t, err)
	addr, err := record.GetAddress()
	require.NoError(t, err)
	require.Equal(t, testfactory.TestAccAddr, addr.String())
	require.Equal(t, testnode.TestAddress(), addr)
}
