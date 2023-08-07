package main

import (
	"auth/internal/application"
	"auth/pkg/infra/logger"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer cancel()

	l, err := logger.New()
	if err != nil {
		log.Fatalf("logger initialization failed: %s", err.Error())
	}
	app := application.New(l)
	err = app.Start()
	if err != nil {
		l.Sugar().Fatalf("app not started: %s", err.Error())
	}

	<-ctx.Done()
	stopCtx, stopCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer stopCancel()
	app.Stop(stopCtx)
}
