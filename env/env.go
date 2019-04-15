package env

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/fusor/cpma/internal/config"
	"github.com/fusor/cpma/internal/sftpclient"
	log "github.com/sirupsen/logrus"
)

// Info describes application environment and configuration
type Info struct {
	Source    string          `mapstructure:"Source"`
	SSH       sftpclient.Info `mapstructure:"SSHCreds"`
	OutputDir string          `mapstructure:"OutputDir"`
}

// New returns a instance of the application settings.
func New() *Info {
	var info Info

	if err := config.Config().Unmarshal(&info); err != nil {
		log.Fatalf("unable to parse configuration: %v", err)
	}

	log.Debugln(spew.Sdump(info))
	return &info
}
