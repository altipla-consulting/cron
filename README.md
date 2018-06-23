
# cron

[![GoDoc](https://godoc.org/github.com/altipla-consulting/cron?status.svg)](https://godoc.org/github.com/altipla-consulting/cron)
[![Build Status](https://travis-ci.org/altipla-consulting/cron.svg?branch=master)](https://travis-ci.org/altipla-consulting/cron)

Cron runner.


### Install

```shell
go get github.com/altipla-consulting/cron
```

This library has the following dependencies:

* [github.com/altipla-consulting/sentry](github.com/altipla-consulting/sentry)
* [github.com/gorhill/cronexpr](github.com/gorhill/cronexpr)
* [github.com/juju/errors](github.com/juju/errors)
* [github.com/julienschmidt/httprouter](github.com/julienschmidt/httprouter)
* [github.com/sirupsen/logrus](github.com/sirupsen/logrus)
* [golang.org/x/net/context](golang.org/x/net/context)


### Contributing

You can make pull requests or create issues in GitHub. Any code you send should be formatted using `gofmt`.


### Running tests

Run the tests

```shell
make test
```


### License

[MIT License](LICENSE)
