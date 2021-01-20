package ctxlog

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.StandardLogger().Level = logrus.DebugLevel
}

type Context interface {
	// Parent returns the PID for the current actors parent
	Parent() *actor.PID

	// Self returns the PID for the current actor
	Self() *actor.PID

	// Actor returns the actor associated with this context
	Actor() actor.Actor

	ActorSystem() *actor.ActorSystem
}

func Infof(c Context, format string, args ...interface{}) {
	logrus.Infof(logMsgf(c, format), args...)
}

func Info(c Context, args ...interface{}) {
	logrus.Info(logMsg(c, args)...)
}

func Debugf(c Context, format string, args ...interface{}) {
	logrus.Debugf(logMsgf(c, format), args...)
}

func Debug(c Context, args ...interface{}) {
	logrus.Debug(logMsg(c, args)...)
}

func logMsg(c Context, args ...interface{}) []interface{} {
	a := []interface{}{contextMsg(c)}
	a = append(a, args...)
	return a
}

func logMsgf(c Context, format string) string {
	f := fmt.Sprintf("pid[%s] - ", c.Self())
	f += format
	return f
}

func contextMsg(c Context) string {
	return fmt.Sprintf("pid[%s] - ", c.Self())
}
