package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/stream module sentinel errors
var (
	ErrInvalidSigner = sdkerrors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")

	ErrDuplicate                                 = sdkerrors.Register(ModuleName, 2, "duplicate")
	ErrEmpty                                     = sdkerrors.Register(ModuleName, 6, "empty")
	ErrNoValidators                              = sdkerrors.Register(ModuleName, 12, "no bonded validators in active set")
	ErrInvalidValAddress                         = sdkerrors.Register(ModuleName, 13, "invalid validator address in current valset %v")
	ErrInvalidEVMAddress                         = sdkerrors.Register(ModuleName, 14, "discovered invalid EVM address stored for validator %v")
	ErrInvalidValset                             = sdkerrors.Register(ModuleName, 15, "generated invalid valset")
	ErrAttestationNotValsetRequest               = sdkerrors.Register(ModuleName, 16, "attestation is not a valset request")
	ErrAttestationNotFound                       = sdkerrors.Register(ModuleName, 18, "attestation not found")
	ErrNilAttestation                            = sdkerrors.Register(ModuleName, 22, "nil attestation")
	ErrUnmarshalllAttestation                    = sdkerrors.Register(ModuleName, 26, "couldn't unmarshall attestation from store")
	ErrNonceHigherThanLatestAttestationNonce     = sdkerrors.Register(ModuleName, 27, "the provided nonce is higher than the latest attestation nonce")
	ErrNoValsetBeforeNonceOne                    = sdkerrors.Register(ModuleName, 28, "there is no valset before attestation nonce 1")
	ErrDataCommitmentNotGenerated                = sdkerrors.Register(ModuleName, 29, "no data commitment has been generated for the provided height")
	ErrDataCommitmentNotFound                    = sdkerrors.Register(ModuleName, 30, "data commitment not found")
	ErrLatestAttestationNonceStillNotInitialized = sdkerrors.Register(ModuleName, 31, "the latest attestation nonce has still not been defined in store")
	ErrInvalidDataCommitmentWindow               = sdkerrors.Register(ModuleName, 32, "invalid data commitment window")
	ErrEarliestAvailableNonceStillNotInitialized = sdkerrors.Register(ModuleName, 33, "the earliest available nonce after pruning has still not been defined in store")
	ErrRequestedNonceWasPruned                   = sdkerrors.Register(ModuleName, 34, "the requested nonce has been pruned")
	ErrUnknownAttestationType                    = sdkerrors.Register(ModuleName, 35, "unknown attestation type")
	ErrEVMAddressNotHex                          = sdkerrors.Register(ModuleName, 36, "the provided evm address is not a valid hex address")
	ErrEVMAddressAlreadyExists                   = sdkerrors.Register(ModuleName, 37, "the provided evm address already exists")
	ErrEVMAddressNotFound                        = sdkerrors.Register(ModuleName, 38, "EVM address not found")
)
