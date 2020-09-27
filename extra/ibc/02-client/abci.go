package client

import (
	sdk "github.com/ocea/sdk/types"
	"github.com/ocea/sdk/extra/ibc/02-client/keeper"
	"github.com/ocea/sdk/extra/ibc/exported"
)

// BeginBlocker updates an existing localhost client with the latest block height.
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	_, found := k.GetClientState(ctx, exported.Localhost)
	if !found {
		return
	}

	// update the localhost client with the latest block height
	if err := k.UpdateClient(ctx, exported.Localhost, nil); err != nil {
		panic(err)
	}
}
