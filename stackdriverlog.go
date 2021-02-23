package superr

import (
	"cloud.google.com/go/errorreporting"
	"cloud.google.com/go/logging"
	"github.com/sirupsen/logrus"
)

type StackDriverHook struct {
	client      *logging.Client
	errorClient *errorreporting.Client
	logger      *logging.Logger
}

var logLevelMappings = map[logrus.Level]logging.Severity{
	logrus.TraceLevel: logging.Default,
	logrus.DebugLevel: logging.Debug,
	logrus.InfoLevel:  logging.Info,
	logrus.WarnLevel:  logging.Warning,
	logrus.ErrorLevel: logging.Error,
	logrus.FatalLevel: logging.Critical,
	logrus.PanicLevel: logging.Critical,
}
