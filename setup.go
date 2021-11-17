package migrationmanager

import (
	"context"

	"github.com/uptrace/bun"
)

//
// SetupTable creates a table if not existing.
//
func SetupTable(connect func() *bun.DB) error {
	conn := connect()
	defer conn.Close()
	m := Migration{}
	if _, err := conn.NewCreateTable().Model(&m).Exec(context.Background()); err != nil {
		return err
	}
	return nil
}
