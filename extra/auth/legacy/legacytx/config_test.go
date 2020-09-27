package legacytx_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/ocea/sdk/codec"
	cryptoAmino "github.com/ocea/sdk/crypto/codec"
	"github.com/ocea/sdk/testutil/testdata"
	sdk "github.com/ocea/sdk/types"
	"github.com/ocea/sdk/extra/auth/legacy/legacytx"
	"github.com/ocea/sdk/extra/auth/testutil"
)

func testCodec() *codec.LegacyAmino {
	cdc := codec.NewLegacyAmino()
	sdk.RegisterLegacyAminoCodec(cdc)
	cryptoAmino.RegisterCrypto(cdc)
	cdc.RegisterConcrete(&testdata.TestMsg{}, "ocea-sdk/Test", nil)
	return cdc
}

func TestStdTxConfig(t *testing.T) {
	cdc := testCodec()
	txGen := legacytx.StdTxConfig{Cdc: cdc}
	suite.Run(t, testutil.NewTxConfigTestSuite(txGen))
}
