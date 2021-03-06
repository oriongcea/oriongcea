package types

import (
	"errors"
	"fmt"

	host "github.com/ocea/sdk/extra/ibc/24-host"
)

// NewPacketAckCommitment creates a new PacketAckCommitment instance.
func NewPacketAckCommitment(portID, channelID string, seq uint64, hash []byte) PacketAckCommitment {
	return PacketAckCommitment{
		PortId:    portID,
		ChannelId: channelID,
		Sequence:  seq,
		Hash:      hash,
	}
}

// Validate performs basic validation of fields returning an error upon any
// failure.
func (pa PacketAckCommitment) Validate() error {
	if len(pa.Hash) == 0 {
		return errors.New("hash bytes cannot be empty")
	}
	return validateGenFields(pa.PortId, pa.ChannelId, pa.Sequence)
}

// NewPacketSequence creates a new PacketSequences instance.
func NewPacketSequence(portID, channelID string, seq uint64) PacketSequence {
	return PacketSequence{
		PortId:    portID,
		ChannelId: channelID,
		Sequence:  seq,
	}
}

// Validate performs basic validation of fields returning an error upon any
// failure.
func (ps PacketSequence) Validate() error {
	return validateGenFields(ps.PortId, ps.ChannelId, ps.Sequence)
}

// NewGenesisState creates a GenesisState instance.
func NewGenesisState(
	channels []IdentifiedChannel, acks, commitments []PacketAckCommitment,
	sendSeqs, recvSeqs, ackSeqs []PacketSequence,
) GenesisState {
	return GenesisState{
		Channels:         channels,
		Acknowledgements: acks,
		Commitments:      commitments,
		SendSequences:    sendSeqs,
		RecvSequences:    recvSeqs,
		AckSequences:     ackSeqs,
	}
}

// DefaultGenesisState returns the ibc channel submodule's default genesis state.
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Channels:         []IdentifiedChannel{},
		Acknowledgements: []PacketAckCommitment{},
		Commitments:      []PacketAckCommitment{},
		SendSequences:    []PacketSequence{},
		RecvSequences:    []PacketSequence{},
		AckSequences:     []PacketSequence{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	for i, channel := range gs.Channels {
		if err := channel.ValidateBasic(); err != nil {
			return fmt.Errorf("invalid channel %v channel index %d: %w", channel, i, err)
		}
	}

	for i, ack := range gs.Acknowledgements {
		if err := ack.Validate(); err != nil {
			return fmt.Errorf("invalid acknowledgement %v ack index %d: %w", ack, i, err)
		}
	}

	for i, commitment := range gs.Commitments {
		if err := commitment.Validate(); err != nil {
			return fmt.Errorf("invalid commitment %v index %d: %w", commitment, i, err)
		}
	}

	for i, ss := range gs.SendSequences {
		if err := ss.Validate(); err != nil {
			return fmt.Errorf("invalid send sequence %v index %d: %w", ss, i, err)
		}
	}

	for i, rs := range gs.RecvSequences {
		if err := rs.Validate(); err != nil {
			return fmt.Errorf("invalid receive sequence %v index %d: %w", rs, i, err)
		}
	}

	for i, as := range gs.AckSequences {
		if err := as.Validate(); err != nil {
			return fmt.Errorf("invalid acknowledgement sequence %v index %d: %w", as, i, err)
		}
	}

	return nil
}

func validateGenFields(portID, channelID string, sequence uint64) error {
	if err := host.PortIdentifierValidator(portID); err != nil {
		return fmt.Errorf("invalid port Id: %w", err)
	}
	if err := host.ChannelIdentifierValidator(channelID); err != nil {
		return fmt.Errorf("invalid channel Id: %w", err)
	}
	if sequence == 0 {
		return errors.New("sequence cannot be 0")
	}
	return nil
}
