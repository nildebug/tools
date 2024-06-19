// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package log

import (
	"context"
	"fmt"
	"path/filepath"

	"go.uber.org/zap/zapcore"
)

type Level int8

const (
	LogLevelDubug = 5
	LogLevelInfo  = 4
	LogLevelWarn  = 3
	LogLevelError = 2
	LogLevelFatal = 1
	LogLevelPanic = 0
)

var (
	pkgLogger   Logger
	sp          = string(filepath.Separator)
	logLevelMap = map[int]zapcore.Level{
		5: zapcore.DebugLevel,
		4: zapcore.InfoLevel,
		3: zapcore.WarnLevel,
		2: zapcore.ErrorLevel,
		1: zapcore.FatalLevel,
		0: zapcore.PanicLevel,
	}
	lzap *ZapLogger
)

// InitFromConfig
//
//	@Description: 初始化日志
//	@param loggerPrefixName 日志文件名前缀
//	@param moduleName 模块名称
//	@param logLevel 日志等级 0(panic) 1(fatal) 2(error) 3(warn) 4(info) 5(debug)
//	@param isStdout  是否输出到控制台
//	@param isJson 	是否输出为json格式
//	@param logLocation 日志文件路径
//	@param rotateCount 日志文件保留数量
//	@param
//	@return error
func InitFromConfig(
	loggerPrefixName, moduleName string,
	logLevel int,
	isStdout bool,
	isJson bool,
	logLocation string,
	rotateCount uint,
) error {
	l, err := NewZapLogger(loggerPrefixName, moduleName, logLevel, isStdout, isJson, logLocation, rotateCount)
	if err != nil {
		return err
	}
	lzap = l
	pkgLogger = l.WithCallDepth(2)
	if isJson {
		pkgLogger = pkgLogger.WithName(moduleName)
	}
	return nil
}

// SetLogLevel 动态设置日志等级
//
// 日志等级 0(panic) 1(fatal) 2(error) 3(warn) 4(info) 5(debug)
func SetLogLevel(level int) error {
	if lzap == nil {
		fmt.Errorf("lzap is nil")
	}
	return lzap.SetLogLevel(level)
}

func ZDebug(ctx context.Context, msg string, keysAndValues ...interface{}) {
	if pkgLogger == nil {
		return
	}
	pkgLogger.Debug(ctx, msg, keysAndValues...)
}

func ZInfo(ctx context.Context, msg string, keysAndValues ...interface{}) {
	if pkgLogger == nil {
		return
	}
	pkgLogger.Info(ctx, msg, keysAndValues...)
}

func ZWarn(ctx context.Context, msg string, err error, keysAndValues ...interface{}) {
	if pkgLogger == nil {
		return
	}
	pkgLogger.Warn(ctx, msg, err, keysAndValues...)
}

func ZError(ctx context.Context, msg string, err error, keysAndValues ...interface{}) {
	if pkgLogger == nil {
		return
	}
	pkgLogger.Error(ctx, msg, err, keysAndValues...)
}
