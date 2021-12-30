package dal

import (
	"database/sql"
	e "mowplow/api/v2/errors"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	DB         *sql.DB
	DBType     string
	DBName     string
	DBUser     string
	DBPassword string
}

func NewDB() *DB {
	// Hard coding for now. eventually will be replaced with
	// dynamic customer specific creds.
	db := DB{
		DBType:     "mysql",
		DBName:     "mowplow",
		DBUser:     "mowplow",
		DBPassword: "",
	}
	return &db
}

func (d *DB) OpenDB() error {
	db, err := sql.Open(d.DBType, d.DBUser+":"+d.DBPassword+"@/"+d.DBName)
	if err != nil {
		e.LogError(err)
	}

	d.DB = db
	return err
}

func (d *DB) CloseDB() error {
	err := d.DB.Close()
	return err
}
