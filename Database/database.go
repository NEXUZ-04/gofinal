package database

import (
	"database/sql"

	customer "github.com/NEXUZ-04/gofinal/Model"
	_ "github.com/lib/pq"
)

type DB struct {
	Table string
	db    *sql.DB
}

func (d *DB) Connect(URL string) error {
	var err error

	d.db, err = sql.Open("postgres", URL)
	if err != nil {
		return err
	}

	return nil
}

func (d *DB) Abort() error {
	var err error

	err = d.db.Close()
	if err != nil {
		return err
	}

	return nil
}

func (d *DB) CreateTB() error {

	_, err := d.db.Exec(`CREATE TABLE IF NOT EXISTS ` + d.Table + ` (id SERIAL PRIMARY KEY,name TEXT,email TEXT, status TEXT);`)
	if err != nil {
		return err
	}

	return nil
}

func (d *DB) QueryAll() ([]customer.Profile, error) {

	var stmt *sql.Stmt
	var rows *sql.Rows
	var err error

	stmt, err = d.db.Prepare("SELECT id, name, email, status FROM " + d.Table)
	if err != nil {
		return nil, err
	}

	rows, err = stmt.Query()
	if err != nil {
		return nil, err
	}

	var pfs = []customer.Profile{}
	for rows.Next() {

		var pf customer.Profile
		var id int
		var name string
		var email string
		var status string

		err = rows.Scan(&id, &name, &email, &status)
		if err != nil {
			return nil, err
		}

		pf = customer.Profile{id, name, email, status}
		pfs = append(pfs, pf)
	}

	return pfs, nil
}

func (d *DB) Query(rowId int) (customer.Profile, error) {

	var stmt *sql.Stmt
	var row *sql.Row
	var err error

	stmt, err = d.db.Prepare("SELECT id, name, email, status FROM " + d.Table + " WHERE id=$1")
	if err != nil {
		return customer.Profile{}, err
	}

	row = stmt.QueryRow(rowId)
	if err != nil {
		return customer.Profile{}, err
	}

	var pf customer.Profile
	var id int
	var name string
	var email string
	var status string

	err = row.Scan(&id, &name, &email, &status)
	if err != nil {
		return customer.Profile{}, err
	}

	pf = customer.Profile{id, name, email, status}
	return pf, nil
}

func (d *DB) Insert(p customer.Profile) (customer.Profile, error) {

	var row *sql.Row
	var err error

	row = d.db.QueryRow("INSERT INTO "+d.Table+" (name, email, status) values ($1, $2, $3) RETURNING id", p.Name, p.Email, p.Status)
	err = row.Scan(&p.ID)
	if err != nil {
		return customer.Profile{}, err
	}

	return p, nil
}

func (d *DB) Update(p customer.Profile) (customer.Profile, error) {

	var stmt *sql.Stmt
	var err error

	stmt, err = d.db.Prepare("UPDATE " + d.Table + " SET name=$2, email=$3, status=$4 WHERE id=$1;")
	if err != nil {
		return customer.Profile{}, err
	}

	if _, err = stmt.Exec(p.ID, p.Name, p.Email, p.Status); err != nil {
		return customer.Profile{}, err
	}

	return p, nil
}

func (d *DB) Delete(rowId int) error {

	var stmt *sql.Stmt
	var err error

	stmt, err = d.db.Prepare("DELETE FROM " + d.Table + " WHERE id=$1;")
	if err != nil {
		return err
	}

	if _, err = stmt.Exec(rowId); err != nil {
		return err
	}

	return nil
}
