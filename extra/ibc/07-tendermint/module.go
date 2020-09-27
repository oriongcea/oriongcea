package tendermint

import (
	"github.com/spf13/cobra"

	"github.com/ocea/sdk/extra/ibc/07-tendermint/client/cli"
	"github.com/ocea/sdk/extra/ibc/07-tendermint/types"
)

// Name returns the IBC client name
func Name() string {
	return types.SubModuleName
}

// GetTxCmd returns the root tx command for the IBC client
func GetTxCmd() *cobra.Command {
	return cli.NewTxCmd()
}
