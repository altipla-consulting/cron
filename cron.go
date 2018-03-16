package cron

import (
	"net/http"
	"reflect"
	"runtime"
	"time"

	"github.com/gorhill/cronexpr"
	"github.com/juju/errors"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

var funcs map[string]Fn

func Handler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	job := ps.ByName("job")
	funcs[job](r.Context())
}

type Fn func(ctx context.Context) error

func Daily(fn Fn) {
	Schedule("0 0 1 * * * *", fn)
}

func Schedule(schedule string, fn Fn) {
	name := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	logger := log.WithFields(log.Fields{
		"name":     name,
		"schedule": schedule,
	})

	if IsLocal() {
		logger.Info("Cron configured")
		funcs[name] = fn
		return
	}

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
