package main

import (
	"atlas-quest/database"
	"atlas-quest/logger"
	"atlas-quest/quest"
	"atlas-quest/rest"
	"atlas-quest/tracing"
	"atlas-quest/wz"
	"context"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const serviceName = "atlas-quest"

func main() {
	l := logger.CreateLogger(serviceName)
	l.Infoln("Starting main service.")

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	tc, err := tracing.InitTracer(l)(serviceName)
	if err != nil {
		l.WithError(err).Fatal("Unable to initialize tracer.")
	}
	defer func(tc io.Closer) {
		err := tc.Close()
		if err != nil {
			l.WithError(err).Errorf("Unable to close tracer.")
		}
	}(tc)

	wzDir := os.Getenv("WZ_DIR")
	wz.GetFileCache().Init(wzDir)

	err = quest.GetCache().Init()
	if err != nil {
		l.WithError(err).Errorf("Unable to load quest cache.")
	}

	db := database.Connect(l, database.SetMigrations())

	rest.CreateService(l, db, ctx, wg, "/ms/quest", quest.InitResource)

	// trap sigterm or interrupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Block until a signal is received.
	sig := <-c
	l.Infof("Initiating shutdown with signal %s.", sig)
	cancel()
	wg.Wait()
	l.Infoln("Service shutdown.")
}
