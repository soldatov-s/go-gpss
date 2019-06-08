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
	restaurant := gpss.NewPipeline("Restaurant  Simulation")
	// 1. Create the Generator and Queue of Visitors, create a Hole
	visitorsG := gpss.NewGenerator("Visitors", 10, 5, 0, 0, nil)
	// visitors_g := NewGenerator("Visitors", 0, 0, 0, 1, nil)
	out := gpss.NewHole("Out")
	visitorsQ := gpss.NewQueue("Visitors queue")
	// 2. Create the Check for checking size the Queue of Visitors
	CheckQueueHndl := func(obj *gpss.Check, transact *gpss.Transaction) bool {
		queue := obj.Pipe.GetObjByName("Visitors queue")
		if queue.(gpss.IQueue).GetLength() >= 6 {
			return false
		}
		return true
	}
	check_queue := gpss.NewCheck("Check size of Visitors queue", CheckQueueHndl, out)
	// 3. Create are Hostess
	hostes1F := gpss.NewFacility("hostess 1", 5, 3)
	hostes2F := gpss.NewFacility("Hostess 2", 5, 3)
	// 4. Create are Tables
	cntTables := 24
	tablesIN := make([]gpss.IBaseObj, cntTables)
	tablesOUT := make([]gpss.IBaseObj, cntTables)
	for i := 0; i < cntTables; i++ {
		table_name := fmt.Sprintf("Table %d", i+1)
		tablesIN[i], tablesOUT[i] = gpss.NewBifacility(table_name)
	}
	// 5. Check that we have empty table
	CheckEmptyTableHndl := func(obj *gpss.Check, transact *gpss.Transaction) bool {
		for i := 0; i < cntTables; i++ {
			table_name := fmt.Sprintf("Table %d", i+1)
			table := obj.Pipe.GetObjByName(table_name).(gpss.IFacility)
			if table.IsEmpty() {
				return true
			}
		}
		return false
	}
	checkEmptyTable := gpss.NewCheck("Check empty table", CheckEmptyTableHndl, nil)
	// 6. Create the Queues and Facilities for waiters
	cntWaiters := 8
	waitersQueue := make([]gpss.IBaseObj, cntWaiters)
	waitersFacility := make([]gpss.IBaseObj, cntWaiters)
	for i := 0; i < cntWaiters; i++ {
		queueName := fmt.Sprintf("Queue waiter %d", i+1)
		waitersQueue[i] = gpss.NewQueue(queueName)
		facilityName := fmt.Sprintf("Waiter %d", i+1)
		waitersFacility[i] = gpss.NewFacility(facilityName, 5, 3)

	}
	// 7. Create the Split for splitting Visitors order to dishes
	// Maybe 1 or 5 dishes, includes bar
	dishesSP := gpss.NewSplit("Selected dishes", 3, 2, nil)
	// 8. Create the Check that transact is a coocked dishes. If false, this transact is
	// an order from tables, needs split to dishes
	dishState := gpss.Parameter{Name: "State", Value: "After kitchen"}
	checkIsCooked := gpss.NewCheck("After kitchen?", nil, dishesSP, dishState)
	// 9. Create the Queue and Facility for cooks and barmans. Each cook cooking
	// only one type of dishes: meat, sushi, salats, dessert
	cook1Q := gpss.NewQueue("Queue of orders to cook 1 (meat dishes)")
	cook2Q := gpss.NewQueue("Queue of orders to cook 2 (fish dishes)")
	cook3Q := gpss.NewQueue("Queue of orders to cook 3 (salats)")
	cook4Q := gpss.NewQueue("Queue of orders to cook 4 (dessert)")
	barQ := gpss.NewQueue("Queue of orders to bar")
	cook1Q := gpss.NewFacility("Cook 1 (meat dishes)", 15, 5)
	cook2Q := gpss.NewFacility("Cook 2 (sushi)", 7, 3)
	cook3F := gpss.NewFacility("Cook 3 (salats)", 10, 4)
	cook4F := gpss.NewFacility("Cook 4 (dessert)", 20, 5)
	barman1F := gpss.NewFacility("Barman 1", 4, 2)
	barman2F := gpss.NewFacility("Barman 2", 4, 2)
	// 10. Create the Assign that dish is cooked, includes bar
	assign := gpss.NewAssign("After kitchen", dishState)
	// 11. Create the Checks for checking that this dishes for this table
	// id_table is a number of first table in group which served by waiter
	checkTbNumber := func(id_table int) gpss.HandleCheckingFunc {
		return func(obj *gpss.Check, transact *gpss.Transaction) bool {
			for i := 0; i < 3; i++ {
				tableName := fmt.Sprintf("Table %d", i+id_table)
				if transact.GetParameter("Facility").(string) == tableName {
					return true
				}
			}
			return false
		}
	}
	checkTb := make([]gpss.IBaseObj, cntWaiters)
	for i := 0; i < cntWaiters; i++ {
		checkName := fmt.Sprintf("Is it order for table %d, %d, %d?", i*3+1, i*3+2, i*3+3)
		checkTb[i] = gpss.NewCheck(checkName, checkTbNumber(i*3+1), nil)
	}
	// 12. Create the Advance for eating simulation
	visitors_eating := gpss.NewAdvance("Visitors eating", 45, 10)
	// 13. Create the Aggregate for aggregate dishes to order
	aggregate := gpss.NewAggregate("Aggregate dishes")
	// 14. Create the Advance for payment simulation
	visitors_pays := gpss.NewAdvance("Visitors pays", 5, 2)
	// 15. Append objects to a pipeline
	restaurant.Append(visitorsG, check_queue)
	restaurant.Append(check_queue, visitorsQ)
	restaurant.Append(visitorsQ, checkEmptyTable)
	restaurant.Append(checkEmptyTable, hostes1F, hostes2F)
	restaurant.AppendISlice(hostes1F, tablesIN)
	restaurant.AppendISlice(hostes2F, tablesIN)

	for i := 0; i < cntWaiters; i++ {
		for j := 0; j < 3; j++ {
			restaurant.Append(tablesIN[j+i*3], waitersQueue[i])
		}
	}

	for i := 0; i < cntWaiters; i++ {
		restaurant.Append(waitersQueue[i], waitersFacility[i])
	}

	restaurant.AppendMultiple(waitersFacility, checkIsCooked)
	restaurant.Append(dishesSP, cook1Q, cook2Q, cook3Q, cook4Q, barQ)
	restaurant.Append(cook1Q, cook1Q)
	restaurant.Append(cook2Q, cook2Q)
	restaurant.Append(cook3Q, cook3F)
	restaurant.Append(cook4Q, cook4F)
	restaurant.Append(barQ, barman1F, barman2F)
	restaurant.Append(cook1Q, assign)
	restaurant.Append(cook2Q, assign)
	restaurant.Append(cook3F, assign)
	restaurant.Append(cook4F, assign)
	restaurant.Append(barman1F, assign)
	restaurant.Append(barman2F, assign)
	restaurant.AppendISlice(assign, checkTb)
	for i := 0; i < cntWaiters; i++ {
		restaurant.Append(checkTb[i], waitersQueue[i])
	}
	restaurant.Append(checkIsCooked, visitors_eating)
	restaurant.Append(visitors_eating, aggregate)
	restaurant.Append(aggregate, visitors_pays)
	restaurant.AppendISlice(visitors_pays, tablesOUT)
	restaurant.AppendMultiple(tablesOUT, out)
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
