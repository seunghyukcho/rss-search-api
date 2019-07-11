package rsserver

import (
	"database/sql"
	"fmt"
)

type DB struct {
	Connection *sql.DB
}

func (db *DB) Open(name, address, id, password string) (err error) {
	dbInfo := fmt.Sprintf("%s:%s@tcp(%s)/%s", id, password, address, name)
	db.Connection, err = sql.Open("mysql", dbInfo)

	return err
}

func (db *DB) Close() (err error) {
	err = db.Connection.Close()
	return err
}
