package v040_test

import (
	"testing"

	"github.com/ocea/sdk/codec"
	cryptocodec "github.com/ocea/sdk/crypto/codec"
	sdk "github.com/ocea/sdk/types"
	v038auth "github.com/ocea/sdk/extra/auth/legacy/v0_38"
	v039auth "github.com/ocea/sdk/extra/auth/legacy/v0_39"
	v038bank "github.com/ocea/sdk/extra/bank/legacy/v0_38"
	v040bank "github.com/ocea/sdk/extra/bank/legacy/v0_40"

	"github.com/stretchr/testify/require"
)

func TestMigrate(t *testing.T) {
	v040Codec := codec.NewLegacyAmino()
	cryptocodec.RegisterCrypto(v040Codec)
	v039auth.RegisterLegacyAminoCodec(v040Codec)

	coins := sdk.NewCoins(sdk.NewInt64Coin("stake", 50))
	addr1, _ := sdk.AccAddressFromBech32("ocea1xxkueklal9vejv9unqu80w9vptyepfa95pd53u")
	acc1 := v038auth.NewBaseAccount(addr1, coins, nil, 1, 0)

	addr2, _ := sdk.AccAddressFromBech32("ocea15v50ymp6n5dn73erkqtmq0u8adpl8d3ujv2e74")
	vaac := v038auth.NewContinuousVestingAccountRaw(
		v038auth.NewBaseVestingAccount(
			v038auth.NewBaseAccount(addr2, coins, nil, 1, 0), coins, nil, nil, 3160620846,
		),
		1580309972,
	)

	bankGenState := v038bank.GenesisState{
		SendEnabled: true,
	}
	authGenState := v039auth.GenesisState{
		Accounts: v038auth.GenesisAccounts{acc1, vaac},
	}

	migrated := v040bank.Migrate(bankGenState, authGenState)
	expected := `{
  "send_enabled": true,
  "balances": [
    {
      "address": "ocea1xxkueklal9vejv9unqu80w9vptyepfa95pd53u",
      "coins": [
        {
          "denom": "stake",
          "amount": "50"
        }
      ]
    },
    {
      "address": "ocea15v50ymp6n5dn73erkqtmq0u8adpl8d3ujv2e74",
      "coins": [
        {
          "denom": "stake",
          "amount": "50"
        }
      ]
    }
  ]
}`

	bz, err := v040Codec.MarshalJSONIndent(migrated, "", "  ")
	require.NoError(t, err)
	require.Equal(t, expected, string(bz))
}
