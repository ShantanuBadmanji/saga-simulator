package env

type envVar string

const (
	SqliteDbDir = envVar("SQLITE_DB_DIR")
	Port        = envVar("PORT")
)
