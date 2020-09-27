package crisis

import (
	"time"

	"github.com/ocea/sdk/telemetry"
	sdk "github.com/ocea/sdk/types"
	"github.com/ocea/sdk/extra/crisis/keeper"
	"github.com/ocea/sdk/extra/crisis/types"
)

// check all registered invariants
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	if k.InvCheckPeriod() == 0 || ctx.BlockHeight()%int64(k.InvCheckPeriod()) != 0 {
		// skip running the invariant check
		return
	}
	k.AssertInvariants(ctx)
}
