package logger

import (
	"bytes"
	"io"
	"os"

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
