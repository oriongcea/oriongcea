package sub

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

var genesisHash,_ = hex.DecodeString("66E9C66B3A40EAEADFAE13C3C2A7A50F")
var genesisTx = []byte("{from:\"O\",to:\"O5W9xKvnk13PuakHVVne82YjKjZs1cfjFwEhE5jgCnTyz\",amount:5000000}")
var genesisTxHash,_ = hex.DecodeString("707afb1e20b935cf96059ca5ca22f31e")
var genesisTime,_ = time.Parse("Jan 2 15:04:05 +0800 MST 2006","Oct 9 8:00:00 +0800 MST 2020")
// MainNetParams defines the network parameters for the main Bitcoin network.
var MainParams = Params{
	Name:        "mainnet",
	DefaultPort: "8990",
	DNSSeeds: []DNSSeed{
		{"seed.taurus.com", true},
	},

	// Chain parameters
	GenesisBlock: &types.Block{
		Header:     types.Header{
			Version:  version.Consensus{
				Block: 1,
				App:   1,
			},
			ChainID:            "taurus",
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
		{30*1440,550000},
		{60*1440,605000},
		{90*1440,665500},
		{120*1440,732050},
		{150*1440,805255},
		{180*1440,708624},
		{210*1440,765314},
		{240*1440,826539},
		{270*1440,892662},
		{300*1440,964075},
		{330*1440,1041201},
		{360*1440,1124497},
		{390*1440,1214457},
		{420*1440,1311614},
		{450*1440,1416543},
		{480*1440,1529866},
		{510*1440,1652256},
		{540*1440,1115272},
		{570*1440,1171036},
		{600*1440,1229588},
		{630*1440,1291067},
		{660*1440,1355621},
		{690*1440,1423402},
		{720*1440,1494572},
		{750*1440,1569301},
		{780*1440,1647766},
		{810*1440,1730154},
		{840*1440,1816662},
		{870*1440,1907495},
		{900*1440,2002870},
		{930*1440,2103013},
		{960*1440,2208164},
		{990*1440,2318572},
		{1020*1440,2434501},
		{1050*1440,2556226},
		{1080*1440,2684037},
		{1110*1440,2818239},
		{1140*1440,2959151},
		{1170*1440,3107108},
		{1200*1440,3262464},
		{1230*1440,3425587},
		{1260*1440,2158120},
		{1290*1440,2222863},
		{1320*1440,2289549},
		{1350*1440,2358236},
		{1380*1440,2428983},
		{1410*1440,2501852},
		{1440*1440,2576908},
		{1470*1440,2654215},
		{1500*1440,2733841},
		{1530*1440,2815857},
		{1560*1440,2900332},
		{1590*1440,2987342},
		{1620*1440,3076963},
		{1650*1440,3169272},
		{1680*1440,3264350},
		{1710*1440,3362280},
		{1740*1440,3463149},
		{1770*1440,3567043},
		{1800*1440,3674055},
		{1830*1440,3784276},
		{1860*1440,3897804},
		{1890*1440,4014739},
		{1920*1440,4135181},
		{1950*1440,4259236},
		{1980*1440,4387013},
		{2010*1440,4518624},
		{2040*1440,4654183},
		{2070*1440,4793808},
		{2100*1440,4937622},
		{2130*1440,5085751},
		{2160*1440,5238323},
		{2190*1440,5395473},
		{2220*1440,5557337},
		{2250*1440,5724058},
		{2280*1440,5895779},
		{2310*1440,6072653},
		{2340*1440,4169888},
		{2370*1440,4253286},
		{2400*1440,4338352},
		{2430*1440,4425119},
		{2460*1440,4513621},
		{2490*1440,4603893},
		{2520*1440,4695971},
		{2550*1440,4789891},
		{2580*1440,4885688},
		{2610*1440,4983402},
		{2640*1440,5083070},
		{2670*1440,5184732},
		{2700*1440,5288426},
		{2730*1440,5394195},
		{2760*1440,5502079},
		{2790*1440,5612120},
		{2820*1440,5724363},
		{2850*1440,5838850},
		{2880*1440,5955627},
		{2910*1440,6074740},
		{2940*1440,6196234},
		{2970*1440,6320159},
		{3000*1440,6446562},
		{3030*1440,6575494},
		{3060*1440,6707004},
		{3090*1440,6841144},
		{3120*1440,6977966},
		{3150*1440,7117526},
		{3180*1440,7259876},
		{3210*1440,7405074},
		{3240*1440,7553175},
		{3270*1440,7704239},
		{3300*1440,7858324},
		{3330*1440,8015490},
		{3360*1440,8175800},
		{3390*1440,8339316},
		{3420*1440,8506102},
		{3450*1440,8676224},
		{3480*1440,8849749},
		{3510*1440,9026744},
		{3540*1440,9207279},
		{3570*1440,9391424},
		{3600*1440,9579253},
		{3630*1440,9770838},
		{3660*1440,1687242},
		{3690*1440, 0},

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
