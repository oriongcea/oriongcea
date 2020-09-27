package cli

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/ocea/sdk/client"
	"github.com/ocea/sdk/client/flags"
	"github.com/ocea/sdk/client/tx"
	"github.com/ocea/sdk/extra/crisis/types"
)

// NewTxCmd returns a root CLI command handler for all extra/crisis transaction commands.
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Crisis transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(NewMsgVerifyInvariantTxCmd())

	return txCmd
}

// NewMsgVerifyInvariantTxCmd returns a CLI command handler for creating a
// MsgVerifyInvariant transaction.
func NewMsgVerifyInvariantTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "invariant-broken [module-name] [invariant-route]",
		Short: "Submit proof that an invariant broken to halt the chain",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			moduleName, route := args[0], args[1]
			if moduleName == "" {
				return errors.New("invalid module name")
			}
			if route == "" {
				return errors.New("invalid invariant route")
			}

			senderAddr := clientCtx.GetFromAddress()

			msg := types.NewMsgVerifyInvariant(senderAddr, moduleName, route)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
