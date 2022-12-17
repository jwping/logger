package logger

import "golang.org/x/exp/slog"

// func TestNewLogger(t *testing.T) {
// 	// fmt.Printf("试试测试\n")
// 	logger := NewLogger(JSON, slog.HandlerOptions{})
// 	logger.Info("测试输出")
// }

func Example() {
	handlerPlus := NewHandler(JSON, slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})
	handlerPlus.AddOutFileForLevel(slog.LevelInfo, "./output1.log", "./output2.log")
	logger := slog.New(handlerPlus)

	logger.Info("Example success!")
	// Output:
	// {"time":"xxxxxxx","level":"INFO","source":"/home/jwping/gopath/logger/logger_test.go:19","msg":"Example success!"}
}
