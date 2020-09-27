package evidence

import (
	sdk "github.com/ocea/sdk/types"
	sdkerrors "github.com/ocea/sdk/types/errors"
	"github.com/ocea/sdk/extra/evidence/exported"
	"github.com/ocea/sdk/extra/evidence/keeper"
	"github.com/ocea/sdk/extra/evidence/types"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case exported.MsgSubmitEvidence:
			return handleMsgSubmitEvidence(ctx, k, msg)

		default:

			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
		}
	}
}

func handleMsgSubmitEvidence(ctx sdk.Context, k keeper.Keeper, msg exported.MsgSubmitEvidence) (*sdk.Result, error) {
	evidence := msg.GetEvidence()
	if err := k.SubmitEvidence(ctx, evidence); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.GetSubmitter().String()),
		),
	)

	return &sdk.Result{
		Data:   evidence.Hash(),
		Events: ctx.EventManager().ABCIEvents(),
	}, nil
}
