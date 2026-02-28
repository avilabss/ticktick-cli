package logger

import (
	"log/slog"
	"os"
)

// LevelTrace is a custom level below Debug for -vvv output.
const LevelTrace = slog.Level(-8)

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

	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slogLevel,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey && a.Value.String() == "DEBUG-4" {
				a.Value = slog.StringValue("TRACE")
			}
			return a
		},
	})

	slog.SetDefault(slog.New(handler))
}

// Trace logs at trace level (-vvv). Raw data and internals.
func Trace(msg string, args ...any) {
	slog.Log(nil, LevelTrace, msg, args...)
}
