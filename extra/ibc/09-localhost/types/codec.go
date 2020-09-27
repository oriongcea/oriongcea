package types

import (
	"github.com/ocea/sdk/codec"
	codectypes "github.com/ocea/sdk/codec/types"
	"github.com/ocea/sdk/extra/ibc/exported"
)

// RegisterInterfaces register the ibc interfaces submodule implementations to protobuf
// Any.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*exported.ClientState)(nil),
		&ClientState{},
	)
}

var (
	// SubModuleCdc references the global extra/ibc/09-localhost module codec.
	// The actual codec used for serialization should be provided to extra/ibc/09-localhost and
	// defined at the application level.
	SubModuleCdc = codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
)
