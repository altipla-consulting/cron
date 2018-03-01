package cron

import (
	"reflect"
	"runtime"
	"time"

	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type Fn func(ctx context.Context) error
type Nexter func(t time.Time) time.Time

func Daily(fn Fn) {
	run(fn, "daily", func(t time.Time) time.Time {
		return time.Date(t.Year(), t.Month(), t.Day()+1, 1, 0, 0, 0, time.UTC)
	})
}

func run(fn Fn, schedule string, next Nexter) {
	logger := log.WithFields(log.Fields{
		"name":     runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name(),
		"schedule": schedule,
	})

	logger.Info("Schedule cron")
	go func() {
		n := time.Now()
		logger.WithFields(log.Fields{"time": next(n)}).Info("Waiting next run")
		time.Sleep(next(n).Sub(n))

		if err := fn(context.Background()); err != nil {
			logger.WithFields(log.Fields{"err": err.Error(), "stack": errors.ErrorStack(err)}).Error("Failed to run cron")
		}
	}()
}
