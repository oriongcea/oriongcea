package v040

import (
	v039auth "github.com/ocea/sdk/extra/auth/legacy/v0_39"
	v038bank "github.com/ocea/sdk/extra/bank/legacy/v0_38"
)

// Migrate accepts exported v0.39 extra/auth and v0.38 extra/bank genesis state and
// migrates it to v0.40 extra/bank genesis state. The migration includes:
//
// - Moving balances from extra/auth to extra/bank genesis state.
func Migrate(bankGenState v038bank.GenesisState, authGenState v039auth.GenesisState) GenesisState {
	balances := make([]Balance, len(authGenState.Accounts))
	for i, acc := range authGenState.Accounts {
		balances[i] = Balance{
			Address: acc.GetAddress(),
			Coins:   acc.GetCoins(),
		}
	}

	return NewGenesisState(bankGenState.SendEnabled, balances)
}
