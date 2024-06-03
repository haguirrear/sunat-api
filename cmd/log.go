package cmd

import (
	"os"

	"github.com/haguirrear/sunatapi/pkg/logger"
)

func GetLogger() *logger.Logger {
	var level logger.LogLevel

	switch VerboseCount {
	case 0:
		level = logger.WarnLevel
	case 1:
		level = logger.DebugLevel
	default:
		level = logger.TraceLevel
	}

	return logger.NewLogger(os.Stderr, level)
}
