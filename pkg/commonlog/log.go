package log

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"

	//"git.internal.yunify.com/manage/common/config"

	"github.com/rifflock/lfshook"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

// 封装logrus.Fields
type Fields logrus.Fields

func GetLogger() *logrus.Logger {
	return logger
}

type File struct {
	Name       string
	Path       string
	MaxHour    int64
	RotateHour int64
}
type Options struct {
	Level   string
	LogPath string
	File    *File
}

func Init(options *Options) {
	initLogger(options)
}
func initLogger(options *Options) {
	SetLogLevelFromStr(options.Level)

	var errorWriter, infoWriter io.Writer
	if options.File != nil && options.File.Path != "" {
		logPath := options.File.Path
		if mode, err := os.Stat(logPath); err != nil {
			if os.IsNotExist(err) {
				if err := os.MkdirAll(logPath, os.ModePerm); err != nil {
					panic(err)
				}
			} else {
				if !mode.IsDir() {
					panic(fmt.Sprintf("log path [%v] is not a directory", logPath))
				}
			}
		}
		var infoLog, errorLog string
		if options.File.Name != "" {
			infoLog = fmt.Sprintf("%s/%s.log", logPath, options.File.Name)
			errorLog = fmt.Sprintf("%s/%s.log.wf", logPath, options.File.Name)
		} else {
			infoLog = logPath + "/info.log"
			errorLog = logPath + "/error.log"
		}
		// 日志分割
		var err error

		logFile := options.File
		errorWriter, err = rotatelogs.New(
			errorLog+".%Y%m%d%H%M",
			rotatelogs.WithLinkName(errorLog),                                        // 生成软链，指向最新日志文件
			rotatelogs.WithMaxAge(time.Duration(logFile.MaxHour)*time.Hour),          // 文件最大保存时间
			rotatelogs.WithRotationTime(time.Duration(logFile.RotateHour)*time.Hour), // 日志切割时间间隔
		)
		if err != nil {
			logger.Errorf("config local file system logger error. %v", errors.WithStack(err))
		}
		infoWriter, err = rotatelogs.New(
			infoLog+".%Y%m%d%H%M",
			rotatelogs.WithLinkName(infoLog),                                         // 生成软链，指向最新日志文件
			rotatelogs.WithMaxAge(time.Duration(logFile.MaxHour)*time.Hour),          // 文件最大保存时间
			rotatelogs.WithRotationTime(time.Duration(logFile.RotateHour)*time.Hour), // 日志切割时间间隔
		)
		if err != nil {
			logger.Errorf("config local file system logger error. %v", errors.WithStack(err))
		}
		formatter := &logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.283+0800",
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "content",
				logrus.FieldKeyFunc:  "caller",
			},
		}
		lfHook := lfshook.NewHook(lfshook.WriterMap{
			logrus.TraceLevel: infoWriter, // 为不同级别设置不同的输出目的
			logrus.DebugLevel: infoWriter,
			logrus.InfoLevel:  infoWriter,
			logrus.WarnLevel:  infoWriter,
			logrus.ErrorLevel: errorWriter,
			logrus.FatalLevel: errorWriter,
			logrus.PanicLevel: errorWriter,
		}, formatter)
		logger.AddHook(lfHook)
		//logger.SetOutput(ioutil.Discard)
	} else {
		infoWriter = os.Stdout
		errorWriter = os.Stderr

		formatter := &logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02T15:04:05.283+0800",
		}
		SetLogFormatter(formatter)
	}
}

func SetLogLevel(level logrus.Level) {
	logger.Level = level
}

func SetLogFormatter(formatter logrus.Formatter) {
	logger.Formatter = formatter
}

func SetLogLevelFromStr(level string) {
	switch l := strings.ToLower(level); l {
	case "trace":
		SetLogLevel(logrus.TraceLevel)
	case "debug":
		SetLogLevel(logrus.DebugLevel)
	case "info":
		SetLogLevel(logrus.InfoLevel)
	case "warn":
		SetLogLevel(logrus.WarnLevel)
	case "error":
		SetLogLevel(logrus.ErrorLevel)
	case "fatal":
		SetLogLevel(logrus.FatalLevel)
	default:
		SetLogLevel(logrus.InfoLevel)
	}
}

// Trace
func Trace(args ...interface{}) {
	if logger.Level >= logrus.TraceLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Trace(args...)
	}
}

// Trace
func Tracef(format string, args ...interface{}) {
	arg := fmt.Sprintf(format, args...)
	if logger.Level >= logrus.TraceLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Trace(arg)
	}
}

// 带有field的Trace
func TraceWithFields(l interface{}, f Fields) {
	if logger.Level >= logrus.TraceLevel {
		entry := logger.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(2)
		entry.Trace(l)
	}
}

// Debug
func Debug(args ...interface{}) {
	if logger.Level >= logrus.DebugLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Debug(args...)
	}
}

// Debugf
func Debugf(format string, args ...interface{}) {
	arg := fmt.Sprintf(format, args...)
	if logger.Level >= logrus.DebugLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Debug(arg)
	}
}

// 带有field的Debug
func DebugWithFields(l interface{}, f Fields) {
	if logger.Level >= logrus.DebugLevel {
		entry := logger.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(2)
		entry.Debug(l)
	}
}

// Info
func Info(args ...interface{}) {
	if logger.Level >= logrus.InfoLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Info(args...)
	}
}

// Info
func Infof(format string, args ...interface{}) {
	arg := fmt.Sprintf(format, args...)
	if logger.Level >= logrus.InfoLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Info(arg)
	}
}

// 带有field的Info
func InfoWithFields(l interface{}, f Fields) {
	if logger.Level >= logrus.InfoLevel {
		entry := logger.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(2)
		entry.Info(l)
	}
}

// Warn
func Warn(args ...interface{}) {
	if logger.Level >= logrus.WarnLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Warn(args...)
	}
}
func Warnf(format string, args ...interface{}) {
	arg := fmt.Sprintf(format, args...)
	if logger.Level >= logrus.WarnLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Warn(arg)
	}
}

// 带有Field的Warn
func WarnWithFields(l interface{}, f Fields) {
	if logger.Level >= logrus.WarnLevel {
		entry := logger.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(2)
		entry.Warn(l)
	}
}

// Error
func Error(args ...interface{}) {
	if logger.Level >= logrus.ErrorLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Error(args...)
	}
}

// Error
func Errorf(format string, args ...interface{}) {
	arg := fmt.Sprintf(format, args...)
	if logger.Level >= logrus.ErrorLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Error(arg)
	}
}

// 带有Fields的Error
func ErrorWithFields(l interface{}, f Fields) {
	if logger.Level >= logrus.ErrorLevel {
		entry := logger.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(2)
		entry.Error(l)
	}
}

// Fatal
func Fatal(args ...interface{}) {
	if logger.Level >= logrus.FatalLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Fatal(args...)
	}
}

// 带有Field的Fatal
func FatalWithFields(l interface{}, f Fields) {
	if logger.Level >= logrus.FatalLevel {
		entry := logger.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(2)
		entry.Fatal(l)
	}
}

// Panic
func Panic(args ...interface{}) {
	if logger.Level >= logrus.PanicLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Panic(args...)
	}
}

// 带有Field的Panic
func PanicWithFields(l interface{}, f Fields) {
	if logger.Level >= logrus.PanicLevel {
		entry := logger.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(2)
		entry.Panic(l)
	}
}

func fileInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}

func NewWriter() io.Writer {
	return logger.Writer()
}
