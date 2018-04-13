package cron

func AddSentry(in string) Option {
	return func(runner *Runner) {
		runner.dsn = in
	}
}
