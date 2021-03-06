package keeper

import (
	"fmt"

	sdk "github.com/ocea/sdk/types"
	"github.com/ocea/sdk/extra/bank/types"
)

// InitGenesis initializes the bank module's state from a given genesis state.
func (k BaseKeeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) {
	k.SetParams(ctx, genState.Params)

	var totalSupply sdk.Coins

	genState.Balances = types.SanitizeGenesisBalances(genState.Balances)
	for _, balance := range genState.Balances {
		addr, err := sdk.AccAddressFromBech32(balance.Address)
		if err != nil {
			panic(err)
		}

		if err := k.ValidateBalance(ctx, addr); err != nil {
			panic(err)
		}

		if err := k.SetBalances(ctx, addr, balance.Coins); err != nil {
			panic(fmt.Errorf("error on setting balances %w", err))
		}

		totalSupply = totalSupply.Add(balance.Coins...)
	}

	if genState.Supply.Empty() {
		genState.Supply = totalSupply
	}

	k.SetSupply(ctx, types.NewSupply(genState.Supply))
}

// ExportGenesis returns the bank module's genesis state.
func (k BaseKeeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return types.NewGenesisState(
		k.GetParams(ctx),
		k.GetAccountsBalances(ctx),
		k.GetSupply(ctx).GetTotal(),
		k.GetAllDenomMetaData(ctx),
	)
}
