package postgres

import (
	"WB_TEST/config"
	"fmt"
	_ "github.com/jackc/pgx/stdlib" // pgx driver
	"github.com/jmoiron/sqlx"
)

// Return PostgresSQL instance
func NewConnToPostgres(c *config.Config) (*sqlx.DB, error) {
	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		c.PostgresConfig.PostgresqlHost,
		c.PostgresConfig.PostgresqlPort,
		c.PostgresConfig.PostgresqlUser,
		c.PostgresConfig.PostgresqlDbname,
		c.PostgresConfig.PostgresqlPassword,
	)

	db, err := sqlx.Connect(c.PostgresConfig.PgDriver, dataSourceName)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
