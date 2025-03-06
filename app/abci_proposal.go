package app

import (
	"bytes"
	"fmt"

	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/log"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	// blocksdkabci "github.com/skip-mev/block-sdk/v2/abci"

	"github.com/sunriselayer/sunrise/x/da/keeper"
	"github.com/sunriselayer/sunrise/x/da/types"
)

var metadataUriSplitter = []byte{0x4D, 0x45, 0x54, 0x41, 0x44, 0x41, 0x54, 0x41} // ASCII for "METADATA"

type ProposalHandler struct {
	logger                 log.Logger
	keeper                 keeper.Keeper
	ModuleManager          *module.Manager
	DefaultProposalHandler *baseapp.DefaultProposalHandler
	// DefaultProposalHandler *blocksdkabci.ProposalHandler
}

func NewProposalHandler(
	logger log.Logger,
	keeper keeper.Keeper,
	ModuleManager *module.Manager,
	proposalHandler *baseapp.DefaultProposalHandler,
	// proposalHandler *blocksdkabci.ProposalHandler,
) *ProposalHandler {
	return &ProposalHandler{
		logger:                 logger,
		keeper:                 keeper,
		ModuleManager:          ModuleManager,
		DefaultProposalHandler: proposalHandler,
	}
}

func (h *ProposalHandler) PrepareProposal() sdk.PrepareProposalHandler {
	return func(ctx sdk.Context, req *abci.PrepareProposalRequest) (*abci.PrepareProposalResponse, error) {
		defaultHandler := h.DefaultProposalHandler.PrepareProposalHandler()
		defaultResponse, err := defaultHandler(ctx, req)
		if err != nil {
			return nil, err
		}

		proposalTxs := defaultResponse.Txs

		verifiedData, err := h.keeper.GetSpecificStatusData(ctx, types.Status_STATUS_VERIFIED)
		if err != nil {
			return nil, err
		}

		if len(verifiedData) > 0 {
			proposalTxs = append(proposalTxs, metadataUriSplitter)

			for _, data := range verifiedData {
				metadataUri := &types.MetadataUriWrapper{
					MetadataUri: data.MetadataUri,
				}

				metadataUriBz, err := metadataUri.Marshal()
				if err != nil {
					return nil, fmt.Errorf("failed to marshal metadata uri: %w", err)
				}

				proposalTxs = append(proposalTxs, metadataUriBz)
			}
		}

		return &abci.PrepareProposalResponse{Txs: proposalTxs}, nil
	}
}

func (h *ProposalHandler) ProcessProposal() sdk.ProcessProposalHandler {
	return func(ctx sdk.Context, req *abci.ProcessProposalRequest) (*abci.ProcessProposalResponse, error) {
		splitterIndex := -1
		for i, txBytes := range req.Txs {
			if bytes.Equal(txBytes, metadataUriSplitter) {
				splitterIndex = i
				break
			}
		}

		if splitterIndex != -1 {
			for i := splitterIndex + 1; i < len(req.Txs); i++ {
				var metadataUri types.MetadataUriWrapper
				err := metadataUri.Unmarshal(req.Txs[i])
				if err != nil {
					h.logger.Error("failed to unmarshal metadata uri", "error", err)
					continue
				}

				if metadataUri.MetadataUri == "" {
					h.logger.Error("metadata uri is empty")
					continue
				}

				data, found, err := h.keeper.GetPublishedData(ctx, metadataUri.MetadataUri)
				if err != nil {
					return nil, err
				}
				if !found {
					h.logger.Error("published data not found", "metadata uri", metadataUri.MetadataUri)
					continue
				}
				if data.Status != types.Status_STATUS_VERIFIED {
					return &abci.ProcessProposalResponse{Status: abci.PROCESS_PROPOSAL_STATUS_REJECT}, nil
				}
			}
		}
		defaultHandler := h.DefaultProposalHandler.ProcessProposalHandler()
		return defaultHandler(ctx, req)
	}
}

func (h *ProposalHandler) PreBlocker(ctx sdk.Context, req *abci.FinalizeBlockRequest) error {
	ctx = ctx.WithEventManager(sdk.NewEventManager())
	for _, moduleName := range h.ModuleManager.OrderPreBlockers {
		if module, ok := h.ModuleManager.Modules[moduleName].(appmodule.HasPreBlocker); ok {
			err := module.PreBlock(ctx)
			if err != nil {
				return err
			}
		}
	}

	splitterIndex := -1
	for i, txBytes := range req.Txs {
		if bytes.Equal(txBytes, metadataUriSplitter) {
			splitterIndex = i
			break
		}
	}

	if splitterIndex != -1 {
		for i := splitterIndex + 1; i < len(req.Txs); i++ {
			var metadataUri types.MetadataUriWrapper
			err := metadataUri.Unmarshal(req.Txs[i])
			if err != nil {
				h.logger.Error("failed to unmarshal metadata uri", "error", err)
				continue
			}

			if metadataUri.MetadataUri == "" {
				h.logger.Error("metadata uri is empty")
				continue
			}

			data, found, err := h.keeper.GetPublishedData(ctx, metadataUri.MetadataUri)
			if err != nil {
				return err
			}
			if !found {
				continue
			}

			err = h.keeper.DeletePublishedData(ctx, data)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
