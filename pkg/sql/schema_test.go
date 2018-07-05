package sql

import (
	"testing"
)

func TestConnectionString(t *testing.T) {
	str := ConnectionURL()
	if str != "postgres://:@localhost:5432/?sslmode=disable" {
		t.Fatal("Unexpected connection string", str)
	}
}
