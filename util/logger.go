package util

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
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
		fmt.Fprintf(b, "[%s] [%s] %s %s", timestamp, entry.Level, fileVal, entry.Message)
	} else {
		fmt.Fprintf(b, "[%s] [%s] %s", timestamp, entry.Level, entry.Message)
	}
	// 输出额外字段
	if len(entry.Data) > 0 {
		t.formatEntries(b, entry.Data)
	}
	b.WriteByte('\n')
	return b.Bytes(), nil
}

type entry struct {
	k string
	v any
}

// formatEntries 格式化并排序日志字段，保证稳定输出顺序
func (t *LogFormatter) formatEntries(w io.Writer, data logrus.Fields) {
	es := make([]entry, 0, len(data))
	for k, v := range data {
		es = append(es, entry{k: k, v: v})
	}
	sort.Slice(es, func(i, j int) bool { return strings.Compare(es[i].k, es[j].k) < 0 })
	for _, e := range es {
		fmt.Fprintf(w, " -%s[%v]", e.k, e.v)
	}
}

// DailyRotateWriter 按日期自动轮换的文件写入器
type DailyRotateWriter struct {
	logDir          string
	levelPrefix     string
	currentDate     string
	currentFilePath string
	file            *os.File
	mu              sync.Mutex
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

	// 检查日志文件
	if err := w.checkFile(); err != nil {
		return 0, err
	}

	return w.file.Write(p)
}

func (w *DailyRotateWriter) checkFile() error {
	currentDate := time.Now().Format("20060102_15")

	// 检查日期是否变化
	if currentDate != w.currentDate {
		return w.rotate()
	}

	// 检查当前文件是否还存在
	opened, err := w.file.Stat()
	if err != nil {
		return err
	}
	stat, err := os.Stat(w.currentFilePath)
	if err != nil || !os.SameFile(stat, opened) {
		return w.rotate()
	}

	return nil
}

// rotate 轮换日志文件
func (w *DailyRotateWriter) rotate() error {
	currentDate := time.Now().Format("20060102_15")
	filename := fmt.Sprintf("%s_%s.log", w.levelPrefix, currentDate)
	filePath := filepath.Join(w.logDir, filename)

	// 关闭旧文件
	if w.file != nil {
		if err := w.file.Close(); err != nil {
			return fmt.Errorf("关闭旧日志文件失败: %w", err)
		}
	}

	// 打开新文件
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("打开日志文件失败: %w", err)
	}

	w.file = file
	w.currentDate = currentDate
	w.currentFilePath = filePath
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
	errorWriter, err := NewDailyRotateWriter(logDir, "error")
	if err != nil {
		return nil, fmt.Errorf("创建 error 级别的日志写入器失败: %w", err)
	}
	hook.writers[logrus.PanicLevel] = errorWriter
	hook.writers[logrus.FatalLevel] = errorWriter
	hook.writers[logrus.ErrorLevel] = errorWriter

	warnWriter, err := NewDailyRotateWriter(logDir, "warn")
	if err != nil {
		return nil, fmt.Errorf("创建 warn 级别的日志写入器失败: %w", err)
	}
	hook.writers[logrus.WarnLevel] = warnWriter

	infoWriter, err := NewDailyRotateWriter(logDir, "info")
	if err != nil {
		return nil, fmt.Errorf("创建 info 级别的日志写入器失败: %w", err)
	}
	hook.writers[logrus.InfoLevel] = infoWriter

	debugWriter, err := NewDailyRotateWriter(logDir, "debug")
	if err != nil {
		return nil, fmt.Errorf("创建 debug 级别的日志写入器失败: %w", err)
	}
	hook.writers[logrus.DebugLevel] = debugWriter

	traceWriter, err := NewDailyRotateWriter(logDir, "trace")
	if err != nil {
		return nil, fmt.Errorf("创建 trace 级别的日志写入器失败: %w", err)
	}
	hook.writers[logrus.TraceLevel] = traceWriter

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
	case "trace":
		return logrus.TraceLevel
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

func FiberLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		method := c.Method()
		path := c.Path()
		status := c.Response().StatusCode()

		err := c.Next()
		latency := float64(time.Since(start).Nanoseconds()) / 1000000.0

		logrus.WithField("body", string(c.Response().Body())).Tracef("%d %8.2fms %s %s", status, latency, method, path)
		return err
	}
}
