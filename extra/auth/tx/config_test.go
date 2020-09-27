package tx

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/ocea/sdk/codec"
	codectypes "github.com/ocea/sdk/codec/types"
	"github.com/ocea/sdk/std"
	"github.com/ocea/sdk/testutil/testdata"
	sdk "github.com/ocea/sdk/types"
	"github.com/ocea/sdk/extra/auth/testutil"
)

func TestGenerator(t *testing.T) {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	interfaceRegistry.RegisterImplementations((*sdk.Msg)(nil), &testdata.TestMsg{})
	protoCodec := codec.NewProtoCodec(interfaceRegistry)
	suite.Run(t, testutil.NewTxConfigTestSuite(NewTxConfig(protoCodec, DefaultSignModes)))
}
