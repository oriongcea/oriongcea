package simulation

import (
	"math/rand"

	simtypes "github.com/ocea/sdk/types/simulation"
	"github.com/ocea/sdk/extra/ibc/02-client/types"
)

// GenClientGenesis returns the default client genesis state.
func GenClientGenesis(_ *rand.Rand, _ []simtypes.Account) types.GenesisState {
	return types.DefaultGenesisState()
}
