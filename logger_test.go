package logger

import (
	"golang.org/x/exp/slog"
)

func ExampleHandler() {
	handlerPlus := NewHandler(JSON, slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})
	err := handlerPlus.AddOutFileForLevel(slog.LevelInfo, "./output1.log", "./output2.log")
	if err != nil {
		return
	}

	log := slog.New(handlerPlus)

	log.Info("ExampleHandler success!")
	// Output:
	// {"time":"xxxxxx","level":"INFO","source":"/home/jwping/gopath/logger/logger_test.go:19","msg":"Example success!"}
}

func ExampleLogger() {
	log := NewLogger(Options{
		Lt:        JSON,
		Level:     LevelDebug,
		AddSource: true,
	})

	err := log.AddOutFileForLevel(slog.LevelInfo, "./output3.log", "./output4.log")
	if err != nil {
		log.Error("AddOutFileForLevel faild", err)
		return
	}

	log.Info("ExampleLogger success!")
	// Output:
	// {"time":"xxxxxx","level":"INFO","source":"/home/jwping/gopath/logger/logger_test.go:36","msg":"ExampleLogger success!"}
}
