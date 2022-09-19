package repository

import (
	"github.com/alloba/TheLibrarian/database"
	"gorm.io/gorm"
	"testing"
	"time"
)

const integrationTestDbPath = "../../../out/library_integration_test.db"

var testBook1 = &database.Book{
	ID:           "testBook1Id",
	Name:         "testBook1",
	DateCreated:  time.Now(),
	DateModified: time.Now(),
}

func deleteTestBooks(db *gorm.DB) {
	db.Where("id like ?", "test%").Delete(&database.Book{})
}

func TestBookRepo_SaveOne(t *testing.T) {
	db := database.Connect(integrationTestDbPath)
	defer deleteTestBooks(db)

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		book *database.Book
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"expectSuccess", fields{db}, args{testBook1}, false},
		{"expectFailureOnDupe", fields{db}, args{testBook1}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := BookRepo{
				db: tt.fields.db,
			}
			if err := repo.CreateOne(tt.args.book); (err != nil) != tt.wantErr {
				t.Errorf("CreateOne() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBookRepo_FindOneByName(t *testing.T) {
	db := database.Connect(integrationTestDbPath)
	defer deleteTestBooks(db)

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		bookname string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"expectSuccess", fields{db}, args{"testBook1"}, false},
		{"expectFailureOnNotFound", fields{db}, args{"namedoesntexist"}, true},
	}

	repo := BookRepo{db: db}
	_ = repo.CreateOne(testBook1)
	defer deleteTestBooks(db)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := repo.FindOneByName(tt.args.bookname); (err != nil) != tt.wantErr {
				t.Errorf("FindOneByName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

}
