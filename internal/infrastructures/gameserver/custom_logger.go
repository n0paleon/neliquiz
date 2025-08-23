package gameserver

import (
	"github.com/sirupsen/logrus"
	"github.com/topfreegames/pitaya/v3/pkg/logger/interfaces"
)

type logger struct {
}

func (l logger) Fatal(format ...interface{}) {

}

func (l logger) Fatalf(format string, args ...interface{}) {

}

func (l logger) Fatalln(args ...interface{}) {

}

func (l logger) Debug(args ...interface{}) {

}

func (l logger) Debugf(format string, args ...interface{}) {

}

func (l logger) Debugln(args ...interface{}) {

}

func (l logger) Error(args ...interface{}) {

}

func (l logger) Errorf(format string, args ...interface{}) {

}

func (l logger) Errorln(args ...interface{}) {

}

func (l logger) Info(args ...interface{}) {

}

func (l logger) Infof(format string, args ...interface{}) {

}

func (l logger) Infoln(args ...interface{}) {

}

func (l logger) Warn(args ...interface{}) {

}

func (l logger) Warnf(format string, args ...interface{}) {

}

func (l logger) Warnln(args ...interface{}) {

}

func (l logger) Panic(args ...interface{}) {

}

func (l logger) Panicf(format string, args ...interface{}) {

}

func (l logger) Panicln(args ...interface{}) {

}

func (l logger) WithFields(fields map[string]interface{}) interfaces.Logger {
	return l
}

func (l logger) WithField(key string, value interface{}) interfaces.Logger {
	return l
}

func (l logger) WithError(err error) interfaces.Logger {
	return l
}

func (l logger) GetInternalLogger() any {
	return logrus.New()
}

func newCustomLogger() interfaces.Logger {
	return logger{}
}
