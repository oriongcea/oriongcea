package simulation

// DONTCOVER

import (
	"fmt"
	"math/rand"

	"github.com/ocea/sdk/extra/simulation"

	simtypes "github.com/ocea/sdk/types/simulation"
	"github.com/ocea/sdk/extra/auth/types"
)

const (
	keyMaxMemoCharacters = "MaxMemoCharacters"
	keyTxSigLimit        = "TxSigLimit"
	keyTxSizeCostPerByte = "TxSizeCostPerByte"
)

// ParamChanges defines the parameters that can be modified by param change proposals
// on the simulation
func ParamChanges(r *rand.Rand) []simtypes.ParamChange {
	return []simtypes.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, keyMaxMemoCharacters,
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GenMaxMemoChars(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, keyTxSigLimit,
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GenTxSigLimit(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, keyTxSizeCostPerByte,
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GenTxSizeCostPerByte(r))
			},
		),
	}
}
