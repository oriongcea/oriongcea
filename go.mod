go 1.15

module github.com/ocea/sdk

require (
	github.com/99designs/keyring v1.1.6
	github.com/armon/go-metrics v0.3.4
	github.com/bgentry/speakeasy v0.1.0
	github.com/btcsuite/btcd v0.21.0-beta
	github.com/btcsuite/btcutil v1.0.2
	github.com/confio/ics23/go v0.0.0-20200817220745-f173e6211efb
	github.com/enigmampc/btcutil v1.0.3-0.20200723161021-e2fb6adb2a25
	github.com/gogo/gateway v1.1.0
	github.com/gogo/protobuf v1.3.1
	github.com/golang/mock v1.4.4
	github.com/golang/protobuf v1.4.2
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.1
	github.com/grpc-ecosystem/grpc-gateway v1.15.0
	github.com/hashicorp/golang-lru v0.5.4
	github.com/magiconair/properties v1.8.4
	github.com/mattn/go-isatty v0.0.12
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/otiai10/copy v1.2.0
	github.com/pelletier/go-toml v1.8.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.7.1
	github.com/prometheus/common v0.14.0
	github.com/rakyll/statik v0.1.7
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/cast v1.3.1
	github.com/spf13/cobra v1.0.0
	github.com/spf13/jwalterweatherman v1.1.0 // indirect; indirects
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.6.1
	github.com/tendermint/btcd v0.1.1
	github.com/tendermint/crypto v0.0.0-20191022145703-50d29ede1e15
	github.com/tendermint/go-amino v0.16.0
	github.com/tendermint/tendermint v0.34.0-rc4
	github.com/tendermint/tm-db v0.6.2
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a
	google.golang.org/genproto v0.0.0-20200825200019-8632dd797987
	google.golang.org/grpc v1.32.0
	google.golang.org/protobuf v1.25.0
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/yaml.v2 v2.3.0
	github.com/btcsuite/btcd v0.0.0-20190115013929-ed77733ec07d
    github.com/davecgh/go-spew v1.1.1 // indirect
    github.com/pkg/errors v0.8.1 // indirect
    github.com/stretchr/testify v1.3.0
    github.com/zondax/hid v0.9.0 // indirect
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4