package cfg

type AppConfig struct {
	IsDev bool
}

type AppConfigOption func(*AppConfig)

func NewAppConfig(IsDev bool, opts ...AppConfigOption) *AppConfig {
	cfg := AppConfig{
		IsDev: IsDev,
	}

	for _, o := range opts {
		o(&cfg)
	}

	return &cfg
}

func (cfg *AppConfig) OverideConfigs(opts ...AppConfigOption) {
	for _, o := range opts {
		o(cfg)
	}
}
