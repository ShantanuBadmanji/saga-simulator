package env

type envVar string

const (
	SQLITE_DB_DIR = envVar("SQLITE_DB_DIR")
	PORT          = envVar("PORT")
	LOG_DIR       = envVar("LOG_DIR")
)
