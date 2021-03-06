package keeper_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/ocea/sdk/simapp"
	store "github.com/ocea/sdk/store/types"
	"github.com/ocea/sdk/extra/upgrade/keeper"
	"github.com/ocea/sdk/extra/upgrade/types"
)

type KeeperTestSuite struct {
	suite.Suite

	homeDir string
	app     *simapp.SimApp
}

func (s *KeeperTestSuite) SetupTest() {
	app := simapp.Setup(false)
	homeDir := filepath.Join(s.T().TempDir(), "x_upgrade_keeper_test")
	app.UpgradeKeeper = keeper.NewKeeper( // recreate keeper in order to use a custom home path
		make(map[int64]bool), app.GetKey(types.StoreKey), app.AppCodec(), homeDir,
	)
	s.T().Log("home dir:", homeDir)
	s.homeDir = homeDir
	s.app = app
}

func (s *KeeperTestSuite) TestReadUpgradeInfoFromDisk() {
	// require no error when the upgrade info file does not exist
	_, err := s.app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	s.Require().NoError(err)

	expected := store.UpgradeInfo{
		Name:   "test_upgrade",
		Height: 100,
	}

	// create an upgrade info file
	s.Require().NoError(s.app.UpgradeKeeper.DumpUpgradeInfoToDisk(expected.Height, expected.Name))

	ui, err := s.app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	s.Require().NoError(err)
	s.Require().Equal(expected, ui)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
