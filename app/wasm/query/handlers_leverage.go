package query

import (
	"fmt"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	lvtypes "github.com/umee-network/umee/v2/x/leverage/types"
)

// HandleLeverageParams handles the get the x/leverage module's parameters.
func (q UmeeQuery) HandleLeverageParams(
	ctx sdk.Context,
	qs lvtypes.QueryServer,
) ([]byte, error) {
	resp, err := qs.Params(sdk.WrapSDKContext(ctx), &lvtypes.QueryParams{})
	if err != nil {
		return nil, wasmvmtypes.UnsupportedRequest{Kind: fmt.Sprintf("error %+v to assigned query Leverage Params", err)}
	}

	return MarshalResponse(resp)
}

// HandleRegisteredTokens handles the get all registered tokens query and response.
func (q UmeeQuery) HandleRegisteredTokens(
	ctx sdk.Context,
	qs lvtypes.QueryServer,
) ([]byte, error) {
	resp, err := qs.RegisteredTokens(sdk.WrapSDKContext(ctx), &lvtypes.QueryRegisteredTokens{})
	if err != nil {
		return nil, wasmvmtypes.UnsupportedRequest{Kind: fmt.Sprintf("error %+v to assigned query Registered Tokens", err)}
	}

	return MarshalResponse(resp)
}

// HandleMarketSummary queries a base asset's current borrowing and supplying conditions.
func (q UmeeQuery) HandleMarketSummary(
	ctx sdk.Context,
	qs lvtypes.QueryServer,
) ([]byte, error) {
	resp, err := qs.MarketSummary(sdk.WrapSDKContext(ctx), q.MarketSummary)
	if err != nil {
		return nil, wasmvmtypes.UnsupportedRequest{Kind: fmt.Sprintf("error %+v to assigned query Market Summary", err)}
	}

	return MarshalResponse(resp)
}

// HandleAccountBalances queries an account's current supply, collateral, and borrow positions.
func (q UmeeQuery) HandleAccountBalances(
	ctx sdk.Context,
	qs lvtypes.QueryServer,
) ([]byte, error) {
	resp, err := qs.AccountBalances(sdk.WrapSDKContext(ctx), q.AccountBalances)
	if err != nil {
		return nil, wasmvmtypes.UnsupportedRequest{Kind: fmt.Sprintf("error %+v to assigned query Account Balances", err)}
	}

	return MarshalResponse(resp)
}

// HandleAccountSummary queries USD values representing an account's total
// positions and borrowing limits. It requires oracle prices to return successfully.
func (q UmeeQuery) HandleAccountSummary(
	ctx sdk.Context,
	qs lvtypes.QueryServer,
) ([]byte, error) {
	resp, err := qs.AccountSummary(sdk.WrapSDKContext(ctx), q.AccountSummary)
	if err != nil {
		return nil, wasmvmtypes.UnsupportedRequest{Kind: fmt.Sprintf("error %+v to assigned query Account Summary", err)}
	}

	return MarshalResponse(resp)
}

// HandleLiquidationTargets queries a list of all borrower account addresses eligible for liquidation.
func (q UmeeQuery) HandleLiquidationTargets(
	ctx sdk.Context,
	qs lvtypes.QueryServer,
) ([]byte, error) {
	resp, err := qs.LiquidationTargets(sdk.WrapSDKContext(ctx), q.LiquidationTargets)
	if err != nil {
		return nil, wasmvmtypes.UnsupportedRequest{Kind: fmt.Sprintf("error %+v to assigned query Liquidation Targets", err)}
	}

	return MarshalResponse(resp)
}
