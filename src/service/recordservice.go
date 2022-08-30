package service

import (
	"fmt"
	"github.com/alloba/TheLibrarian/database/repository"
	"runtime"
	"strings"
)

type RecordService struct {
	repo *repository.RecordRepo
}

func NewRecordService(repo *repository.RecordRepo) *RecordService {
	return &RecordService{repo: repo}
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
