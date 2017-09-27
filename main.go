package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"./iwnet"
)

func main() {
	// create log file
	logFile, err := os.OpenFile("IWNetServerInGo.log", os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	// logger will write into te file AND into stdout
	mw := io.MultiWriter(os.Stdout, logFile)
	logger := log.New(mw, "", log.Ldate|log.Ltime)

	fmt.Print("IWNetServerInGo by gmb\n\n")

	// start the IPserver
	IPServer := iwnet.NewIPServer(logger)
	IPServer.StartAndRun()
	defer IPServer.Close()

	// start the HTTPserver
	HTTPServer := iwnet.NewHTTPServer(logger)
	HTTPServer.StartAndRunHTTP()

	// wait for Ctrl + C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // Ctrl + C Win/*nix
	go func() {
		select {
		case <-c:
			return
		}
	}()
}
