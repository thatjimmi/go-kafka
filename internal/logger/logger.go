package logger

import (
	"go.uber.org/zap"
)

var l *zap.Logger

func init() {
	var err error
	l, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
}

func New() *zap.Logger {
	return l
}
