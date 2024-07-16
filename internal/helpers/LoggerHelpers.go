package helpers

import "go.uber.org/zap"

var TLog = Log()

func Log() *zap.SugaredLogger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	return logger.Sugar()
}
