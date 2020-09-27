package types

import (
	"github.com/ocea/sdk/codec"
	sdk "github.com/ocea/sdk/types"
	sdkerrors "github.com/ocea/sdk/types/errors"
	clienttypes "github.com/ocea/sdk/extra/ibc/02-client/types"
	host "github.com/ocea/sdk/extra/ibc/24-host"
	"github.com/ocea/sdk/extra/ibc/exported"
)

// GetConsensusState retrieves the consensus state from the client prefixed
// store. An error is returned if the consensus state does not exist.
func GetConsensusState(store sdk.KVStore, cdc codec.BinaryMarshaler, height exported.Height) (*ConsensusState, error) {
	bz := store.Get(host.KeyConsensusState(height))
	if bz == nil {
		return nil, sdkerrors.Wrapf(
			clienttypes.ErrConsensusStateNotFound,
			"consensus state does not exist for height %d", height,
		)
	}

	var consensusStateI exported.ConsensusState
	if err := codec.UnmarshalAny(cdc, &consensusStateI, bz); err != nil {
		return nil, sdkerrors.Wrapf(clienttypes.ErrInvalidConsensus, "unmarshal error: %v", err)
	}

	consensusState, ok := consensusStateI.(*ConsensusState)
	if !ok {
		return nil, sdkerrors.Wrapf(
			clienttypes.ErrInvalidConsensus,
			"invalid consensus type %T, expected %T", consensusState, &ConsensusState{},
		)
	}

	return consensusState, nil
}
