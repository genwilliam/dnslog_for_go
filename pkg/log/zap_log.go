package log

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Zap  = zap.NewNop()
	once sync.Once
)

func getProjectDir() string {
	// 1. 获取可执行文件路径
	exePath, _ := os.Executable()
	exeDir := filepath.Dir(exePath)

	// 2. 获取当前工作目录
	wd, _ := os.Getwd()

	// 3. 检测是否是 go run 产生的临时目录
	// /var/.../T/go-buildxxxx/exe/main
	if strings.Contains(exeDir, os.TempDir()) {
		// go run：使用工作目录（项目根目录）
		return wd
	}

	// go build：使用可执行文件目录
	return exeDir
}

func InitZapLogger() {
	once.Do(func() {

		projectDir := getProjectDir()
		logDir := filepath.Join(projectDir, "logs")
		os.MkdirAll(logDir, 0755)

		logFile := filepath.Join(logDir, time.Now().Format("2006-01-02")+".log")

		encoderConfig := zapcore.EncoderConfig{
			TimeKey:      "time",
			LevelKey:     "level",
			CallerKey:    "caller",
			MessageKey:   "msg",
			EncodeLevel:  zapcore.CapitalLevelEncoder,
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		}

		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		fileEncoder := zapcore.NewJSONEncoder(encoderConfig)

		fileWriter, _ := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

		core := zapcore.NewTee(
			zapcore.NewCore(fileEncoder, zapcore.AddSync(fileWriter), zapcore.InfoLevel),
			zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
		)

		Zap = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	})
}

func Sync() {
	if Zap != nil {
		_ = Zap.Sync()
	}
}

func Info(msg string, fields ...zap.Field)  { Zap.Info(msg, fields...) }
func Warn(msg string, fields ...zap.Field)  { Zap.Warn(msg, fields...) }
func Error(msg string, fields ...zap.Field) { Zap.Error(msg, fields...) }
func Debug(msg string, fields ...zap.Field) { Zap.Debug(msg, fields...) }
func Fatal(msg string, fields ...zap.Field) { Zap.Fatal(msg, fields...) }
