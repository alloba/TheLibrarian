package repository

import (
	"github.com/alloba/TheLibrarian/database"
	"gorm.io/gorm"
	"testing"
	"time"
)

func TestBookRepo_Exists(t *testing.T) {
	db := database.Connect(integrationTestDbPath)
	testBook := getTestBook("testBook1")

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		BookId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{"exists", fields{db}, args{BookId: testBook.ID}, true, false},
		{"notExists", fields{db}, args{BookId: "nonexistBookid"}, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanUpTestBooks(db)
			defer cleanUpTestBooks(db)

			repo := BookRepo{
				db: tt.fields.db,
			}

			if err := repo.SaveOne(testBook); err != nil {
				t.Fatalf("failed setup")
			}

			got, err := repo.Exists(tt.args.BookId)
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

func TestBookRepo_SaveOne(t *testing.T) {
	db := database.Connect(integrationTestDbPath)
	testBook := getTestBook("testId")

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		Book *database.Book
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "newBookSaved", fields: fields{db: db}, args: args{testBook}, wantErr: false},
		{name: "existingBookFailedSave", fields: fields{db: db}, args: args{testBook}, wantErr: true},
	}
	cleanUpTestBooks(db)
	defer cleanUpTestBooks(db)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := BookRepo{
				db: tt.fields.db,
			}
			if err := repo.SaveOne(tt.args.Book); (err != nil) != tt.wantErr {
				t.Errorf("SaveOne() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBookRepo_DeleteOne(t *testing.T) {
	db := database.Connect(integrationTestDbPath)
	testBook := getTestBook("testBook1")

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		BookId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"successDelete", fields{db}, args{testBook.ID}, false},
		{"successNoActionTaken", fields{db}, args{"fakeidnothere"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanUpTestBooks(db)
			defer cleanUpTestBooks(db)
			repo := BookRepo{
				db: tt.fields.db,
			}
			if err := repo.SaveOne(testBook); err != nil {
				t.Error(logTrace(err))
			}

			if err := repo.DeleteOne(tt.args.BookId); (err != nil) != tt.wantErr {
				t.Errorf("DeleteOne() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func getTestBook(id string) *database.Book {
	return &database.Book{
		ID:           id,
		Name:         "testname",
		DateCreated:  time.Now(),
		DateModified: time.Now(),
	}
}
