package app

import (
	"fmt"

	"cosmossdk.io/log"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"

	// blocksdkabci "github.com/skip-mev/block-sdk/v2/abci"

	"github.com/sunriselayer/sunrise/x/da/keeper"
	"github.com/sunriselayer/sunrise/x/da/types"
)

type ProposalHandler struct {
	logger                 log.Logger
	keeper                 keeper.Keeper
	DefaultProposalHandler *baseapp.DefaultProposalHandler
	// DefaultProposalHandler *blocksdkabci.ProposalHandler
}

func NewProposalHandler(
	logger log.Logger,
	keeper keeper.Keeper,
	proposalHandler *baseapp.DefaultProposalHandler,
	// proposalHandler *blocksdkabci.ProposalHandler,
) *ProposalHandler {
	return &ProposalHandler{
		logger:                 logger,
		keeper:                 keeper,
		DefaultProposalHandler: proposalHandler,
	}
}

func (h *ProposalHandler) ProcessProposal() sdk.ProcessProposalHandler {
	return func(ctx sdk.Context, req *abci.ProcessProposalRequest) (*abci.ProcessProposalResponse, error) {
		txs := req.Txs
		verifiedData, err := h.keeper.GetSpecificStatusDataBeforeTime(ctx, types.Status_STATUS_VERIFIED, ctx.BlockTime().Unix())
		if err != nil {
			return nil, err
		}

		for _, data := range verifiedData {
			metadataUri := &types.MetadataUriWrapper{
				MetadataUri: data.MetadataUri,
			}

			metadataUriBz, err := metadataUri.Marshal()
			if err != nil {
				return nil, fmt.Errorf("failed to marshal metadata uri: %w", err)
			}

			txs = append(txs, metadataUriBz)
		}

		defaultReq := *req
		defaultReq.Txs = txs
		defaultHandler := h.DefaultProposalHandler.ProcessProposalHandler()
		return defaultHandler(ctx, &defaultReq)
	}
}
