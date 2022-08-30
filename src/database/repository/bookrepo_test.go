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
			if err := repo.CreateOne(tt.args.book); (err != nil) != tt.wantErr {
				t.Errorf("CreateOne() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBookRepo_FindOneByName(t *testing.T) {
	db := database.Connect(integrationTestDbPath)
	repo := BookRepo{db: db}
	testBook := getTestBook("testBookAssociatedWithEdition")
	defer wipeTestDatabase(db)

	err := repo.CreateOne(testBook)
	if err != nil {
		t.Fatalf("couldnt save book into db")
	}

	res, err := repo.FindOneByName(testBook.Name)
	if err != nil {
		t.Fatalf("couldnt find book - %v", err.Error())
	}
	t.Logf("%#v", res)

}
