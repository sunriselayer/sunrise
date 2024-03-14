package testnode

import (
	"fmt"
	"os"
	"path"
	"strings"

	"cosmossdk.io/log"
	"github.com/cometbft/cometbft/node"
	"github.com/cometbft/cometbft/rpc/client/local"
	srvconfig "github.com/cosmos/cosmos-sdk/server/config"
	srvgrpc "github.com/cosmos/cosmos-sdk/server/grpc"
	srvtypes "github.com/cosmos/cosmos-sdk/server/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// noOpCleanup is a function that conforms to the cleanup function signature and
// performs no operation.
var noOpCleanup = func() error { return nil }

// StartNode starts the tendermint node along with a local core rpc client. The
// rpc is returned via the client.Context. The function returned should be
// called during cleanup to teardown the node, core client, along with canceling
// the internal context.Context in the returned Context.
func StartNode(tmNode *node.Node, cctx Context) (Context, func() error, error) {
	if err := tmNode.Start(); err != nil {
		return cctx, noOpCleanup, err
	}

	coreClient := local.New(tmNode)

	cctx.Context = cctx.WithClient(coreClient)
	cleanup := func() error {
		err := tmNode.Stop()
		if err != nil {
			return err
		}
		tmNode.Wait()
		if err = removeDir(path.Join([]string{cctx.HomeDir, "config"}...)); err != nil {
			return err
		}
		if err = removeDir(path.Join([]string{cctx.HomeDir, tmNode.Config().DBPath}...)); err != nil {
			return err
		}
		return nil
	}

	return cctx, cleanup, nil
}

// StartGRPCServer starts the grpc server using the provided application and
// config. A grpc client connection to that server is also added to the client
// context. The returned function should be used to shutdown the server.
func StartGRPCServer(app srvtypes.Application, appCfg *srvconfig.Config, cctx Context) (Context, func() error, error) {
	emptycleanup := func() error { return nil }
	// Add the tx service in the gRPC router.
	fmt.Println("before RegisterTxService")
	app.RegisterTxService(cctx.Context)

	// Add the tendermint queries service in the gRPC router.
	fmt.Println("before RegisterTendermintService")
	app.RegisterTendermintService(cctx.Context)

	nodeGRPCAddr := strings.Replace(appCfg.GRPC.Address, "0.0.0.0", "localhost", 1)
	fmt.Println("before grpc.Dial")
	conn, err := grpc.Dial(nodeGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return Context{}, emptycleanup, err
	}

	fmt.Println("after grpc.Dial")
	cctx.Context = cctx.WithGRPCClient(conn)
	fmt.Println("after WithGRPCClient")

	fmt.Println("before NewGRPCServer")
	// Run maxSendMsgSize := cfg.MaxSendMsgSize to gogoreflection.Register(grpcSrv)
	grpcSrv, err := srvgrpc.NewGRPCServer(cctx.Context, app, appCfg.GRPC)
	if err != nil {
		return Context{}, emptycleanup, err
	}
	fmt.Println("before StartGRPCServer")

	err = srvgrpc.StartGRPCServer(cctx.rootCtx, log.NewNopLogger().With("module", "grpc-server"), appCfg.GRPC, grpcSrv)
	fmt.Println("after StartGRPCServer")
	if err != nil {
		return Context{}, emptycleanup, err
	}

	return cctx, func() error {
		// grpcSrv.Stop()
		return nil
	}, nil
}

// DefaultAppConfig wraps the default config described in the server
func DefaultAppConfig() *srvconfig.Config {
	return srvconfig.DefaultConfig()
}

// removeDir removes the directory `rootDir`.
// The main use of this is to reduce the flakiness of the CI when it's unable to delete
// the config folder of the tendermint node.
// This will manually go over the files contained inside the provided `rootDir`
// and delete them one by one.
func removeDir(rootDir string) error {
	dir, err := os.ReadDir(rootDir)
	if err != nil {
		return err
	}
	for _, d := range dir {
		path := path.Join([]string{rootDir, d.Name()}...)
		err := os.RemoveAll(path)
		if err != nil {
			return err
		}
	}
	return os.RemoveAll(rootDir)
}
