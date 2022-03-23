// Copyright (c) 2017 Uber Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package log

import (
	"context"
	"github.com/tkeel-io/rule-util/pkg/logfield"
	"strings"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type registeredLogger struct {
	logger       Factory
	isRegistered bool
}

var (
	globalLogger registeredLogger = registeredLogger{noopFactory, false}
	callbacks                     = []Callback{}
)

// Logger is a simplified abstraction of the zap.Logger
type Factory interface {
	Bg() Logger
	For(ctx context.Context) Logger
	With(fields ...zapcore.Field) Factory
	withOptions(opts ...zap.Option) Factory
	Check(lvl zapcore.Level, msg string) *zapcore.CheckedEntry
}

// Factory is the default logging wrapper that can create
// logger instances either for a given Context or context-less.
type factory struct {
	logger *zap.Logger
}

type Callback func(logger Factory)

func InitLogger(source, appID string) (logger Factory, level Level) {
	logger, level = NewProduction(
		logf.Fields(logf.String("source", source)),
		logf.Fields(logf.String("appID", appID)),
		logf.AddStacktrace(FatalLevel),
	)
	SetGlobalLogger(logger)
	return
}

func SetGlobalLogger(logger Factory) {
	globalLogger = registeredLogger{logger, true}
	for _, cb := range callbacks {
		cb(logger)
	}
}

func Register(cb Callback) {
	callbacks = append(callbacks, cb)
}

func GlobalLogger() Factory {
	return globalLogger.logger
}

func globalLoggerWarp() Factory {
	return globalLogger.logger.withOptions(
		zap.AddCallerSkip(1))
}

func For(ctx context.Context) Logger {
	return globalLogger.logger.For(ctx)
}

// New creates a new Factory.
func New(options ...Option) Factory {
	opts := []Option{
		zap.AddCallerSkip(1),
	}
	zapLogger, _ := zap.NewDevelopment(
		append(opts, options...)...,
	)

	return factory{logger: zapLogger}
}

// New creates a new Factory.
func NewDevelopment(options ...Option) (Factory, Level) {

	logcfg := zap.NewDevelopmentConfig()

	alevel := zap.NewAtomicLevel()
	logcfg.Level = alevel

	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.MessageKey = "content"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logcfg.EncoderConfig = encoderConfig

	opts := []Option{
		zap.AddCallerSkip(1),
	}

	zapLogger, _ := logcfg.Build(
		append(opts, options...)...,
	)

	return factory{logger: zapLogger}, alevel
}

// New creates a new Factory.
func NewProduction(options ...Option) (Factory, Level) {

	logcfg := zap.NewProductionConfig()

	alevel := zap.NewAtomicLevel()
	logcfg.Level = alevel

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.MessageKey = "content"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logcfg.EncoderConfig = encoderConfig

	opts := []Option{
		zap.AddCallerSkip(1),
	}

	zapLogger, _ := logcfg.Build(
		append(opts, options...)...,
	)

	return factory{logger: zapLogger}, alevel
}

// NewFactory creates a new Factory.
func NewFactory(logger *zap.Logger) Factory {
	return factory{logger: logger}
}

// Bg creates a context-unaware logger.
func (b factory) Bg() Logger {
	return logger(b)
}

// For returns a context-aware Logger. If the context
// contains an OpenTracing span, all logging calls are also
// echo-ed into the span.
func (b factory) For(ctx context.Context) Logger {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		// TODO for Jaeger span extract trace/span IDs as fields
		return spanLogger{span: span, logger: b.logger}
	}
	return b.Bg()
}

// With creates a child logger, and optionally adds some context fields to that logger.
func (b factory) With(fields ...zapcore.Field) Factory {
	return factory{logger: b.logger.With(fields...)}
}

func (b factory) withOptions(opts ...zap.Option) Factory {
	return factory{logger: b.logger.WithOptions(opts...)}
}

func GetLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn":
		return WarnLevel
	case "error":
		return ErrorLevel
	case "dpanic":
		return DPanicLevel
	case "panic":
		return PanicLevel
	case "fatal":
		return FatalLevel
	}
	return ErrorLevel
}

func (b factory) Check(lvl zapcore.Level, msg string) *zapcore.CheckedEntry {
	return b.logger.Check(lvl, msg)
}
