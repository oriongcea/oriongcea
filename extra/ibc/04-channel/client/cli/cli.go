package cli

import (
	"github.com/spf13/cobra"

	"github.com/ocea/sdk/client"
	"github.com/ocea/sdk/extra/ibc/04-channel/types"
)

// GetQueryCmd returns the query commands for IBC channels
func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        types.SubModuleName,
		Short:                      "IBC channel query subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	queryCmd.AddCommand(
		GetCmdQueryChannels(),
		GetCmdQueryChannel(),
		GetCmdQueryConnectionChannels(),
		GetCmdQueryChannelClientState(),
		GetCmdQueryPacketCommitment(),
		GetCmdQueryPacketCommitments(),
		GetCmdQueryUnreceivedPackets(),
		GetCmdQueryUnrelayedAcks(),
		GetCmdQueryNextSequenceReceive(),
		// TODO: next sequence Send ?
	)

	return queryCmd
}

// NewTxCmd returns a CLI command handler for all extra/ibc channel transaction commands.
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.SubModuleName,
		Short:                      "IBC channel transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		NewChannelOpenInitCmd(),
		NewChannelOpenTryCmd(),
		NewChannelOpenAckCmd(),
		NewChannelOpenConfirmCmd(),
		NewChannelCloseInitCmd(),
		NewChannelCloseConfirmCmd(),
	)

	return txCmd
}
