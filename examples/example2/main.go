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
	p := NewPipeline("Water Closet Simulation", false)
	g := NewGenerator("Office", 0, 0, 0, 10, nil)
	a1 := NewAdvance("Wanted to use the toilet", 90, 60)
	a2 := NewAdvance("Path to WC", 5, 3)
	q := NewQueue("Queue to the WC")
	f1 := NewFacility("WC1", 15, 10)
	f2 := NewFacility("WC2", 15, 10)
	a3 := NewAdvance("Path from WC", 5, 3)
	p.Append(g, a1)
	p.Append(a1, a2)
	p.Append(a2, q)
	p.Append(q, f1, f2)
	p.Append(f1, a3)
	p.Append(f2, a3)
	p.Append(a3, a1)
	p.Start(540)

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
