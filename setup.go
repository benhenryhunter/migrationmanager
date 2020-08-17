package migrationmanager

import (
	"github.com/dickmanben/migrationmanager/db"
	"github.com/go-pg/pg/v10/orm"
)

//
// SetupTable creates a table if not existing.
//
func SetupTable() error {
	conn := db.Connect()
	defer conn.Close()
	m := Migration{}
	if err := conn.Model(&m).CreateTable(&orm.CreateTableOptions{IfNotExists: true}); err != nil {
		return err
	}
	return nil
}
