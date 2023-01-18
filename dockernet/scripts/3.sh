### LIQ STAKE 
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

source ${SCRIPT_DIR}/../config.sh

$DREDGER_MAIN_CMD tx stakeibc redeem-stake 89 GAIA $GAIA_RECEIVER_ADDRESS --from ${DREDGER_VAL_PREFIX}1 -y
