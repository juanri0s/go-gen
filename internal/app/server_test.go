package app

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `HEALTHY 1.0.0`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func Test_initGit(t *testing.T) {
	type args struct {
		p string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "empty path error",
			args:    args{p: ""},
			wantErr: true,
		},
		{
			name:    "invalid path error",
			args:    args{p: "abc"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := initGit(tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("initGit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_initMod(t *testing.T) {
	type args struct {
		p string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "empty path error",
			args:    args{p: ""},
			wantErr: true,
		},
		{
			name:    "invalid path error",
			args:    args{p: "abc"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := initMod(tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("initMod() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_setRepoURL(t *testing.T) {
	type args struct {
		p   string
		url string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "empty path error",
			args:    args{p: "", url: "abc"},
			wantErr: true,
		},
		{
			name:    "invalid path error",
			args:    args{p: "abc", url: "abc"},
			wantErr: true,
		},
		{
			name:    "empty url error",
			args:    args{p: "/", url: ""},
			wantErr: true,
		},
		{
			name:    "invalid url error",
			args:    args{p: "abc", url: "abc"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := setRepoURL(tt.args.p, tt.args.url); (err != nil) != tt.wantErr {
				t.Errorf("setRepoURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
