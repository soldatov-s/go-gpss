// examples
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	. "github.com/soldatov-s/go-gpss"
)

func main() {
	// Test small cafe with barista and cook. Used split and aggregate components
	p := NewPipeline("Cafe Simulation", false)
	g := NewGenerator("Visitors", 18, 6, 0, 0, nil)
	q := NewQueue("Visitors queue")
	orders_f := NewFacility("Order Acceptance", 5, 3)
	split := NewSplit("Split orders", 1, 1, nil)
	barista_q := NewQueue("Queue of orders to barista")
	barista_f := NewFacility("Barista", 5, 2)
	cook_q := NewQueue("Queue of orders to cook")
	cook_f := NewFacility("Cook", 10, 5)
	aggregate := NewAggregate("Aggregate orders")
	h := NewHole("Out")
	p.Append(g, q)
	p.Append(q, orders_f)
	p.Append(orders_f, split)
	p.Append(split, barista_q, cook_q)
	p.Append(barista_q, barista_f)
	p.Append(cook_q, cook_f)
	p.Append(barista_f, aggregate)
	p.Append(cook_f, aggregate)
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
	p.PrintReport()
	fmt.Println("Exit program")
}
