package database

import (
	"database/sql"
)

type DB struct {
	db *sql.DB
}

func (d *DB) Connect() error {
	var err error
	d.db, err = sql.Open("postgres", "postgres://oeenwpjw:zQqtsCaL5VoY3x-NX8dbzfH7unkJ9Lb0@hansken.db.elephantsql.com:5432/oeenwpjw")
	if err != nil {
		return err
	}

	return nil
}

func (d *DB) abort() error {
	var err error
	err = d.db.Close()
	if err != nil {
		return err
	}

	return nil
}