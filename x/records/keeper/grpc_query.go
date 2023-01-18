package keeper

import (
	"github.com/TessorNetwork/dredger/v4/x/records/types"
)

var _ types.QueryServer = Keeper{}
