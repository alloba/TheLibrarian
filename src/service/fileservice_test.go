package service

import (
	"testing"
)

func TestFileService_createFileContainer(t *testing.T) {
	type fields struct {
		archiveBasePath string
	}
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    func(container *FileContainer) bool
		wantErr bool
	}{
		{name: "fileContainer",
			fields: fields{archiveBasePath: "../../out/filebin"},
			args:   args{path: "./util.go"},
			want: func(container *FileContainer) bool {
				if container.OriginName != "util.go" || container.IsDir || container.Hash == "" {
					return false
				}
				return true
			},
			wantErr: false,
		},
		{name: "dirContainer",
			fields: fields{archiveBasePath: "../../out/filebin"},
			args:   args{path: "./"},
			want: func(container *FileContainer) bool {
				if container.OriginName != "service" || !container.IsDir || container.Hash != "" {
					return false
				}
				return true
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := FileService{
				archiveBasePath: tt.fields.archiveBasePath,
			}
			got, err := service.createFileContainer(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("createFileContainer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.want(got) {
				t.Errorf("want conditions not satisfied")
			}
		})
	}
}

func TestFileService_WriteContainerToArchive_MANUALVERIFY(t *testing.T) {
	service := NewFileService("../../out/filebin")
	container, err := service.createFileContainer("./fileservice_test.go")
	if err != nil {
		t.Fatalf("failed init process - container creation - %v", err.Error())
	}

	err = service.WriteContainerToArchive(container)
	if err != nil {
		t.Fatalf("failed writing process - %v", err.Error())
	}
}
