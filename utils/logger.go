package utils

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

var Logger *zap.Logger

func InitializedLogger() *zap.Logger {
	fileLoggerConfig := zap.NewProductionEncoderConfig()
	fileLoggerConfig.MessageKey = "message"
	fileLoggerConfig.LevelKey = "level"
	fileLoggerConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	fileLoggerConfig.TimeKey = "timestamp"
	fileLoggerConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	fileLoggerConfig.CallerKey = "caller"
	fileLoggerConfig.EncodeCaller = zapcore.ShortCallerEncoder
	fileLoggerConfig.FunctionKey = "func"
	logFile, _ := os.OpenFile("logs/errors.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	core := zapcore.NewTee(
		// logger to record in warn level (including errors) to masterdata.log
		zapcore.NewCore(
			zapcore.NewJSONEncoder(fileLoggerConfig),
			zapcore.AddSync(logFile),
			zapcore.WarnLevel,
		),
		// logger to record in debug level in terminal
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
			zapcore.AddSync(os.Stdout),
			zapcore.DebugLevel,
		),
	)

	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}

// RequestLog logs all requests that occurs when service is running
func RequestLog(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ec echo.Context) error {
		err := next(ec)
		if err != nil {
			ec.Error(err)
		}

		request := ec.Request()
		response := ec.Response()

		fields := []zapcore.Field{
			zap.Int("status", response.Status),
			zap.String("latency", time.Since(time.Now()).String()),
			zap.String("method", request.Method),
			zap.String("uri", request.RequestURI),
			zap.String("remote_ip", ec.RealIP()),
		}

		statusCode := response.Status
		switch {
		case statusCode >= 500:
			Logger.Error("Internal Server Error", fields...)
		case statusCode >= 400:
			Logger.Warn("Client-side Error", fields...)
		case statusCode >= 300:
			Logger.Info("Redirection", fields...)
		default:
			Logger.Info("Success", fields...)
		}

		return nil
	}
}
