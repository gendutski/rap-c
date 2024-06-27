package middleware

import (
	"net/http"
	"os"
	"path/filepath"
	"rap-c/app/entity"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

const (
	logPath        string = "log"
	errorLogFile   string = "error.log"
	warningLogFile string = "warning.log"
)

func SetLog(enableWarnFileLog bool) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:      true,
		LogStatus:   true,
		LogError:    true,
		HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger := logrus.New()

			logger.SetFormatter(&logrus.TextFormatter{
				FullTimestamp: true,
			})
			if v.Error == nil {
				logger.WithFields(logrus.Fields{
					"URI":    v.URI,
					"status": v.Status,
				}).Info("request")
			} else {
				// create error log file
				errorLog, err := os.OpenFile(filepath.Join(logPath, errorLogFile), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
				if err != nil {
					logger.Errorf("Failed to create error file hook: %v", err)
				} else {
					defer errorLog.Close()

					logger.AddHook(&FileHook{
						writer:    errorLog,
						logLevels: []logrus.Level{logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel},
					})
				}

				// create warning log file
				if enableWarnFileLog {
					warnLog, err := os.OpenFile(filepath.Join(logPath, warningLogFile), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
					if err != nil {
						logger.Errorf("Failed to create warning file hook: %v", err)
					} else {
						defer warnLog.Close()

						logger.AddHook(&FileHook{
							writer:    warnLog,
							logLevels: []logrus.Level{logrus.WarnLevel},
						})
					}
				}

				logrusFields := logrus.Fields{
					"URI":    v.URI,
					"status": v.Status,
					"error":  v.Error,
				}
				var message interface{} = "request error"

				if _err, ok := v.Error.(*echo.HTTPError); ok {
					message = _err.Message
					logrusFields["error"] = _err.Internal
					if cErr, ok := _err.Internal.(*entity.InternalError); ok {
						logrusFields["code"] = cErr.Code
					}
				}

				if v.Status < http.StatusInternalServerError {
					logger.WithFields(logrusFields).Warn(message)
				} else {
					logger.WithFields(logrusFields).Error(message)
				}
			}
			return nil
		},
	})
}

// logrus file hook
type FileHook struct {
	writer    *os.File
	logLevels []logrus.Level
}

// Fire write log ke file
func (hook *FileHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	_, err = hook.writer.Write([]byte(line))
	return err
}

// Levels returns the levels logged by this hook
func (hook *FileHook) Levels() []logrus.Level {
	return hook.logLevels
}
