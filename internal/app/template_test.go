package app

import (
	"testing"
)

func Test_makeDirForEntry(t *testing.T) {
	type args struct {
		m Metadata
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "error on empty metadata",
			args:    args{Metadata{}},
			wantErr: true,
		},
		{
			name:    "error on invalid project path",
			args:    args{Metadata{ProjectPath: "", Entrypoint: "abc"}},
			wantErr: true,
		},
		{
			name:    "error on invalid entrypoint",
			args:    args{Metadata{ProjectPath: "abc", Entrypoint: ""}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := makeDirForEntrypoint(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("makeDirForEntry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_addFileFromTemplate(t *testing.T) {
	type args struct {
		fType string
		m     Metadata
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "error on empty file type",
			args:    args{"", Metadata{}},
			wantErr: true,
		},
		{
			name:    "error on invalid or unsupported file type",
			args:    args{"unsupported", Metadata{}},
			wantErr: true,
		},
		{
			name:    "error on empty project path",
			args:    args{"unsupported", Metadata{ProjectPath: ""}},
			wantErr: true,
		},
		{
			name:    "error on invalid project path",
			args:    args{"unsupported", Metadata{ProjectPath: "abc"}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := addFileFromTemplate(tt.args.fType, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("addFileFromTemplate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
