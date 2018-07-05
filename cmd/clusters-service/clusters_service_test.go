package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestList(t *testing.T) {
	t.Fatal("This test fails")
}
