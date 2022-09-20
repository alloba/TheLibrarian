package service

import (
	"fmt"
	"github.com/alloba/TheLibrarian/database"
	"github.com/alloba/TheLibrarian/database/repository"
	"gorm.io/gorm"
	"testing"
)

const integrationTestDbPath = "../../out/library_integration_test.db"

var db = database.Connect(integrationTestDbPath)
var fileService = NewFileService("../../out/filebin")

func TestRecordService_CreateRecord(t *testing.T) {
	type fields struct {
		RecordRepo *repository.RecordRepo
	}
	type args struct {
		fileContainer *FileContainer
	}

	recordRepo := repository.NewRecordRepo(db)
	fileContainer, err := fileService.CreateFileContainer("RecordService_test.go") //this file.
	if err != nil {
		panic(fmt.Errorf("failed during test init -- cannot create fileContainer - %v", err))
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "successCreate", fields: fields{RecordRepo: recordRepo}, args: args{fileContainer: fileContainer}, wantErr: false},
		{name: "failDupe", fields: fields{RecordRepo: recordRepo}, args: args{fileContainer: fileContainer}, wantErr: true},
	}

	defer deleteTestRecords(db)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := RecordService{
				RecordRepo: tt.fields.RecordRepo,
			}
			if err := service.CreateRecord(tt.args.fileContainer); (err != nil) != tt.wantErr {
				t.Errorf("CreateRecord() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRecordService_DeleteRecord(t *testing.T) {
	type fields struct {
		RecordRepo *repository.RecordRepo
	}
	type args struct {
		fileHash string
	}

	recordRepo := repository.NewRecordRepo(db)
	fileContainer, err := fileService.CreateFileContainer("RecordService_test.go") //this file.
	if err != nil {
		panic(fmt.Errorf("failed during test init -- cannot create fileContainer - %v", err))
	}
	initService := RecordService{RecordRepo: recordRepo}
	err = initService.CreateRecord(fileContainer)
	if err != nil {
		panic(fmt.Errorf("couln not init test in deleteRecord - %v", err))
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "successDelete", fields: fields{RecordRepo: recordRepo}, args: args{fileHash: fileContainer.Hash}, wantErr: false},
		{name: "successDeleteNotExist", fields: fields{RecordRepo: recordRepo}, args: args{fileHash: "blah blah"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := RecordService{
				RecordRepo: tt.fields.RecordRepo,
			}
			if err := service.DeleteRecord(tt.args.fileHash); (err != nil) != tt.wantErr {
				t.Errorf("DeleteRecord() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRecordService_GetRecord(t *testing.T) {
	type fields struct {
		RecordRepo *repository.RecordRepo
	}
	type args struct {
		fileHash string
	}

	recordRepo := repository.NewRecordRepo(db)
	fileContainer, err := fileService.CreateFileContainer("RecordService_test.go") //this file.
	if err != nil {
		panic(fmt.Errorf("failed during test init -- cannot create fileContainer - %v", err))
	}
	initService := RecordService{RecordRepo: recordRepo}
	err = initService.CreateRecord(fileContainer)
	if err != nil {
		panic(fmt.Errorf("couln not init test in deleteRecord - %v", err))
	}

	defer deleteTestRecords(db)
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *database.Record
		wantErr bool
	}{
		{name: "successGet", fields: fields{RecordRepo: recordRepo}, args: args{fileHash: fileContainer.Hash}, want: nil, wantErr: false},
		{name: "failNotFound", fields: fields{RecordRepo: recordRepo}, args: args{fileHash: "blahblah"}, want: nil, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := RecordService{
				RecordRepo: tt.fields.RecordRepo,
			}
			got, err := service.GetRecord(tt.args.fileHash)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil && err != nil && !tt.wantErr {
				t.Errorf("GetRecord() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func deleteTestRecords(db *gorm.DB) {
	db.Where("1 <> 2").Delete(&database.Record{})
}
