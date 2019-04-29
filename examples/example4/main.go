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
	// TODO: adds all object from pic04.jpg
	restaurant := NewPipeline("Restaurant  Simulation", false)
	visitors_g := NewGenerator("Visitors", 18, 6, 0, 0, nil)
	out := NewHole("Out")
	visitors_q := NewQueue("Visitors queue")
	CheckQueueHndl := func(obj *Check, transact ITransaction) bool {
		var queue IQueue
		queue = obj.GetPipeline().GetObjByName("Check size of Visitors queue").(IQueue)
		if queue.GetLength() >= 6 {
			return false
		} else {
			return true
		}
	}
	check_queue := NewCheck("Check size of Visitors queue", CheckQueueHndl, out)
	hostes1_f := NewFacility("Hostes 1", 5, 3)
	hostes2_f := NewFacility("Hostes 2", 5, 3)
	tb1_in, tb1_out := NewBifacility("Table 1")
	tb2_in, tb2_out := NewBifacility("Table 2")
	tb3_in, tb3_out := NewBifacility("Table 3")
	tb4_in, tb4_out := NewBifacility("Table 4")
	tb5_in, tb5_out := NewBifacility("Table 5")
	tb6_in, tb6_out := NewBifacility("Table 6")
	tb7_in, tb7_out := NewBifacility("Table 7")
	tb8_in, tb8_out := NewBifacility("Table 8")
	dishes_sp := NewSplit("Selected dishes", 3, 2, nil)
	cook1_q := NewQueue("Queue of orders to cook 1 (meat dishes)")
	cook2_q := NewQueue("Queue of orders to cook 2 (fish dishes)")
	cook3_q := NewQueue("Queue of orders to cook 3 (salats)")
	cook4_q := NewQueue("Queue of orders to cook 4 (dessert)")
	bar_q := NewQueue("Queue of orders to bar")
	cook1_f := NewFacility("Cook 1 (meat dishes)", 15, 5)
	cook2_f := NewFacility("Cook 2 (fish dishes)", 15, 5)
	cook3_f := NewFacility("Cook 3 (salats)", 15, 5)
	cook4_f := NewFacility("Cook 4 (dessert)", 15, 5)
	bar_f := NewFacility("Bar", 15, 5)
	aggregate := NewAggregate("Aggregate orders")
	visitors_eating := NewAdvance("Visitors eating", 45, 10)
	visitors_pays := NewAdvance("Visitors pays", 5, 2)

	restaurant.Append(visitors_g, check_queue)
	restaurant.Append(check_queue, visitors_q)
	restaurant.Append(visitors_q, hostes1_f, hostes2_f)
	restaurant.Append(hostes1_f, tb1_in, tb2_in, tb3_in, tb4_in, tb5_in,
		tb7_in, tb8_in)
	restaurant.Append(hostes2_f, tb1_in, tb2_in, tb3_in, tb4_in, tb5_in,
		tb7_in, tb8_in)
	restaurant.Append(tb1_in, dishes_sp)
	restaurant.Append(tb2_in, dishes_sp)
	restaurant.Append(tb3_in, dishes_sp)
	restaurant.Append(tb4_in, dishes_sp)
	restaurant.Append(tb5_in, dishes_sp)
	restaurant.Append(tb6_in, dishes_sp)
	restaurant.Append(tb7_in, dishes_sp)
	restaurant.Append(tb8_in, dishes_sp)
	restaurant.Append(dishes_sp, cook1_q, cook2_q, cook3_q, cook4_q, bar_q)
	restaurant.Append(cook1_q, cook1_f)
	restaurant.Append(cook2_q, cook2_f)
	restaurant.Append(cook3_q, cook3_f)
	restaurant.Append(cook4_q, cook4_f)
	restaurant.Append(bar_q, bar_f)
	restaurant.Append(cook1_f, aggregate)
	restaurant.Append(cook2_f, aggregate)
	restaurant.Append(cook3_f, aggregate)
	restaurant.Append(cook4_f, aggregate)
	restaurant.Append(cook4_f, aggregate)
	restaurant.Append(bar_f, aggregate)
	restaurant.Append(aggregate, visitors_eating)
	restaurant.Append(visitors_eating, visitors_pays)
	restaurant.Append(visitors_pays, tb1_out, tb2_out, tb3_out, tb4_out, tb5_out,
		tb6_out, tb7_out, tb8_out)
	restaurant.Append(tb1_out, out)
	restaurant.Append(tb2_out, out)
	restaurant.Append(tb3_out, out)
	restaurant.Append(tb4_out, out)
	restaurant.Append(tb5_out, out)
	restaurant.Append(tb6_out, out)
	restaurant.Append(tb7_out, out)
	restaurant.Append(tb8_out, out)
	restaurant.Append(out)
	restaurant.Start(480)

	// Exit handler
	exit := make(chan struct{})
	closeSignal := make(chan os.Signal)
	signal.Notify(closeSignal, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-closeSignal
		close(exit)
	}()

	go func() {
		<-restaurant.Done
		close(exit)
	}()

	// Exit app if chan is closed
	<-exit
	restaurant.PrintReport()
	fmt.Println("Exit program")
}
