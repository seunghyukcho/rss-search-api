package handle

import "database/sql"

func Transaction(tx *sql.Tx, e error) (err error) {
	if err := tx.Rollback(); err != nil {
		return err
	}

	return e
}
