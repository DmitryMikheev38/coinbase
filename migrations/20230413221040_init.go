package migrations

import (
	"database/sql"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upInit, downInit)
}

func upInit(tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.Exec(`
	CREATE TABLE ticks (
	    timestamp bigint NOT NULL,
		symbol varchar(7) NOT NULL,
		best_bid double NOT NULL,
		best_ask double NOT NULL
	);`)

	return err
}

func downInit(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	_, err := tx.Exec(`DROP TABLE ticks;`)
	return err
}
