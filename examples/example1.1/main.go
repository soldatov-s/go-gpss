// examples
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/soldatov-s/go-gpss/objects"
)

var (
	exit chan struct{}
)

func doneHandler(p *objects.Pipeline) {
	p.Report()
	close(exit)
}

func signalLoop() {
	closeSignal := make(chan os.Signal)
	signal.Notify(closeSignal, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-closeSignal
		close(exit)
	}()
}

func main() {
	// Init exit chan
	exit = make(chan struct{})

	// Build pipeline
	// Generator -> Queue -> Facility -> Hole
	p := objects.NewPipeline("Barbershop").
		AddObject(objects.NewGenerator("Clients", 18, 6, 0, 0, nil)).
		AddObject(objects.NewQueue("Chairs")).
		AddObject(objects.NewFacility("Master", 16, 4)).
		AddObject(objects.NewHole("Out"))
	// Start simulation
	p.Start(480)

	<-p.Done
	p.Report()

	// Exit
	fmt.Println("Exit program")
}
