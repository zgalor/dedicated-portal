package sql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate"
	// Register the migrate postgresl driver
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
)

// EnsureSchema makes sure our DB's schema matches that defined by schemaPath.
// It migrates it if needed and returns error if it can't do so.
func EnsureSchema(schemaPath, connectionUrl string) error {
	waitForDatabase(connectionUrl)
	m, err := migrate.New(
		fmt.Sprintf("file:///%s", schemaPath),
		connectionUrl,
	)
	if err != nil {
		return fmt.Errorf("migrate.New: %v", err)
	}
	defer m.Close()

	err = outputVersion(m)
	if err != nil {
		return fmt.Errorf("outputVersion: %v", err)
	}
	runMigration(m)
	if err != nil {
		return fmt.Errorf("migrate: %v", err)
	}
	outputVersion(m)
	if err != nil {
		return fmt.Errorf("outputVersion: %v", err)
	}
	return nil
}

func waitForDatabase(connectionUrl string) {
	for {
		err := tryConnect(connectionUrl)
		if err != nil {
			fmt.Println("Could not connect: ", err)
			time.Sleep(500 * time.Millisecond)
		} else {
			return
		}
	}
}

func tryConnect(connectionURL string) error {
	db, err := sql.Open("postgres", connectionURL)
	defer db.Close()
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("sql.Open: %v...", err))
	}
	err = db.Ping()
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("db.Ping: %v...", err))
	}
	return nil
}

func outputVersion(m *migrate.Migrate) error {
	ver, dirty, err := m.Version()
	if err == migrate.ErrNilVersion {
		fmt.Print("Current version [0], dirty [false]\n")
	} else if err != nil {
		return fmt.Errorf("m.Version: %v", err)
	} else {
		fmt.Printf("Current version [%v], dirty [%v]\n", ver, dirty)
	}
	return nil
}

func runMigration(m *migrate.Migrate) error {
	err := m.Up()
	if err == migrate.ErrNoChange {
		fmt.Printf("Already fully migrated!\n")
	} else if err != nil {
		return fmt.Errorf("m.Up: %v", err)
	} else {
		fmt.Print("Migrated Successfully!\n")
	}
	return nil
}
