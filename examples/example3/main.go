// examples
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/soldatov-s/go-gpss"
)

func main() {
	// Test small cafe with barista and cook. Used split and aggregate components
	p := gpss.NewPipeline("Cafe Simulation")
	g := gpss.NewGenerator("Visitors", 18, 6, 0, 0, nil)
	q := gpss.NewQueue("Visitors queue")
	ordersF := gpss.NewFacility("Order Acceptance", 5, 3)
	split := gpss.NewSplit("Split orders", 1, 1, nil)
	baristaQ := gpss.NewQueue("Queue of orders to barista")
	baristaF := gpss.NewFacility("Barista", 5, 2)
	cookQ := gpss.NewQueue("Queue of orders to cook")
	cookF := gpss.NewFacility("Cook", 10, 5)
	aggregate := gpss.NewAggregate("Aggregate orders")
	h := gpss.NewHole("Out")
	p.Append(g, q)
	p.Append(q, ordersF)
	p.Append(ordersF, split)
	p.Append(split, baristaQ, cookQ)
	p.Append(baristaQ, baristaF)
	p.Append(cookQ, cookF)
	p.Append(baristaF, aggregate)
	p.Append(cookF, aggregate)
	p.Append(aggregate, h)
	p.Append(h)
	p.Start(480)

	// Exit handler
	exit := make(chan struct{})
	closeSignal := make(chan os.Signal)
	signal.Notify(closeSignal, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-closeSignal
		close(exit)
	}()

	go func() {
		<-p.Done
		close(exit)
	}()

	// Exit app if chan is closed
	<-exit
	p.Report()
	fmt.Println("Exit program")
}
