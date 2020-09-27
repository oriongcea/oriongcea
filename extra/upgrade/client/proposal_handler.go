package client

import (
	govclient "github.com/ocea/sdk/extra/gov/client"
	"github.com/ocea/sdk/extra/upgrade/client/cli"
	"github.com/ocea/sdk/extra/upgrade/client/rest"
)

var ProposalHandler = govclient.NewProposalHandler(cli.NewCmdSubmitUpgradeProposal, rest.ProposalRESTHandler)
var CancelProposalHandler = govclient.NewProposalHandler(cli.NewCmdSubmitCancelUpgradeProposal, rest.ProposalCancelRESTHandler)
