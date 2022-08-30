package domain

import (
	"github.com/alloba/TheLibrarian/database"
	"gorm.io/gorm"
	"testing"
	"time"
)

var dbPath = "../../out/library_integration_test.db"

func TestRecordRepo_SaveOne(t *testing.T) {
	db := database.Connect(dbPath)
	cleanUp(db)
	defer cleanUp(db)

	testRecord := &database.Record{
		ID:               "testHash",
		FilePointer:      "filepointer",
		Name:             "filename",
		Extension:        ".fileExt",
		DateFileModified: time.Now(),
		DateCreated:      time.Now(),
		DateModified:     time.Now(),
	}

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		record *database.Record
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "newRecordSaved", fields: fields{db: db}, args: args{testRecord}, wantErr: false},
		{name: "existingRecordFailedSave", fields: fields{db: db}, args: args{testRecord}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := RecordRepo{
				db: tt.fields.db,
			}
			if err := repo.SaveOne(tt.args.record); (err != nil) != tt.wantErr {
				t.Errorf("SaveOne() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func cleanUp(db *gorm.DB) {
	err := db.Where("id like ?", "test%").Delete(&database.Record{}).Error
	if err != nil {
		panic(logTrace(err))
	}
}
