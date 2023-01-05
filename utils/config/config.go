package config

type Config struct {
	Postgres struct {
		Driver    string `mapstructure:"driver"`
		Source    string `mapstructure:"source"`
		Migration struct {
			MigrateUrl string `mapstructure:"migrate_url"`
		} `mapstructure:"migration"`
	}
}
