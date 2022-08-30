package domain

import (
	"fmt"
	"github.com/alloba/TheLibrarian/database"
	"gorm.io/gorm"
	"runtime"
	"strings"
)

type RecordRepo struct {
	db *gorm.DB
}

func NewRecordRepo(db *gorm.DB) *RecordRepo {
	return &RecordRepo{db: db}
}

type RecordService struct {
	repo *RecordRepo
}

func NewRecordService(repo *RecordRepo) *RecordService {
	return &RecordService{repo: repo}
}

func (repo RecordRepo) Exists(recordId string) (bool, error) {
	var exists bool
	err := repo.db.Model(&database.Record{}).
		Select("count(*) > 0").
		Where("id = ?", recordId).
		Find(&exists).Error

	if err != nil {
		return false, logTrace(err)
	}
	return exists, nil
}

func (repo RecordRepo) SaveOne(record *database.Record) error {
	if err := repo.db.Create(record).Error; err != nil {
		return logTrace(err)
	}
	return nil
}

func (repo RecordRepo) DeleteOne(recordId string) error {
	exist, err := repo.Exists(recordId)
	if err != nil {
		return logTrace(err)
	}

	if !exist {
		return nil
	}

	err = repo.db.Where("id = ?", recordId).Delete(&database.Record{}).Error
	if err != nil {
		return logTrace(err)
	}
	return nil
}

func logTrace(err error) error {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return fmt.Errorf("ERR_CANNOT_TRACE_CALLER: %v", err.Error())
	}

	fullName := fmt.Sprintf("%s", runtime.FuncForPC(pc).Name())
	methodName := strings.Split(fullName, "/")[len(strings.Split(fullName, "/"))-1]

	return fmt.Errorf("%v: %v", methodName, err.Error())
}
