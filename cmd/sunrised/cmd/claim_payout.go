package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"cosmossdk.io/math"
	"github.com/spf13/cobra"

	"github.com/sunriselayer/sunrise/app/claimpayout"
)

// ClaimPayoutCmd checks a Keplr ADR-036 signature and records an external payout.
// The hot-wallet key stays outside this command. After reserve, send the tokens,
// then record the transaction hash with complete. release drops a failed send.
func ClaimPayoutCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim-payout",
		Short: "Verify a claim signature and record an external USDC or Cosmos payout",
	}
	cmd.AddCommand(claimPayoutVerifyCmd(), claimPayoutReserveCmd(), claimPayoutCompleteCmd(), claimPayoutReleaseCmd())
	return cmd
}

func claimPayoutVerifyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify",
		Short: "Check a signature against claims.json without recording a payout",
		RunE: func(cmd *cobra.Command, _ []string) error {
			auth, pub, sig, ledgerPath, claimsPath, err := readClaimFlags(cmd)
			if err != nil {
				return err
			}
			storePath, _ := cmd.Flags().GetString("store")
			available, payout, err := claimpayout.Check(ledgerPath, claimsPath, storePath, auth, pub, sig, time.Now().UTC())
			if err != nil {
				return err
			}
			cmd.Printf("claimant %s asset %s payout %s available %s\n", auth.Claimant, auth.Asset, payout, available)
			return nil
		},
	}
	addClaimFlags(cmd)
	cmd.Flags().String("store", "", "Optional payout record used to subtract already reserved amounts")
	return cmd
}

func claimPayoutReserveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reserve",
		Short: "Hold a nonce after the signature, balance, and caps check",
		RunE: func(cmd *cobra.Command, _ []string) error {
			auth, pub, sig, ledgerPath, claimsPath, err := readClaimFlags(cmd)
			if err != nil {
				return err
			}
			storePath, err := requiredString(cmd, "store")
			if err != nil {
				return err
			}
			policy, err := readPolicy(cmd)
			if err != nil {
				return err
			}
			if err := claimpayout.Reserve(storePath, ledgerPath, claimsPath, auth, pub, sig, policy, time.Now().UTC()); err != nil {
				return err
			}
			cmd.Printf("reserved nonce %d for %s %s to %s\n", auth.Nonce, auth.Amount, auth.Asset, auth.Destination)
			return nil
		},
	}
	addClaimFlags(cmd)
	cmd.Flags().String("store", "", "Payout record file")
	cmd.Flags().String("max-per-tx", "", "Maximum amount for this asset in one payout")
	cmd.Flags().String("max-daily", "", "Maximum amount for this asset during the UTC day")
	return cmd
}

func claimPayoutCompleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "complete",
		Short: "Record the hot-wallet transaction hash for a reserved nonce",
		RunE: func(cmd *cobra.Command, _ []string) error {
			storePath, err := requiredString(cmd, "store")
			if err != nil {
				return err
			}
			nonce, err := cmd.Flags().GetUint64("nonce")
			if err != nil {
				return fmt.Errorf("claim-payout complete: nonce: %w", err)
			}
			txHash, err := requiredString(cmd, "tx-hash")
			if err != nil {
				return err
			}
			if err := claimpayout.Complete(storePath, nonce, txHash, time.Now().UTC()); err != nil {
				return err
			}
			cmd.Printf("consumed nonce %d with %s\n", nonce, txHash)
			return nil
		},
	}
	cmd.Flags().String("store", "", "Payout record file")
	cmd.Flags().Uint64("nonce", 0, "Nonce from the signed authorization")
	cmd.Flags().String("tx-hash", "", "Hot-wallet transaction hash")
	return cmd
}

func claimPayoutReleaseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "release",
		Short: "Drop a reservation after the external send fails",
		RunE: func(cmd *cobra.Command, _ []string) error {
			storePath, err := requiredString(cmd, "store")
			if err != nil {
				return err
			}
			nonce, err := cmd.Flags().GetUint64("nonce")
			if err != nil {
				return fmt.Errorf("claim-payout release: nonce: %w", err)
			}
			if err := claimpayout.Release(storePath, nonce); err != nil {
				return err
			}
			cmd.Printf("released nonce %d\n", nonce)
			return nil
		},
	}
	cmd.Flags().String("store", "", "Payout record file")
	cmd.Flags().Uint64("nonce", 0, "Nonce from the signed authorization")
	return cmd
}

func addClaimFlags(cmd *cobra.Command) {
	cmd.Flags().String("ledger", "", "claims.json file whose bytes are hashed into the signature")
	cmd.Flags().String("claims", "", "claims.json used for balances; defaults to --ledger")
	cmd.Flags().String("authorization", "", "Signed authorization JSON")
	cmd.Flags().String("pubkey", "", "Compressed secp256k1 pubkey hex")
	cmd.Flags().String("signature", "", "64-byte signature hex")
}

func readClaimFlags(cmd *cobra.Command) (claimpayout.Authorization, []byte, []byte, string, string, error) {
	ledgerPath, err := requiredString(cmd, "ledger")
	if err != nil {
		return claimpayout.Authorization{}, nil, nil, "", "", err
	}
	claimsPath, err := cmd.Flags().GetString("claims")
	if err != nil {
		return claimpayout.Authorization{}, nil, nil, "", "", fmt.Errorf("claim-payout: claims: %w", err)
	}
	if claimsPath == "" {
		claimsPath = ledgerPath
	}
	authPath, err := requiredString(cmd, "authorization")
	if err != nil {
		return claimpayout.Authorization{}, nil, nil, "", "", err
	}
	raw, err := os.ReadFile(authPath)
	if err != nil {
		return claimpayout.Authorization{}, nil, nil, "", "", fmt.Errorf("claim-payout: read authorization %s: %w", authPath, err)
	}
	var auth claimpayout.Authorization
	if err := json.Unmarshal(raw, &auth); err != nil {
		return claimpayout.Authorization{}, nil, nil, "", "", fmt.Errorf("claim-payout: authorization %s: %w", authPath, err)
	}
	pubHex, err := requiredString(cmd, "pubkey")
	if err != nil {
		return claimpayout.Authorization{}, nil, nil, "", "", err
	}
	sigHex, err := requiredString(cmd, "signature")
	if err != nil {
		return claimpayout.Authorization{}, nil, nil, "", "", err
	}
	pub, err := claimpayout.DecodeHex("pubkey", pubHex)
	if err != nil {
		return claimpayout.Authorization{}, nil, nil, "", "", err
	}
	sig, err := claimpayout.DecodeHex("signature", sigHex)
	if err != nil {
		return claimpayout.Authorization{}, nil, nil, "", "", err
	}
	return auth, pub, sig, ledgerPath, claimsPath, nil
}

func readPolicy(cmd *cobra.Command) (claimpayout.Policy, error) {
	perTx, err := requiredString(cmd, "max-per-tx")
	if err != nil {
		return claimpayout.Policy{}, err
	}
	daily, err := requiredString(cmd, "max-daily")
	if err != nil {
		return claimpayout.Policy{}, err
	}
	perTxAmount, ok := math.NewIntFromString(perTx)
	if !ok {
		return claimpayout.Policy{}, fmt.Errorf("claim-payout: max-per-tx %q", perTx)
	}
	dailyAmount, ok := math.NewIntFromString(daily)
	if !ok {
		return claimpayout.Policy{}, fmt.Errorf("claim-payout: max-daily %q", daily)
	}
	return claimpayout.Policy{MaxPerTx: perTxAmount, MaxDaily: dailyAmount}, nil
}

func requiredString(cmd *cobra.Command, name string) (string, error) {
	value, err := cmd.Flags().GetString(name)
	if err != nil {
		return "", fmt.Errorf("claim-payout: %s: %w", name, err)
	}
	if value == "" {
		return "", fmt.Errorf("claim-payout: --%s is required", name)
	}
	return value, nil
}
