package cli

import (
	"github.com/spf13/cobra"

	"github.com/ocea/sdk/client"
	"github.com/ocea/sdk/extra/ibc/02-client/types"
)

// GetQueryCmd returns the query commands for IBC clients
func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        types.SubModuleName,
		Short:                      "IBC client query subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	queryCmd.AddCommand(
		GetCmdQueryClientStates(),
		GetCmdQueryClientState(),
		GetCmdQueryConsensusStates(),
		GetCmdQueryConsensusState(),
		GetCmdQueryHeader(),
		GetCmdNodeConsensusState(),
	)

	return queryCmd
}