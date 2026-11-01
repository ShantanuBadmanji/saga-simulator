package cfg

type AppConfig struct {
	Dev bool
}

type AppConfigOption func(*AppConfig)

func NewAppConfig(dev bool, opts ...AppConfigOption) *AppConfig {
	cfg := AppConfig{
		Dev: dev,
	}

	for _, o := range opts {
		o(&cfg)
	}

	return &cfg
}

func (cfg *AppConfig) Apply(opts ...AppConfigOption) {
	for _, o := range opts {
		o(cfg)
	}
}
