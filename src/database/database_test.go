package database

import (
	"testing"
	"github.com/alloba/TheLibrarian/configuration"
)

var dbPath = "../../out/library_integration_test.db"
var config = configuration.NewLibraryConfig(dbPath, "")

func TestConnect(t *testing.T) {
	t.Run("successfulConnect", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("fail Connect: %v", r)
			}
		}()
		Connect(dbPath)
	})
	t.Run("failedConnectBadPath", func(t *testing.T) {
		defer func() { _ = recover() }()
		Connect("lkjajkla;sdf./asdfasdf/.a/dsf/a.sd.f/a./df././/")
		t.Fatalf("unreachable")
	})

}
