package main

import (
	"github.com/fairytale5571/crypto_page/pkg/app"
	"log"
	"os"
	"os/signal"
)

func main() {
	a, err := app.New()
	if err != nil {
		log.Fatalf("cant create app %s", err.Error())
		return
	}

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	sig := <-stop
	a.Logger.Infof("Close: received %v", sig.String())

	err = a.DB.Close()
	if err != nil {
		a.Logger.Errorf("Close: error close database: %v", err)
		return
	}
	a.Server.Stop()
	log.Fatalf("Graceful shutdown\n************************************************************************\n\n")

}
