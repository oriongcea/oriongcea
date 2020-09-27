package types

import (
	"github.com/gogo/protobuf/grpc"

	client "github.com/ocea/sdk/extra/ibc/02-client"
	clienttypes "github.com/ocea/sdk/extra/ibc/02-client/types"
	connection "github.com/ocea/sdk/extra/ibc/03-connection"
	connectiontypes "github.com/ocea/sdk/extra/ibc/03-connection/types"
	channel "github.com/ocea/sdk/extra/ibc/04-channel"
	channeltypes "github.com/ocea/sdk/extra/ibc/04-channel/types"
)

// QueryServer defines the IBC interfaces that the gRPC query server must implement
type QueryServer interface {
	clienttypes.QueryServer
	connectiontypes.QueryServer
	channeltypes.QueryServer
}

// RegisterQueryService registers each individual IBC submodule query service
func RegisterQueryService(server grpc.Server, queryService QueryServer) {
	client.RegisterQueryService(server, queryService)
	connection.RegisterQueryService(server, queryService)
	channel.RegisterQueryService(server, queryService)
}