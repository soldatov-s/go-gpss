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
	restaurant := NewPipeline("Restaurant  Simulation")
	// 1. Create the Generator and Queue of Visitors, create a Hole
	visitors_g := NewGenerator("Visitors", 10, 5, 0, 0, nil)
	// visitors_g := NewGenerator("Visitors", 0, 0, 0, 1, nil)
	out := NewHole("Out")
	visitors_q := NewQueue("Visitors queue")
	// 2. Create the Check for checking size the Queue of Visitors
	CheckQueueHndl := func(obj *Check, transact *Transaction) bool {
		queue := obj.Pipe.GetObjByName("Visitors queue")
		if queue.(IQueue).GetLength() >= 6 {
			return false
		}
		return true
	}
	check_queue := NewCheck("Check size of Visitors queue", CheckQueueHndl, out)
	// 3. Create are Hostess
	hostes1_f := NewFacility("hostess 1", 5, 3)
	hostes2_f := NewFacility("Hostess 2", 5, 3)
	// 4. Create are Tables
	cnt_tables := 24
	tables_in := make([]IBaseObj, cnt_tables)
	tables_out := make([]IBaseObj, cnt_tables)
	for i := 0; i < cnt_tables; i++ {
		table_name := fmt.Sprintf("Table %d", i+1)
		tables_in[i], tables_out[i] = NewBifacility(table_name)
	}
	// 5. Check that we have empty table
	CheckEmptyTableHndl := func(obj *Check, transact *Transaction) bool {
		for i := 0; i < cnt_tables; i++ {
			table_name := fmt.Sprintf("Table %d", i+1)
			table := obj.Pipe.GetObjByName(table_name).(IFacility)
			if table.IsEmpty() {
				return true
			}
		}
		return false
	}
	check_empty_table := NewCheck("Check empty table", CheckEmptyTableHndl, nil)
	// 6. Create the Queues and Facilities for waiters
	cnt_waiters := 8
	waiters_queue := make([]IBaseObj, cnt_waiters)
	waiters_facility := make([]IBaseObj, cnt_waiters)
	for i := 0; i < cnt_waiters; i++ {
		queue_name := fmt.Sprintf("Queue waiter %d", i+1)
		waiters_queue[i] = NewQueue(queue_name)
		facility_name := fmt.Sprintf("Waiter %d", i+1)
		waiters_facility[i] = NewFacility(facility_name, 5, 3)

	}
	// 7. Create the Split for splitting Visitors order to dishes
	// Maybe 1 or 5 dishes, includes bar
	dishes_sp := NewSplit("Selected dishes", 3, 2, nil)
	// 8. Create the Check that transact is a coocked dishes. If false, this transact is
	// an order from tables, needs split to dishes
	dish_state := Parameter{Name: "State", Value: "After kitchen"}
	check_is_cooked := NewCheck("After kitchen?", nil, dishes_sp, dish_state)
	// 9. Create the Queue and Facility for cooks and barmans. Each cook cooking
	// only one type of dishes: meat, sushi, salats, dessert
	cook1_q := NewQueue("Queue of orders to cook 1 (meat dishes)")
	cook2_q := NewQueue("Queue of orders to cook 2 (fish dishes)")
	cook3_q := NewQueue("Queue of orders to cook 3 (salats)")
	cook4_q := NewQueue("Queue of orders to cook 4 (dessert)")
	bar_q := NewQueue("Queue of orders to bar")
	cook1_f := NewFacility("Cook 1 (meat dishes)", 15, 5)
	cook2_f := NewFacility("Cook 2 (sushi)", 7, 3)
	cook3_f := NewFacility("Cook 3 (salats)", 10, 4)
	cook4_f := NewFacility("Cook 4 (dessert)", 20, 5)
	barman1_f := NewFacility("Barman 1", 4, 2)
	barman2_f := NewFacility("Barman 2", 4, 2)
	// 10. Create the Assign that dish is cooked, includes bar
	assign := NewAssign("After kitchen", dish_state)
	// 11. Create the Checks for checking that this dishes for this table
	// id_table is a number of first table in group which served by waiter
	check_tb_number := func(id_table int) HandleCheckingFunc {
		return func(obj *Check, transact *Transaction) bool {
			for i := 0; i < 3; i++ {
				table_name := fmt.Sprintf("Table %d", i+id_table)
				if transact.GetParameter("Facility").(string) == table_name {
					return true
				}
			}
			return false
		}
	}
	check_tb := make([]IBaseObj, cnt_waiters)
	for i := 0; i < cnt_waiters; i++ {
		check_name := fmt.Sprintf("Is it order for table %d, %d, %d?", i*3+1, i*3+2, i*3+3)
		check_tb[i] = NewCheck(check_name, check_tb_number(i*3+1), nil)
	}
	// 12. Create the Advance for eating simulation
	visitors_eating := NewAdvance("Visitors eating", 45, 10)
	// 13. Create the Aggregate for aggregate dishes to order
	aggregate := NewAggregate("Aggregate dishes")
	// 14. Create the Advance for payment simulation
	visitors_pays := NewAdvance("Visitors pays", 5, 2)
	// 15. Append objects to a pipeline
	restaurant.Append(visitors_g, check_queue)
	restaurant.Append(check_queue, visitors_q)
	restaurant.Append(visitors_q, check_empty_table)
	restaurant.Append(check_empty_table, hostes1_f, hostes2_f)
	restaurant.AppendISlice(hostes1_f, tables_in)
	restaurant.AppendISlice(hostes2_f, tables_in)

	for i := 0; i < cnt_waiters; i++ {
		for j := 0; j < 3; j++ {
			restaurant.Append(tables_in[j+i*3], waiters_queue[i])
		}
	}

	for i := 0; i < cnt_waiters; i++ {
		restaurant.Append(waiters_queue[i], waiters_facility[i])
	}

	restaurant.AppendMultiple(waiters_facility, check_is_cooked)
	restaurant.Append(dishes_sp, cook1_q, cook2_q, cook3_q, cook4_q, bar_q)
	restaurant.Append(cook1_q, cook1_f)
	restaurant.Append(cook2_q, cook2_f)
	restaurant.Append(cook3_q, cook3_f)
	restaurant.Append(cook4_q, cook4_f)
	restaurant.Append(bar_q, barman1_f, barman2_f)
	restaurant.Append(cook1_f, assign)
	restaurant.Append(cook2_f, assign)
	restaurant.Append(cook3_f, assign)
	restaurant.Append(cook4_f, assign)
	restaurant.Append(barman1_f, assign)
	restaurant.Append(barman2_f, assign)
	restaurant.AppendISlice(assign, check_tb)
	for i := 0; i < cnt_waiters; i++ {
		restaurant.Append(check_tb[i], waiters_queue[i])
	}
	restaurant.Append(check_is_cooked, visitors_eating)
	restaurant.Append(visitors_eating, aggregate)
	restaurant.Append(aggregate, visitors_pays)
	restaurant.AppendISlice(visitors_pays, tables_out)
	restaurant.AppendMultiple(tables_out, out)
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
	restaurant.Report()
	fmt.Println("Exit program")
}
