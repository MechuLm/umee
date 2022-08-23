package app

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/group"
	"github.com/cosmos/cosmos-sdk/x/nft"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	bech32ibckeeper "github.com/osmosis-labs/bech32-ibc/x/bech32ibc/keeper"
	bech32ibctypes "github.com/osmosis-labs/bech32-ibc/x/bech32ibc/types"
	leveragetypes "github.com/umee-network/umee/v2/x/leverage/types"
	oracletypes "github.com/umee-network/umee/v2/x/oracle/types"
)

const UpgradeV3_0Plan = "v1.0-v3.0"

func (app UmeeApp) RegisterUpgradeHandlers() {
	// v3 upgrade handler performs upgrade from v1->v3
	app.UpgradeKeeper.SetUpgradeHandler(
		UpgradeV3_0Plan,
		func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			ctx.Logger().Info("Upgrade handler execution", "name", UpgradeV3_0Plan)
			err := setupBech32ibcKeeper(&app.bech32IbcKeeper, ctx)
			if err != nil {
				return nil, sdkerrors.Wrapf(
					err, "Calypso %q Upgrade: Unable to upgrade, bech32ibc module not initialized", UpgradeV3_0Plan)
			}

			ctx.Logger().Info("Upgrade handler execution finished, running migrations", "name", UpgradeV3_0Plan)
			return app.mm.RunMigrations(ctx, app.configurator, fromVM)
		})

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}

	if upgradeInfo.Name == UpgradeV3_0Plan && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := storetypes.StoreUpgrades{
			Added: []string{
				group.ModuleName,
				nft.ModuleName,
				bech32ibctypes.ModuleName,
				oracletypes.ModuleName,
				leveragetypes.ModuleName,
			},
		}
		// configure store loader that checks if version == upgradeHeight and applies store upgrades
		app.SetStoreLoader(
			upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}
}

// Sets up bech32ibc module by setting the native account prefix to "umee".
// Failing to set the native prefix will cause a chain halt on init genesis or
// in the firstBeginBlocker assertions.
func setupBech32ibcKeeper(bech32IbcKeeper *bech32ibckeeper.Keeper, ctx sdk.Context) error {
	return bech32IbcKeeper.SetNativeHrp(ctx, sdk.GetConfig().GetBech32AccountAddrPrefix())
}