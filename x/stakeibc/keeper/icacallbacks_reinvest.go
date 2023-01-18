package keeper

import (
	"fmt"

	"github.com/TessorNetwork/dredger/utils"
	epochtypes "github.com/TessorNetwork/dredger/x/epochs/types"
	icacallbackstypes "github.com/TessorNetwork/dredger/x/icacallbacks/types"
	recordstypes "github.com/TessorNetwork/dredger/x/records/types"
	"github.com/TessorNetwork/dredger/x/stakeibc/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	channeltypes "github.com/cosmos/ibc-go/v5/modules/core/04-channel/types"
	"github.com/golang/protobuf/proto" //nolint:staticcheck
)

// Marshalls reinvest callback arguments
func (k Keeper) MarshalReinvestCallbackArgs(ctx sdk.Context, reinvestCallback types.ReinvestCallback) ([]byte, error) {
	out, err := proto.Marshal(&reinvestCallback)
	if err != nil {
		k.Logger(ctx).Error(fmt.Sprintf("MarshalReinvestCallbackArgs %v", err.Error()))
		return nil, err
	}
	return out, nil
}

// Unmarshalls reinvest callback arguments into a ReinvestCallback struct
func (k Keeper) UnmarshalReinvestCallbackArgs(ctx sdk.Context, reinvestCallback []byte) (*types.ReinvestCallback, error) {
	unmarshalledReinvestCallback := types.ReinvestCallback{}
	if err := proto.Unmarshal(reinvestCallback, &unmarshalledReinvestCallback); err != nil {
		k.Logger(ctx).Error(fmt.Sprintf("UnmarshalReinvestCallbackArgs %s", err.Error()))
		return nil, err
	}
	return &unmarshalledReinvestCallback, nil
}

// ICA Callback after reinvestment
//   If successful:
//      * Creates a new DepositRecord with the reinvestment amount
//   If timeout/failure:
//      * Does nothing
func ReinvestCallback(k Keeper, ctx sdk.Context, packet channeltypes.Packet, ackResponse *icacallbackstypes.AcknowledgementResponse, args []byte) error {
	// Fetch callback args
	reinvestCallback, err := k.UnmarshalReinvestCallbackArgs(ctx, args)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrUnmarshalFailure, fmt.Sprintf("Unable to unmarshal reinvest callback args: %s", err.Error()))
	}
	chainId := reinvestCallback.HostZoneId
	k.Logger(ctx).Info(utils.LogICACallbackWithHostZone(chainId, ICACallbackID_Reinvest, "Starting reinvest callback"))

	// Check for timeout (ack nil)
	// No action is necessary on a timeout
	if ackResponse.Status == icacallbackstypes.AckResponseStatus_TIMEOUT {
		k.Logger(ctx).Error(utils.LogICACallbackStatusWithHostZone(chainId, ICACallbackID_Reinvest,
			icacallbackstypes.AckResponseStatus_TIMEOUT, packet))
		return nil
	}

	// Check for a failed transaction (ack error)
	// No action is necessary on a failure
	if ackResponse.Status == icacallbackstypes.AckResponseStatus_FAILURE {
		k.Logger(ctx).Error(utils.LogICACallbackStatusWithHostZone(chainId, ICACallbackID_Reinvest,
			icacallbackstypes.AckResponseStatus_FAILURE, packet))
		return nil
	}

	k.Logger(ctx).Info(utils.LogICACallbackStatusWithHostZone(chainId, ICACallbackID_Reinvest,
		icacallbackstypes.AckResponseStatus_SUCCESS, packet))

	// Get the current dredger epoch number
	dredgerEpochTracker, found := k.GetEpochTracker(ctx, epochtypes.DREDGER_EPOCH)
	if !found {
		k.Logger(ctx).Error("failed to find epoch")
		return sdkerrors.Wrapf(types.ErrInvalidLengthEpochTracker, "no number for epoch (%s)", epochtypes.DREDGER_EPOCH)
	}

	// Create a new deposit record so that rewards are reinvested
	record := recordstypes.DepositRecord{
		Amount:             reinvestCallback.ReinvestAmount.Amount,
		Denom:              reinvestCallback.ReinvestAmount.Denom,
		HostZoneId:         reinvestCallback.HostZoneId,
		Status:             recordstypes.DepositRecord_DELEGATION_QUEUE,
		Source:             recordstypes.DepositRecord_WITHDRAWAL_ICA,
		DepositEpochNumber: dredgerEpochTracker.EpochNumber,
	}
	k.RecordsKeeper.AppendDepositRecord(ctx, record)

	return nil
}
