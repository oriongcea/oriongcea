package codec

import (
	"github.com/tendermint/tendermint/crypto"
	tmed25519 "github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/sr25519"

	"github.com/ocea/sdk/codec"
	"github.com/ocea/sdk/crypto/keys/ed25519"
	kmultisig "github.com/ocea/sdk/crypto/keys/multisig"
	"github.com/ocea/sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/ocea/sdk/crypto/types"
)

var amino *codec.LegacyAmino

func init() {
	amino = codec.NewLegacyAmino()
	RegisterCrypto(amino)
}

// RegisterCrypto registers all crypto dependency types with the provided Amino
// codec.
func RegisterCrypto(cdc *codec.LegacyAmino) {
	// TODO We now register both Tendermint's PubKey and our own PubKey. In the
	// long-term, we should move away from Tendermint's PubKey, and delete this
	// first line.
	cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	cdc.RegisterInterface((*cryptotypes.PubKey)(nil), nil)
	cdc.RegisterConcrete(sr25519.PubKey{},
		sr25519.PubKeyName, nil)
	// TODO Same as above, for ED25519
	cdc.RegisterConcrete(tmed25519.PubKey{},
		tmed25519.PubKeyName, nil)
	cdc.RegisterConcrete(&ed25519.PubKey{},
		ed25519.PubKeyName, nil)
	cdc.RegisterConcrete(&secp256k1.PubKey{},
		secp256k1.PubKeyName, nil)
	cdc.RegisterConcrete(&kmultisig.LegacyAminoPubKey{},
		kmultisig.PubKeyAminoRoute, nil)

	cdc.RegisterInterface((*crypto.PrivKey)(nil), nil)
	cdc.RegisterConcrete(sr25519.PrivKey{},
		sr25519.PrivKeyName, nil)
	// TODO Same as above
	cdc.RegisterConcrete(tmed25519.PrivKey{},
		tmed25519.PrivKeyName, nil)
	cdc.RegisterConcrete(&ed25519.PrivKey{},
		ed25519.PrivKeyName, nil)
	cdc.RegisterConcrete(&secp256k1.PrivKey{},
		secp256k1.PrivKeyName, nil)
}

// PrivKeyFromBytes unmarshals private key bytes and returns a PrivKey
func PrivKeyFromBytes(privKeyBytes []byte) (privKey crypto.PrivKey, err error) {
	err = amino.UnmarshalBinaryBare(privKeyBytes, &privKey)
	return
}

// PubKeyFromBytes unmarshals public key bytes and returns a PubKey
func PubKeyFromBytes(pubKeyBytes []byte) (pubKey crypto.PubKey, err error) {
	err = amino.UnmarshalBinaryBare(pubKeyBytes, &pubKey)
	return
}