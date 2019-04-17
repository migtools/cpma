package log

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

func init() {
	f, err := os.OpenFile("cpma.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatalf("error opening file: %v", err)
	}

	mw := io.MultiWriter(os.Stdout, f)
	logrus.SetOutput(mw)
	logrus.SetLevel(logrus.InfoLevel)

}
