package env

import (
	"log"

	"github.com/fusor/cpma/internal/sftpclient"
	"github.com/spf13/viper"
)

// Info structures the application settings.
type Info struct {
	SFTP       sftpclient.Info `mapstructure:"Source"`
	OutputPath string          `mapstructure:"outputPath"`
}

// New returns a instance of the application settings.
func New() *Info {
	var info Info

	if err := viper.Unmarshal(&info); err != nil {
		log.Fatalf("unable to parse configuration: %v", err)
	}

	return &info
}
