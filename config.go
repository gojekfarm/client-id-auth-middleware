package clientauth

type Config struct {
	dbDriver  string
	dbConnURL string
}

func NewConfig(dbDriver, dbConnURL string) *Config {
	return &Config{
		dbDriver:  dbDriver,
		dbConnURL: dbConnURL,
	}
}
