#!/bin/bash

echo "-----------------------"
echo "## Add new CosmWasm CW20 contract"

SCRIPTPATH="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
ARTIFACTS_PATH=$SCRIPTPATH/../../artifacts/
CHAIN_ID="umee-local-beta-testnet"

RESP=$(umeed tx wasm store $ARTIFACTS_PATH/cw20_base.wasm --chain-id $CHAIN_ID --from alice --keyring-backend test --gas 100000000 -y)
CODE_ID=$(echo "$RESP" | jq -r '.logs[0].events[1].attributes[-1].value')
ALICE_ADDR="$(umeed keys show alice --keyring-backend=test -a)"
BOB_ADDR="$(umeed keys show bob --keyring-backend=test -a)"

echo "* Code id: $CODE_ID"
echo "* Download code"
TMPDIR=$(mktemp -t wasmdXXXXXX)
umeed q wasm code "$CODE_ID" "$TMPDIR"
echo "-----------------------"
echo "## List code"
umeed query wasm list-code --chain-id $CHAIN_ID


echo "-----------------------"
echo "## Create new contract instance"
INIT="{\"decimals\": 2, \"name\":\"CashName\", \"symbol\": \"SYMBOL\", \"initial_balances\":[{\"address\": \"$ALICE_ADDR\", \"amount\": \"64534\"}]}"
echo "----------INIT: $INIT"
(umeed tx wasm instantiate "$CODE_ID" "$INIT" --admin="$(umeed keys show alice -a --keyring-backend=test)" \
  --from alice --keyring-backend test --amount="100uumee" --label "test-cw20-rafilx" \
  --gas 1000000 -y --chain-id $CHAIN_ID)


echo "Contract Information"

CONTRACT_ADDR=$(umeed query wasm list-contract-by-code "$CODE_ID" -o json --chain-id $CHAIN_ID  | jq -r '.contracts[-1]')
echo "* Contract address: $CONTRACT_ADDR"
echo "### Query all"
RESP=$(umeed query wasm contract-state all "$CONTRACT_ADDR" -o json --chain-id $CHAIN_ID )
echo "$RESP" | jq
echo "### Query smart - check balance"
umeed query wasm contract-state smart "$CONTRACT_ADDR" "{\"balance\":{\"address\": \"$ALICE_ADDR\"}}" -o json --chain-id $CHAIN_ID  | jq
echo "### Query raw"
KEY=$(echo "$RESP" | jq -r ".models[0].key")
umeed query wasm contract-state raw "$CONTRACT_ADDR" "$KEY" -o json --chain-id $CHAIN_ID  | jq

echo "-----------------------"
echo "## Execute contract $CONTRACT_ADDR"
AMOUNT_TO_TRANSFER=123
MSG_TRANSFER="{\"transfer\": {\"recipient\": \"$BOB_ADDR\", \"amount\": \"$AMOUNT_TO_TRANSFER\"}}"
echo "## Sending $AMOUNT_TO_TRANSFER from $ALICE_ADDR to $BOB_ADDR"
umeed tx wasm execute "$CONTRACT_ADDR" "$MSG_TRANSFER" --from alice --keyring-backend test -o json --chain-id $CHAIN_ID -y | jq


echo "-----------------------"
echo "## Check balance of bob $BOB_ADDR in the contract $CONTRACT_ADDR"
BOBS_BALANCE_DATA="$(umeed query wasm contract-state smart "$CONTRACT_ADDR" "{\"balance\":{\"address\": \"$BOB_ADDR\"}}" -o json --chain-id $CHAIN_ID)"
echo $BOBS_BALANCE_DATA | jq

BOB_BALANCE=$(echo $BOBS_BALANCE_DATA | jq -r '.data.balance')
echo "Bob balance $BOB_BALANCE"

echo "-----------------------"

if [ "$AMOUNT_TO_TRANSFER" != "$BOB_BALANCE" ]; then
  echo "Amounts are not equal $AMOUNT_TO_TRANSFER != $BOB_BALANCE"
  exit 1
fi
echo "Amounts are equal! =D $AMOUNT_TO_TRANSFER == $BOB_BALANCE"