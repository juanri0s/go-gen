package app

import (
	"html/template"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
)

// Example TODO:
type Example struct {
	ServicePackage string
	Copyright      bool
	DefaultImports string
}

func testTemplate() {
	var one = Example{
		"main", true, `"fmt", "os"`,
	}

	templateText, err := ioutil.ReadFile("./generator/templates/main.gotmpl")
	if err != nil {
		log.Error("parsing: %s", err)
	}
	// Create a template, add the function map, and parse the text.
	tmpl, err := template.New("Main").Parse(string(templateText))
	if err != nil {
		log.Error("parsing: %s", err)
	}

	f, err := os.Create("./generator/templates/testing.go")
	if err != nil {
		log.Error("create file: ", err)
		return
	}

	// Run the template to verify the output.
	err = tmpl.Execute(f, one)
	if err != nil {
		log.Error("execution: %s", err)
	}

	f.Close()
}
