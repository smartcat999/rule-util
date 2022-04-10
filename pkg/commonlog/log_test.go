package log

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

var _ = myInit()

func myInit() interface{} {
	//initLogger()
	return nil
}

func TestExample(t *testing.T) {
	//Debug("Debug")
	//Info("Info")
	//Warn("Warn")
	//Error("Error")
	//Fatal("Fatal")
	//Panic("Panic")
	DebugWithFields("Create Action| Input message invalid", Fields{"animal": "walrus"})
	InfoWithFields("Create Action| Input message invalid", Fields{"animal": "walrus"})
	SetLogLevel(logrus.DebugLevel)
	DebugWithFields("Create Action| Input message invalid", Fields{"animal": "walrus"})

	//logrus.WithFields(logrus.Fields{
	//	"animal": "walrus",
	//}).Info("A walrus appears")
}

func TestWarnf(t *testing.T) {
	initLogger(nil)
	go func() {
		for i := 0; i < 1000; i++ {
			Warnf("helelel%s:%d", "-hexing111", i)
		}
	}()
	go func() {
		for i := 0; i < 1000; i++ {
			Errorf("helelel%s:%d", "-hexing111", i)
		}
	}()
	time.Sleep(2 * time.Second)
}
func TestColor(t *testing.T) {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	logrus.Warnf("hello world %s", "sdsds")
}
