package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(service string) (*zap.SugaredLogger, error) {

	config := zap.NewDevelopmentConfig()

	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.InitialFields = map[string]interface{}{
		"service": service,
	}

	log, err := config.Build(zap.WithCaller(true))
	if err != nil {
		return nil, err
	}

	return log.Sugar(), nil
}
