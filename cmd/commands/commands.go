package commands

import (
	"flag"
	"fmt"
)

func generateFromSpec() {
	fPackage := flag.String("package", "main", "package test")
	fCopyright := flag.Bool("copyright", true, "copyright test")
	fImports := flag.String("imports", "", "import test")

	flag.Parse()

	fmt.Printf("textPtr: %s, metricPtr: %t, uniquePtr: %s\n", *fPackage, *fCopyright, *fImports)
}
