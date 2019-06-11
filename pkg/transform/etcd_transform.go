package transform

import (
	"errors"
	"fmt"
	"gopkg.in/go-ini/ini.v1"
	"strings"

	"github.com/fusor/cpma/pkg/env"
	"github.com/fusor/cpma/pkg/io"
	"github.com/sirupsen/logrus"
)

// ETCDExtraction holds ETCD data extracted from OCP3
type ETCDExtraction struct {
	TLSCipherSuites string
	ClientPort      string
}

// ETCDTransform is an ETCD specific transform
type ETCDTransform struct {
}

// Transform converts data collected from an OCP3 into a useful output
func (e ETCDExtraction) Transform() ([]Output, error) {
	logrus.Info("ETCDTransform::Transform")
	reports, err := e.buildReportOutput()
	if err != nil {
		return nil, err
	}
	outputs := []Output{reports}
	return outputs, nil
}

func (e ETCDExtraction) buildReportOutput() (Output, error) {
	reportOutput := ReportOutput{
		Component: ETCDComponentName,
	}

	var TLSConfidence = "green"
	var TLSMessage = "No Custom TLS Cipher Suites were found"
	var clientPortConfidence = "green"

	if e.ClientPort != "2379" {
		clientPortConfidence = "red"
	}

	if e.TLSCipherSuites != "" {
		TLSConfidence = "red"
		TLSMessage = fmt.Sprintf("The Openshift 4 ETCD Cluster is not configurable. TLS Cipher Suite configuration was detected, %v", e.TLSCipherSuites)
	}

	reportOutput.Reports = append(reportOutput.Reports,
		Report{
			Name:       "ETCD Client Port",
			Kind:       "Configuration",
			Supported:  false,
			Confidence: clientPortConfidence,
			Comment:    fmt.Sprintf("The Openshift 4 ETCD Cluster is not configurable and uses port 2379. Your Openshift 3 Cluster is using port %v", e.ClientPort),
		})

	reportOutput.Reports = append(reportOutput.Reports,
		Report{
			Name:       "ETCD TLS Cipher Suites",
			Kind:       "Configuration",
			Supported:  false,
			Confidence: TLSConfidence,
			Comment:    TLSMessage,
		})

	return reportOutput, nil
}

// Extract collects ETCD configuration from an OCP3 cluster
func (e ETCDTransform) Extract() (Extraction, error) {
	logrus.Info("ETCDTransform::Extract")
	content, err := io.FetchFile(env.Config().GetString("ETCDConfigFile"))
	if err != nil {
		return nil, err
	}

	ETCDConfig, err := ini.Load(content)
	if err != nil {
		return nil, err
	}

	var extraction ETCDExtraction
	portArray := strings.Split(ETCDConfig.Section("").Key("ETCD_LISTEN_CLIENT_URLS").String(), ":")
	extraction.ClientPort = portArray[len(portArray)-1]
	extraction.TLSCipherSuites = ETCDConfig.Section("").Key("ETCD_CIPHER_SUITES").String()

	return extraction, nil
}

// Validate confirms we have recieved good ETCD configuration data during Extract
func (e ETCDExtraction) Validate() error {

	if e.ClientPort == "" {
		err := errors.New("ETCD Client Port could not be determined")
		return err
	}

	return nil
}

// Name returns a human readable name for the transform
func (e ETCDTransform) Name() string {
	return "ETCD"
}
