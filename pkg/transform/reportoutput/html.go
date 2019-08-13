package reportoutput

import (
	"encoding/json"
	"html/template"
	"os"
	"path/filepath"

	rice "github.com/GeertJohan/go.rice"
	"github.com/fusor/cpma/pkg/env"
	"k8s.io/apimachinery/pkg/api/resource"
)

// Output reads report stucture, generates html using go templates and writes it to a file
func htmlOutput(report ReportOutput) error {
	path := filepath.Join(env.Config().GetString("WorkDir"), "report.html")

	f, err := os.Create(path)
	defer f.Close()

	if err != nil {
		return err
	}

	htmlTemplate, err := parseTemplates()
	if err != nil {
		return err
	}

	err = htmlTemplate.Execute(f, report)
	if err != nil {
		return err
	}

	return nil
}

func parseTemplates() (*template.Template, error) {
	templateBox, err := rice.FindBox("staticpage")
	if err != nil {
		return nil, err
	}

	bootstrapCSS, err := templateBox.String("css/bootstrap.min.css")
	if err != nil {
		return nil, err
	}

	stylesCSS, err := templateBox.String("css/styles.css")
	if err != nil {
		return nil, err
	}

	bootstrapJS, err := templateBox.String("js/bootstrap.min.js")
	if err != nil {
		return nil, err
	}

	jqueryJS, err := templateBox.String("js/jquery-3.3.1.slim.min.js")
	if err != nil {
		return nil, err
	}

	popperJS, err := templateBox.String("js/popper.min.js")
	if err != nil {
		return nil, err
	}

	helpersTemplateString, err := templateBox.String("templates/helpers.gohtml")
	if err != nil {
		return nil, err
	}

	nodesTemplateString, err := templateBox.String("templates/nodes.gohtml")
	if err != nil {
		return nil, err
	}

	quotasTemplateString, err := templateBox.String("templates/quotas.gohtml")
	if err != nil {
		return nil, err
	}

	namespacesTemplateString, err := templateBox.String("templates/namespaces.gohtml")
	if err != nil {
		return nil, err
	}

	pvsTemplateString, err := templateBox.String("templates/pvs.gohtml")
	if err != nil {
		return nil, err
	}

	clusterReportTemplateString, err := templateBox.String("templates/cluster-report.gohtml")
	if err != nil {
		return nil, err
	}

	mainTemplateString, err := templateBox.String("templates/main.gohtml")
	if err != nil {
		return nil, err
	}

	htmlTemplate := template.Must(template.New("html").Parse(helpersTemplateString))

	htmlTemplate = template.Must(htmlTemplate.Parse(nodesTemplateString))

	htmlTemplate = template.Must(htmlTemplate.Funcs(template.FuncMap{
		"formatQuantity": func(q resource.Quantity) string {
			json, _ := json.Marshal(q)
			return string(json)
		},
	}).Parse(quotasTemplateString))

	htmlTemplate = template.Must(htmlTemplate.Parse(namespacesTemplateString))

	htmlTemplate = template.Must(htmlTemplate.Parse(pvsTemplateString))

	htmlTemplate = template.Must(htmlTemplate.Parse(clusterReportTemplateString))

	htmlTemplate = template.Must(htmlTemplate.Funcs(template.FuncMap{
		"bootstrapCSS": func() template.CSS {
			return template.CSS(bootstrapCSS)
		},
		"stylesCSS": func() template.CSS {
			return template.CSS(stylesCSS)
		},
		"bootstrapJS": func() template.JS {
			return template.JS(bootstrapJS)
		},
		"jqueryJS": func() template.JS {
			return template.JS(jqueryJS)
		},
		"popperJS": func() template.JS {
			return template.JS(popperJS)
		},
	}).Parse(mainTemplateString))

	return htmlTemplate, nil
}
