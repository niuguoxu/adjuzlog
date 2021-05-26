package adjuzLog

import (
	"fmt"
	"runtime"
	log "github.com/sirupsen/logrus"
	"time"
	"path"
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"os"
)

// ConfigLocalFilesystemLogger set logger
func ConfigLocalFilesystemLogger(logPath string, logFileName string, maxAge time.Duration, rotationTime time.Duration) {
	hostName, _ := os.Hostname()
	baseLogPaht := path.Join(logPath, logFileName)
	// fileName := "%Y%m%d_"+fmt.Sprintf("%s_", logFileName)+hostName+".log"
	// baseLogPaht := path.Join(logPath, fileName)
	writer, err := rotatelogs.New(
		baseLogPaht+"_"+"%Y%m%d_"+hostName+".log",
		rotatelogs.WithLinkName(baseLogPaht), // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(maxAge), // 文件最大保存时间
		rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
	)
	errWriter, err := rotatelogs.New(
		baseLogPaht+"_"+"%Y%m%d_"+hostName+"_err.log",
		rotatelogs.WithLinkName(baseLogPaht), // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(maxAge), // 文件最大保存时间
		rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
	)
	if err != nil {
		log.Errorf("config local file system logger error. %+v", errors.WithStack(err))
	}
	lfHook := lfshook.NewHook(lfshook.WriterMap{
		// log.DebugLevel: writer, // 为不同级别设置不同的输出目的
		log.InfoLevel:  writer,
		log.WarnLevel:  errWriter,
		log.ErrorLevel: errWriter,
		log.FatalLevel: errWriter,
		log.PanicLevel: errWriter,
	}, &log.JSONFormatter{TimestampFormat:"2006-01-02 15:04.999"})

	cxHook := ContextHook{}
	log.AddHook(cxHook)
	log.AddHook(lfHook)
	log.SetFormatter(&log.JSONFormatter{TimestampFormat:"2006-01-02 15:04.999"})
}

// ContextHook ...
type ContextHook struct{}

// Levels ...
func (hook ContextHook) Levels() []log.Level {
	return log.AllLevels
}

// Fire ...
func (hook ContextHook) Fire(entry *log.Entry) error {
	if pc, file, line, ok := runtime.Caller(10); ok {
		funcName := runtime.FuncForPC(pc).Name()
		entry.Data["source"] = fmt.Sprintf("%s:%v:%s", path.Base(file), line, path.Base(funcName))
	}

	return nil
}
