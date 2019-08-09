package reportoutput

import (
	"encoding/json"

	"github.com/fusor/cpma/pkg/io"
	"github.com/pkg/errors"
)

func jsonOutput(r ReportOutput) error {
	jsonFile := "report.json"

	jsonReports, err := json.MarshalIndent(r, "", " ")
	if err != nil {
		return errors.Wrap(err, "unable to marshal reports")
	}

	if err := io.WriteFile(jsonReports, jsonFile); err != nil {
		return errors.Wrapf(err, "unable to write to report file: %s", jsonFile)
	}

	return nil
}
