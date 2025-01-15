package logging

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

var (
	logger *logrus.Logger
)

// GetLogger returns the initialized logrus.Logger instance or panics if the logger has not been initialized.
func GetLogger() *logrus.Logger {
	if logger == nil {
		panic(fmt.Errorf("logger not initialized"))
	}
	return logger
}

func init() {
	logger = logrus.New()
}
