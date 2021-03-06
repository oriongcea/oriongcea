package keeper

import (
	"bytes"
	"fmt"
	"reflect"

	sdk "github.com/ocea/sdk/types"
	sdkerrors "github.com/ocea/sdk/types/errors"
	clienttypes "github.com/ocea/sdk/extra/ibc/02-client/types"
	"github.com/ocea/sdk/extra/ibc/03-connection/types"
	commitmenttypes "github.com/ocea/sdk/extra/ibc/23-commitment/types"
	"github.com/ocea/sdk/extra/ibc/exported"
)

// ConnOpenInit initialises a connection attempt on chain A.
//
// NOTE: Identifiers are checked on msg validation.
func (k Keeper) ConnOpenInit(
	ctx sdk.Context,
	connectionID, // identifier
	clientID string,
	counterparty types.Counterparty, // desiredCounterpartyConnectionIdentifier, counterpartyPrefix, counterpartyClientIdentifier
) error {
	_, found := k.GetConnection(ctx, connectionID)
	if found {
		return types.ErrConnectionExists
	}

	// connection defines chain A's ConnectionEnd
	connection := types.NewConnectionEnd(types.INIT, clientID, counterparty, types.GetCompatibleEncodedVersions())
	k.SetConnection(ctx, connectionID, connection)

	if err := k.addConnectionToClient(ctx, clientID, connectionID); err != nil {
		return err
	}

	k.Logger(ctx).Info(fmt.Sprintf("connection %s state updated: NONE -> INIT", connectionID))
	return nil
}

// ConnOpenTry relays notice of a connection attempt on chain A to chain B (this
// code is executed on chain B).
//
// NOTE:
//  - Here chain A acts as the counterparty
//  - Identifiers are checked on msg validation
func (k Keeper) ConnOpenTry(
	ctx sdk.Context,
	connectionID string, // desiredIdentifier
	counterparty types.Counterparty, // counterpartyConnectionIdentifier, counterpartyPrefix and counterpartyClientIdentifier
	clientID string, // clientID of chainA
	clientState exported.ClientState, // clientState that chainA has for chainB
	counterpartyVersions []string, // supported versions of chain A
	proofInit []byte, // proof that chainA stored connectionEnd in state (on ConnOpenInit)
	proofClient []byte, // proof that chainA stored a light client of chainB
	proofConsensus []byte, // proof that chainA stored chainB's consensus state at consensus height
	proofHeight exported.Height, // height at which relayer constructs proof of A storing connectionEnd in state
	consensusHeight exported.Height, // latest height of chain B which chain A has stored in its chain B client
) error {
	selfHeight := clienttypes.GetSelfHeight(ctx)
	if consensusHeight.GTE(selfHeight) {
		return sdkerrors.Wrapf(
			sdkerrors.ErrInvalidHeight,
			"consensus height is greater than or equal to the current block height (%s >= %s)", consensusHeight, selfHeight,
		)
	}

	// validate client parameters of a chainB client stored on chainA
	if err := k.clientKeeper.ValidateSelfClient(ctx, clientState); err != nil {
		return err
	}

	expectedConsensusState, found := k.clientKeeper.GetSelfConsensusState(ctx, consensusHeight)
	if !found {
		return clienttypes.ErrSelfConsensusStateNotFound
	}

	// expectedConnection defines Chain A's ConnectionEnd
	// NOTE: chain A's counterparty is chain B (i.e where this code is executed)
	prefix := k.GetCommitmentPrefix()
	expectedCounterparty := types.NewCounterparty(clientID, connectionID, commitmenttypes.NewMerklePrefix(prefix.Bytes()))
	expectedConnection := types.NewConnectionEnd(types.INIT, counterparty.ClientId, expectedCounterparty, counterpartyVersions)

	// chain B picks a version from Chain A's available versions that is compatible
	// with Chain B's supported IBC versions. PickVersion will select the intersection
	// of the supported versions and the counterparty versions.
	version, err := types.PickVersion(counterpartyVersions)
	if err != nil {
		return err
	}

	// connection defines chain B's ConnectionEnd
	connection := types.NewConnectionEnd(types.TRYOPEN, clientID, counterparty, []string{version})

	// Check that ChainA committed expectedConnectionEnd to its state
	if err := k.VerifyConnectionState(
		ctx, connection, proofHeight, proofInit, counterparty.ConnectionId,
		expectedConnection,
	); err != nil {
		return err
	}

	// Check that ChainA stored the clientState provided in the msg
	if err := k.VerifyClientState(ctx, connection, proofHeight, proofClient, clientState); err != nil {
		return err
	}

	// Check that ChainA stored the correct ConsensusState of chainB at the given consensusHeight
	if err := k.VerifyClientConsensusState(
		ctx, connection, proofHeight, consensusHeight, proofConsensus, expectedConsensusState,
	); err != nil {
		return err
	}

	// If connection already exists for connectionID, ensure that the existing connection's
	// counterparty is chainA and connection is on INIT stage.
	// Check that existing connection versions for initialized connection is equal to compatible
	// versions for this chain.
	previousConnection, found := k.GetConnection(ctx, connectionID)
	if found && !(previousConnection.State == types.INIT &&
		previousConnection.Counterparty.ConnectionId == counterparty.ConnectionId &&
		bytes.Equal(previousConnection.Counterparty.Prefix.Bytes(), counterparty.Prefix.Bytes()) &&
		previousConnection.ClientId == clientID &&
		previousConnection.Counterparty.ClientId == counterparty.ClientId &&
		reflect.DeepEqual(previousConnection.Versions, types.GetCompatibleEncodedVersions())) {
		return sdkerrors.Wrap(types.ErrInvalidConnection, "cannot relay connection attempt")
	}

	// store connection in chainB state
	if err := k.addConnectionToClient(ctx, clientID, connectionID); err != nil {
		return sdkerrors.Wrapf(err, "failed to add connection with ID %s to client with ID %s", connectionID, clientID)
	}

	k.SetConnection(ctx, connectionID, connection)
	k.Logger(ctx).Info(fmt.Sprintf("connection %s state updated: %s -> TRYOPEN ", connectionID, previousConnection.State))
	return nil
}

// ConnOpenAck relays acceptance of a connection open attempt from chain B back
// to chain A (this code is executed on chain A).
//
// NOTE: Identifiers are checked on msg validation.
func (k Keeper) ConnOpenAck(
	ctx sdk.Context,
	connectionID string,
	clientState exported.ClientState, // client state for chainA on chainB
	encodedVersion string, // version that ChainB chose in ConnOpenTry
	proofTry []byte, // proof that connectionEnd was added to ChainB state in ConnOpenTry
	proofClient []byte, // proof of client state on chainB for chainA
	proofConsensus []byte, // proof that chainB has stored ConsensusState of chainA on its client
	proofHeight exported.Height, // height that relayer constructed proofTry
	consensusHeight exported.Height, // latest height of chainA that chainB has stored on its chainA client
) error {
	// Check that chainB client hasn't stored invalid height
	selfHeight := clienttypes.GetSelfHeight(ctx)
	if consensusHeight.GTE(selfHeight) {
		return sdkerrors.Wrapf(
			sdkerrors.ErrInvalidHeight,
			"consensus height is greater than or equal to the current block height (%s >= %s)", consensusHeight, selfHeight,
		)
	}

	// Retrieve connection
	connection, found := k.GetConnection(ctx, connectionID)
	if !found {
		return sdkerrors.Wrap(types.ErrConnectionNotFound, connectionID)
	}

	// Verify the provided version against the previously set connection state
	switch {
	// connection on ChainA must be in INIT or TRYOPEN
	case connection.State != types.INIT && connection.State != types.TRYOPEN:
		return sdkerrors.Wrapf(
			types.ErrInvalidConnectionState,
			"connection state is not INIT or TRYOPEN (got %s)", connection.State.String(),
		)

	// if the connection is INIT then the provided version must be supproted
	case connection.State == types.INIT && !types.IsSupportedVersion(encodedVersion):
		return sdkerrors.Wrapf(
			types.ErrInvalidConnectionState,
			"connection state is in INIT but the provided encoded version is not supported %s", encodedVersion,
		)

	// if the connection is in TRYOPEN then the encoded version must be the only set version in the
	// retreived connection state.
	case connection.State == types.TRYOPEN && (len(connection.Versions) != 1 || connection.Versions[0] != encodedVersion):
		return sdkerrors.Wrapf(
			types.ErrInvalidConnectionState,
			"connection state is in TRYOPEN but the provided encoded version (%s) is not set in the previous connection %s", encodedVersion, connection,
		)
	}

	// validate client parameters of a chainA client stored on chainB
	if err := k.clientKeeper.ValidateSelfClient(ctx, clientState); err != nil {
		return err
	}

	// Retrieve chainA's consensus state at consensusheight
	expectedConsensusState, found := k.clientKeeper.GetSelfConsensusState(ctx, consensusHeight)
	if !found {
		return clienttypes.ErrSelfConsensusStateNotFound
	}

	prefix := k.GetCommitmentPrefix()
	expectedCounterparty := types.NewCounterparty(connection.ClientId, connectionID, commitmenttypes.NewMerklePrefix(prefix.Bytes()))
	expectedConnection := types.NewConnectionEnd(types.TRYOPEN, connection.Counterparty.ClientId, expectedCounterparty, []string{encodedVersion})

	// Ensure that ChainB stored expected connectionEnd in its state during ConnOpenTry
	if err := k.VerifyConnectionState(
		ctx, connection, proofHeight, proofTry, connection.Counterparty.ConnectionId,
		expectedConnection,
	); err != nil {
		return err
	}

	// Check that ChainB stored the clientState provided in the msg
	if err := k.VerifyClientState(ctx, connection, proofHeight, proofClient, clientState); err != nil {
		return err
	}

	// Ensure that ChainB has stored the correct ConsensusState for chainA at the consensusHeight
	if err := k.VerifyClientConsensusState(
		ctx, connection, proofHeight, consensusHeight, proofConsensus, expectedConsensusState,
	); err != nil {
		return err
	}

	k.Logger(ctx).Info(fmt.Sprintf("connection %s state updated: %s -> OPEN ", connectionID, connection.State))

	// Update connection state to Open
	connection.State = types.OPEN
	connection.Versions = []string{encodedVersion}
	k.SetConnection(ctx, connectionID, connection)
	return nil
}

// ConnOpenConfirm confirms opening of a connection on chain A to chain B, after
// which the connection is open on both chains (this code is executed on chain B).
//
// NOTE: Identifiers are checked on msg validation.
func (k Keeper) ConnOpenConfirm(
	ctx sdk.Context,
	connectionID string,
	proofAck []byte, // proof that connection opened on ChainA during ConnOpenAck
	proofHeight exported.Height, // height that relayer constructed proofAck
) error {
	// Retrieve connection
	connection, found := k.GetConnection(ctx, connectionID)
	if !found {
		return sdkerrors.Wrap(types.ErrConnectionNotFound, connectionID)
	}

	// Check that connection state on ChainB is on state: TRYOPEN
	if connection.State != types.TRYOPEN {
		return sdkerrors.Wrapf(
			types.ErrInvalidConnectionState,
			"connection state is not TRYOPEN (got %s)", connection.State.String(),
		)
	}

	prefix := k.GetCommitmentPrefix()
	expectedCounterparty := types.NewCounterparty(connection.ClientId, connectionID, commitmenttypes.NewMerklePrefix(prefix.Bytes()))
	expectedConnection := types.NewConnectionEnd(types.OPEN, connection.Counterparty.ClientId, expectedCounterparty, connection.Versions)

	// Check that connection on ChainA is open
	if err := k.VerifyConnectionState(
		ctx, connection, proofHeight, proofAck, connection.Counterparty.ConnectionId,
		expectedConnection,
	); err != nil {
		return err
	}

	// Update ChainB's connection to Open
	connection.State = types.OPEN
	k.SetConnection(ctx, connectionID, connection)
	k.Logger(ctx).Info(fmt.Sprintf("connection %s state updated: TRYOPEN -> OPEN ", connectionID))
	return nil
}
