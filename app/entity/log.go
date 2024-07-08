package entity

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

const (
	storagePath    string = "storage"
	logPath        string = "log"
	errorLogFile   string = "error.log"
	warningLogFile string = "warning.log"
)

type RapCLog struct {
	log               *logrus.Logger
	uri               string
	method            string
	status            int
	message           string
	err               error
	enableWarnFileLog bool
}

func InitLog(uri, method, message string, status int, err error, enableWarnFileLog bool) RapCLog {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	return RapCLog{
		log:               logger,
		uri:               uri,
		method:            method,
		status:            status,
		message:           message,
		err:               err,
		enableWarnFileLog: enableWarnFileLog,
	}
}

func (e RapCLog) Log() {
	logrusFields := logrus.Fields{
		"URI":    e.uri,
		"Method": e.method,
		"Status": e.status,
		"Error":  e.err,
	}

	if e.err == nil {
		delete(logrusFields, "Error")
		e.log.WithFields(logrusFields).Info(e.message)
	} else {
		// create error log file
		errorLog, err := os.OpenFile(filepath.Join(storagePath, logPath, errorLogFile), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			e.log.Errorf("Failed to create error file hook: %v", err)
		} else {
			defer errorLog.Close()

			e.log.AddHook(&FileHook{
				writer:    errorLog,
				logLevels: []logrus.Level{logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel},
			})
		}

		// create warning log file
		if e.enableWarnFileLog {
			warnLog, err := os.OpenFile(filepath.Join(storagePath, logPath, warningLogFile), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				e.log.Errorf("Failed to create warning file hook: %v", err)
			} else {
				defer warnLog.Close()

				e.log.AddHook(&FileHook{
					writer:    warnLog,
					logLevels: []logrus.Level{logrus.WarnLevel},
				})
			}
		}

		var message interface{} = fmt.Sprintf("%s error", e.message)
		if _err, ok := e.err.(*echo.HTTPError); ok {
			message = _err.Message
			logrusFields["error"] = _err.Internal
			if cErr, ok := _err.Internal.(*InternalError); ok {
				logrusFields["code"] = cErr.Code
			}
		}

		if e.status < http.StatusInternalServerError {
			e.log.WithFields(logrusFields).Warn(message)
		} else {
			e.log.WithFields(logrusFields).Error(message)
		}
	}
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
