package config

import "time"

type Config struct {
	Postgres struct {
		Driver     string `mapstructure:"driver"`
		Source     string `mapstructure:"source"`
		MigrateUrl string `mapstructure:"migrate_url"`
	}
	Server struct {
		GrpcServerAddress string `mapstructure:"grpc_server_address"`
		HttpServerAddress string `mapstructure:"http_server_address"`
	}
	Token struct {
		TokenSymmetricKey   string        `mapstructure:"token_symmetric_key"`
		AccessTokenDuration time.Duration `mapstructure:"access_token_duration"`
	}
}
