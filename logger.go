package logger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"

	"golang.org/x/exp/slog"
)

type LogType int

const (
	TEXT LogType = 1 << iota
	JSON
)

type HandlerPlus struct {
	slog.Handler

	buffer  *bytes.Buffer
	outList map[slog.Level][]io.Writer
	ch      chan interface{}
}

func (h *HandlerPlus) Handle(r slog.Record) error {
	<-h.ch

	var stdWriter io.Writer
	switch r.Level {
	case slog.LevelError:
		stdWriter = os.Stderr
	default:
		stdWriter = os.Stdout
	}

	outList := make([]io.Writer, len(h.outList[r.Level]))
	copy(outList, h.outList[r.Level])
	outList = append(outList, stdWriter)
	ioWriter := io.MultiWriter(outList...)

	err := h.Handler.Handle(r)
	if err != nil {
		// l.Handler.Handle(slog.Record{
		// 	Time: time.Now(),
		// 	Level: slog.LevelError,
		// 	Message: "l.Handler.Handle err: " + err.Error(),
		// })
		return err
	}
	ioWriter.Write(h.buffer.Bytes())
	h.buffer.Reset()
	h.ch <- nil

	return nil
}

func (h *HandlerPlus) WithAttrs(attr []slog.Attr) slog.Handler {
	return &HandlerPlus{
		Handler: h.Handler.WithAttrs(attr),
		buffer:  bytes.NewBuffer(make([]byte, 10)),
		outList: make(map[slog.Level][]io.Writer),
		ch:      h.ch,
	}
}

func (h *HandlerPlus) WithGroup(name string) slog.Handler {
	return &HandlerPlus{
		Handler: h.Handler.WithGroup(name),
		buffer:  bytes.NewBuffer(make([]byte, 10)),
		outList: make(map[slog.Level][]io.Writer),
		ch:      h.ch,
	}
}

func (h *HandlerPlus) AddOutFileForLevel(level slog.Level, files ...string) error {
	for _, file := range files {
		fd, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			// l.Handler.Handle(slog.Record{
			// 	Time:    time.Now(),
			// 	Level:   slog.LevelError,
			// 	Message: "os openfile faild:" + err.Error(),
			// })
			return err
		}

		h.outList[level] = append(h.outList[level], fd)
	}

	return nil
}

func NewHandler(lt LogType, opts slog.HandlerOptions) *HandlerPlus {

	handlerPlus := &HandlerPlus{
		buffer:  bytes.NewBuffer(make([]byte, 10)),
		outList: make(map[slog.Level][]io.Writer),
		ch:      make(chan interface{}, 1),
	}

	var handler slog.Handler
	switch lt {
	case TEXT:
		handler = opts.NewTextHandler(handlerPlus.buffer)
	case JSON:
		handler = opts.NewJSONHandler(handlerPlus.buffer)
	}

	handlerPlus.Handler = handler

	handlerPlus.ch <- nil
	return handlerPlus
}

type Level int

const (
	LevelDebug Level = -4
	LevelInfo  Level = 0
	LevelWarn  Level = 4
	LevelError Level = 8
)

type Options struct {
	Lt LogType
	// When AddSource is true, the handler adds a ("source", "file:line")
	// attribute to the output indicating the source code position of the log
	// statement. AddSource is false by default to skip the cost of computing
	// this information.
	AddSource bool

	// Level reports the minimum record level that will be logged.
	// The handler discards records with lower levels.
	// If Level is nil, the handler assumes LevelInfo.
	// The handler calls Level.Level for each record processed;
	// to adjust the minimum level dynamically, use a LevelVar.
	Level Level
}

type Logger struct {
	slog.Logger
}

func NewLogger(opts Options) *Logger {
	return &Logger{
		Logger: *slog.New(NewHandler(opts.Lt, slog.HandlerOptions{
			Level:     slog.Level(opts.Level),
			AddSource: opts.AddSource,
		})),
	}
}

func (l *Logger) AddOutFileForLevel(level slog.Level, files ...string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	rv := reflect.ValueOf(l.Logger.Handler())
	rverr := rv.MethodByName("AddOutFileForLevel").CallSlice([]reflect.Value{reflect.ValueOf(slog.LevelInfo), reflect.ValueOf(files)})
	err, ok := rverr[0].Interface().(error)
	if ok {
		return
	}
	err = nil
	return
}
