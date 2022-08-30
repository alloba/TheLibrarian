package repository

import (
	"github.com/alloba/TheLibrarian/database"
	"gorm.io/gorm"
	"testing"
	"time"
)

func TestEditionRepo_Exists(t *testing.T) {
	db := database.Connect(integrationTestDbPath)
	repoManager := NewRepoManager(db)
	testBook := getTestBook("testBookAssociatedWithEdition")
	testEdition := getTestEdition("testEdition1", testBook.ID)
	defer wipeTestDatabase(db)

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		EditionId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{"exists", fields{db}, args{EditionId: testEdition.ID}, true, false},
		{"notExists", fields{db}, args{EditionId: "nonexistEditionid"}, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := EditionRepo{
				db: tt.fields.db,
			}
			defer wipeTestDatabase(db)

			if err := repoManager.Book.SaveOne(testBook); err != nil {
				t.Fatalf("failed setup")
			}
			if err := repo.SaveOne(testEdition); err != nil {
				t.Fatalf("failed setup")
			}

			got, err := repo.Exists(tt.args.EditionId)
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

func TestEditionRepo_SaveOne(t *testing.T) {
	db := database.Connect(integrationTestDbPath)
	repoManager := NewRepoManager(db)
	testBook := getTestBook("testBookAssociatedWithEdition")
	testEdition := getTestEdition("testEdition1", testBook.ID)

	repoManager.Book.SaveOne(testBook)
	defer wipeTestDatabase(db)

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		Edition *database.Edition
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "newEditionSaved", fields: fields{db: db}, args: args{testEdition}, wantErr: false},
		{name: "existingEditionFailedSave", fields: fields{db: db}, args: args{testEdition}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			repo := EditionRepo{
				db: tt.fields.db,
			}
			if err := repo.SaveOne(tt.args.Edition); (err != nil) != tt.wantErr {
				t.Errorf("SaveOne() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEditionRepo_DeleteOne(t *testing.T) {
	db := database.Connect(integrationTestDbPath)
	repoManager := NewRepoManager(db)
	testBook := getTestBook("testBookAssociatedWithEdition")
	testEdition := getTestEdition("testEdition1", testBook.ID)
	defer wipeTestDatabase(db)

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		EditionId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"successDelete", fields{db}, args{testEdition.ID}, false},
		{"successNoActionTaken", fields{db}, args{"fakeidnothere"}, false},
	}
	repoManager.Book.SaveOne(testBook)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer repoManager.Edition.DeleteOne(tt.args.EditionId)
			repo := EditionRepo{
				db: tt.fields.db,
			}
			if err := repo.SaveOne(testEdition); err != nil {
				t.Error(logTrace(err))
			}

			if err := repo.DeleteOne(tt.args.EditionId); (err != nil) != tt.wantErr {
				t.Errorf("DeleteOne() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func getTestEdition(id string, bookId string) *database.Edition {

	return &database.Edition{
		ID:            id,
		EditionNumber: 0,
		BookID:        bookId,
		DateCreated:   time.Now(),
		DateModified:  time.Now(),
	}
}
