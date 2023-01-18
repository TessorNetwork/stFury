package keeper

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	strideapp "github.com/TessorNetwork/dredger/app"
	"github.com/TessorNetwork/dredger/x/records/keeper"
)

func RecordsKeeper(t testing.TB) (*keeper.Keeper, sdk.Context) {
	app := strideapp.InitStrideTestApp(true)
	recordKeeper := app.RecordsKeeper
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: 1, ChainID: "dredger-1", Time: time.Now().UTC()})

	return &recordKeeper, ctx
}
