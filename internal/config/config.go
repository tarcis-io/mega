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
	LogFormat string
)

const (
	LogFormatText LogFormat = "text"
	LogFormatJSON LogFormat = "json"
)

type (
	LogOutput string
)

const (
	LogOutputStdout LogOutput = "stdout"
	LogOutputStderr LogOutput = "stderr"
)
