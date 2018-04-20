package cron

import (
	"net/http"
	"path/filepath"
	"reflect"
	"runtime"
	"time"

	"github.com/altipla-consulting/sentry"
	"github.com/gorhill/cronexpr"
	"github.com/juju/errors"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type Fn func(ctx context.Context) error

type Runner struct {
	funcs map[string]Fn
	dsn   string
}

type Option func(runner *Runner)

func NewRunner(opts ...Option) *Runner {
	runner := &Runner{
		funcs: map[string]Fn{},
	}

	for _, opt := range opts {
		opt(runner)
	}

	return runner
}

func (runner *Runner) Handler() func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		job := ps.ByName("job")
		runner.funcs[job](r.Context())
	}
}

func (runner *Runner) Daily(fn Fn) {
	runner.Schedule("0 0 1 * * * *", fn)
}

func (runner *Runner) Schedule(schedule string, fn Fn) {
	name := filepath.Base(runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name())
	logger := log.WithFields(log.Fields{
		"name":     name,
		"schedule": schedule,
	})

	logger.Info("Configuring cron")

	if IsLocal() {
		logger.Info("Cron configured")
		runner.funcs[name] = fn
		return
	}

	var client *sentry.Client
	if runner.dsn != "" {
		client = sentry.NewClient(runner.dsn)
	}

	expr := cronexpr.MustParse(schedule)

	go func() {
		next := expr.Next(time.Now())
		logger.WithField("next-run", next).Info("Schedule cron")
		time.Sleep(next.Sub(time.Now()))

		ctx := context.Background()
		if client != nil {
			ctx = sentry.WithContext(ctx)
		}

		if err := fn(ctx); err != nil {
			logger.WithFields(log.Fields{"err": err.Error(), "stack": errors.ErrorStack(err)}).Error("Failed to run cron")
			if client != nil {
				client.ReportInternal(ctx, err)
			}
		}
	}()
}
