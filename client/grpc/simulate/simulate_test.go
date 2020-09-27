package simulate_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/ocea/sdk/baseapp"
	"github.com/ocea/sdk/client"
	"github.com/ocea/sdk/client/grpc/simulate"
	"github.com/ocea/sdk/client/tx"
	codectypes "github.com/ocea/sdk/codec/types"
	"github.com/ocea/sdk/simapp"
	"github.com/ocea/sdk/testutil/testdata"
	sdk "github.com/ocea/sdk/types"
	txtypes "github.com/ocea/sdk/types/tx"
	"github.com/ocea/sdk/types/tx/signing"
	authsigning "github.com/ocea/sdk/extra/auth/signing"
	authtypes "github.com/ocea/sdk/extra/auth/types"
	banktypes "github.com/ocea/sdk/extra/bank/types"
)

type IntegrationTestSuite struct {
	suite.Suite

	app         *simapp.SimApp
	clientCtx   client.Context
	queryClient simulate.SimulateServiceClient
	sdkCtx      sdk.Context
}

func (s *IntegrationTestSuite) SetupSuite() {
	app := simapp.Setup(true)
	sdkCtx := app.BaseApp.NewContext(true, tmproto.Header{})

	app.AccountKeeper.SetParams(sdkCtx, authtypes.DefaultParams())
	app.BankKeeper.SetParams(sdkCtx, banktypes.DefaultParams())

	// Set up TxConfig.
	encodingConfig := simapp.MakeEncodingConfig()
	clientCtx := client.Context{}.WithTxConfig(encodingConfig.TxConfig)

	// Create new simulation server.
	srv := simulate.NewSimulateServer(app.BaseApp.Simulate, encodingConfig.InterfaceRegistry)

	queryHelper := baseapp.NewQueryServerTestHelper(sdkCtx, app.InterfaceRegistry())
	simulate.RegisterSimulateServiceServer(queryHelper, srv)
	queryClient := simulate.NewSimulateServiceClient(queryHelper)

	s.app = app
	s.clientCtx = clientCtx
	s.queryClient = queryClient
	s.sdkCtx = sdkCtx
}

func (s IntegrationTestSuite) TestSimulateService() {
	// Create an account with some funds.
	priv1, _, addr1 := testdata.KeyTestPubAddr()
	acc1 := s.app.AccountKeeper.NewAccountWithAddress(s.sdkCtx, addr1)
	err := acc1.SetAccountNumber(0)
	s.Require().NoError(err)
	s.app.AccountKeeper.SetAccount(s.sdkCtx, acc1)
	s.app.BankKeeper.SetBalances(s.sdkCtx, addr1, sdk.Coins{
		sdk.NewInt64Coin("atom", 10000000),
	})

	// Create a test extra/bank MsgSend.
	coins := sdk.NewCoins(sdk.NewInt64Coin("atom", 10))
	_, _, addr2 := testdata.KeyTestPubAddr()
	msg := banktypes.NewMsgSend(addr1, addr2, coins)
	feeAmount := testdata.NewTestFeeAmount()
	gasLimit := testdata.NewTestGasLimit()
	memo := "foo"
	accSeq, accNum := uint64(0), uint64(0)

	// Create a txBuilder.
	txBuilder := s.clientCtx.TxConfig.NewTxBuilder()
	txBuilder.SetMsgs(msg)
	txBuilder.SetMemo(memo)
	txBuilder.SetFeeAmount(feeAmount)
	txBuilder.SetGasLimit(gasLimit)
	// 1st round: set empty signature
	sigV2 := signing.SignatureV2{
		PubKey: priv1.PubKey(),
		Data: &signing.SingleSignatureData{
			SignMode:  s.clientCtx.TxConfig.SignModeHandler().DefaultMode(),
			Signature: nil,
		},
	}
	txBuilder.SetSignatures(sigV2)
	// 2nd round: actually sign
	sigV2, err = tx.SignWithPrivKey(
		s.clientCtx.TxConfig.SignModeHandler().DefaultMode(),
		authsigning.SignerData{ChainID: s.sdkCtx.ChainID(), AccountNumber: accNum, Sequence: accSeq},
		txBuilder, priv1, s.clientCtx.TxConfig, accSeq,
	)
	txBuilder.SetSignatures(sigV2)

	any, ok := txBuilder.(codectypes.IntoAny)
	s.Require().True(ok)
	cached := any.AsAny().GetCachedValue()
	txTx, ok := cached.(*txtypes.Tx)
	s.Require().True(ok)
	res, err := s.queryClient.Simulate(
		context.Background(),
		&simulate.SimulateRequest{Tx: txTx},
	)
	s.Require().NoError(err)

	// Check the result and gas used are correct.
	s.Require().Equal(len(res.GetResult().GetEvents()), 4) // 1 transfer, 3 messages.
	s.Require().True(res.GetGasInfo().GetGasUsed() > 0)    // Gas used sometimes change, just check it's not empty.
}

func TestSimulateTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
