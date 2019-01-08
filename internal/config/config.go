package config

import (
	"github.com/yarnaid/vm_syn/internal/logger"
)

// Config stores all application data
type Config struct {
	Logger *logger.Logger
}
