package rest

import (
	"github.com/gorilla/mux"

	"github.com/ocea/sdk/client"
)

// RegisterRoutes registers minting module REST handlers on the provided router.
func RegisterRoutes(clientCtx client.Context, r *mux.Router) {
	registerQueryRoutes(clientCtx, r)
}
