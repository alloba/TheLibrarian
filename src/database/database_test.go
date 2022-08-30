package database

import (
	"testing"
)

var dbPath = "../../out/library_integration_test.db"

func TestConnect(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("fail Connect: %v", r)
		}
	}()
	Connect(dbPath)
}
