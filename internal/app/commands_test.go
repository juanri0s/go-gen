package app

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func TestMetadata_new(t *testing.T) {
	var m *Metadata
	def := Metadata{
		Description:  "A default service for auth0",
		Entrypoint:   "default-service",
		HasCopyright: true,
		HasGitIgnore: true,
		HasLicense:   true,
		Imports:      DefaultImports,
		IsPrivate:    true,
		MainBranch:   "main",
		Name:         "default-repo",
		Owner:        "default-owner",
		ProjectPath:  "",
		Version:      "1.0.0",
	}
	tests := []struct {
		name string
		m    *Metadata
		want Metadata
	}{
		{
			name: "default metadata is accuraterlyreturned",
			m:    m,
			want: def,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.new(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Metadata.new() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generateServiceFromDefault(t *testing.T) {
	type args struct {
		token string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "Error on empty token",
			args:    args{""},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Error on invalid or expired token",
			args:    args{"abc"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateServiceFromDefault(tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateServiceFromDefault() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("generateServiceFromDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generateServiceFromFile(t *testing.T) {
	type args struct {
		f     string
		token string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "Error on empty file",
			args:    args{"", "abc"},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Error on invalid or empty token",
			args:    args{"abc", ""},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Error on invalid file",
			args:    args{"abc", "abc"},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Error on invalid or unsupported file type",
			args:    args{"./test/assets/test.txt", "abc"},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Error on invalid json marshal",
			args:    args{"./test/assets/test.json", "abc"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateServiceFromFile(tt.args.f, tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateServiceFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("generateServiceFromFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerator_generate(t *testing.T) {
	tests := []struct {
		name    string
		g       *Generator
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.g.generate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Generator.generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Generator.generate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCli_Run(t *testing.T) {
	s := ""
	app := &cli.App{
		Action: func(c *cli.Context) error {
			s = s + c.Args().First()
			return nil
		},
	}

	err := app.Run([]string{"command", "foo"})
	assert.Equal(t, err, nil)
	err = app.Run([]string{"command", "bar"})
	assert.Equal(t, err, nil)
	assert.Equal(t, s, "foobar")
}
