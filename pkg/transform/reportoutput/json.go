package reportoutput

import (
	"encoding/json"

	"github.com/fusor/cpma/pkg/io"
	"github.com/pkg/errors"
)

func jsonOutput(r ReportOutput) {
	jsonReports, err := json.MarshalIndent(r, "", " ")
	if err != nil {
		panic(errors.Wrap(err, "unable to marshal reports"))
	}

	if err := io.WriteFile(jsonReports, jsonFileName); err != nil {
		panic(errors.Wrapf(err, "unable to write to report file: %s", jsonFileName))
	}
}
