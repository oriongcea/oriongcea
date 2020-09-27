package simulation

import (
	"math/rand"

	simtypes "github.com/ocea/sdk/types/simulation"
	"github.com/ocea/sdk/extra/ibc/03-connection/types"
)

// GenConnectionGenesis returns the default connection genesis state.
func GenConnectionGenesis(_ *rand.Rand, _ []simtypes.Account) types.GenesisState {
	return types.DefaultGenesisState()
}
