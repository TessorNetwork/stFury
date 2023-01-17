package keeper

import (
	"github.com/TessorNetwork/dredger/v4/x/stakeibc/types"
)

var _ types.QueryServer = Keeper{}
