package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/skeris/flat-grabber/app"
	"go.uber.org/zap"
)

func main() {

	defer time.Sleep(1500 * time.Millisecond)

	ctx, _ := context.WithCancel(context.Background())

	nedSvc, err := app.New(ctx, getEnv())
	if err != nil {
		panic(err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	logger := nedSvc.GetLogger()

	select {
	case <-stop:
		//defer cancel()
		logger.Info("Application was interrupted.")
	case err := <-nedSvc.GetErr():
		//defer cancel()
		logger.Panic("A fatal error occured", zap.Error(err))
	}
}
