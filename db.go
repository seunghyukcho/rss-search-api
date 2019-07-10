package main

import (
	"database/sql"
)

type Env struct {
	db *sql.DB
}
