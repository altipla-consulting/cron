package main

import (
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	"github.com/altipla-consulting/cron"
)

func main() {
	cron.Daily(Sync)
	time.Sleep(10 * time.Second)
}

func Sync(ctx context.Context) error {
	log.Println("here")
	return nil
}
