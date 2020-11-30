package hrpc

import (
	"fmt"
	"github.com/kamva/hexa"
	"google.golang.org/grpc/grpclog"
)

// logger implements the gRPC logger v2
type logger struct {
	logger hexa.Logger
	v      int
}

func (l *logger) Info(args ...interface{}) {
	l.logger.Info(fmt.Sprintln(args...))
}

func (l *logger) Infoln(args ...interface{}) {
	l.logger.Info(fmt.Sprintln(args...))
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(format, fmt.Sprintln(args...)))
}

func (l *logger) Warning(args ...interface{}) {
	l.logger.Warn(fmt.Sprintln(args...))
}

func (l *logger) Warningln(args ...interface{}) {
	l.logger.Warn(fmt.Sprintln(args...))
}

func (l *logger) Warningf(format string, args ...interface{}) {
	l.logger.Warn(fmt.Sprintf(format, fmt.Sprintln(args...)))
}

func (l *logger) Error(args ...interface{}) {
	l.logger.Error(fmt.Sprintln(args...))
}

func (l *logger) Errorln(args ...interface{}) {
	l.logger.Error(fmt.Sprintln(args...))
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.logger.Error(fmt.Sprintf(format, fmt.Sprintln(args...)))
}

func (l *logger) Fatal(args ...interface{}) {
	l.logger.Error(fmt.Sprintln(args...))
}

func (l *logger) Fatalln(args ...interface{}) {
	l.logger.Error(fmt.Sprintln(args...))
}

func (l *logger) Fatalf(format string, args ...interface{}) {
	l.logger.Error(fmt.Sprintf(format, fmt.Sprintln(args...)))
}

func (l *logger) V(level int) bool {
	return level <= l.v
}

// NewLogger returns new instance of the gRPC Logger v2
func NewLogger(l hexa.Logger, level int) grpclog.LoggerV2 {
	// Detect log v
	return &logger{logger: l, v: level}
}

var _ grpclog.LoggerV2 = &logger{}
