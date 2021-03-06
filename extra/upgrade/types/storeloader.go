package types

import (
	"github.com/ocea/sdk/baseapp"
	store "github.com/ocea/sdk/store/types"
	sdk "github.com/ocea/sdk/types"
)

// UpgradeStoreLoader is used to prepare baseapp with a fixed StoreLoader
// pattern. This is useful for custom upgrade loading logic.
func UpgradeStoreLoader(upgradeHeight int64, storeUpgrades *store.StoreUpgrades) baseapp.StoreLoader {
	return func(ms sdk.CommitMultiStore) error {
		if upgradeHeight == ms.LastCommitID().Version {
			// Check if the current commit version and upgrade height matches
			if len(storeUpgrades.Renamed) > 0 || len(storeUpgrades.Deleted) > 0 {
				return ms.LoadLatestVersionAndUpgrade(storeUpgrades)
			}
		}

		// Otherwise load default store loader
		return baseapp.DefaultStoreLoader(ms)
	}
}
