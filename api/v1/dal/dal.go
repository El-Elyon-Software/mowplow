package dal

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	DB         *sql.DB
	DBType     string
	DBName     string
	DBUser     string
	DBPassword string
}

type CRUD interface {
	Create()
	Retrieve()
	RetrieveFilter()
	Update()
	Delete()
	DeleteFilter()
}

func (d *DB) OpenDB() error {
	db, err := sql.Open(d.DBType, d.DBUser+":"+d.DBPassword+"@/"+d.DBName)
	d.DB = db
	return err
}
