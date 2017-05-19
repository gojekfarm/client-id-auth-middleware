package clientauth

type Config struct {
	dbDriver     string
	dbConnURL    string
	HeaderConfig *HeaderConfig
}

type HeaderConfig struct {
	ClientIDName string
	PassKeyName  string
}

func NewConfig(dbDriver, dbConnURL string) *Config {
	return &Config{
		dbDriver:     dbDriver,
		dbConnURL:    dbConnURL,
		HeaderConfig: &HeaderConfig{ClientIDName: "Client-ID", PassKeyName: "Pass-Key"},
	}
}

func (cfg *Config) SetHeaderConfig(clientIDName, passKeyName string) {
	cfg.HeaderConfig.ClientIDName = clientIDName
	cfg.HeaderConfig.PassKeyName = passKeyName
}
