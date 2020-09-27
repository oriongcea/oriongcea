package rest

import (
	"github.com/gorilla/mux"

	"github.com/ocea/sdk/client"
)

func RegisterHandlers(clientCtx client.Context, r *mux.Router) {
	registerQueryRoutes(clientCtx, r)
	registerTxHandlers(clientCtx, r)
}
