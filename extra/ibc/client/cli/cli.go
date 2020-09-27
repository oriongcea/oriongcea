package cli

import (
	"github.com/spf13/cobra"

	"github.com/ocea/sdk/client"
	ibcclient "github.com/ocea/sdk/extra/ibc/02-client"
	connection "github.com/ocea/sdk/extra/ibc/03-connection"
	channel "github.com/ocea/sdk/extra/ibc/04-channel"
	tendermint "github.com/ocea/sdk/extra/ibc/07-tendermint"
	host "github.com/ocea/sdk/extra/ibc/24-host"
	"github.com/ocea/sdk/extra/ibc/light-clients/solomachine"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	ibcTxCmd := &cobra.Command{
		Use:                        host.ModuleName,
		Short:                      "IBC transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	ibcTxCmd.AddCommand(
		solomachine.GetTxCmd(),
		tendermint.GetTxCmd(),
		connection.GetTxCmd(),
		channel.GetTxCmd(),
	)

	return ibcTxCmd
}

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	// Group ibc queries under a subcommand
	ibcQueryCmd := &cobra.Command{
		Use:                        host.ModuleName,
		Short:                      "Querying commands for the IBC module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	ibcQueryCmd.AddCommand(
		ibcclient.GetQueryCmd(),
		connection.GetQueryCmd(),
		channel.GetQueryCmd(),
	)

	return ibcQueryCmd
}
