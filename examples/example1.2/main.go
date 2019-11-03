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
	// Generator -> Queue -> ...
	p := objects.NewPipeline("Barbershop", doneHandler).
		AddObject(objects.NewGenerator("Clients", 18, 6, 0, 0, nil)).
		AddObject(objects.NewQueue("Chairs"))

	// Create BiFacility
	fIN, fOUT := objects.NewBifacility("Master")

	// Build pipeline
	// ... -> BiFacilityIn -> Advance -> BeFacilityOut -> Hole
	p.
		AddObject(fIN).
		AddObject(objects.NewAdvance("Master work", 16, 4)).
		AddObject(fOUT).
		AddObject(objects.NewHole("Out"))
	// Start simulation
	p.Start(480)

	// Signal handler
	signalLoop()

	// Exit
	<-exit
	fmt.Println("Exit program")
}
