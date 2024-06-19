package log

import (
	"context"
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/nildebug/tools/mcontext"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

type ZapLogger struct {
	zap *zap.SugaredLogger
	//level            zapcore.Level
	loggerName       string
	loggerPrefixName string
	atomicLevel      zap.AtomicLevel
}

func NewZapLogger(
	loggerPrefixName, loggerName string,
	logLevel int,
	isStdout bool,
	isJson bool,
	logLocation string,
	rotateCount uint,
) (*ZapLogger, error) {
	atomicLevel := zap.NewAtomicLevelAt(logLevelMap[logLevel])
	zapConfig := zap.Config{
		Level: atomicLevel,
		// EncoderConfig: zap.NewProductionEncoderConfig(),
		// InitialFields:     map[string]interface{}{"PID": os.Getegid()},
		DisableStacktrace: true,
	}
	if isJson {
		zapConfig.Encoding = "json"
	} else {
		zapConfig.Encoding = "console"
	}
	// if isStdout {
	// 	zapConfig.OutputPaths = append(zapConfig.OutputPaths, "stdout", "stderr")
	// }
	zl := &ZapLogger{atomicLevel: atomicLevel, loggerName: loggerName, loggerPrefixName: loggerPrefixName}
	opts, err := zl.cores(isStdout, isJson, logLocation, rotateCount)
	if err != nil {
		return nil, err
	}
	l, err := zapConfig.Build(opts)
	if err != nil {
		return nil, err
	}
	zl.zap = l.Sugar()

	return zl, nil
}

func (l *ZapLogger) cores(isStdout bool, isJson bool, logLocation string, rotateCount uint) (zap.Option, error) {
	c := zap.NewProductionEncoderConfig()
	c.EncodeTime = l.timeEncoder
	c.EncodeDuration = zapcore.SecondsDurationEncoder
	c.MessageKey = "msg"
	c.LevelKey = "level"
	c.TimeKey = "time"
	c.CallerKey = "caller"
	c.NameKey = "logger"
	c.FunctionKey = "func"
	var fileEncoder zapcore.Encoder
	writer, err := l.getWriter(logLocation, rotateCount)
	if err != nil {
		return nil, err
	}
	var cores []zapcore.Core
	// if logLocation == "" && !isStdout {
	// 	return nil, errors.New("log storage location is empty and not stdout")
	// }
	if logLocation != "" {
		//写入日志时使用JSON格式
		c.EncodeLevel = zapcore.CapitalLevelEncoder
		fileEncoder = zapcore.NewJSONEncoder(c)
		fileEncoder.AddInt("PID", os.Getpid())
		cores = []zapcore.Core{
			zapcore.NewCore(fileEncoder, writer, l.atomicLevel),
		}
	}

	//调试输出 使用console
	if isStdout {
		//输出使用console
		c.EncodeLevel = zapcore.CapitalLevelEncoder
		c.EncodeCaller = l.customCallerEncoder
		fileEncoder = zapcore.NewConsoleEncoder(c)
		cores = append(cores, zapcore.NewCore(fileEncoder, zapcore.Lock(os.Stdout), l.atomicLevel))
		// cores = append(cores, zapcore.NewCore(fileEncoder, zapcore.Lock(os.Stderr), zap.NewAtomicLevelAt(l.level)))
	}
	return zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewTee(cores...)
	}), nil
}

// SetLogLevel 动态设置日志等级
//
// 日志等级 0(panic) 1(fatal) 2(error) 3(warn) 4(info) 5(debug)
func (l *ZapLogger) SetLogLevel(logLevel int) error {
	level, ok := logLevelMap[logLevel]
	if !ok {
		return fmt.Errorf("invalid log level: %d", logLevel)
	}
	l.atomicLevel.SetLevel(level) // 动态设置日志等级
	return nil
}

func (l *ZapLogger) customCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	//s := "[file:" + caller.FullPath() + "]"
	s := caller.FullPath() //使用完整路径
	//caller.Function
	// color, ok := _levelToColor[l.level]
	// if !ok {
	// 	color = _levelToColor[zapcore.ErrorLevel]
	// }
	enc.AppendString(s)
}

func (l *ZapLogger) timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	layout := "2006-01-02 15:04:05.000"
	type appendTimeEncoder interface {
		AppendTimeLayout(time.Time, string)
	}
	if enc, ok := enc.(appendTimeEncoder); ok {
		enc.AppendTimeLayout(t, layout)
		return
	}
	enc.AppendString(t.Format(layout))
}

func (l *ZapLogger) getWriter(logLocation string, rorateCount uint) (zapcore.WriteSyncer, error) {
	logf, err := rotatelogs.New(logLocation+sp+l.loggerPrefixName+".%Y-%m-%d.log",
		rotatelogs.WithRotationCount(rorateCount), //日志保留个数
		rotatelogs.WithRotationTime(24*time.Hour), //日志旋转时间
	)
	if err != nil {
		return nil, err
	}
	return zapcore.AddSync(logf), nil
}

func (l *ZapLogger) ToZap() *zap.SugaredLogger {
	return l.zap
}

func (l *ZapLogger) Debug(ctx context.Context, msg string, keysAndValues ...interface{}) {
	keysAndValues = l.kvAppend(ctx, keysAndValues)
	l.zap.Debugw(msg, keysAndValues...)

}

func (l *ZapLogger) Info(ctx context.Context, msg string, keysAndValues ...interface{}) {
	keysAndValues = l.kvAppend(ctx, keysAndValues)
	l.zap.Infow(msg, keysAndValues...)
}

func (l *ZapLogger) Warn(ctx context.Context, msg string, err error, keysAndValues ...interface{}) {
	if err != nil {
		keysAndValues = append(keysAndValues, "error", err.Error())
	}
	keysAndValues = l.kvAppend(ctx, keysAndValues)
	l.zap.Warnw(msg, keysAndValues...)
}

func (l *ZapLogger) Error(ctx context.Context, msg string, err error, keysAndValues ...interface{}) {
	if err != nil {
		keysAndValues = append(keysAndValues, "error", err.Error())
	}
	keysAndValues = l.kvAppend(ctx, keysAndValues)
	l.zap.Errorw(msg, keysAndValues...)
}

func (l *ZapLogger) kvAppend(ctx context.Context, keysAndValues []interface{}) []interface{} {
	if ctx == nil {
		return keysAndValues
	}
	key := mcontext.GetCenterKey(ctx)
	if key != "" {
		keysAndValues = append([]interface{}{mcontext.CenterKey, key}, keysAndValues...)
	}
	return keysAndValues
}

func (l *ZapLogger) WithValues(keysAndValues ...interface{}) Logger {
	dup := *l
	dup.zap = l.zap.With(keysAndValues...)
	return &dup
}

func (l *ZapLogger) WithName(name string) Logger {
	dup := *l
	dup.zap = l.zap.Named(name)
	return &dup
}

func (l *ZapLogger) WithCallDepth(depth int) Logger {
	dup := *l
	dup.zap = l.zap.WithOptions(zap.AddCallerSkip(depth))
	return &dup
}
