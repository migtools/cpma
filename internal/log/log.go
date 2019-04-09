package log

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

type customFormatter struct {
	log.TextFormatter
}

func (f *customFormatter) Format(entry *log.Entry) ([]byte, error) {
	_, err := f.TextFormatter.Format(entry)

	time := entry.Time.Format("04/10 15:04:05")
	level := strings.ToUpper(entry.Level.String())

	str := fmt.Sprintf("%s %5s %s", time, level, entry.Message)
	return []byte(str), err
}

func init() {
	f, err := os.OpenFile("cpma.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	// TODO: Replace with flag
	//mw := io.MultiWriter(os.Stdout, f)
	log.SetOutput(f)
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&customFormatter{})
}
