package klog

import "context"

type Level int

const (
	TRACE Level = iota
	DEBUG
	INFO
	WARN
	ERR
)

type FullLogger interface {
	TraceCtx(ctx context.Context, v ...interface{})
	DebugCtx(ctx context.Context, v ...interface{})
	InfoCtx(ctx context.Context, v ...interface{})
	WarnCtx(ctx context.Context, v ...interface{})
	ErrCtx(ctx context.Context, v ...interface{})
	Trace(v ...interface{})
	Debug(v ...interface{})
	Info(v ...interface{})
	Warn(v ...interface{})
	Err(v ...interface{})
	Tracef(format string, v ...interface{})
	Debugf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Errf(format string, v ...interface{})
}
