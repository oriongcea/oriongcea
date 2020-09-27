package client

import (
	"github.com/ocea/sdk/extra/distribution/client/cli"
	"github.com/ocea/sdk/extra/distribution/client/rest"
	govclient "github.com/ocea/sdk/extra/gov/client"
)

// ProposalHandler is the community spend proposal handler.
var (
	ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
)
