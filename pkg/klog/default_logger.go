package klog

import (
	"context"
	"log"
	"os"
)

type defaultLogger struct {
	logLevel Level
	stdLog   *log.Logger
}

func (d *defaultLogger) setLogLevel(l Level) {
	d.logLevel = l
}

func (d *defaultLogger) setLogger(l *log.Logger) {
	d.stdLog = l
}

var logger FullLogger = &defaultLogger{
	logLevel: INFO,
	stdLog:   log.New(os.Stderr, "[Krpc]", log.LstdFlags|log.Lshortfile|log.Lmicroseconds),
}

func TraceCtx(ctx context.Context, v ...interface{}) {
	logger.TraceCtx(ctx, v...)
}
func DebugCtx(ctx context.Context, v ...interface{}) {
	logger.DebugCtx(ctx, v...)
}
func InfoCtx(ctx context.Context, v ...interface{}) {
	logger.InfoCtx(ctx, v...)
}
func WarnCtx(ctx context.Context, v ...interface{}) {
	logger.WarnCtx(ctx, v...)
}
func ErrCtx(ctx context.Context, v ...interface{}) {
	logger.ErrCtx(ctx, v...)
}
func Trace(v ...interface{}) {
	logger.Trace(v...)
}
func Debug(v ...interface{}) {
	logger.Debug(v...)
}
func Info(v ...interface{}) {
	logger.Info(v...)
}
func Warn(v ...interface{}) {
	logger.Warn(v...)
}
func Err(v ...interface{}) {
	logger.Err(v...)
}

func Tracef(format string, v ...interface{}) {
	logger.Tracef(format, v...)
}
func Debugf(format string, v ...interface{}) {
	logger.Debugf(format, v...)
}
func Infof(format string, v ...interface{}) {
	logger.Infof(format, v...)
}
func Warnf(format string, v ...interface{}) {
	logger.Warnf(format, v...)
}
func Errf(format string, v ...interface{}) {
	logger.Errf(format, v...)
}

func (d *defaultLogger) TraceCtx(ctx context.Context, v ...interface{}) {
	if d.logLevel > TRACE {
		return
	}
	d.stdLog.Printf("%v\n", v...)
}

func (d *defaultLogger) DebugCtx(ctx context.Context, v ...interface{}) {
	if d.logLevel > DEBUG {
		return
	}
	d.stdLog.Printf("%v\n", v...)
}

func (d *defaultLogger) InfoCtx(ctx context.Context, v ...interface{}) {
	if d.logLevel > INFO {
		return
	}
	d.stdLog.Printf("%v\n", v...)
}

func (d *defaultLogger) WarnCtx(ctx context.Context, v ...interface{}) {
	if d.logLevel > WARN {
		return
	}
	d.stdLog.Printf("%v\n", v...)
}
func (d *defaultLogger) ErrCtx(ctx context.Context, v ...interface{}) {
	if d.logLevel > ERR {
		return
	}
	d.stdLog.Printf("%v\n", v...)
}
func (d *defaultLogger) Trace(v ...interface{}) {
	if d.logLevel > TRACE {
		return
	}
	d.stdLog.Printf("%v\n", v...)
}
func (d *defaultLogger) Debug(v ...interface{}) {
	if d.logLevel > DEBUG {
		return
	}
	d.stdLog.Printf("%v\n", v...)
}
func (d *defaultLogger) Info(v ...interface{}) {
	if d.logLevel > INFO {
		return
	}
	d.stdLog.Printf("%v\n", v...)
}
func (d *defaultLogger) Warn(v ...interface{}) {
	if d.logLevel > WARN {
		return
	}
	d.stdLog.Printf("%v\n", v...)
}
func (d *defaultLogger) Err(v ...interface{}) {
	if d.logLevel > ERR {
		return
	}
	d.stdLog.Printf("%v\n", v...)
}

func (d *defaultLogger) Tracef(format string, v ...interface{}) {
	if d.logLevel > TRACE {
		return
	}
	d.stdLog.Printf(format, v...)
}
func (d *defaultLogger) Debugf(format string, v ...interface{}) {
	if d.logLevel > DEBUG {
		return
	}
	d.stdLog.Printf(format, v...)
}
func (d *defaultLogger) Infof(format string, v ...interface{}) {
	if d.logLevel > INFO {
		return
	}
	d.stdLog.Printf(format, v...)
}
func (d *defaultLogger) Warnf(format string, v ...interface{}) {
	if d.logLevel > WARN {
		return
	}
	d.stdLog.Printf(format, v...)
}
func (d *defaultLogger) Errf(format string, v ...interface{}) {
	if d.logLevel > ERR {
		return
	}
	d.stdLog.Printf(format, v...)
}
