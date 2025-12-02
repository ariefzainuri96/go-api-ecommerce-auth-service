package logger

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

func NewLogger() *zap.Logger {
    config := zap.NewProductionConfig()

    // Format timestamps to ISO8601 for easier reading
    config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

    logger, err := config.Build()
    if err != nil {
        panic(err)
    }

    return logger
}
