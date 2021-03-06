package cli_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/suite"
	tmcli "github.com/tendermint/tendermint/libs/cli"

	"github.com/ocea/sdk/client"
	"github.com/ocea/sdk/client/flags"
	clitestutil "github.com/ocea/sdk/testutil/cli"
	"github.com/ocea/sdk/testutil/network"
	sdk "github.com/ocea/sdk/types"
	sdkerrors "github.com/ocea/sdk/types/errors"
	"github.com/ocea/sdk/types/query"
	"github.com/ocea/sdk/extra/bank/client/cli"
	banktestutil "github.com/ocea/sdk/extra/bank/client/testutil"
	"github.com/ocea/sdk/extra/bank/types"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	cfg := network.DefaultConfig()
	cfg.NumValidators = 1

	s.cfg = cfg
	s.network = network.New(s.T(), cfg)

	_, err := s.network.WaitForHeight(1)
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func (s *IntegrationTestSuite) TestGetBalancesCmd() {
	val := s.network.Validators[0]

	testCases := []struct {
		name      string
		args      []string
		expectErr bool
		respType  proto.Message
		expected  proto.Message
	}{
		{"no address provided", []string{}, true, nil, nil},
		{
			"total account balance",
			[]string{
				val.Address.String(),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
				fmt.Sprintf("--%s=1", flags.FlagHeight),
			},
			false,
			&types.QueryAllBalancesResponse{},
			&types.QueryAllBalancesResponse{
				Balances: sdk.NewCoins(
					sdk.NewCoin(fmt.Sprintf("%stoken", val.Moniker), s.cfg.AccountTokens),
					sdk.NewCoin(s.cfg.BondDenom, s.cfg.StakingTokens.Sub(s.cfg.BondedTokens)),
				),
				Pagination: &query.PageResponse{},
			},
		},
		{
			"total account balance of a specific denom",
			[]string{
				val.Address.String(),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
				fmt.Sprintf("--%s=%s", cli.FlagDenom, s.cfg.BondDenom),
				fmt.Sprintf("--%s=1", flags.FlagHeight),
			},
			false,
			&sdk.Coin{},
			NewCoin(s.cfg.BondDenom, s.cfg.StakingTokens.Sub(s.cfg.BondedTokens)),
		},
		{
			"total account balance of a bogus denom",
			[]string{
				val.Address.String(),
				fmt.Sprintf("--%s=foobar", cli.FlagDenom),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			false,
			&sdk.Coin{},
			NewCoin("foobar", sdk.ZeroInt()),
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.GetBalancesCmd()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, tc.args)

			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().NoError(val.ClientCtx.JSONMarshaler.UnmarshalJSON(out.Bytes(), tc.respType))
				s.Require().Equal(tc.expected.String(), tc.respType.String())
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGetCmdQueryTotalSupply() {
	val := s.network.Validators[0]

	testCases := []struct {
		name      string
		args      []string
		expectErr bool
		respType  proto.Message
		expected  proto.Message
	}{
		{
			name: "total supply",
			args: []string{
				fmt.Sprintf("--%s=1", flags.FlagHeight),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			respType: &types.QueryTotalSupplyResponse{},
			expected: &types.QueryTotalSupplyResponse{
				Supply: sdk.NewCoins(
					sdk.NewCoin(fmt.Sprintf("%stoken", val.Moniker), s.cfg.AccountTokens),
					sdk.NewCoin(s.cfg.BondDenom, s.cfg.StakingTokens.Add(sdk.NewInt(10))),
				)},
		},
		{
			name: "total supply of a specific denomination",
			args: []string{
				fmt.Sprintf("--%s=1", flags.FlagHeight),
				fmt.Sprintf("--%s=%s", cli.FlagDenom, s.cfg.BondDenom),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			respType: &sdk.Coin{},
			expected: &sdk.Coin{s.cfg.BondDenom, s.cfg.StakingTokens.Add(sdk.NewInt(10))},
		},
		{
			name: "total supply of a bogus denom",
			args: []string{
				fmt.Sprintf("--%s=1", flags.FlagHeight),
				fmt.Sprintf("--%s=foobar", cli.FlagDenom),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			respType: &sdk.Coin{},
			expected: &sdk.Coin{"foobar", sdk.ZeroInt()},
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.GetCmdQueryTotalSupply()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(out.Bytes(), tc.respType))
				s.Require().Equal(tc.expected, tc.respType)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestNewSendTxCmdGenOnly() {
	val := s.network.Validators[0]

	clientCtx := val.ClientCtx

	ctx := context.Background()
	ctx = context.WithValue(ctx, client.ClientContextKey, &clientCtx)

	from := val.Address
	to := val.Address
	amount := sdk.NewCoins(
		sdk.NewCoin(fmt.Sprintf("%stoken", val.Moniker), sdk.NewInt(10)),
		sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10)),
	)
	args := []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
		fmt.Sprintf("--%s=true", flags.FlagGenerateOnly),
	}

	bz, err := banktestutil.MsgSendExec(clientCtx, from, to, amount, args...)
	s.Require().NoError(err)
	tx, err := s.cfg.TxConfig.TxJSONDecoder()(bz.Bytes())
	s.Require().NoError(err)
	s.Require().Equal([]sdk.Msg{types.NewMsgSend(from, to, amount)}, tx.GetMsgs())
}

func (s *IntegrationTestSuite) TestNewSendTxCmd() {
	val := s.network.Validators[0]

	testCases := []struct {
		name         string
		from, to     sdk.AccAddress
		amount       sdk.Coins
		args         []string
		expectErr    bool
		respType     proto.Message
		expectedCode uint32
	}{
		{
			"valid transaction",
			val.Address,
			val.Address,
			sdk.NewCoins(
				sdk.NewCoin(fmt.Sprintf("%stoken", val.Moniker), sdk.NewInt(10)),
				sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10)),
			),
			[]string{
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
			},
			false,
			&sdk.TxResponse{},
			0,
		},
		{
			"not enough fees",
			val.Address,
			val.Address,
			sdk.NewCoins(
				sdk.NewCoin(fmt.Sprintf("%stoken", val.Moniker), sdk.NewInt(10)),
				sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10)),
			),
			[]string{
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(1))).String()),
			},
			false,
			&sdk.TxResponse{},
			sdkerrors.ErrInsufficientFee.ABCICode(),
		},
		{
			"not enough gas",
			val.Address,
			val.Address,
			sdk.NewCoins(
				sdk.NewCoin(fmt.Sprintf("%stoken", val.Moniker), sdk.NewInt(10)),
				sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10)),
			),
			[]string{
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
				"--gas=10",
			},
			false,
			&sdk.TxResponse{},
			sdkerrors.ErrOutOfGas.ABCICode(),
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			clientCtx := val.ClientCtx

			bz, err := banktestutil.MsgSendExec(clientCtx, tc.from, tc.to, tc.amount, tc.args...)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)

				s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), tc.respType), bz.String())
				txResp := tc.respType.(*sdk.TxResponse)
				s.Require().Equal(tc.expectedCode, txResp.Code)
			}
		})
	}
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func NewCoin(denom string, amount sdk.Int) *sdk.Coin {
	coin := sdk.NewCoin(denom, amount)
	return &coin
}
