package app

import (
	"html/template"
	"io/ioutil"
	"os"
	"strings"
)

// DefaultImports TODO
const DefaultImports = `"fmt", "os"`

func setupService(m Metadata) error {
	err := initMod(m.ProjectPath)
	if err != nil {
		return err
	}

	if m.HasGitIgnore {
		err = templateService("gitignore", m)
		if err != nil {
			return err
		}
	}

	if m.HasLicense {
		err = templateService("license", m)
		if err != nil {
			return err
		}
	}

	err = makeDir(m)
	if err != nil {
		return err
	}

	err = templateService("main", m)
	if err != nil {
		return err
	}

	err = templateService("docker", m)
	if err != nil {
		return err
	}

	return nil
	// TODO: create /cmd/entry, create /internal/app, run go mod init, get .gitignore, get license, get dockerfile
}

func makeDir(m Metadata) error {
	entryPath := m.ProjectPath + "/cmd" + m.Entrypoint
	_, err := os.Stat(entryPath)
	if os.IsNotExist(err) {
		err := os.MkdirAll(entryPath, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func templateService(fType string, m Metadata) error {
	var tmplPath string
	var path string
	var file string
	switch strings.ToLower(fType) {
	case "main":
		tmplPath = "internal/app/templates/simple-main.gotmpl"
		path = m.ProjectPath + "/cmd" + m.Entrypoint + "/"
		file = "main.go"
		break
	case "docker":
		tmplPath = "internal/app/templates/docker.tmpl"
		path = m.ProjectPath
		file = "TEST.TXT"
		break
	case "license":
		tmplPath = "internal/app/templates/license.tmpl"
		path = m.ProjectPath
		file = "LICENSE"
		break
	case "gitignore":
		tmplPath = "internal/app/templates/gitignore.tmpl"
		path = m.ProjectPath
		file = ".gitignore"
		break
	default:
		tmplPath = ""
		path = ""
		file = ""
		break
	}

	templateF, err := ioutil.ReadFile(tmplPath)
	if err != nil {
		return err
	}
	// Create a template, add the function map, and parse the text.
	tmpl, err := template.New("").Parse(string(templateF))
	if err != nil {
		return err
	}

	f, err := os.Create(path + file)
	if err != nil {
		return err
	}

	// Run the template to verify the output.
	err = tmpl.Execute(f, m)
	if err != nil {
		return err
	}

	f.Close()
	return nil
}
