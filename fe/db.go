package fe

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/stdlib" // postgres driver
)

var password string

const (
	driver   = "pgx"
	user     = "fe"
	database = "fatal-encounters"
)

func init() {
	password = os.Getenv("FE_DB_PW")
}

// OpenDB ...
func OpenDB() (*sql.DB, error) {
	source := fmt.Sprintf("user=%s password=%s database=%s sslmode=disable", user, password, database)
	return sql.Open(driver, source)
}
