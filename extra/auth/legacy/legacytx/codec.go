package legacytx

import (
	"github.com/ocea/sdk/codec"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(StdTx{}, "ocea-sdk/StdTx", nil)
}
