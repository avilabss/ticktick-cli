package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync"
	"time"
)

// LevelTrace is a custom level below Debug for -vvv output.
const LevelTrace = slog.Level(-8)

type handler struct {
	level slog.Level
	mu    *sync.Mutex
	out   io.Writer
}

func (h *handler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *handler) Handle(_ context.Context, r slog.Record) error {
	levelStr := r.Level.String()
	if r.Level == LevelTrace {
		levelStr = "TRACE"
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	timestamp := r.Time.Format(time.TimeOnly)
	fmt.Fprintf(h.out, "%s [%-5s] %s", timestamp, levelStr, r.Message)

	r.Attrs(func(a slog.Attr) bool {
		fmt.Fprintf(h.out, " %s=%v", a.Key, a.Value)
		return true
	})

	fmt.Fprintln(h.out)
	return nil
}

func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *handler) WithGroup(name string) slog.Handler {
	return h
}

// SetVerbosity configures slog based on the verbosity level.
//
//	0: errors/warnings only (slog.LevelWarn)
//	1 (-v): info level (slog.LevelInfo)
//	2 (-vv): debug level (slog.LevelDebug)
//	3 (-vvv): trace level (LevelTrace)
func SetVerbosity(level int) {
	slogLevel := slog.LevelWarn
	switch {
	case level >= 3:
		slogLevel = LevelTrace
	case level >= 2:
		slogLevel = slog.LevelDebug
	case level >= 1:
		slogLevel = slog.LevelInfo
	}

	h := &handler{
		level: slogLevel,
		mu:    &sync.Mutex{},
		out:   os.Stderr,
	}

	slog.SetDefault(slog.New(h))
}

// Trace logs at trace level (-vvv). Raw data and internals.
func Trace(msg string, args ...any) {
	slog.Log(context.TODO(), LevelTrace, msg, args...)
}
