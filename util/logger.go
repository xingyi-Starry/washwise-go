package util

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	defaultLogDir   = "logs"
	defaultLogLevel = logrus.InfoLevel
)

type LogFormatter struct{}

// Format 实现Formatter(entry *logrus.Entry) ([]byte, error)接口
func (t *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	//自定义日期格式
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	if entry.HasCaller() {
		//自定义文件路径
		fileVal := fmt.Sprintf("%s:%d", entry.Caller.File, entry.Caller.Line)
		//自定义输出格式
		fmt.Fprintf(b, "[%s] [%s] %s %s\n", timestamp, entry.Level, fileVal, entry.Message)
	} else {
		fmt.Fprintf(b, "[%s] [%s] %s\n", timestamp, entry.Level, entry.Message)
	}
	// 输出额外字段
	for k, v := range entry.Data {
		fmt.Fprintf(b, " -%s[%v]", k, v)
	}
	return b.Bytes(), nil
}

// DailyRotateWriter 按日期自动轮换的文件写入器
type DailyRotateWriter struct {
	logDir      string
	levelPrefix string
	currentDate string
	file        *os.File
	mu          sync.Mutex
}

// NewDailyRotateWriter 创建一个新的按日期轮换的写入器
func NewDailyRotateWriter(logDir, levelPrefix string) (*DailyRotateWriter, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("创建日志目录失败: %w", err)
	}

	writer := &DailyRotateWriter{
		logDir:      logDir,
		levelPrefix: levelPrefix,
	}

	if err := writer.rotate(); err != nil {
		return nil, err
	}

	return writer, nil
}

// Write 实现 io.Writer 接口
func (w *DailyRotateWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	// 检查是否需要轮换日志文件
	currentDate := time.Now().Format("20060102_15")
	if currentDate != w.currentDate {
		if err := w.rotate(); err != nil {
			return 0, err
		}
	}

	return w.file.Write(p)
}

// rotate 轮换日志文件
func (w *DailyRotateWriter) rotate() error {
	currentDate := time.Now().Format("20060102_15")
	filename := fmt.Sprintf("%s_%s.log", w.levelPrefix, currentDate)
	filepath := filepath.Join(w.logDir, filename)

	// 关闭旧文件
	if w.file != nil {
		if err := w.file.Close(); err != nil {
			return fmt.Errorf("关闭旧日志文件失败: %w", err)
		}
	}

	// 打开新文件
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("打开日志文件失败: %w", err)
	}

	w.file = file
	w.currentDate = currentDate
	return nil
}

// Close 关闭文件
func (w *DailyRotateWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.file != nil {
		return w.file.Close()
	}
	return nil
}

// LevelFileHook 不同级别的日志写入不同文件的Hook
type LevelFileHook struct {
	writers map[logrus.Level]*DailyRotateWriter
	logDir  string
}

// NewLevelFileHook 创建一个新的级别文件Hook
func NewLevelFileHook(logDir string) (*LevelFileHook, error) {
	if logDir == "" {
		logDir = defaultLogDir
	}

	hook := &LevelFileHook{
		writers: make(map[logrus.Level]*DailyRotateWriter),
		logDir:  logDir,
	}

	// 为每个日志级别创建写入器
	levels := map[logrus.Level]string{
		// 合并 panic, fatal, error 到 error 级别
		logrus.PanicLevel: "error",
		logrus.FatalLevel: "error",
		logrus.ErrorLevel: "error",

		logrus.WarnLevel: "warn",
		logrus.InfoLevel: "info",

		// 合并 debug, trace 到 debug 级别
		logrus.DebugLevel: "debug",
		logrus.TraceLevel: "debug",
	}

	for level, prefix := range levels {
		writer, err := NewDailyRotateWriter(logDir, prefix)
		if err != nil {
			return nil, fmt.Errorf("创建 %s 级别的日志写入器失败: %w", prefix, err)
		}
		hook.writers[level] = writer
	}

	return hook, nil
}

// Levels 返回这个Hook处理的日志级别
func (hook *LevelFileHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire 当日志事件发生时触发
func (hook *LevelFileHook) Fire(entry *logrus.Entry) error {
	writer, ok := hook.writers[entry.Level]
	if !ok {
		return nil
	}

	// 格式化日志
	line, err := entry.Logger.Formatter.Format(entry)
	if err != nil {
		return fmt.Errorf("格式化日志失败: %w", err)
	}

	// 写入对应级别的文件
	_, err = writer.Write(line)
	return err
}

// Close 关闭所有文件
func (hook *LevelFileHook) Close() error {
	for _, writer := range hook.writers {
		if err := writer.Close(); err != nil {
			return err
		}
	}
	return nil
}

// LogConfig 日志配置
type LogConfig struct {
	Level string // 日志等级
	Dir   string // 日志目录
}

// ParseLogLevel 解析日志等级字符串
func ParseLogLevel(level string) logrus.Level {
	switch strings.ToLower(level) {
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn", "warning":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal":
		return logrus.FatalLevel
	case "panic":
		return logrus.PanicLevel
	default:
		return defaultLogLevel
	}
}

// InitLogger 使用配置初始化日志系统
func InitLogger(config LogConfig) {
	// 设置基本配置
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&LogFormatter{})

	// 设置日志等级
	level := ParseLogLevel(config.Level)
	logrus.SetLevel(level)

	// 设置日志目录
	logDir := config.Dir
	if logDir == "" {
		logDir = defaultLogDir
	}

	// 创建文件Hook
	fileHook, err := NewLevelFileHook(logDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "创建文件Hook失败: %v\n", err)
		return
	}

	// 添加Hook
	logrus.AddHook(fileHook)

	// 屏蔽标准输出
	logrus.SetOutput(io.Discard)
}
