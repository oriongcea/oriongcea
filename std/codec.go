package std

import (
	"github.com/ocea/sdk/codec"
	"github.com/ocea/sdk/codec/types"
	cryptocodec "github.com/ocea/sdk/crypto/codec"
	sdk "github.com/ocea/sdk/types"
	txtypes "github.com/ocea/sdk/types/tx"
	vesting "github.com/ocea/sdk/extra/auth/vesting/types"
)

// RegisterLegacyAminoCodec registers types with the Amino codec.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	sdk.RegisterLegacyAminoCodec(cdc)
	cryptocodec.RegisterCrypto(cdc)
	vesting.RegisterLegacyAminoCodec(cdc)
}

// RegisterInterfaces registers Interfaces from sdk/types, vesting, crypto, tx.
func RegisterInterfaces(interfaceRegistry types.InterfaceRegistry) {
	sdk.RegisterInterfaces(interfaceRegistry)
	txtypes.RegisterInterfaces(interfaceRegistry)
	cryptocodec.RegisterInterfaces(interfaceRegistry)
	vesting.RegisterInterfaces(interfaceRegistry)
}
