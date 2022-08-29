package fileutil

import (
	"testing"
)

func TestGetAllNestedFilePaths(t *testing.T) {
	type args struct {
		dirPath string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "recursiveWalk",
			args: args{"../"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAllNestedFilePaths(tt.args.dirPath)
			if err != nil {
				t.Errorf("GetAllNestedFilePaths() error = %v", err.Error())
				return
			}
			t.Logf("%#v", got)
		})
	}
}
