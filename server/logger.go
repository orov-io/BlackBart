package server

import (
	stackdriver "github.com/TV4/logrus-stackdriver-formatter"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

// setLogger is called by the service. It adds the logger if their option is sets
// when you initialize the service.
func setLogger(options *LoggerOptions) {

	if options.Format != nil {
		log.Formatter = options.Format
	}

	log.Level = options.Level

	log.Infof("Start to loggin with %v config", options.Env)
}

// GetLogger returns a logger ready to log to stackdriver
func GetLogger() *logrus.Logger {
	return log
}

func localLogging() (logrus.Level, logrus.Formatter) {
	return logrus.TraceLevel, nil
}

func nonProdServerLogging() (logrus.Level, logrus.Formatter) {
	return logrus.DebugLevel, getLoggerForServer()
}

func prodServerLogging() (logrus.Level, logrus.Formatter) {
	return logrus.InfoLevel, getLoggerForServer()
}

func getLoggerForServer() logrus.Formatter {
	return stackdriver.NewFormatter(
		stackdriver.WithService(GetEnvOrDefaultString(serviceNameKey)),
		stackdriver.WithVersion(GetEnvOrDefaultString(serviceVersionKey)),
	)
}

func getDefaultLoggerConfByEnv(env string) (logrus.Level, logrus.Formatter) {
	switch env {
	case "local":
		return localLogging()
	case "dev":
		return nonProdServerLogging()
	case "pre":
		return nonProdServerLogging()
	case "prod":
		return prodServerLogging()
	}

	return nonProdServerLogging()
}
