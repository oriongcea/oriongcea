package client

import (
	govclient "github.com/ocea/sdk/extra/gov/client"
	"github.com/ocea/sdk/extra/params/client/cli"
	"github.com/ocea/sdk/extra/params/client/rest"
)

// ProposalHandler is the param change proposal handler.
var ProposalHandler = govclient.NewProposalHandler(cli.NewSubmitParamChangeProposalTxCmd, rest.ProposalRESTHandler)
