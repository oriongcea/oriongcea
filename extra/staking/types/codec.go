package types

import (
	"github.com/ocea/sdk/codec"
	"github.com/ocea/sdk/codec/types"
	cryptocodec "github.com/ocea/sdk/crypto/codec"
	sdk "github.com/ocea/sdk/types"
)

// RegisterLegacyAminoCodec registers the necessary extra/staking interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateValidator{}, "ocea-sdk/MsgCreateValidator", nil)
	cdc.RegisterConcrete(&MsgEditValidator{}, "ocea-sdk/MsgEditValidator", nil)
	cdc.RegisterConcrete(&MsgDelegate{}, "ocea-sdk/MsgDelegate", nil)
	cdc.RegisterConcrete(&MsgUndelegate{}, "ocea-sdk/MsgUndelegate", nil)
	cdc.RegisterConcrete(&MsgBeginRedelegate{}, "ocea-sdk/MsgBeginRedelegate", nil)
}

// RegisterInterfaces registers the extra/staking interfaces types with the interface registry
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateValidator{},
		&MsgEditValidator{},
		&MsgDelegate{},
		&MsgUndelegate{},
		&MsgBeginRedelegate{},
	)
}

var (
	amino = codec.NewLegacyAmino()

	// ModuleCdc references the global extra/staking module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding as Amino is
	// still used for that purpose.
	//
	// The actual codec used for serialization should be provided to extra/staking and
	// defined at the application level.
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}
