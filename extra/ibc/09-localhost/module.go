package localhost

import (
	"github.com/ocea/sdk/extra/ibc/09-localhost/types"
)

// Name returns the IBC client name
func Name() string {
	return types.SubModuleName
}
