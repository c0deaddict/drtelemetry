package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/c0deaddict/drtelemetry/telemetry"
	"github.com/c0deaddict/drtelemetry/ui"
)

func main() {
	fmt.Println("**********************************************")
	fmt.Println("**** DiRT Rally Telemetry Overlay by zh32 ****")
	fmt.Println("**********************************************")

	flag.Parse()
	log.SetFlags(0)

	dataChannel, quit := telemetry.RunServer()
	go ui.ListenAndServe(dataChannel)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	sig := <-c

	fmt.Printf("captured %v, stopping profiler and exiting..", sig)
	close(quit)
	os.Exit(1)
}
