package common

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()
var levelStr = "debug"

func init() {
	Logger.SetReportCaller(true)
	if v := os.Getenv("LOG_LEVEL"); v != "" {
		levelStr = strings.ToLower(v)
	}
	LoggerLevel, err := logrus.ParseLevel(levelStr)
	if err != nil {
		LoggerLevel = logrus.DebugLevel
	}
	Logger.SetLevel(LoggerLevel)
	Logger.SetFormatter(&MyFormatter{})
}

type MyFormatter struct{}

func (m *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	var newLog string

	//HasCaller()为true才会有调用信息
	if entry.HasCaller() {
		fName := filepath.Base(entry.Caller.File)
		newLog = fmt.Sprintf("[%s] [%s] [%s:%d %s] %s\n",
			timestamp, entry.Level, fName, entry.Caller.Line, entry.Caller.Function, entry.Message)
	} else {
		newLog = fmt.Sprintf("[%s] [%s] %s\n", timestamp, entry.Level, entry.Message)
	}

	b.WriteString(newLog)
	return b.Bytes(), nil
}
