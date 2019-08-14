package reportoutput

import (
	"encoding/json"

	"github.com/fusor/cpma/pkg/io"
	"github.com/pkg/errors"
)

func jsonOutput(r ReportOutput) {
	jsonFile := "report.json"

	jsonReports, err := json.MarshalIndent(r, "", " ")
	if err != nil {
		panic(errors.Wrap(err, "unable to marshal reports"))
	}

	if err := io.WriteFile(jsonReports, jsonFile); err != nil {
		panic(errors.Wrapf(err, "unable to write to report file: %s", jsonFile))
	}
}
