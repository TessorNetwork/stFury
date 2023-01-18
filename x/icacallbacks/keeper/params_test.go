package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	testkeeper "github.com/TessorNetwork/dredger/testutil/keeper"
	"github.com/TessorNetwork/dredger/x/icacallbacks/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.IcacallbacksKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
