package gitee_utils

import (
	"io"
	"os"
	"path"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

var (
	logPath = "src/log"
	logFile = "strategy.log"
)
var LogInstance = logrus.New()

// 日志初始化
func init() {
	// 打开文件
	logFileName := path.Join(logPath, logFile)

	fileWriter, err := os.OpenFile(logFileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}

	mw := io.MultiWriter(os.Stdout, fileWriter)

	// 使用滚动压缩方式记录日志
	rolling(logFileName)
	// 设置日志输出JSON格式
	LogInstance.SetFormatter(&logrus.TextFormatter{})
	// 设置日志记录级别
	LogInstance.SetLevel(logrus.DebugLevel)

	LogInstance.SetOutput(mw)
}

// 日志滚动设置
func rolling(logFile string) {
	// 设置输出
	LogInstance.SetOutput(&lumberjack.Logger{
		Filename:   logFile, //日志文件位置
		MaxSize:    50,      // 单文件最大容量,单位是MB
		MaxBackups: 3,       // 最大保留过期文件个数
		MaxAge:     5,       // 保留过期文件的最大时间间隔,单位是天
		Compress:   true,    // 是否需要压缩滚动日志, 使用的 gzip 压缩
	})
}
