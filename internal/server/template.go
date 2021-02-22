package server

import (
	"io"
	"os"
	"path/filepath"
	"text/template"
)

func parseIndexTemplate(uiDir string) (indexTemplate *template.Template, err error) {
	templateFilepath := filepath.Join(uiDir, "index.html")

	file, err := os.Open(templateFilepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	if err := file.Close(); err != nil {
		return nil, err
	}

	indexTemplate = template.New("index.html")
	return indexTemplate.Parse(string(b))
}
