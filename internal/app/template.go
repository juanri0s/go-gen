package app

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// DefaultImports represents the default imports that we want the initial service to have.
var DefaultImports = []string{"fmt"}

// setupService sets up a service based on the configuration metadata.
func setupService(m Metadata) error {
	err := initMod(m.ProjectPath)
	if err != nil {
		return err
	}

	if m.HasGitIgnore {
		err := addFileFromTemplate("gitignore", m)
		if err != nil {
			return err
		}
	}

	if m.HasLicense {
		err := addFileFromTemplate("license", m)
		if err != nil {
			return err
		}
	}

	err = makeDirForEntrypoint(m)
	if err != nil {
		return err
	}

	err = addFileFromTemplate("main", m)
	if err != nil {
		return err
	}

	err = addFileFromTemplate("docker", m)
	if err != nil {
		return err
	}

	return nil
}

// makeDirForEntrypoint sets up the main logic directory with entrypoint.
func makeDirForEntrypoint(m Metadata) error {
	if m.ProjectPath == "" {
		return fmt.Errorf("Invalid project path from metadata")
	}
	if m.Entrypoint == "" {
		return fmt.Errorf("Invalid entrypoint from metadata")
	}

	entryPath := m.ProjectPath + "/cmd/" + m.Entrypoint

	log.Printf("creating project entry path for %s", entryPath)

	_, err := os.Stat(entryPath)
	if os.IsNotExist(err) {
		err := os.MkdirAll(entryPath, 0755)
		if err != nil {
			return fmt.Errorf("unable to make directory for %s - %w", entryPath, err)
		}
	}
	return nil
}

// addFileFromTemplate templates a file based on the fileType and matching template.
func addFileFromTemplate(fType string, m Metadata) error {
	var tmplPath string
	var path string
	var file string
	switch strings.ToLower(fType) {
	case "main":
		tmplPath = "internal/app/templates/simple-main.gotmpl"
		path = m.ProjectPath + "/cmd/" + m.Entrypoint + "/"
		file = "main.go"
		break
	case "docker":
		tmplPath = "internal/app/templates/dockerfile.tmpl"
		path = m.ProjectPath + "/"
		file = "Dockerfile"
		break
	case "license":
		tmplPath = "internal/app/templates/license.tmpl"
		path = m.ProjectPath + "/"
		file = "LICENSE"
		break
	case "gitignore":
		tmplPath = "internal/app/templates/gitignore.tmpl"
		path = m.ProjectPath + "/"
		file = ".gitignore"
		break
	default:
		return fmt.Errorf("Invalid or unsupported file type")
	}

	templateF, err := ioutil.ReadFile(tmplPath)
	if err != nil {
		return fmt.Errorf("unable to read template file from %s  - %w", tmplPath, err)
	}
	// Create a template and parse the text.
	tmpl, err := template.New("Main").Parse(string(templateF))
	if err != nil {
		return fmt.Errorf("unable to parse template file - %w", err)
	}

	f, err := os.Create(path + file)
	if err != nil {
		return fmt.Errorf("unable to create file from path %s - %w", path+file, err)
	}

	// Run the template to verify the output.
	err = tmpl.Execute(f, m)
	if err != nil {
		return fmt.Errorf("unable to execute templating - %w", err)
	}

	f.Close()
	return nil
}
