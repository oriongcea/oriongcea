package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/ocea/sdk/types"
)

func TestMsgUnjailGetSignBytes(t *testing.T) {
	addr := sdk.AccAddress("abcd")
	msg := NewMsgUnjail(sdk.ValAddress(addr))
	bytes := msg.GetSignBytes()
	require.Equal(
		t,
		`{"type":"ocea-sdk/MsgUnjail","value":{"address":"oceavaloper1v93xxeqhg9nn6"}}`,
		string(bytes),
	)
}
