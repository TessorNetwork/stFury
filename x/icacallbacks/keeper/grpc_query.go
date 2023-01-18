package keeper

import (
	"github.com/TessorNetwork/dredger/x/icacallbacks/types"
)

var _ types.QueryServer = Keeper{}
