package transform_test

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/fusor/cpma/pkg/transform"
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
	expectedReport := transform.ComponentReport{
		Component: "ETCD",
	}

	expectedReport.Reports = append(expectedReport.Reports,
		transform.Report{
			Name:       "ETCD Client Port",
			Kind:       "Configuration",
			Supported:  false,
			Confidence: 2,
			Comment:    "The Openshift 4 ETCD Cluster is not configurable and uses port 2379. Your Openshift 3 Cluster is using port 2379",
		})
	expectedReport.Reports = append(expectedReport.Reports,
		transform.Report{
			Name:       "ETCD TLS Cipher Suites",
			Kind:       "Configuration",
			Supported:  false,
			Confidence: 0,
			Comment:    "The Openshift 4 ETCD Cluster is not configurable. TLS Cipher Suite configuration was detected, TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256",
		})

	expectedReportOutput := transform.ReportOutput{
		ComponentReports: []transform.ComponentReport{expectedReport},
	}

	testCases := []struct {
		name            string
		expectedReports transform.ReportOutput
	}{
		{
			name:            "transform crio extraction",
			expectedReports: expectedReportOutput,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualReportsChan := make(chan transform.ReportOutput)

			// Override flush method
			transform.ReportOutputFlush = func(reports transform.ReportOutput) error {
				actualReportsChan <- reports
				return nil
			}

			go func() {
				transformOutput, err := testExtraction.Transform()
				if err != nil {
					t.Error(err)
				}
				for _, output := range transformOutput {
					output.Flush()
				}
			}()

			actualReports := <-actualReportsChan
			assert.Equal(t, actualReports, tc.expectedReports)
		})

	}
}
