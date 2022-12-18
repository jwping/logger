package logger

import (
	"golang.org/x/exp/slog"
)

func ExampleHandler() {
	handlerPlus := NewHandler(JSON, slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})
	err := handlerPlus.AddOutFileForLevel(LevelInfo, "./output1.log", "./output2.log")
	if err != nil {
		return
	}

	log := slog.New(handlerPlus)

	log.Info("ExampleHandler success!")
	// Output:
	// {"time":"xxxxxx","level":"INFO","source":"/home/jwping/gopath/logger/logger_test.go:19","msg":"Example success!"}
}

// recommend
func ExampleLogger() {
	log := NewLogger(Options{
		Lt:        JSON,
		Level:     LevelDebug,
		AddSource: true,
	})

	err := log.AddOutFileForLevel(LevelDebug, "./output3.log", "./output4.log")
	if err != nil {
		log.Error("AddOutFileForLevel faild", err)
		return
	}

	log.Info("ExampleLogger success!")

	// If you need a new logger, otherwise use the example above
	log, err = log.WithAttrs([]slog.Attr{
		{
			Key:   "TestKey1",
			Value: slog.AnyValue("TestValue1"),
		},
		{
			Key:   "TestKey2",
			Value: slog.AnyValue("TestValue2"),
		},
	})
	if err != nil {
		log.Error("WithAttrs faild", err)
		return
	}
	log.Debug("ExampleLogger WithAttrs Output!")

	// Output:
	// {"time":"xxxxxx","level":"INFO","source":"/home/jwping/gopath/logger/logger_test.go:38","msg":"ExampleLogger success!"}
	// {"time":"xxxxxx","level":"DEBUG","source":"/home/jwping/gopath/logger/logger_test.go:54","msg":"ExampleLogger WithAttrs Output!","TestKey1":"TestValue1","TestKey2":"TestValue2"}
}
