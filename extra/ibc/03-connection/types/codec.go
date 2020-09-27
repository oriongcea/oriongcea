package types

import (
	"github.com/ocea/sdk/codec"
	codectypes "github.com/ocea/sdk/codec/types"
	sdk "github.com/ocea/sdk/types"
	"github.com/ocea/sdk/extra/ibc/exported"
)

// RegisterInterfaces register the ibc interfaces submodule implementations to protobuf
// Any.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterInterface(
		"ocea_sdk.ibc.v1.connection.ConnectionI",
		(*exported.ConnectionI)(nil),
	)
	registry.RegisterInterface(
		"ocea_sdk.ibc.v1.connection.CounterpartyConnectionI",
		(*exported.CounterpartyConnectionI)(nil),
	)
	registry.RegisterImplementations(
		(*exported.ConnectionI)(nil),
		&ConnectionEnd{},
	)
	registry.RegisterImplementations(
		(*exported.CounterpartyConnectionI)(nil),
		&Counterparty{},
	)
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgConnectionOpenInit{},
		&MsgConnectionOpenTry{},
		&MsgConnectionOpenAck{},
		&MsgConnectionOpenConfirm{},
	)
}

var (
	// SubModuleCdc references the global extra/ibc/03-connection module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding.
	//
	// The actual codec used for serialization should be provided to extra/ibc/03-connectionl and
	// defined at the application level.
	SubModuleCdc = codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
)
