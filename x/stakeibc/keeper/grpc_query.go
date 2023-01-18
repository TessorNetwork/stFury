package keeper

import (
	"github.com/TessorNetwork/dredger/x/stakeibc/types"
)

var _ types.QueryServer = Keeper{}
