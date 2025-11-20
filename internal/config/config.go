package config

type (
	// LogLevel represents the severity or verbosity of the log records.
	LogLevel string
)

const (
	// LogLevelDebug
	LogLevelDebug LogLevel = "debug"

	// LogLevelInfo
	LogLevelInfo LogLevel = "info"

	// LogLevelWarn
	LogLevelWarn LogLevel = "warn"

	// LogLevelError
	LogLevelError LogLevel = "error"
)

type (
	// LogFormat represents the encoding style of the log records.
	LogFormat string
)

const (
	// LogFormatText renders log records as human-readable plain text.
	LogFormatText LogFormat = "text"

	// LogFormatJSON renders log records as structured JSON objects.
	LogFormatJSON LogFormat = "json"
)

type (
	// LogOutput represents the destination stream for log records.
	LogOutput string
)

const (
	// LogOutputStdout writes log records to the standard output stream.
	LogOutputStdout LogOutput = "stdout"

	// LogOutputStderr writes log records to the standard error stream.
	LogOutputStderr LogOutput = "stderr"
)
