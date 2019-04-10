package log

import (
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

func init() {
	f, err := os.OpenFile("cpma.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	mw := io.MultiWriter(os.Stdout, f)
	log.SetOutput(mw)
	log.SetLevel(log.InfoLevel)

	log.Println("CPMA Log started")
}
