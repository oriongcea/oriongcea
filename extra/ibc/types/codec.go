package types

import (
	codectypes "github.com/ocea/sdk/codec/types"
	clienttypes "github.com/ocea/sdk/extra/ibc/02-client/types"
	connectiontypes "github.com/ocea/sdk/extra/ibc/03-connection/types"
	channeltypes "github.com/ocea/sdk/extra/ibc/04-channel/types"
	ibctmtypes "github.com/ocea/sdk/extra/ibc/07-tendermint/types"
	localhosttypes "github.com/ocea/sdk/extra/ibc/09-localhost/types"
	commitmenttypes "github.com/ocea/sdk/extra/ibc/23-commitment/types"
	solomachinetypes "github.com/ocea/sdk/extra/ibc/light-clients/solomachine/types"
)

// RegisterInterfaces registers extra/ibc interfaces into protobuf Any.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	clienttypes.RegisterInterfaces(registry)
	connectiontypes.RegisterInterfaces(registry)
	channeltypes.RegisterInterfaces(registry)
	solomachinetypes.RegisterInterfaces(registry)
	ibctmtypes.RegisterInterfaces(registry)
	localhosttypes.RegisterInterfaces(registry)
	commitmenttypes.RegisterInterfaces(registry)
}
