package config

type (
	LogLevel string
)

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

type (
	// LogFormat represents the encoding style of the log records.
	LogFormat string
)

const (
	// LogFormatText renders logs as human-readable plain text.
	LogFormatText LogFormat = "text"

	// LogFormatJSON renders logs as structured JSON objects.
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
