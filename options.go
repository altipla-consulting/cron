package cron

func WithSentry(dsn string) Option {
	return func(runner *Runner) {
		runner.dsn = dsn
	}
}
