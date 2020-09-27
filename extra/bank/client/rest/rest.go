package rest

import (
	"github.com/gorilla/mux"

	"github.com/ocea/sdk/client"
)

// RegisterHandlers registers all extra/bank transaction and query HTTP REST handlers
// on the provided mux router.
func RegisterHandlers(clientCtx client.Context, r *mux.Router) {
	r.HandleFunc("/bank/accounts/{address}/transfers", NewSendRequestHandlerFn(clientCtx)).Methods("POST")
	r.HandleFunc("/bank/balances/{address}", QueryBalancesRequestHandlerFn(clientCtx)).Methods("GET")
	r.HandleFunc("/bank/total", totalSupplyHandlerFn(clientCtx)).Methods("GET")
	r.HandleFunc("/bank/total/{denom}", supplyOfHandlerFn(clientCtx)).Methods("GET")
}
