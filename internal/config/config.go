package config

type (
	// LogLevel represents the severity or verbosity of the log records.
	LogLevel string
)

const (
	// LogLevelDebug captures detailed information, typically useful for development
	// and debugging purposes.
	LogLevelDebug LogLevel = "debug"

	// LogLevelInfo captures general information about the application's operation.
	LogLevelInfo LogLevel = "info"

	// LogLevelWarn captures non-critical events or potentially harmful situations.
	LogLevelWarn LogLevel = "warn"

	// LogLevelError captures critical events or errors that require immediate
	// attention.
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
	// LogOutputStdout writes log records to the standard output stream (stdout).
	LogOutputStdout LogOutput = "stdout"

	// LogOutputStderr writes log records to the standard error stream (stderr).
	LogOutputStderr LogOutput = "stderr"
)
