/*
 * Copyright (C) 2019 Yunify, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this work except in compliance with the License.
 * You may obtain a copy of the License in the LICENSE file, or at:
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package log

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Option = zap.Option
type Level = zap.AtomicLevel

var (
	DebugLevel  = zapcore.DebugLevel
	InfoLevel   = zapcore.InfoLevel
	WarnLevel   = zapcore.WarnLevel
	ErrorLevel  = zapcore.ErrorLevel
	DPanicLevel = zapcore.DPanicLevel
	PanicLevel  = zapcore.PanicLevel
	FatalLevel  = zapcore.FatalLevel
)

type NoopFactory struct {
	logger *zap.Logger
}

type NoopLogger struct{}

var noopFactory = &NoopFactory{
	logger: zap.NewNop(),
}

var noopLogger = &NoopLogger{}

func (f *NoopFactory) Bg() Logger {
	return noopLogger
}

func (f *NoopFactory) For(ctx context.Context) Logger {
	return noopLogger

}

func (f *NoopFactory) With(fields ...zapcore.Field) Factory {
	return f
}
func (f *NoopFactory) withOptions(opts ...zap.Option) Factory {
	return f
}

func (*NoopLogger) Debug(msg string, fields ...zapcore.Field) {
}

func (*NoopLogger) Info(msg string, fields ...zapcore.Field) {
}

func (*NoopLogger) Warn(msg string, fields ...zapcore.Field) {
}

func (*NoopLogger) Error(msg string, fields ...zapcore.Field) {
	fmt.Printf("Error: %v\n", msg, fields)
}

func (*NoopLogger) Fatal(msg string, fields ...zapcore.Field) {
	panic(fmt.Sprintln("Fatal: %v", msg, fields))
}

func (l *NoopLogger) Panic(msg string, fields ...zapcore.Field) {
	panic(fmt.Sprintln("Panic: %v", msg, fields))
}

func (l *NoopFactory) Check(lvl zapcore.Level, msg string) *zapcore.CheckedEntry {
	return l.logger.Check(lvl, msg)
}

func (*NoopLogger) Flush() {
}

func (*NoopLogger) With(fields ...zapcore.Field) Logger {
	return noopLogger
}
