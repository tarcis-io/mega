// Package config loads and provides the application configuration.
package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type (
	// LogLevel represents the severity or verbosity of log records.
	LogLevel string
)

const (
	// LogLevelDebug captures detailed information, typically useful for development
	// and debugging.
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
	// LogFormat represents the encoding style of log records.
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

const (
	// EnvLogLevel specifies the environment variable name for configuring the
	// [LogLevel].
	//
	// Expected values:
	//
	//  - [LogLevelDebug]
	//  - [LogLevelInfo]
	//  - [LogLevelWarn]
	//  - [LogLevelError]
	//
	// Default: [DefaultLogLevel]
	EnvLogLevel = "LOG_LEVEL"

	// EnvLogFormat specifies the environment variable name for configuring the
	// [LogFormat].
	//
	// Expected values:
	//
	//  - [LogFormatText]
	//  - [LogFormatJSON]
	//
	// Default: [DefaultLogFormat]
	EnvLogFormat = "LOG_FORMAT"

	// EnvLogOutput specifies the environment variable name for configuring the
	// [LogOutput].
	//
	// Expected values:
	//
	//  - [LogOutputStdout]
	//  - [LogOutputStderr]
	//  - file path (e.g. /var/log/app.log)
	//
	// Default: [DefaultLogOutput]
	EnvLogOutput = "LOG_OUTPUT"
)

const (
	// DefaultLogLevel defines the default [LogLevel].
	DefaultLogLevel = LogLevelInfo

	// DefaultLogFormat defines the default [LogFormat].
	DefaultLogFormat = LogFormatText

	// DefaultLogOutput defines the default [LogOutput].
	DefaultLogOutput = LogOutputStdout
)

type (
	// Config represents the immutable application configuration.
	Config struct {
		logLevel  LogLevel
		logFormat LogFormat
		logOutput LogOutput
	}
)

// LogLevel returns the configured severity or verbosity of log records.
func (cfg *Config) LogLevel() LogLevel {
	return cfg.logLevel
}

// LogFormat returns the configured encoding style of log records.
func (cfg *Config) LogFormat() LogFormat {
	return cfg.logFormat
}

// LogOutput returns the configured destination stream for log records.
func (cfg *Config) LogOutput() LogOutput {
	return cfg.logOutput
}

type (
	// loader
	loader struct {
		errs []error
	}
)

// newLoader
func newLoader() *loader {
	return &loader{}
}

// logLevel
func (l *loader) logLevel() LogLevel {
	env := getEnv(EnvLogLevel, string(DefaultLogLevel))
	switch val := LogLevel(strings.ToLower(env)); val {
	case LogLevelDebug, LogLevelInfo, LogLevelWarn, LogLevelError:
		return val
	}
	l.appendError(fmt.Errorf("invalid log level (%s) got=%q", EnvLogLevel, env))
	return ""
}

// logFormat
func (l *loader) logFormat() LogFormat {
	env := getEnv(EnvLogFormat, string(DefaultLogFormat))
	switch val := LogFormat(strings.ToLower(env)); val {
	case LogFormatText, LogFormatJSON:
		return val
	}
	l.appendError(fmt.Errorf("invalid log format (%s) got=%q", EnvLogFormat, env))
	return ""
}

// logOutput
func (l *loader) logOutput() LogOutput {
	env := getEnv(EnvLogOutput, string(DefaultLogOutput))
	switch val := LogOutput(strings.ToLower(env)); val {
	case LogOutputStdout, LogOutputStderr:
		return val
	}
	if env == "" {
		l.appendError(fmt.Errorf("invalid log output (%s) got=%q", EnvLogOutput, env))
		return ""
	}
	return LogOutput(env)
}

// appendError
func (l *loader) appendError(err error) {
	l.errs = append(l.errs, err)
}

// Err
func (l *loader) Err() error {
	if len(l.errs) == 0 {
		return nil
	}
	return errors.Join(l.errs...)
}

// getEnv
func getEnv(key, defaultValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultValue
}
