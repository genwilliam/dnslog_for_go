package example

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Zap *zap.Logger // 声明全局变量，封装info等方法

func exampleDemo() { // 通过zap.NewExample()创建一个logger，用于演示
	logger := zap.NewExample()
	defer logger.Sync()
	logger.Info("hello, example")
}
func productionDemo() { // 通过zap.NewProduction()创建一个logger，用于生产环境
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("hello, production")
}
func developmentDemo() { // 通过zap.NewDevelopment()创建一个logger，用于开发环境
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	logger.Info("hello, development")
}

// 你也可以通过config来控制logger的行为
func configDemo() { // 通过zap.Config创建一个logger，用于配置
	encoderConfig := zapcore.EncoderConfig{
		// 在这里配置编码器
		TimeKey:        "time",                           // 时间键
		LevelKey:       "level",                          // 级别键
		NameKey:        "logger",                         // 名称键
		CallerKey:      "caller",                         // 调用者键
		MessageKey:     "msg",                            // 消息键
		StacktraceKey:  "stacktrace",                     // 堆栈跟踪键
		LineEnding:     zapcore.DefaultLineEnding,        // 行结束符
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, // 级别编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,       // 时间编码器
		EncodeDuration: zapcore.StringDurationEncoder,    // 持续时间编码器
		EncodeCaller:   zapcore.ShortCallerEncoder,       // 调用者编码器
		FunctionKey:    "function",                       // 函数键
		// ...
	}
	// 后续可以使用encoderConfig来创建一个编码器
	// encoder := zapcore.NewJSONEncoder(encoderConfig)
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig), // 使用控制台编码器
		zapcore.AddSync(os.Stdout),               // 输出到标准输出
		zapcore.DebugLevel,                       // 日志级别
	)
	zapp := zap.New(core, zap.AddCaller()) // 创建logger,这样你就可以控制logger的行为了
	// 用法
	zapp.Info("hello, config")

}

// Info 你也可以封装info,debug等方法,方便后续的调用
func Info(msg string, fields ...zap.Field) {
	if Zap != nil {
		Zap.Info(msg, fields...)
	}
}
func Debug(msg string, fields ...zap.Field) {
	if Zap != nil {
		Zap.Debug(msg, fields...)
	}
}

// ……
