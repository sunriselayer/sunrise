package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/sunriselayer/sunrise/app/shutdownledger"
)

// ShutdownLedgerCmd writes holder files from a pre-upgrade state export.
// It does not open the chain database, so it can run after the chain is shut down.
func ShutdownLedgerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "shutdown-ledger",
		Short: "Build holder files from a pre-upgrade state JSON file",
		RunE: func(cmd *cobra.Command, _ []string) error {
			input, err := cmd.Flags().GetString("input")
			if err != nil {
				return fmt.Errorf("shutdown-ledger: read input flag: %w", err)
			}
			output, err := cmd.Flags().GetString("output-dir")
			if err != nil {
				return fmt.Errorf("shutdown-ledger: read output-dir flag: %w", err)
			}
			if input == "" || output == "" {
				return fmt.Errorf("shutdown-ledger: --input and --output-dir are required")
			}
			if err := shutdownledger.Build(input, output); err != nil {
				return err
			}
			cmd.Printf("wrote holder files to %s\n", output)
			return nil
		},
	}
	cmd.Flags().String("input", "", "Path to pre-upgrade-state.json")
	cmd.Flags().String("output-dir", "", "Directory for bank.json, positions.json, lockup.json, in_flight.json, staking.json, and claims.json")
	return cmd
}
