package transform_test

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/fusor/cpma/pkg/env"
	"github.com/fusor/cpma/pkg/transform"
	"github.com/fusor/cpma/pkg/transform/reportoutput"
	"github.com/stretchr/testify/assert"
	"gopkg.in/go-ini/ini.v1"
)

var testExtraction = func() transform.ETCDExtraction {
	file := "testdata/etcd.conf"
	content, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("Error loading test file: %s\n", file)
	}

	ETCDConfig, err := ini.Load(content)
	if err != nil {
		fmt.Printf("Error loading ini content from file: %s\n", file)
	}

	var extraction transform.ETCDExtraction
	portArray := strings.Split(ETCDConfig.Section("").Key("ETCD_LISTEN_CLIENT_URLS").String(), ":")
	extraction.ClientPort = portArray[len(portArray)-1]
	extraction.TLSCipherSuites = ETCDConfig.Section("").Key("ETCD_CIPHER_SUITES").String()

	return extraction
}()

func TestETCDExtractionTransform(t *testing.T) {
	expectedReport := reportoutput.ComponentReport{
		Component: "ETCD",
	}

	expectedReport.Reports = append(expectedReport.Reports,
		reportoutput.Report{
			Name:       "ETCD Client Port",
			Kind:       "Configuration",
			Supported:  false,
			Confidence: 2,
			Comment:    "The Openshift 4 ETCD Cluster is not configurable and uses port 2379. Your Openshift 3 Cluster is using port 2379",
		})
	expectedReport.Reports = append(expectedReport.Reports,
		reportoutput.Report{
			Name:       "ETCD TLS Cipher Suites",
			Kind:       "Configuration",
			Supported:  false,
			Confidence: 0,
			Comment:    "The Openshift 4 ETCD Cluster is not configurable. TLS Cipher Suite configuration was detected, TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256",
		})

	expectedReportOutput := reportoutput.ReportOutput{
		ComponentReports: []reportoutput.ComponentReport{expectedReport},
	}

	testCases := []struct {
		name            string
		expectedReports reportoutput.ReportOutput
	}{
		{
			name:            "transform etcd extraction",
			expectedReports: expectedReportOutput,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualReportsChan := make(chan reportoutput.ReportOutput)
			transform.FinalReportOutput = transform.Report{}

			// Override flush method
			transform.ReportOutputFlush = func(reports transform.Report) error {
				actualReportsChan <- reports.Report
				return nil
			}

			go func() {
				env.Config().Set("Reporting", true)
				env.Config().Set("Manifests", true)

				_, err := testExtraction.Transform()
				if err != nil {
					t.Error(err)
				}
				transform.FinalReportOutput.Flush()
			}()

			actualReports := <-actualReportsChan
			assert.Equal(t, actualReports.ComponentReports, tc.expectedReports.ComponentReports)
		})

	}
}
