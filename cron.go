package cron

import (
	"reflect"
	"runtime"
	"time"

	"github.com/gorhill/cronexpr"
	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type Fn func(ctx context.Context) error

func Daily(fn Fn) {
	Schedule("0 0 1 * * * *", fn)
}

func Schedule(schedule string, fn Fn) {
	logger := log.WithFields(log.Fields{
		"name":     runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name(),
		"schedule": schedule,
	})

	expr := cronexpr.MustParse(schedule)

	go func() {
		next := expr.Next(time.Now())
		logger.WithFields(log.Fields{"next-run": next}).Info("Schedule cron")
		time.Sleep(next.Sub(time.Now()))

		if err := fn(context.Background()); err != nil {
			logger.WithFields(log.Fields{"err": err.Error(), "stack": errors.ErrorStack(err)}).Error("Failed to run cron")
		}
	}()
}
