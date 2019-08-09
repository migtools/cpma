package reportoutput

import (
	"html/template"
	"os"
	"path/filepath"

	"github.com/fusor/cpma/pkg/env"
	"github.com/fusor/cpma/pkg/transform/reportoutput/templates"
)

var htmlTemplate = template.Must(template.New("html").Parse(templates.ReportHTML))

// Output reads report stucture, generates html using go templates and writes it to a file
func htmlOutput(report ReportOutput) error {
	path := filepath.Join(env.Config().GetString("WorkDir"), "report.html")

	f, err := os.Create(path)
	defer f.Close()

	if err != nil {
		return err
	}

	err = htmlTemplate.Execute(f, report)
	if err != nil {
		return err
	}

	return nil
}
