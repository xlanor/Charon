package logger

import (
	"go.uber.org/zap"
)

func Sugar() *zap.SugaredLogger {
	return zap.L().Sugar()
}
