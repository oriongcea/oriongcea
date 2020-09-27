package keyring

import (
	"github.com/ocea/sdk/codec"
	cryptocodec "github.com/ocea/sdk/crypto/codec"
	"github.com/ocea/sdk/crypto/hd"
)

// CryptoCdc defines the codec required for keys and info
var CryptoCdc *codec.LegacyAmino

func init() {
	CryptoCdc = codec.NewLegacyAmino()
	cryptocodec.RegisterCrypto(CryptoCdc)
	RegisterLegacyAminoCodec(CryptoCdc)
	CryptoCdc.Seal()
}

// RegisterLegacyAminoCodec registers concrete types and interfaces on the given codec.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*Info)(nil), nil)
	cdc.RegisterConcrete(hd.BIP44Params{}, "crypto/keys/hd/BIP44Params", nil)
	cdc.RegisterConcrete(localInfo{}, "crypto/keys/localInfo", nil)
	cdc.RegisterConcrete(ledgerInfo{}, "crypto/keys/ledgerInfo", nil)
	cdc.RegisterConcrete(offlineInfo{}, "crypto/keys/offlineInfo", nil)
	cdc.RegisterConcrete(multiInfo{}, "crypto/keys/multiInfo", nil)
}
