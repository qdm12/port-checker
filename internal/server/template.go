package server

import (
	"os"
	"path/filepath"
	"text/template"
)

func parseTemplate(uiDir string) (template *template.Template, err error) {
	templateFilepath := filepath.Join(uiDir, "index.html")

	file, err := os.Open(templateFilepath)
	if err != nil {
		return nil, err
	}
	if err := file.Close(); err != nil {
		return nil, err
	}

	return template.ParseFiles(templateFilepath)
}
