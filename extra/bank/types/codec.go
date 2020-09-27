package types

import (
	"github.com/ocea/sdk/codec"
	"github.com/ocea/sdk/codec/types"
	cryptocodec "github.com/ocea/sdk/crypto/codec"
	sdk "github.com/ocea/sdk/types"
	"github.com/ocea/sdk/extra/bank/exported"
)

// RegisterLegacyAminoCodec registers the necessary extra/bank interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*exported.SupplyI)(nil), nil)
	cdc.RegisterConcrete(&Supply{}, "ocea-sdk/Supply", nil)
	cdc.RegisterConcrete(&MsgSend{}, "ocea-sdk/MsgSend", nil)
	cdc.RegisterConcrete(&MsgMultiSend{}, "ocea-sdk/MsgMultiSend", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSend{},
		&MsgMultiSend{},
	)

	registry.RegisterInterface(
		"ocea.bank.v1beta1.SupplyI",
		(*exported.SupplyI)(nil),
		&Supply{},
	)
}

var (
	amino = codec.NewLegacyAmino()

	// ModuleCdc references the global extra/bank module codec. Note, the codec should
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
