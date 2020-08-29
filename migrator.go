package migrationmanager

import (
	"fmt"
	"strings"

	"github.com/go-pg/pg/v10"
)

//
// MigrateUp runs the up function for migrations
//
func MigrateUp(migrations []Migration, connect func() *pg.DB) (bool, error) {
	previousMigrations, err := getMigrationsThatHaveBeenRan(connect)
	if err != nil {
		fmt.Printf("Unable to retrived previous migrations: %s\n", err)
		return false, err
	}
	count := 0
	for _, migration := range migrations {
		if existingMigration := getExistingMigration(migration.Name, previousMigrations); existingMigration != nil {
			continue
		}
		count++
		fmt.Printf("Running Migration: %s\n", migration.Name)
		if err := migration.Up(); err != nil {
			return false, fmt.Errorf("migration %s has failed with error: %s", migration.Name, err)
		}
		if err := addMigration(migration, connect); err != nil {
			return false, fmt.Errorf("unable to insert migration: %s", err)
		}
		fmt.Printf("Finished Migration: %s\n", migration.Name)
	}
	fmt.Printf("Successfully ran %v new migrations\n", count)
	return true, nil
}

//
// MigrateDown runs the down scripts for migrations
//
func MigrateDown(migrations []Migration, connect func() *pg.DB) (bool, error) {
	previousMigrations, err := getMigrationsThatHaveBeenRan(connect)
	if err != nil {
		fmt.Printf("Unable to retrived previous migrations: %s\n", err)
		return false, err
	}

	count := 0
	for _, migration := range migrations {
		existingMigration := getExistingMigration(migration.Name, previousMigrations)
		if existingMigration == nil {
			continue
		}
		count++
		fmt.Printf("Reverting Migration: %s\n", migration.Name)
		if err := migration.Down(); err != nil {
			return false, fmt.Errorf("migration %s has failed with error: %s", migration.Name, err)
		}
		if err := removeMigration(existingMigration, connect); err != nil {
			return false, fmt.Errorf("unable to insert migration: %s", err)
		}
		fmt.Printf("Finished Migration: %s\n", migration.Name)
	}
	fmt.Printf("Successfully reverted %v migrations\n", count)
	return true, nil
}

func getExistingMigration(name string, previousMigrations []Migration) *Migration {
	for _, previousMigration := range previousMigrations {
		if name == previousMigration.Name {
			return &previousMigration
		}
	}
	return nil
}

func getMigrationsThatHaveBeenRan(connect func() *pg.DB) ([]Migration, error) {
	conn := connect()
	defer conn.Close()
	migrations := []Migration{}
	if err := conn.Model(&migrations).Order("created_at DESC").Select(); err != nil {
		if strings.Contains(err.Error(), "ERROR #42P01") {
			if err := SetupTable(connect); err != nil {
				return nil, err
			}
			return getMigrationsThatHaveBeenRan(connect)
		}
		return nil, err
	}
	return migrations, nil
}

func addMigration(migration Migration, connect func() *pg.DB) error {
	conn := connect()
	defer conn.Close()
	if _, err := conn.Model(&migration).Insert(&migration); err != nil {
		return err
	}
	return nil
}

func removeMigration(migration *Migration, connect func() *pg.DB) error {
	conn := connect()
	defer conn.Close()
	if _, err := conn.Model(migration).WherePK().Delete(); err != nil {
		return err
	}
	return nil
}
