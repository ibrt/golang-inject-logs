package logz

import (
	"github.com/getsentry/sentry-go"
	"github.com/ibrt/golang-errors/errorz"
	"github.com/sirupsen/logrus"
)

// Level describes a level.
type Level string

func (l Level) toLogrus() logrus.Level {
	switch l {
	case Debug:
		return logrus.DebugLevel
	case Info:
		return logrus.InfoLevel
	case Warning:
		return logrus.WarnLevel
	case Error:
		return logrus.ErrorLevel
	default:
		panic(errorz.Errorf("unknown level: %v", errorz.A(l), errorz.SkipPackage()))
	}
}

func (l Level) toSentry() sentry.Level {
	switch l {
	case Debug:
		return sentry.LevelDebug
	case Info:
		return sentry.LevelInfo
	case Warning:
		return sentry.LevelWarning
	case Error:
		return sentry.LevelError
	default:
		panic(errorz.Errorf("unknown level: %v", errorz.A(l), errorz.SkipPackage()))
	}
}

func levelFromSentry(l sentry.Level) Level {
	switch l {
	case sentry.LevelFatal, sentry.LevelError:
		return Error
	case sentry.LevelWarning:
		return Warning
	case sentry.LevelInfo:
		return Info
	case sentry.LevelDebug:
		return Debug
	default:
		panic(errorz.Errorf("unknown level: %v", errorz.A(l), errorz.SkipPackage()))
	}
}

// Known levels.
const (
	Debug   Level = "debug"
	Info    Level = "info"
	Warning Level = "warning"
	Error   Level = "error"
)
