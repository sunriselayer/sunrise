package testfactory_test

import (
	"testing"

	"github.com/sunrise-zone/sunrise-app/app/encoding"
	util "github.com/sunrise-zone/sunrise-app/test/util"
	"github.com/sunrise-zone/sunrise-app/test/util/testfactory"
	"github.com/sunrise-zone/sunrise-app/test/util/testnode"

	"github.com/stretchr/testify/require"
)

func TestTestAccount(t *testing.T) {
	enc := encoding.MakeConfig(util.ModuleBasics)
	kr := testfactory.TestKeyring(enc.Codec)
	record, err := kr.Key(testfactory.TestAccName)
	require.NoError(t, err)
	addr, err := record.GetAddress()
	require.NoError(t, err)
	require.Equal(t, testfactory.TestAccAddr, addr.String())
	require.Equal(t, testnode.TestAddress(), addr)
}
