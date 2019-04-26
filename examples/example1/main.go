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
	// Test Barbershop
	p := NewPipeline("Barbershop", false)
	g := NewGenerator("Clients", 18, 6, 0, 0, nil)
	q := NewQueue("Chairs")
	f := NewFacility("Master", 16, 4)
	h := NewHole("Out")
	p.Append(g, q)
	p.Append(q, f)
	p.Append(f, h)
	p.Append(h)
	p.Start(480)

	// Comment code before and uncoment next code for test Barbershop with
	// bifacility
	// p := NewPipeline("Barbershop", false)
	// g := NewGenerator("Clients", 18, 6, 0, 0, nil)
	// q := NewQueue("Chairs")
	// f_in, f_out := NewBifacility("Master")
	// a := NewAdvance("Master work", 16, 4)
	// h := NewHole("Out")
	// p.Append(g, q)
	// p.Append(q, f_in)
	// p.Append(f_in, a)
	// p.Append(a, f_out)
	// p.Append(f_out, h)
	// p.Append(h)
	// p.Start(480)

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
