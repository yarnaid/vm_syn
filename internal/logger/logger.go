package logger

import (
	"github.com/sirupsen/logrus"
)

// Logger is internal logger interface
type Logger interface {
	logrus.FieldLogger
}
