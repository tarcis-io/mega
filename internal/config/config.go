// Package config loads and provides the application configuration.
package config

type (
	// LogLevel represents the severity level of log messages.
	LogLevel string
)

const (
	// LogLevelDebug defines the debug log level.
	// Debug messages are typically used for detailed information during
	// development.
	LogLevelDebug LogLevel = "debug"

	// LogLevelInfo defines the info log level.
	// Info messages are used to provide general information about the
	// application's execution.
	LogLevelInfo LogLevel = "info"

	// LogLevelWarn defines the warning log level.
	// Warning messages are used to indicate potential issues or unexpected
	// situations.
	LogLevelWarn LogLevel = "warn"

	// LogLevelError defines the error log level.
	// Error messages are used to report errors or exceptional conditions.
	LogLevelError LogLevel = "error"
)

type (
	// LogFormat represents the output format of log messages.
	LogFormat string
)

const (
	// LogFormatText defines the text-based log format.
	// Text-based formats are human-readable and typically used for console
	// output.
	LogFormatText LogFormat = "text"

	// LogFormatJSON defines the JSON-based log format.
	// JSON-based formats are machine-readable and typically used for
	// structured logging.
	LogFormatJSON LogFormat = "json"
)

type (
	// LogOutput represents the destination for log messages.
	LogOutput string
)

const (
	// LogOutputStdout defines the standard log output stream (stdout).
	// Stdout is typically used for general, non-error messages.
	LogOutputStdout LogOutput = "stdout"

	// LogOutputStderr defines the error log output stream (stderr).
	// Stderr is typically used for error messages.
	LogOutputStderr LogOutput = "stderr"
)

const (
	EnvLogLevel  = "LOG_LEVEL"
	EnvLogFormat = "LOG_FORMAT"
	EnvLogOutput = "LOG_OUTPUT"
)

const (
	DefaultLogLevel  = LogLevelInfo
	DefaultLogFormat = LogFormatText
	DefaultLogOutput = LogOutputStdout
)
