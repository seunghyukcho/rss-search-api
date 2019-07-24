package handle

import "database/sql"

// Transaction is used to handle transaction errors when an error occur.
func Transaction(tx *sql.Tx, e error) (err error) {
	if err := tx.Rollback(); err != nil {
		return err
	}

	return e
}
