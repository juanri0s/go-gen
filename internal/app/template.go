package app

import (
	"html/template"
	"io/ioutil"
	"os"
)

// DefaultImports TODO
const DefaultImports = `"fmt", "os"`

func setupService(m Metadata) (string, error) {
	err := initMod(m.ProjectPath)
	if err != nil {
		return "", err
	}

	if m.HasGitIgnore {
		err = initGitIgnore(m)
		if err != nil {
			return "", err
		}
	}

	if m.HasLicense {
		err = initLicense(m)
		if err != nil {
			return "", err
		}
	}

	entryPath, err := makeDir(m)
	if err != nil {
		return "", err
	}

	return entryPath, nil
	// TODO: create /cmd/entry, create /internal/app, run go mod init, get .gitignore, get license, get dockerfile
}

func makeDir(m Metadata) (string, error) {
	entryPath := m.ProjectPath + "/cmd" + m.Entrypoint
	_, err := os.Stat(entryPath)
	if os.IsNotExist(err) {
		err := os.MkdirAll(entryPath, 0755)
		if err != nil {
			return "", err
		}
	}
	return entryPath, nil
}

func initGitIgnore(m Metadata) error {
	// relative path of template
	const gitignoreT = "internal/app/templates/gitignore.tmpl"
	templateF, err := ioutil.ReadFile(gitignoreT)
	if err != nil {
		return err
	}
	// Create a template, add the function map, and parse the text.
	tmpl, err := template.New("").Parse(string(templateF))
	if err != nil {
		return err
	}

	f, err := os.Create(m.ProjectPath + "/.gitignore")
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

func initLicense(m Metadata) error {
	// relative path of template
	const licenseT = "internal/app/templates/license.tmpl"
	templateF, err := ioutil.ReadFile(licenseT)
	if err != nil {
		return err
	}
	// Create a template, add the function map, and parse the text.
	tmpl, err := template.New("").Parse(string(templateF))
	if err != nil {
		return err
	}

	f, err := os.Create(m.ProjectPath + "/LICENSE")
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

func templateRepo(entryPath string, m Metadata) error {
	// relative path of template
	const mainT = "internal/app/templates/simple-main.gotmpl"
	templateF, err := ioutil.ReadFile(mainT)
	if err != nil {
		return err
	}
	// Create a template, add the function map, and parse the text.
	tmpl, err := template.New("").Parse(string(templateF))
	if err != nil {
		return err
	}

	f, err := os.Create(entryPath + "/main.go")
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
