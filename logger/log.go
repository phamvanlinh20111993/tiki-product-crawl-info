package logger

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"runtime"
	"selfstudy/crawl/product/configuration"
	"selfstudy/crawl/product/util"
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"
)

type PrettyHandlerOptions struct {
	SlogOpts slog.HandlerOptions
}

const (
	LevelTrace = slog.Level(-8)
	LevelFatal = slog.Level(12)
	TargetFile = "FILE"
	TargetCmd  = "STDOUT"
)

var LevelNames = map[slog.Leveler]string{
	LevelTrace:      "TRACE", // -8
	slog.LevelDebug: "DEBUG", // -4
	slog.LevelInfo:  "INFO",  // 0
	slog.LevelWarn:  "WARN",  // 4
	slog.LevelError: "ERROR", //8
	LevelFatal:      "FATAL", // 12
}

var Levels = map[string]slog.Level{
	"TRACE": LevelTrace,
	"FATAL": LevelFatal,
	"DEBUG": slog.LevelDebug,
	"INFO":  slog.LevelInfo,
	"ERROR": slog.LevelError,
	"WARN":  slog.LevelWarn,
}

func slogOption() *slog.HandlerOptions {
	var loggerConfig = configuration.GetLoggerConfig()
	return &slog.HandlerOptions{
		Level:     Levels[strings.ToUpper(loggerConfig.Level)],
		AddSource: loggerConfig.IsAddSource,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				a.Value = slog.StringValue(
					util.TimeToString(a.Value.Time(), util.Format_yyyy_mm_dd_space_hh_dot_mm_dot_ss),
				)
			}
			if a.Key == slog.LevelKey {
				level := a.Value.Any().(slog.Level)
				levelLabel, exists := LevelNames[level]
				if !exists {
					levelLabel = level.String()
				}

				a.Value = slog.StringValue(levelLabel)
			}
			if a.Key == slog.MessageKey {
				a.Key = "message"
			}
			return a
		},
	}
}

var (
	logger *slog.Logger
	once   sync.Once
)

func getLoggerStdoutInstance() *slog.Logger {
	once.Do(func() {
		//opts := PrettyHandlerOptions{
		//	SlogOpts: *slogOpts,
		//}
		//handler := NewPrettyHandler(os.Stdout, opts)
		handler := slog.NewTextHandler(os.Stdout, slogOption())
		logger = slog.New(handler)
	})
	return logger
}

var (
	loggerPathFolder               = util.GetPathSeparator() + "log" + util.GetPathSeparator()
	logFileName             string = ""
	DefaultLogFileExtension        = ".log"
	logFile                 *os.File
	loggerFile              *slog.Logger
	lock                    sync.Mutex
)

func removeOldLogs(path string) {
	entries, err := os.ReadDir(path)
	if err != nil {
		log.Println(err)
	}

	type LogFile struct {
		logFileName string
		cTime       time.Time
	}

	var logFiles []LogFile

	var logFileConfig = configuration.GetLoggerConfig()
	for _, e := range entries {
		if !e.IsDir() {
			fileInfo, err := e.Info()
			if err != nil {
				log.Println(err)
				continue
			}

			if !strings.HasPrefix(fileInfo.Name(), logFileConfig.FilePrefixName) {
				continue
			}

			var cTime time.Time = fileInfo.ModTime()
			if runtime.GOOS == "windows" {
				d := fileInfo.Sys().(*syscall.Win32FileAttributeData)
				cTime = time.Unix(0, d.CreationTime.Nanoseconds())
			}

			logFiles = append(logFiles, LogFile{path + fileInfo.Name(), cTime})
		}
	}

	if logFileConfig.KeepLogDays > len(logFiles) {
		return
	}

	sort.Slice(logFiles, func(i, j int) bool {
		return logFiles[j].cTime.After(logFiles[i].cTime)
	})

	for i := 0; i < len(logFiles)-logFileConfig.KeepLogDays; i++ {
		// Using Remove() function
		e := os.Remove(logFiles[i].logFileName)
		if e != nil {
			log.Println(e)
		}
	}

}

// TODO need to improve, too many logic => slow logging
func getLoggerFileInstance() *slog.Logger {
	var logFileConfig = configuration.GetLoggerConfig()
	var filePattern string = logFileConfig.FilePattern
	if logFileConfig.FilePattern == "" {
		filePattern = util.Format_yyyy_mm_dd
	}

	var logFilePath string = logFileConfig.FilePath
	if logFilePath == "" || logFilePath == strings.TrimSpace(".") {
		logFilePath = util.GetCurrentFolder()
	}
	// C:\Users\Lenovo\AppData\Local\JetBrains\IntelliJIdea2025.3\tmp\GoLand\log
	var path string = logFilePath + loggerPathFolder
	var currLogFileName string = logFileConfig.FilePrefixName + util.CurrentTimeToString(filePattern) + DefaultLogFileExtension

	if currLogFileName == logFileName {
		logFileName = currLogFileName
		return loggerFile
	}

	lock.Lock()
	removeOldLogs(logFilePath)

	logFileName = currLogFileName
	if logFile != nil {
		if err := logFile.Close(); err != nil {
			log.Println(err)
		}
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, 0644)
		log.Println(err)
	}
	var err error
	logFile, err = os.OpenFile(path+logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}

	//handler := NewPrettyHandler(os.Stdout, opts)
	handler := slog.NewTextHandler(logFile, slogOption())
	loggerFile = slog.New(handler)

	lock.Unlock()

	return loggerFile
}

func logCommon(logLevel slog.Level, msg string, args ...any) {
	var msgFormat strings.Builder
	var attrs []slog.Attr

	for _, v := range args {
		attrValue, isSlogAttr := v.(slog.Attr)
		if isSlogAttr {
			attrs = append(attrs, attrValue)
			continue
		}

		errorValue, isError := v.(error)
		if isError {
			msgFormat.WriteString(" ")
			msgFormat.WriteString(errorValue.Error())
			continue
		}

		strValue, isStr := v.(string)
		if isStr {
			msgFormat.WriteString(" ")
			msgFormat.WriteString(strValue)
			continue
		}

		msgFormat.WriteString(" ")
		msgFormat.WriteString(fmt.Sprintf("%v", v))
	}

	for _, val := range configuration.GetLoggerConfig().Target {
		if strings.ToUpper(val) == TargetCmd {
			getLoggerStdoutInstance().LogAttrs(
				context.Background(),
				logLevel,
				msg+msgFormat.String(),
				attrs...)
		}
	}
	for _, val := range configuration.GetLoggerConfig().Target {
		if strings.ToUpper(val) == TargetFile {
			getLoggerFileInstance().LogAttrs(
				context.Background(),
				logLevel,
				msg+msgFormat.String(),
				attrs...)
		}
	}
}

func logInfo(msg string, args ...any) {
	var logLevel = Levels[strings.ToUpper(configuration.GetLoggerConfig().Level)]
	if logLevel <= slog.LevelInfo {
		logCommon(slog.LevelInfo, msg, args...)
	}
}

func logError(msg string, args ...any) {
	var logLevel = Levels[strings.ToUpper(configuration.GetLoggerConfig().Level)]
	if logLevel <= slog.LevelError {
		logCommon(slog.LevelError, msg, args...)
	}
}

func logDebug(msg string, args ...any) {
	var logLevel = Levels[strings.ToUpper(configuration.GetLoggerConfig().Level)]
	if logLevel <= slog.LevelDebug {
		logCommon(slog.LevelDebug, msg, args...)
	}
}

func logWarn(msg string, args ...any) {
	var logLevel = Levels[strings.ToUpper(configuration.GetLoggerConfig().Level)]
	if logLevel <= slog.LevelWarn {
		logCommon(slog.LevelWarn, msg, args...)
	}
}

var LogInfo = logInfo
var LogError = logError
var LogDebug = logDebug
var LogWarn = logWarn
