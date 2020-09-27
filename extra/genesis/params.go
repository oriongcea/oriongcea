package genesis

import (
	"encoding/hex"
	"github.com/tendermint/tendermint/proto/tendermint/types"
	"github.com/tendermint/tendermint/proto/tendermint/version"
	"time"
)

type Checkpoint struct {
	Height int32
	// CoinbaseMaturity is the number of blocks required before newly mined
	// coins (coinbase transactions) can be spent.
	CoinbaseMaturity uint16
}

type DNSSeed struct {
	// Host defines the hostname of the seed.
	Host string

	// HasFiltering defines whether the seed supports filtering
	// by service flags (wire.ServiceFlag).
	HasFiltering bool
}

// Params defines a ocea network by its parameters.  These parameters may be
// used by ocea applications to differentiate networks as well as addresses
// and keys for one network from those intended for use on another network.
type Params struct {
	// Name defines a human-readable identifier for the network.
	Name string

	// DefaultPort defines the default peer-to-peer port for the network.
	DefaultPort string

	// DNSSeeds defines a list of DNS seeds for the network that are used
	// as one method to discover peers.
	DNSSeeds []DNSSeed

	// GenesisBlock defines the first block of the chain.
	GenesisBlock *types.Block

	// GenesisHash is the starting block hash.
	GenesisHash []byte



	// TargetTimePerBlock is the desired amount of time to generate each
	// block.
	TargetTimePerBlock time.Duration

	// RetargetAdjustmentFactor is the adjustment factor used to limit
	// the minimum and maximum amount of adjustment that can occur between
	// difficulty retargets.
	RetargetAdjustmentFactor int64


	// GenerateSupported specifies whether or not CPU mining is allowed.
	GenerateSupported bool

	// Checkpoints ordered from oldest to newest.
	Checkpoints []Checkpoint

	// These fields are related to voting on consensus rule changes as
	// defined by BIP0009.
	//
	// RuleChangeActivationThreshold is the number of blocks in a threshold
	// state retarget window for which a positive vote for a rule change
	// must be cast in order to lock in a rule change. It should typically
	// be 95% for the main network and 75% for test networks.
	//
	// MinerConfirmationWindow is the number of blocks in each threshold
	// state retarget window.
	RuleChangeActivationThreshold uint32
	MinerConfirmationWindow       uint32

	// Mempool parameters
	RelayNonStdTxs bool



}

var genesisHash,_ = hex.DecodeString("8f7645d9a48eaf24a72331ff29d0d2bf")
var genesisTx = []byte("{from:\"O\",to:\"OECu9M4jcXErksXXm9Q6ACpZyHSJgofn2YLLxPcrGfMK\",amount:10000000}")
var genesisTxHash,_ = hex.DecodeString("0fdc1365e72901dcead49bd115c66be3")
var genesisTime,_ = time.Parse("Jan 2 15:04:05 +0800 MST 2006","Jun 16 12:27:59 +0800 MST 2020")
// MainNetParams defines the network parameters for the main Bitcoin network.
var MainParams = Params{
	Name:        "mainnet",
	DefaultPort: "8989",
	DNSSeeds: []DNSSeed{
		{"seed.ocea.com", true},
		{"dnsseed.bluematt.me", true},
		{"seed.orion.com", true},
		{"seed.jonasschnelli.io", false},
	},

	// Chain parameters
	GenesisBlock: &types.Block{
		Header:     types.Header{
			Version:  version.Consensus{
				Block: 1,
				App:   1,
			},
			ChainID:            "ocea",
			Height:             0,
			Time:              genesisTime ,
		},
		Data:       types.Data{
			Txs:  [][]byte{genesisTx},
			Hash: genesisTxHash,
		},
	},
	GenesisHash:              genesisHash,
	TargetTimePerBlock:       time.Minute ,    // 1 minutes
	RetargetAdjustmentFactor: 4,               // 25% less, 400% more
	GenerateSupported:        false,

	// Checkpoints ordered from oldest to newest.
	Checkpoints: []Checkpoint{
		{0, 1000000},
		{30*1440, 1100000},
		{60*1440, 1210000},
		{90*1440, 1331000},
		{120*1440, 1464100},
		{150*1440, 1610510},
		{180*1440, 1417248},
		{210*1440, 1530628},
		{240*1440, 1653079},
		{270*1440, 1785325},
		{300*1440, 1928151},
		{330*1440, 2082403},
		{360*1440, 2248995},
		{390*1440, 2428915},
		{420*1440, 2623228},
		{450*1440, 2833086},
		{480*1440, 3059733},
		{510*1440, 3304512},
		{540*1440, 2230545},
		{570*1440, 2342073},
		{600*1440, 2459176},
		{630*1440, 2582135},
		{660*1440, 2711242},
		{690*1440, 2846804},
		{720*1440, 2989144},
		{750*1440, 3138602},
		{780*1440, 3295532},
		{810*1440, 3460308},
		{840*1440, 3633324},
		{870*1440, 3814990},
		{900*1440, 4005740},
		{930*1440, 4206027},
		{960*1440, 4416328},
		{990*1440, 4637144},
		{1020*1440, 4869002},
		{1050*1440, 5112452},
		{1080*1440, 5368074},
		{1110*1440, 5636478},
		{1140*1440, 5918302},
		{1170*1440, 6214217},
		{1200*1440, 6524928},
		{1230*1440, 6851174},
		{1260*1440, 4316240},
		{1290*1440, 4445727},
		{1320*1440, 4579099},
		{1350*1440, 4716472},
		{1380*1440, 4857966},
		{1410*1440, 5003705},
		{1440*1440, 5153816},
		{1470*1440, 5308431},
		{1500*1440, 5467683},
		{1530*1440, 5631714},
		{1560*1440, 5800665},
		{1590*1440, 5974685},
		{1620*1440, 6153926},
		{1650*1440, 6338544},
		{1680*1440, 6528700},
		{1710*1440, 6724561},
		{1740*1440, 6926298},
		{1770*1440, 7134087},
		{1800*1440, 7348110},
		{1830*1440, 7568553},
		{1860*1440, 7795609},
		{1890*1440, 8029478},
		{1920*1440, 8270362},
		{1950*1440, 8518473},
		{1980*1440, 8774027},
		{2010*1440, 9037248},
		{2040*1440, 9308366},
		{2070*1440, 9587616},
		{2100*1440, 9875245},
		{2130*1440, 10171502},
		{2160*1440, 10476647},
		{2190*1440, 10790947},
		{2220*1440, 11114675},
		{2250*1440, 11448116},
		{2280*1440, 11791559},
		{2310*1440, 12145306},
		{2340*1440, 8339777},
		{2370*1440, 8506572},
		{2400*1440, 8676704},
		{2430*1440, 8850238},
		{2460*1440, 9027242},
		{2490*1440, 9207787},
		{2520*1440, 9391943},
		{2550*1440, 9579782},
		{2580*1440, 9771377},
		{2610*1440, 9966805},
		{2640*1440, 10166141},
		{2670*1440, 10369464},
		{2700*1440, 10576853},
		{2730*1440, 10788390},
		{2760*1440, 11004158},
		{2790*1440, 11224241},
		{2820*1440, 11448726},
		{2850*1440, 11677701},
		{2880*1440, 11911255},
		{2910*1440, 12149480},
		{2940*1440, 12392469},
		{2970*1440, 12640319},
		{3000*1440, 12893125},
		{3030*1440, 13150988},
		{3060*1440, 13414008},
		{3090*1440, 13682288},
		{3120*1440, 13955933},
		{3150*1440, 14235052},
		{3180*1440, 14519753},
		{3210*1440, 14810148},
		{3240*1440, 15106351},
		{3270*1440, 15408478},
		{3300*1440, 15716648},
		{3330*1440, 16030981},
		{3360*1440, 16351600},
		{3390*1440, 16678632},
		{3420*1440, 17012205},
		{3450*1440, 17352449},
		{3480*1440, 17699498},
		{3510*1440, 18053488},
		{3540*1440, 18414558},
		{3570*1440, 18782849},
		{3600*1440, 19158506},
		{3630*1440, 19541676},
		{3660*1440, 3374485},
		{3690*1440, 0}, //ten years
	},

	// Consensus rule change deployments.
	//
	// The miner confirmation window is defined as:
	//   target proof of work timespan / target proof of work spacing
	RuleChangeActivationThreshold: 1916, // 95% of MinerConfirmationWindow
	MinerConfirmationWindow:       2016, //

	// Mempool parameters
	RelayNonStdTxs: false,

}
