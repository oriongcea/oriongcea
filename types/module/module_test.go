package module_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/ocea/sdk/codec/types"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/ocea/sdk/client"
	"github.com/ocea/sdk/codec"
	"github.com/ocea/sdk/tests/mocks"
	sdk "github.com/ocea/sdk/types"
	"github.com/ocea/sdk/types/module"
)

var errFoo = errors.New("dummy")

func TestBasicManager(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	t.Cleanup(mockCtrl.Finish)
	legacyAmino := codec.NewLegacyAmino()
	interfaceRegistry := types.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(interfaceRegistry)

	clientCtx := client.Context{}
	clientCtx = clientCtx.WithLegacyAmino(legacyAmino)
	wantDefaultGenesis := map[string]json.RawMessage{"mockAppModuleBasic1": json.RawMessage(``)}

	mockAppModuleBasic1 := mocks.NewMockAppModuleBasic(mockCtrl)

	mockAppModuleBasic1.EXPECT().Name().AnyTimes().Return("mockAppModuleBasic1")
	mockAppModuleBasic1.EXPECT().DefaultGenesis(gomock.Eq(cdc)).Times(1).Return(json.RawMessage(``))
	mockAppModuleBasic1.EXPECT().ValidateGenesis(gomock.Eq(cdc), gomock.Eq(nil), gomock.Eq(wantDefaultGenesis["mockAppModuleBasic1"])).Times(1).Return(errFoo)
	mockAppModuleBasic1.EXPECT().RegisterRESTRoutes(gomock.Eq(client.Context{}), gomock.Eq(&mux.Router{})).Times(1)
	mockAppModuleBasic1.EXPECT().RegisterLegacyAminoCodec(gomock.Eq(legacyAmino)).Times(1)
	mockAppModuleBasic1.EXPECT().RegisterInterfaces(gomock.Eq(interfaceRegistry)).Times(1)
	mockAppModuleBasic1.EXPECT().GetTxCmd().Times(1).Return(nil)
	mockAppModuleBasic1.EXPECT().GetQueryCmd().Times(1).Return(nil)

	mm := module.NewBasicManager(mockAppModuleBasic1)
	require.Equal(t, mm["mockAppModuleBasic1"], mockAppModuleBasic1)

	mm.RegisterLegacyAminoCodec(legacyAmino)
	mm.RegisterInterfaces(interfaceRegistry)

	require.Equal(t, wantDefaultGenesis, mm.DefaultGenesis(cdc))

	var data map[string]string
	require.Equal(t, map[string]string(nil), data)

	require.True(t, errors.Is(errFoo, mm.ValidateGenesis(cdc, nil, wantDefaultGenesis)))

	mm.RegisterRESTRoutes(client.Context{}, &mux.Router{})

	mockCmd := &cobra.Command{Use: "root"}
	mm.AddTxCommands(mockCmd)

	mm.AddQueryCommands(mockCmd)

	// validate genesis returns nil
	require.Nil(t, module.NewBasicManager().ValidateGenesis(cdc, nil, wantDefaultGenesis))
}

func TestGenesisOnlyAppModule(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	t.Cleanup(mockCtrl.Finish)

	mockModule := mocks.NewMockAppModuleGenesis(mockCtrl)
	mockInvariantRegistry := mocks.NewMockInvariantRegistry(mockCtrl)
	goam := module.NewGenesisOnlyAppModule(mockModule)

	require.True(t, goam.Route().Empty())
	require.Empty(t, goam.QuerierRoute())
	require.Nil(t, goam.LegacyQuerierHandler(nil))

	// no-op
	goam.RegisterInvariants(mockInvariantRegistry)
	goam.BeginBlock(sdk.Context{}, abci.RequestBeginBlock{})
	require.Equal(t, []abci.ValidatorUpdate{}, goam.EndBlock(sdk.Context{}, abci.RequestEndBlock{}))
}

func TestManagerOrderSetters(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	t.Cleanup(mockCtrl.Finish)
	mockAppModule1 := mocks.NewMockAppModule(mockCtrl)
	mockAppModule2 := mocks.NewMockAppModule(mockCtrl)

	mockAppModule1.EXPECT().Name().Times(2).Return("module1")
	mockAppModule2.EXPECT().Name().Times(2).Return("module2")
	mm := module.NewManager(mockAppModule1, mockAppModule2)
	require.NotNil(t, mm)
	require.Equal(t, 2, len(mm.Modules))

	require.Equal(t, []string{"module1", "module2"}, mm.OrderInitGenesis)
	mm.SetOrderInitGenesis("module2", "module1")
	require.Equal(t, []string{"module2", "module1"}, mm.OrderInitGenesis)

	require.Equal(t, []string{"module1", "module2"}, mm.OrderExportGenesis)
	mm.SetOrderExportGenesis("module2", "module1")
	require.Equal(t, []string{"module2", "module1"}, mm.OrderExportGenesis)

	require.Equal(t, []string{"module1", "module2"}, mm.OrderBeginBlockers)
	mm.SetOrderBeginBlockers("module2", "module1")
	require.Equal(t, []string{"module2", "module1"}, mm.OrderBeginBlockers)

	require.Equal(t, []string{"module1", "module2"}, mm.OrderEndBlockers)
	mm.SetOrderEndBlockers("module2", "module1")
	require.Equal(t, []string{"module2", "module1"}, mm.OrderEndBlockers)
}

func TestManager_RegisterInvariants(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	t.Cleanup(mockCtrl.Finish)

	mockAppModule1 := mocks.NewMockAppModule(mockCtrl)
	mockAppModule2 := mocks.NewMockAppModule(mockCtrl)
	mockAppModule1.EXPECT().Name().Times(2).Return("module1")
	mockAppModule2.EXPECT().Name().Times(2).Return("module2")
	mm := module.NewManager(mockAppModule1, mockAppModule2)
	require.NotNil(t, mm)
	require.Equal(t, 2, len(mm.Modules))

	// test RegisterInvariants
	mockInvariantRegistry := mocks.NewMockInvariantRegistry(mockCtrl)
	mockAppModule1.EXPECT().RegisterInvariants(gomock.Eq(mockInvariantRegistry)).Times(1)
	mockAppModule2.EXPECT().RegisterInvariants(gomock.Eq(mockInvariantRegistry)).Times(1)
	mm.RegisterInvariants(mockInvariantRegistry)
}

func TestManager_RegisterRoutes(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	t.Cleanup(mockCtrl.Finish)

	mockAppModule1 := mocks.NewMockAppModule(mockCtrl)
	mockAppModule2 := mocks.NewMockAppModule(mockCtrl)
	mockAppModule1.EXPECT().Name().Times(2).Return("module1")
	mockAppModule2.EXPECT().Name().Times(2).Return("module2")
	mm := module.NewManager(mockAppModule1, mockAppModule2)
	require.NotNil(t, mm)
	require.Equal(t, 2, len(mm.Modules))

	router := mocks.NewMockRouter(mockCtrl)
	handler1, handler2 := sdk.Handler(func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		return nil, nil
	}), sdk.Handler(func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		return nil, nil
	})
	route1 := sdk.NewRoute("route1", handler1)
	route2 := sdk.NewRoute("route2", handler2)
	mockAppModule1.EXPECT().Route().Times(2).Return(route1)
	mockAppModule2.EXPECT().Route().Times(2).Return(route2)
	router.EXPECT().AddRoute(gomock.Any()).Times(2) // Use of Any due to limitations to compare Functions as the sdk.Handler

	queryRouter := mocks.NewMockQueryRouter(mockCtrl)
	mockAppModule1.EXPECT().QuerierRoute().Times(2).Return("querierRoute1")
	mockAppModule2.EXPECT().QuerierRoute().Times(1).Return("")
	handler3 := sdk.Querier(nil)
	mockAppModule1.EXPECT().NewQuerierHandler().Times(1).Return(handler3)
	queryRouter.EXPECT().AddRoute(gomock.Eq("querierRoute1"), gomock.Eq(handler3)).Times(1)

	amino := codec.NewLegacyAmino()
	mm.RegisterRoutes(router, queryRouter, amino)
}

func TestManager_RegisterQueryServices(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	t.Cleanup(mockCtrl.Finish)

	mockAppModule1 := mocks.NewMockAppModule(mockCtrl)
	mockAppModule2 := mocks.NewMockAppModule(mockCtrl)
	mockAppModule1.EXPECT().Name().Times(2).Return("module1")
	mockAppModule2.EXPECT().Name().Times(2).Return("module2")
	mm := module.NewManager(mockAppModule1, mockAppModule2)
	require.NotNil(t, mm)
	require.Equal(t, 2, len(mm.Modules))

	queryRouter := mocks.NewMockServer(mockCtrl)
	mockAppModule1.EXPECT().RegisterQueryService(queryRouter).Times(1)
	mockAppModule2.EXPECT().RegisterQueryService(queryRouter).Times(1)

	mm.RegisterQueryServices(queryRouter)
}

func TestManager_InitGenesis(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	t.Cleanup(mockCtrl.Finish)

	mockAppModule1 := mocks.NewMockAppModule(mockCtrl)
	mockAppModule2 := mocks.NewMockAppModule(mockCtrl)
	mockAppModule1.EXPECT().Name().Times(2).Return("module1")
	mockAppModule2.EXPECT().Name().Times(2).Return("module2")
	mm := module.NewManager(mockAppModule1, mockAppModule2)
	require.NotNil(t, mm)
	require.Equal(t, 2, len(mm.Modules))

	ctx := sdk.Context{}
	interfaceRegistry := types.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(interfaceRegistry)
	genesisData := map[string]json.RawMessage{"module1": json.RawMessage(`{"key": "value"}`)}

	mockAppModule1.EXPECT().InitGenesis(gomock.Eq(ctx), gomock.Eq(cdc), gomock.Eq(genesisData["module1"])).Times(1).Return(nil)
	require.Equal(t, abci.ResponseInitChain{Validators: []abci.ValidatorUpdate(nil)}, mm.InitGenesis(ctx, cdc, genesisData))

	// test panic
	genesisData = map[string]json.RawMessage{
		"module1": json.RawMessage(`{"key": "value"}`),
		"module2": json.RawMessage(`{"key": "value"}`)}
	mockAppModule1.EXPECT().InitGenesis(gomock.Eq(ctx), gomock.Eq(cdc), gomock.Eq(genesisData["module1"])).Times(1).Return([]abci.ValidatorUpdate{{}})
	mockAppModule2.EXPECT().InitGenesis(gomock.Eq(ctx), gomock.Eq(cdc), gomock.Eq(genesisData["module2"])).Times(1).Return([]abci.ValidatorUpdate{{}})
	require.Panics(t, func() { mm.InitGenesis(ctx, cdc, genesisData) })
}

func TestManager_ExportGenesis(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	t.Cleanup(mockCtrl.Finish)

	mockAppModule1 := mocks.NewMockAppModule(mockCtrl)
	mockAppModule2 := mocks.NewMockAppModule(mockCtrl)
	mockAppModule1.EXPECT().Name().Times(2).Return("module1")
	mockAppModule2.EXPECT().Name().Times(2).Return("module2")
	mm := module.NewManager(mockAppModule1, mockAppModule2)
	require.NotNil(t, mm)
	require.Equal(t, 2, len(mm.Modules))

	ctx := sdk.Context{}
	interfaceRegistry := types.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(interfaceRegistry)
	mockAppModule1.EXPECT().ExportGenesis(gomock.Eq(ctx), gomock.Eq(cdc)).Times(1).Return(json.RawMessage(`{"key1": "value1"}`))
	mockAppModule2.EXPECT().ExportGenesis(gomock.Eq(ctx), gomock.Eq(cdc)).Times(1).Return(json.RawMessage(`{"key2": "value2"}`))

	want := map[string]json.RawMessage{
		"module1": json.RawMessage(`{"key1": "value1"}`),
		"module2": json.RawMessage(`{"key2": "value2"}`)}
	require.Equal(t, want, mm.ExportGenesis(ctx, cdc))
}

func TestManager_BeginBlock(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	t.Cleanup(mockCtrl.Finish)

	mockAppModule1 := mocks.NewMockAppModule(mockCtrl)
	mockAppModule2 := mocks.NewMockAppModule(mockCtrl)
	mockAppModule1.EXPECT().Name().Times(2).Return("module1")
	mockAppModule2.EXPECT().Name().Times(2).Return("module2")
	mm := module.NewManager(mockAppModule1, mockAppModule2)
	require.NotNil(t, mm)
	require.Equal(t, 2, len(mm.Modules))

	req := abci.RequestBeginBlock{Hash: []byte("test")}

	mockAppModule1.EXPECT().BeginBlock(gomock.Any(), gomock.Eq(req)).Times(1)
	mockAppModule2.EXPECT().BeginBlock(gomock.Any(), gomock.Eq(req)).Times(1)
	mm.BeginBlock(sdk.Context{}, req)
}

func TestManager_EndBlock(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	t.Cleanup(mockCtrl.Finish)

	mockAppModule1 := mocks.NewMockAppModule(mockCtrl)
	mockAppModule2 := mocks.NewMockAppModule(mockCtrl)
	mockAppModule1.EXPECT().Name().Times(2).Return("module1")
	mockAppModule2.EXPECT().Name().Times(2).Return("module2")
	mm := module.NewManager(mockAppModule1, mockAppModule2)
	require.NotNil(t, mm)
	require.Equal(t, 2, len(mm.Modules))

	req := abci.RequestEndBlock{Height: 10}

	mockAppModule1.EXPECT().EndBlock(gomock.Any(), gomock.Eq(req)).Times(1).Return([]abci.ValidatorUpdate{{}})
	mockAppModule2.EXPECT().EndBlock(gomock.Any(), gomock.Eq(req)).Times(1)
	ret := mm.EndBlock(sdk.Context{}, req)
	require.Equal(t, []abci.ValidatorUpdate{{}}, ret.ValidatorUpdates)

	// test panic
	mockAppModule1.EXPECT().EndBlock(gomock.Any(), gomock.Eq(req)).Times(1).Return([]abci.ValidatorUpdate{{}})
	mockAppModule2.EXPECT().EndBlock(gomock.Any(), gomock.Eq(req)).Times(1).Return([]abci.ValidatorUpdate{{}})
	require.Panics(t, func() { mm.EndBlock(sdk.Context{}, req) })
}