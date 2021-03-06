package reflection_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/ocea/sdk/baseapp"
	"github.com/ocea/sdk/client/grpc/reflection"
	"github.com/ocea/sdk/simapp"
)

type IntegrationTestSuite struct {
	suite.Suite

	queryClient reflection.ReflectionServiceClient
}

func (s *IntegrationTestSuite) SetupSuite() {
	app := simapp.Setup(false)

	srv := reflection.NewReflectionServiceServer(app.InterfaceRegistry())

	sdkCtx := app.BaseApp.NewContext(false, tmproto.Header{})
	queryHelper := baseapp.NewQueryServerTestHelper(sdkCtx, app.InterfaceRegistry())

	reflection.RegisterReflectionServiceServer(queryHelper, srv)
	queryClient := reflection.NewReflectionServiceClient(queryHelper)

	s.queryClient = queryClient
}

func (s IntegrationTestSuite) TestSimulateService() {
	// We will test the following interface for testing.
	var iface = "ocea.evidence.v1beta1.Evidence"

	// Test that "ocea.evidence.v1beta1.Evidence" is included in the
	// interfaces.
	resIface, err := s.queryClient.ListAllInterfaces(
		context.Background(),
		&reflection.ListAllInterfacesRequest{},
	)
	s.Require().NoError(err)
	s.Require().Contains(resIface.GetInterfaceNames(), iface)

	// Test that "ocea.evidence.v1beta1.Evidence" has at least the
	// Equivocation implementations.
	resImpl, err := s.queryClient.ListImplementations(
		context.Background(),
		&reflection.ListImplementationsRequest{InterfaceName: iface},
	)
	s.Require().NoError(err)
	s.Require().Contains(resImpl.GetImplementationMessageNames(), "/ocea.evidence.v1beta1.Equivocation")
}

func TestSimulateTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
