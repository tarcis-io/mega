// Package config loads and provides the application configuration.
package config

import (
	"errors"
	"os"
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
