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
	// Generator -> Queue1 -> Facility1 -> Split ---> ...
	//                                            |
	//                                            --> ...
	p := objects.NewPipeline("Cafe Simulation", doneHandler).
		AddObject(objects.NewGenerator("Visitors", 18, 6, 0, 0, nil)).
		AddObject(objects.NewQueue("Visitors queue")).
		AddObject(objects.NewFacility("Order Acceptance", 5, 3)).
		AddObject(objects.NewSplit("Split orders", 1, 1, nil))

	//  ... -> Split ---> Queue2 -> ...
	//                |
	//                --> Queue3 -> ...
	baristaQ := objects.NewQueue("Queue of orders to barista")
	cookQ := objects.NewQueue("Queue of orders to cook")
	p.AddObject(baristaQ, cookQ)

	//  ... -> Queue2 -> Facility2 -> ...
	//
	//  ... -> Queue3 -> Facility3 -> ...
	baristaF := objects.NewFacility("Barista", 5, 2)
	cookF := objects.NewFacility("Cook", 10, 5)
	baristaQ.LinkObject(baristaF)
	cookQ.LinkObject(cookF)

	//  ... -> Facility2 ---> Aggregate -> Hole
	//                    ^
	//                    |
	//  ... -> Facility3 --
	aggregate := objects.NewAggregate("Aggregate orders")
	baristaF.LinkObject(aggregate)
	cookF.LinkObject(aggregate)
	aggregate.LinkObject(objects.NewHole("Out"))

	// Start simulation
	p.Start(480)

	// Signal handler
	signalLoop()

	// Exit
	<-exit
	fmt.Println("Exit program")
}
