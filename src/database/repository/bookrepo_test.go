package repository

import (
	"github.com/alloba/TheLibrarian/database"
	"gorm.io/gorm"
	"testing"
)

func TestBookRepo_SaveOne(t *testing.T) {
	db := database.Connect(integrationTestDbPath)
	testBook := getTestBook("testBookAssociatedWithEdition")
	defer wipeTestDatabase(db)

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
		{"expectSuccess", fields{db}, args{testBook}, false},
		{"expectFailureOnDupe", fields{db}, args{testBook}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := BookRepo{
				db: tt.fields.db,
			}
			if err := repo.SaveOne(tt.args.book); (err != nil) != tt.wantErr {
				t.Errorf("SaveOne() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
