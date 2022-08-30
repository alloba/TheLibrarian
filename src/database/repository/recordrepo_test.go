package repository

import (
	"github.com/alloba/TheLibrarian/database"
	"gorm.io/gorm"
	"testing"
	"time"
)

func TestRecordRepo_Exists(t *testing.T) {
	db := database.Connect(integrationTestDbPath)
	testRecord := getTestRecord("testRecord1")

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		recordId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{"exists", fields{db}, args{recordId: testRecord.ID}, true, false},
		{"notExists", fields{db}, args{recordId: "nonexistrecordid"}, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanUpTestRecords(db)
			defer cleanUpTestRecords(db)

			repo := RecordRepo{
				db: tt.fields.db,
			}

			if err := repo.SaveOne(testRecord); err != nil {
				t.Fatalf("failed setup")
			}

			got, err := repo.Exists(tt.args.recordId)
			if (err != nil) != tt.wantErr {
				t.Errorf("Exists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Exists() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecordRepo_SaveOne(t *testing.T) {
	db := database.Connect(integrationTestDbPath)
	testRecord := getTestRecord("testhash")

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
	cleanUpTestRecords(db)
	defer cleanUpTestRecords(db)
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

func TestRecordRepo_DeleteOne(t *testing.T) {
	db := database.Connect(integrationTestDbPath)
	testRecord := getTestRecord("testRecord1")

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		recordId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"successDelete", fields{db}, args{testRecord.ID}, false},
		{"successNoActionTaken", fields{db}, args{"fakeidnothere"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanUpTestRecords(db)
			defer cleanUpTestRecords(db)
			repo := RecordRepo{
				db: tt.fields.db,
			}
			if err := repo.SaveOne(testRecord); err != nil {
				t.Error(logTrace(err))
			}

			if err := repo.DeleteOne(tt.args.recordId); (err != nil) != tt.wantErr {
				t.Errorf("DeleteOne() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func getTestRecord(id string) *database.Record {
	return &database.Record{
		ID:               id,
		FilePointer:      "testpointer",
		Name:             "testname",
		Extension:        "test",
		DateFileModified: time.Now(),
		DateCreated:      time.Now(),
		DateModified:     time.Now(),
	}
}
