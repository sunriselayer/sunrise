package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/sunriselayer/sunrise/app/shutdownledger"
)

// ShutdownLedgerCmd writes holder files from a pre-upgrade state export.
// It does not open the chain database, so it can run after the chain is shut down.
// The export does not record its block time, so --snapshot-time supplies the time
// that claims.json uses for lockup vesting.
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
			snapshotTime, err := cmd.Flags().GetString("snapshot-time")
			if err != nil {
				return fmt.Errorf("shutdown-ledger: read snapshot-time flag: %w", err)
			}
			if input == "" || output == "" || snapshotTime == "" {
				return fmt.Errorf("shutdown-ledger: --input, --output-dir, and --snapshot-time are required")
			}
			snapshot, err := time.Parse(time.RFC3339Nano, snapshotTime)
			if err != nil {
				return fmt.Errorf("shutdown-ledger: parse snapshot-time %q: %w", snapshotTime, err)
			}
			if err := shutdownledger.Build(input, output, snapshot); err != nil {
				return err
			}
			cmd.Printf("wrote holder files to %s\n", output)
			return nil
		},
	}
	cmd.Flags().String("input", "", "Path to pre-upgrade-state.json")
	cmd.Flags().String("output-dir", "", "Directory for bank.json, positions.json, lockup.json, in_flight.json, staking.json, and claims.json")
	cmd.Flags().String("snapshot-time", "", "RFC3339 block time at the export's initial_height, used for lockup vesting in claims.json")
	return cmd
}
