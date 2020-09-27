package keeper_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/ocea/sdk/baseapp"
	"github.com/ocea/sdk/codec"
	cryptotypes "github.com/ocea/sdk/crypto/types"
	"github.com/ocea/sdk/simapp"
	sdk "github.com/ocea/sdk/types"
	"github.com/ocea/sdk/extra/ibc/02-client/keeper"
	"github.com/ocea/sdk/extra/ibc/02-client/types"
	ibctmtypes "github.com/ocea/sdk/extra/ibc/07-tendermint/types"
	localhosttypes "github.com/ocea/sdk/extra/ibc/09-localhost/types"
	commitmenttypes "github.com/ocea/sdk/extra/ibc/23-commitment/types"
	"github.com/ocea/sdk/extra/ibc/exported"
	ibctesting "github.com/ocea/sdk/extra/ibc/testing"
	ibctestingmock "github.com/ocea/sdk/extra/ibc/testing/mock"
	stakingtypes "github.com/ocea/sdk/extra/staking/types"
)

const (
	testChainID       = "gaiahub-0"
	testChainIDEpoch1 = "gaiahub-1"

	testClientID  = "gaiachain"
	testClientID2 = "ethbridge"
	testClientID3 = "ethermint"

	height = 5

	trustingPeriod time.Duration = time.Hour * 24 * 7 * 2
	ubdPeriod      time.Duration = time.Hour * 24 * 7 * 3
	maxClockDrift  time.Duration = time.Second * 10
)

var testClientHeight = types.NewHeight(0, 5)
var testClientHeightEpoch1 = types.NewHeight(1, 5)

type KeeperTestSuite struct {
	suite.Suite

	coordinator *ibctesting.Coordinator

	chainA *ibctesting.TestChain
	chainB *ibctesting.TestChain

	cdc            codec.Marshaler
	ctx            sdk.Context
	keeper         *keeper.Keeper
	consensusState *ibctmtypes.ConsensusState
	header         *ibctmtypes.Header
	valSet         *tmtypes.ValidatorSet
	valSetHash     tmbytes.HexBytes
	privVal        tmtypes.PrivValidator
	now            time.Time
	past           time.Time

	queryClient types.QueryClient
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 2)

	suite.chainA = suite.coordinator.GetChain(ibctesting.GetChainID(0))
	suite.chainB = suite.coordinator.GetChain(ibctesting.GetChainID(1))

	isCheckTx := false
	suite.now = time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)
	suite.past = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	now2 := suite.now.Add(time.Hour)
	app := simapp.Setup(isCheckTx)

	suite.cdc = app.AppCodec()
	suite.ctx = app.BaseApp.NewContext(isCheckTx, tmproto.Header{Height: height, ChainID: testClientID, Time: now2})
	suite.keeper = &app.IBCKeeper.ClientKeeper
	suite.privVal = ibctestingmock.NewPV()

	pubKey, err := suite.privVal.GetPubKey()
	suite.Require().NoError(err)

	testClientHeightMinus1 := types.NewHeight(0, height-1)

	validator := tmtypes.NewValidator(pubKey.(cryptotypes.IntoTmPubKey).AsTmPubKey(), 1)
	suite.valSet = tmtypes.NewValidatorSet([]*tmtypes.Validator{validator})
	suite.valSetHash = suite.valSet.Hash()
	suite.header = ibctmtypes.CreateTestHeader(testChainID, testClientHeight, testClientHeightMinus1, now2, suite.valSet, suite.valSet, []tmtypes.PrivValidator{suite.privVal})
	suite.consensusState = ibctmtypes.NewConsensusState(suite.now, commitmenttypes.NewMerkleRoot([]byte("hash")), suite.valSetHash)

	var validators stakingtypes.Validators
	for i := 1; i < 11; i++ {
		privVal := ibctestingmock.NewPV()
		pk, err := privVal.GetPubKey()
		suite.Require().NoError(err)
		val := stakingtypes.NewValidator(sdk.ValAddress(pk.Address()), pk, stakingtypes.Description{})
		val.Status = sdk.Bonded
		val.Tokens = sdk.NewInt(rand.Int63())
		validators = append(validators, val)

		app.StakingKeeper.SetHistoricalInfo(suite.ctx, int64(i), stakingtypes.NewHistoricalInfo(suite.ctx.BlockHeader(), validators))
	}

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, app.IBCKeeper.ClientKeeper)
	suite.queryClient = types.NewQueryClient(queryHelper)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestSetClientState() {
	clientState := ibctmtypes.NewClientState(testChainID, ibctmtypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, types.ZeroHeight(), commitmenttypes.GetSDKSpecs(), false, false)
	suite.keeper.SetClientState(suite.ctx, testClientID, clientState)

	retrievedState, found := suite.keeper.GetClientState(suite.ctx, testClientID)
	suite.Require().True(found, "GetClientState failed")
	suite.Require().Equal(clientState, retrievedState, "Client states are not equal")
}

func (suite *KeeperTestSuite) TestSetClientConsensusState() {
	suite.keeper.SetClientConsensusState(suite.ctx, testClientID, testClientHeight, suite.consensusState)

	retrievedConsState, found := suite.keeper.GetClientConsensusState(suite.ctx, testClientID, testClientHeight)
	suite.Require().True(found, "GetConsensusState failed")

	tmConsState, ok := retrievedConsState.(*ibctmtypes.ConsensusState)
	suite.Require().True(ok)
	suite.Require().Equal(suite.consensusState, tmConsState, "ConsensusState not stored correctly")
}

func (suite *KeeperTestSuite) TestValidateSelfClient() {
	testCases := []struct {
		name        string
		clientState exported.ClientState
		expPass     bool
	}{
		{
			"success",
			ibctmtypes.NewClientState(testChainID, ibctmtypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, testClientHeight, commitmenttypes.GetSDKSpecs(), false, false),
			true,
		},
		{
			"invalid client type",
			localhosttypes.NewClientState(testChainID, testClientHeight),
			false,
		},
		{
			"frozen client",
			&ibctmtypes.ClientState{testChainID, ibctmtypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, testClientHeight, testClientHeight, commitmenttypes.GetSDKSpecs(), false, false},
			false,
		},
		{
			"incorrect chainID",
			ibctmtypes.NewClientState("gaiatestnet", ibctmtypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, testClientHeight, commitmenttypes.GetSDKSpecs(), false, false),
			false,
		},
		{
			"invalid client height",
			ibctmtypes.NewClientState(testChainID, ibctmtypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, types.NewHeight(0, testClientHeight.EpochHeight+10), commitmenttypes.GetSDKSpecs(), false, false),
			false,
		},
		{
			"invalid client epoch",
			ibctmtypes.NewClientState(testChainIDEpoch1, ibctmtypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, testClientHeightEpoch1, commitmenttypes.GetSDKSpecs(), false, false),
			false,
		},
		{
			"invalid proof specs",
			ibctmtypes.NewClientState(testChainID, ibctmtypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, testClientHeight, nil, false, false),
			false,
		},
		{
			"invalid trust level",
			ibctmtypes.NewClientState(testChainID, ibctmtypes.Fraction{0, 1}, trustingPeriod, ubdPeriod, maxClockDrift, testClientHeight, commitmenttypes.GetSDKSpecs(), false, false),
			false,
		},
		{
			"invalid unbonding period",
			ibctmtypes.NewClientState(testChainID, ibctmtypes.DefaultTrustLevel, trustingPeriod, ubdPeriod+10, maxClockDrift, testClientHeight, commitmenttypes.GetSDKSpecs(), false, false),
			false,
		},
		{
			"invalid trusting period",
			ibctmtypes.NewClientState(testChainID, ibctmtypes.DefaultTrustLevel, ubdPeriod+10, ubdPeriod, maxClockDrift, testClientHeight, commitmenttypes.GetSDKSpecs(), false, false),
			false,
		},
	}

	ctx := suite.ctx.WithChainID(testChainID)
	ctx = ctx.WithBlockHeight(height)

	for _, tc := range testCases {
		err := suite.keeper.ValidateSelfClient(ctx, tc.clientState)
		if tc.expPass {
			suite.Require().NoError(err, "expected valid client for case: %s", tc.name)
		} else {
			suite.Require().Error(err, "expected invalid client for case: %s", tc.name)
		}
	}
}

func (suite KeeperTestSuite) TestGetAllClients() {
	clientIDs := []string{
		testClientID2, testClientID3, testClientID,
	}
	expClients := []exported.ClientState{
		ibctmtypes.NewClientState(testChainID, ibctmtypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, types.ZeroHeight(), commitmenttypes.GetSDKSpecs(), false, false),
		ibctmtypes.NewClientState(testChainID, ibctmtypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, types.ZeroHeight(), commitmenttypes.GetSDKSpecs(), false, false),
		ibctmtypes.NewClientState(testChainID, ibctmtypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, types.ZeroHeight(), commitmenttypes.GetSDKSpecs(), false, false),
	}

	for i := range expClients {
		suite.keeper.SetClientState(suite.ctx, clientIDs[i], expClients[i])
	}

	// add localhost client
	localHostClient, found := suite.keeper.GetClientState(suite.ctx, exported.Localhost)
	suite.Require().True(found)
	expClients = append(expClients, localHostClient)

	clients := suite.keeper.GetAllClients(suite.ctx)
	suite.Require().Len(clients, len(expClients))
	suite.Require().Equal(expClients, clients)
}

func (suite KeeperTestSuite) TestGetAllGenesisClients() {
	clientIDs := []string{
		testClientID2, testClientID3, testClientID,
	}
	expClients := []exported.ClientState{
		ibctmtypes.NewClientState(testChainID, ibctmtypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, types.ZeroHeight(), commitmenttypes.GetSDKSpecs(), false, false),
		ibctmtypes.NewClientState(testChainID, ibctmtypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, types.ZeroHeight(), commitmenttypes.GetSDKSpecs(), false, false),
		ibctmtypes.NewClientState(testChainID, ibctmtypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, types.ZeroHeight(), commitmenttypes.GetSDKSpecs(), false, false),
	}

	expGenClients := make([]types.IdentifiedClientState, len(expClients))

	for i := range expClients {
		suite.keeper.SetClientState(suite.ctx, clientIDs[i], expClients[i])
		expGenClients[i] = types.NewIdentifiedClientState(clientIDs[i], expClients[i])
	}

	// add localhost client
	localHostClient, found := suite.keeper.GetClientState(suite.ctx, exported.Localhost)
	suite.Require().True(found)
	expGenClients = append(expGenClients, types.NewIdentifiedClientState(exported.Localhost, localHostClient))

	genClients := suite.keeper.GetAllGenesisClients(suite.ctx)

	suite.Require().Equal(expGenClients, genClients)
}

func (suite KeeperTestSuite) TestGetConsensusState() {
	suite.ctx = suite.ctx.WithBlockHeight(10)
	cases := []struct {
		name    string
		height  types.Height
		expPass bool
	}{
		{"zero height", types.ZeroHeight(), false},
		{"height > latest height", types.NewHeight(0, uint64(suite.ctx.BlockHeight())+1), false},
		{"latest height - 1", types.NewHeight(0, uint64(suite.ctx.BlockHeight())-1), true},
		{"latest height", types.GetSelfHeight(suite.ctx), true},
	}

	for i, tc := range cases {
		tc := tc
		cs, found := suite.keeper.GetSelfConsensusState(suite.ctx, tc.height)
		if tc.expPass {
			suite.Require().True(found, "Case %d should have passed: %s", i, tc.name)
			suite.Require().NotNil(cs, "Case %d should have passed: %s", i, tc.name)
		} else {
			suite.Require().False(found, "Case %d should have failed: %s", i, tc.name)
			suite.Require().Nil(cs, "Case %d should have failed: %s", i, tc.name)
		}
	}
}

func (suite KeeperTestSuite) TestConsensusStateHelpers() {
	// initial setup
	clientState := ibctmtypes.NewClientState(testChainID, ibctmtypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, testClientHeight, commitmenttypes.GetSDKSpecs(), false, false)

	suite.keeper.SetClientState(suite.ctx, testClientID, clientState)
	suite.keeper.SetClientConsensusState(suite.ctx, testClientID, testClientHeight, suite.consensusState)

	nextState := ibctmtypes.NewConsensusState(suite.now, commitmenttypes.NewMerkleRoot([]byte("next")), suite.valSetHash)

	testClientHeightPlus5 := types.NewHeight(0, height+5)

	header := ibctmtypes.CreateTestHeader(testClientID, testClientHeightPlus5, testClientHeight, suite.header.Header.Time.Add(time.Minute),
		suite.valSet, suite.valSet, []tmtypes.PrivValidator{suite.privVal})

	// mock update functionality
	clientState.LatestHeight = header.GetHeight().(types.Height)
	suite.keeper.SetClientConsensusState(suite.ctx, testClientID, header.GetHeight(), nextState)
	suite.keeper.SetClientState(suite.ctx, testClientID, clientState)

	latest, ok := suite.keeper.GetLatestClientConsensusState(suite.ctx, testClientID)
	suite.Require().True(ok)
	suite.Require().Equal(nextState, latest, "Latest client not returned correctly")

	// Should return existing consensusState at latestClientHeight
	lte, ok := suite.keeper.GetClientConsensusStateLTE(suite.ctx, testClientID, types.NewHeight(0, height+3))
	suite.Require().True(ok)
	suite.Require().Equal(suite.consensusState, lte, "LTE helper function did not return latest client state below height: %d", height+3)
}

// 2 clients in total are created on chainA. The first client is updated so it contains an initial consensus state
// and a consensus state at the update height.
func (suite KeeperTestSuite) TestGetAllConsensusStates() {
	clientA, _ := suite.coordinator.SetupClients(suite.chainA, suite.chainB, ibctesting.Tendermint)

	clientState := suite.chainA.GetClientState(clientA)
	expConsensusHeight0 := clientState.GetLatestHeight()
	consensusState0, ok := suite.chainA.GetConsensusState(clientA, expConsensusHeight0)
	suite.Require().True(ok)

	// update client to create a second consensus state
	err := suite.coordinator.UpdateClient(suite.chainA, suite.chainB, clientA, ibctesting.Tendermint)
	suite.Require().NoError(err)

	clientState = suite.chainA.GetClientState(clientA)
	expConsensusHeight1 := clientState.GetLatestHeight()
	suite.Require().True(expConsensusHeight1.GT(expConsensusHeight0))
	consensusState1, ok := suite.chainA.GetConsensusState(clientA, expConsensusHeight1)
	suite.Require().True(ok)

	expConsensus := []exported.ConsensusState{
		consensusState0,
		consensusState1,
	}

	// create second client on chainA
	clientA2, _ := suite.coordinator.SetupClients(suite.chainA, suite.chainB, ibctesting.Tendermint)
	clientState = suite.chainA.GetClientState(clientA2)

	expConsensusHeight2 := clientState.GetLatestHeight()
	consensusState2, ok := suite.chainA.GetConsensusState(clientA2, expConsensusHeight2)
	suite.Require().True(ok)

	expConsensus2 := []exported.ConsensusState{consensusState2}

	expConsensusStates := types.ClientsConsensusStates{
		types.NewClientConsensusStates(clientA, []types.ConsensusStateWithHeight{
			types.NewConsensusStateWithHeight(expConsensusHeight0.(types.Height), expConsensus[0]),
			types.NewConsensusStateWithHeight(expConsensusHeight1.(types.Height), expConsensus[1]),
		}),
		types.NewClientConsensusStates(clientA2, []types.ConsensusStateWithHeight{
			types.NewConsensusStateWithHeight(expConsensusHeight2.(types.Height), expConsensus2[0]),
		}),
	}.Sort()

	consStates := suite.chainA.App.IBCKeeper.ClientKeeper.GetAllConsensusStates(suite.chainA.GetContext())
	suite.Require().Equal(expConsensusStates, consStates, "%s \n\n%s", expConsensusStates, consStates)
}
