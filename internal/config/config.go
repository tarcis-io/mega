// Package config loads and provides the application configuration.
package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

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
	// EnvLogLevel defines the environment variable name for setting the log
	// level.
	EnvLogLevel = "LOG_LEVEL"

	// EnvLogFormat defines the environment variable name for setting the log
	// format.
	EnvLogFormat = "LOG_FORMAT"

	// EnvLogOutput defines the environment variable name for setting the log
	// output.
	EnvLogOutput = "LOG_OUTPUT"
)

const (
	// DefaultLogLevel defines the default log level used if not otherwise
	// specified.
	DefaultLogLevel = LogLevelInfo

	// DefaultLogFormat defines the default log format used if not otherwise
	// specified.
	DefaultLogFormat = LogFormatText

	// DefaultLogOutput defines the default log output used if not otherwise
	// specified.
	DefaultLogOutput = LogOutputStdout
)

type (
	// Config represents the application configuration.
	Config struct {
		// logLevel specifies the severity level of log messages.
		logLevel LogLevel

		// logFormat specifies the output format of log messages.
		logFormat LogFormat

		// logOutput specifies the destination for log messages.
		logOutput LogOutput
	}
)

// New
func New() (*Config, error) {
	l := newLoader()
	cfg := &Config{
		logLevel:  l.logLevel(),
		logFormat: l.logFormat(),
		logOutput: l.logOutput(),
	}
	if err := l.Err(); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return cfg, nil
}

// LogLevel returns the log level configured for the application.
// It can be "debug", "info", "warn", or "error".
func (cfg *Config) LogLevel() LogLevel {
	return cfg.logLevel
}

// LogFormat returns the log format configured for the application.
// It can be "text" or "json".
func (cfg *Config) LogFormat() LogFormat {
	return cfg.logFormat
}

// LogOutput returns the log output configured for the application.
// It can be "stdout", "stderr", or a file path.
func (cfg *Config) LogOutput() LogOutput {
	return cfg.logOutput
}

type (
	// loader is a helper struct for loading the application configuration.
	loader struct {
		// errs holds all the errors encountered during the loading process.
		errs []error
	}
)

// newLoader creates and returns a new loader instance.
func newLoader() *loader {
	return &loader{}
}

// logLevel loads, validates, and returns the log level from the environment
// variables.
// It accepts "debug", "info", "warn", or "error" (case-insensitive).
// It returns the default value if the environment variable is not set.
// If the value is invalid, it records the error and returns an empty LogLevel.
func (l *loader) logLevel() LogLevel {
	env := getEnv(EnvLogLevel, string(DefaultLogLevel))
	switch val := LogLevel(strings.ToLower(env)); val {
	case LogLevelDebug, LogLevelInfo, LogLevelWarn, LogLevelError:
		return val
	}
	l.appendError(fmt.Errorf("invalid log level (%s) got=%q", EnvLogLevel, env))
	return ""
}

// logFormat loads, validates, and returns the log format from the environment
// variables.
// It accepts "text" or "json" (case-insensitive).
// It returns the default value if the environment variable is not set.
// If the value is invalid, it records the error and returns an empty
// LogFormat.
func (l *loader) logFormat() LogFormat {
	env := getEnv(EnvLogFormat, string(DefaultLogFormat))
	switch val := LogFormat(strings.ToLower(env)); val {
	case LogFormatText, LogFormatJSON:
		return val
	}
	l.appendError(fmt.Errorf("invalid log format (%s) got=%q", EnvLogFormat, env))
	return ""
}

// logOutput loads, validates, and returns the log output from the environment
// variables.
// It accepts "stdout", "stderr" (case-insensitive), or a file path.
// It returns the default value if the environment variable is not set.
// If the value is invalid, it records the error and returns an empty
// LogOutput.
func (l *loader) logOutput() LogOutput {
	env := getEnv(EnvLogOutput, string(DefaultLogOutput))
	switch val := LogOutput(strings.ToLower(env)); val {
	case LogOutputStdout, LogOutputStderr:
		return val
	case "":
		l.appendError(fmt.Errorf("invalid log output (%s) got=%q", EnvLogOutput, env))
		return ""
	}
	return LogOutput(env)
}

// appendError adds a new error to the loader's internal slice of errors.
func (l *loader) appendError(err error) {
	l.errs = append(l.errs, err)
}

// Err returns a single error by joining all the errors recorded by the loader.
// If there are no errors recorded, it returns nil.
func (l *loader) Err() error {
	if len(l.errs) == 0 {
		return nil
	}
	return errors.Join(l.errs...)
}

// getEnv retrieves the value of the environment variable with the given name.
// If the variable is not set, it returns the provided default value.
func getEnv(name, defaultValue string) string {
	if val, ok := os.LookupEnv(name); ok {
		return val
	}
	return defaultValue
}
