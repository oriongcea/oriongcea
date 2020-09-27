package types

import (
	"github.com/ocea/sdk/codec"
	cryptocodec "github.com/ocea/sdk/crypto/codec"
)

var (
	amino = codec.NewLegacyAmino()
)

func init() {
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}
