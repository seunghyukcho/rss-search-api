package rsserver

import (
	"database/sql"
	"fmt"
)

type DB struct {
	conn *sql.DB
}

func (db *DB) Open(name, address, id, password string) (err error) {
	dbInfo := fmt.Sprintf("%s:%s@tcp(%s)/%s", id, password, address, name)
	db.conn, err = sql.Open("mysql", dbInfo)

	return err
}

func (db *DB) Close() (err error) {
	err = db.conn.Close()
	return err
}
