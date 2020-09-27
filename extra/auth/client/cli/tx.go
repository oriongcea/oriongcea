package cli

import (
	"github.com/spf13/cobra"

	"github.com/ocea/sdk/client"
	"github.com/ocea/sdk/extra/auth/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Auth transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(
		GetMultiSignCommand(),
		GetSignCommand(),
		GetValidateSignaturesCommand(),
		GetSignBatchCommand(),
	)
	return txCmd
}
