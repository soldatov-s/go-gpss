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
	// Test Barbershop
	p := gpss.NewPipeline("Barbershop")
	g := gpss.NewGenerator("Clients", 18, 6, 0, 0, nil)
	q := gpss.NewQueue("Chairs")
	f := gpss.NewFacility("Master", 16, 4)
	h := gpss.NewHole("Out")
	p.Append(g, q)
	p.Append(q, f)
	p.Append(f, h)
	p.Append(h)
	p.Start(480)

	// Comment code before and uncoment next code for test Barbershop with
	// bifacility
	// p := gpss.NewPipeline("Barbershop", false)
	// g := gpss.NewGenerator("Clients", 18, 6, 0, 0, nil)
	// q := gpss.NewQueue("Chairs")
	// fIN, fOUT := gpss.NewBifacility("Master")
	// a := gpss.NewAdvance("Master work", 16, 4)
	// h := gpss.NewHole("Out")
	// p.Append(g, q)
	// p.Append(q, fIN)
	// p.Append(fIN, a)
	// p.Append(a, fOUT)
	// p.Append(fOUT, h)
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
	p.Report()
	fmt.Println("Exit program")
}
