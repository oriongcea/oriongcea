package types

import (
	"github.com/ocea/sdk/codec"
	codectypes "github.com/ocea/sdk/codec/types"
	sdk "github.com/ocea/sdk/types"
	"github.com/ocea/sdk/extra/ibc/exported"
)

// RegisterInterfaces register the ibc channel submodule interfaces to protobuf
// Any.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterInterface(
		"ocea_sdk.ibc.v1.channel.ChannelI",
		(*exported.ChannelI)(nil),
	)
	registry.RegisterInterface(
		"ocea_sdk.ibc.v1.channel.CounterpartyChannelI",
		(*exported.CounterpartyChannelI)(nil),
	)
	registry.RegisterInterface(
		"ocea_sdk.ibc.v1.channel.PacketI",
		(*exported.PacketI)(nil),
	)
	registry.RegisterImplementations(
		(*exported.ChannelI)(nil),
		&Channel{},
	)
	registry.RegisterImplementations(
		(*exported.CounterpartyChannelI)(nil),
		&Counterparty{},
	)
	registry.RegisterImplementations(
		(*exported.PacketI)(nil),
		&Packet{},
	)
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgChannelOpenInit{},
		&MsgChannelOpenTry{},
		&MsgChannelOpenAck{},
		&MsgChannelOpenConfirm{},
		&MsgChannelCloseInit{},
		&MsgChannelCloseConfirm{},
		&MsgRecvPacket{},
		&MsgAcknowledgement{},
		&MsgTimeout{},
		&MsgTimeoutOnClose{},
	)
}

// SubModuleCdc references the global extra/ibc/04-channel module codec. Note, the codec should
// ONLY be used in certain instances of tests and for JSON encoding.
//
// The actual codec used for serialization should be provided to extra/ibc/04-channel and
// defined at the application level.
var SubModuleCdc = codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
