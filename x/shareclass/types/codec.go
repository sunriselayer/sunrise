package types

import (
	"cosmossdk.io/core/registry"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"

	lockuptypes "github.com/sunriselayer/sunrise/x/accounts/non_voting_delegatable_lockup/v1"
)

func RegisterInterfaces(registrar registry.InterfaceRegistrar) {
	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateValidator{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgNonVotingDelegate{},
		&MsgNonVotingUndelegate{},
		&MsgClaimRewards{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&lockuptypes.MsgInitNonVotingDelegatableLockupAccount{},
		&lockuptypes.MsgInitNonVotingDelegatablePeriodicLockingAccount{},
		&lockuptypes.MsgDelegate{},
		&lockuptypes.MsgUndelegate{},
		&lockuptypes.MsgWithdrawReward{},
		&lockuptypes.MsgSend{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},
	)
	msgservice.RegisterMsgServiceDesc(registrar, &_Msg_serviceDesc)
}
