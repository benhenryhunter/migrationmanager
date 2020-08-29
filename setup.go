package migrationmanager

import (
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

//
// SetupTable creates a table if not existing.
//
func SetupTable(connect func() *pg.DB) error {
	conn := connect()
	defer conn.Close()
	m := Migration{}
	if err := conn.Model(&m).CreateTable(&orm.CreateTableOptions{IfNotExists: true}); err != nil {
		return err
	}
	return nil
}
