package slashing

import (
	sdk "github.com/ocea/sdk/types"
	sdkerrors "github.com/ocea/sdk/types/errors"
	"github.com/ocea/sdk/extra/slashing/keeper"
	"github.com/ocea/sdk/extra/slashing/types"
)

// NewHandler creates an sdk.Handler for all the slashing type messages
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgUnjail:
			return handleMsgUnjail(ctx, msg, k)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
		}
	}
}

// Validators must submit a transaction to unjail itself after
// having been jailed (and thus unbonded) for downtime
func handleMsgUnjail(ctx sdk.Context, msg *types.MsgUnjail, k keeper.Keeper) (*sdk.Result, error) {
	valAddr, valErr := sdk.ValAddressFromBech32(msg.ValidatorAddr)
	if valErr != nil {
		return nil, valErr
	}
	err := k.Unjail(ctx, valAddr)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.ValidatorAddr),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}
