package keeper

import (
	"github.com/TessorNetwork/dredger/x/records/types"
)

var _ types.QueryServer = Keeper{}
