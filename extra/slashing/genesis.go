package slashing

import (
	"github.com/ocea/sdk/extra/genesis"
	"github.com/ocea/sdk/extra/slashing/keeper"
	"github.com/ocea/sdk/extra/slashing/types"
	"github.com/ocea/sdk/extra/staking/exported"
	sdk "github.com/ocea/sdk/types"
)

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, stakingKeeper types.StakingKeeper, data *genesis.GenesisState) {
	stakingKeeper.IterateValidators(ctx,
		func(index int64, validator exported.ValidatorI) bool {
			keeper.AddPubkey(ctx, validator.GetConsPubKey())
			return false
		},
	)

	for _, info := range data.SigningInfos {
		address, err := sdk.ConsAddressFromBech32(info.Address)
		if err != nil {
			panic(err)
		}
		keeper.SetValidatorSigningInfo(ctx, address, info.ValidatorSigningInfo)
	}

	for _, array := range data.MissedBlocks {
		address, err := sdk.ConsAddressFromBech32(array.Address)
		if err != nil {
			panic(err)
		}
		for _, missed := range array.MissedBlocks {
			keeper.SetValidatorMissedBlockBitArray(ctx, address, missed.Index, missed.Missed)
		}
	}

	keeper.SetParams(ctx, data.Params)
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) (data *genesis.GenesisState) {
	params := keeper.GetParams(ctx)
	signingInfos := make([]genesis.SigningInfo, 0)
	missedBlocks := make([]genesis.ValidatorMissedBlocks, 0)
	keeper.IterateValidatorSigningInfos(ctx, func(address sdk.ConsAddress, info types.ValidatorSigningInfo) (stop bool) {
		bechAddr := address.String()
		signingInfos = append(signingInfos, genesis.SigningInfo{
			Address:              bechAddr,
			ValidatorSigningInfo: info,
		})

		localMissedBlocks := keeper.GetValidatorMissedBlocks(ctx, address)

		missedBlocks = append(missedBlocks, genesis.ValidatorMissedBlocks{
			Address:      bechAddr,
			MissedBlocks: localMissedBlocks,
		})

		return false
	})

	return genesis.NewGenesisState(params, signingInfos, missedBlocks)
}
