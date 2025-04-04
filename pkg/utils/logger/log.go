package logger

import (
	"context"
	"github.com/sirupsen/logrus"
	"good-template-go/conf"
	"os"
	"sync"
)

var (
	DefaultLogger    *logrus.Logger
	DefaultBaseEntry *logrus.Entry
	initOnce         sync.Once
)

func Init(name string) {
	initOnce.Do(func() {
		DefaultLogger = logrus.New()

		// Set log level
		logLevel := conf.GetConfig().LoggerLevel
		if l, e := logrus.ParseLevel(logLevel); e == nil {
			DefaultLogger.SetLevel(l)
		}

		// Set log format
		logFormat := conf.GetConfig().LogFormat
		if logFormat == "json" {
			DefaultLogger.SetFormatter(&logrus.JSONFormatter{
				TimestampFormat:  "",
				DisableTimestamp: false,
				DataKey:          "",
				FieldMap:         nil,
				CallerPrettyfier: nil,
				PrettyPrint:      false,
			})
		} else {
			DefaultLogger.SetFormatter(&logrus.TextFormatter{
				ForceColors:      true,
				FullTimestamp:    true,
				PadLevelText:     true,
				ForceQuote:       true,
				QuoteEmptyFields: true,
			})
		}

		DefaultLogger.SetOutput(os.Stdout)
		DefaultBaseEntry = DefaultLogger.WithField("type", name)
	})
}

// LogError ghi log lỗi với thông điệp và chi tiết lỗi
func LogError(log *logrus.Entry, err error, message string) {
	log.WithError(err).Error("*** " + message + " ***")
}

// Tag sets a tag name then returns a log entry ready to write
func Tag(tag string) *logrus.Entry {
	if DefaultBaseEntry == nil {
		Init("function")
	}
	return DefaultBaseEntry.WithField("tag", tag)
}

func WithTag(tag string) *logrus.Entry {
	l := Tag(tag)
	return l
}

func WithCtx(ctx context.Context, tag string) *logrus.Entry {
	l := Tag(tag)
	if requestID, ok := ctx.Value("x-request-id").(string); ok && requestID != "" {
		l = l.WithField("x-request-id", requestID)
	}
	return l
}
