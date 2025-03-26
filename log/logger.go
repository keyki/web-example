package log

import (
	"bytes"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
	"web-example/types"
)

var baseLogger = logrus.New()

type CustomFormatter struct{}

func init() {
	baseLogger.SetFormatter(&CustomFormatter{})
	baseLogger.SetLevel(logrus.InfoLevel)
}

func BaseLogger() *logrus.Entry {
	return baseLogger.WithFields(logrus.Fields{
		"request-id": "application",
	})
}

func Logger(ctx context.Context) *logrus.Entry {
	if ctx == nil {
		return BaseLogger()
	}
	if logger, ok := ctx.Value(types.LogKey).(*logrus.Entry); ok {
		return logger
	}
	return BaseLogger()
}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b bytes.Buffer

	// Timestamp
	b.WriteString(fmt.Sprintf(`"%s" `, entry.Time.Format(time.RFC3339)))

	// Request ID (or any other fields except "msg")
	for key, value := range entry.Data {
		if key != "msg" { // Avoid printing msg here
			b.WriteString(fmt.Sprintf(`%s=%v `, key, value))
		}
	}

	// Finally, append the msg field at the end
	b.WriteString(fmt.Sprintf(`"%s"`, entry.Message))

	b.WriteByte('\n') // End with a newline
	return b.Bytes(), nil
}
