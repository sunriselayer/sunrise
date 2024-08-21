package testnode

// import (
// 	"github.com/sunriselayer/sunrise/app/encoding"
// 	"github.com/sunriselayer/sunrise/pkg/user"
// 	"github.com/sunriselayer/sunrise/test/util/testfactory"

// 	testencoding "github.com/sunriselayer/sunrise/test/util/encoding"
// )

// func NewOfflineSigner() (*user.Signer, error) {
// 	encCfg := encoding.MakeConfig(testencoding.ModuleEncodingRegisters...)
// 	kr, addr := NewKeyring(testfactory.TestAccName)
// 	return user.NewSigner(kr, nil, addr[0], encCfg.TxConfig, testfactory.ChainID, 1, 0)
// }

// func NewSingleSignerFromContext(ctx Context) (*user.Signer, error) {
// 	encCfg := encoding.MakeConfig(testencoding.ModuleEncodingRegisters...)
// 	return user.SetupSingleSigner(ctx.GoContext(), ctx.Keyring, ctx.GRPCClient, encCfg)
// }

// func NewSignerFromContext(ctx Context, acc string) (*user.Signer, error) {
// 	encCfg := encoding.MakeConfig(testencoding.ModuleEncodingRegisters...)
// 	addr := testfactory.GetAddress(ctx.Keyring, acc)
// 	return user.SetupSigner(ctx.GoContext(), ctx.Keyring, ctx.GRPCClient, addr, encCfg)
// }
