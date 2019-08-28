package reportoutput

import (
	"encoding/json"

	"github.com/fusor/cpma/pkg/io"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func jsonOutput(r ReportOutput) {
	jsonReports, err := json.MarshalIndent(r, "", " ")
	if err != nil {
		panic(errors.Wrap(err, "unable to marshal reports"))
	}

	if err := io.WriteFile(jsonReports, jsonFileName); err != nil {
		panic(errors.Wrapf(err, "unable to write to report file: %s", jsonFileName))
	}

	logrus.Infof("Report:Added: %s", jsonFileName)
}
