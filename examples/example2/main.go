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
	p := objects.NewPipeline("Water Closet Simulation").
		AddObject(objects.NewGenerator("Office", 0, 0, 0, 10, nil)).
		AddObject(objects.NewAdvance("Wanted to use the toilet", 90, 60)).
		AddObject(objects.NewAdvance("Path to WC", 5, 3)).
		AddObject(objects.NewQueue("Queue to the WC")).
		AddObject(objects.NewFacility("WC1", 15, 10), objects.NewFacility("WC2", 15, 10)).
		AddObject(objects.NewAdvance("Path from WC", 5, 3)).
		Loop("Wanted to use the toilet")

		// f1 := gpss.NewFacility("WC1", 15, 10)
		// f2 := gpss.NewFacility("WC2", 15, 10)
		// a3 := gpss.NewAdvance("Path from WC", 5, 3)
		// p.Append(g, a1)
		// p.Append(a1, a2)
		// p.Append(a2, q)
		// p.Append(q, f1, f2)
		// p.Append(f1, a3)
		// p.Append(f2, a3)
		// p.Append(a3, a1)

		// Start simulation
	p.Start(540)

	// Signal handler
	signalLoop()

	// Exit
	<-exit
	fmt.Println("Exit program")
}
