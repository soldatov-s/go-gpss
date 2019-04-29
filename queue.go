// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package gpss

import (
	"fmt"
	"sync"
)

type IQueue interface {
	IsObjectAfterMeEmpty(transact ITransaction) bool // Check that after queue exist empty object
	GetLength() int                                  // Get queue length
}

type Queue struct {
	BaseObj
	sum_timequeue   float64 // Sum all transact queue time
	sum_zeroEntries float64 // Sum zero entrise
	sum_Entries     float64 // Sum all entries
	max_content     int     // Max content in queue
	sum_content     float64 // Sum content in queue
}

func NewQueue(name string) *Queue {
	obj := &Queue{}
	obj.BaseObj.Init(name)
	return obj
}

func (obj *Queue) HandleTransact(transact ITransaction) {
	transact.InqQueueTime()
	transact.PrintInfo()
}

// Check that after queue exist empty object
func (obj *Queue) IsObjectAfterMeEmpty(transact ITransaction) bool {
	for _, o := range obj.GetDst() {
		if o.AppendTransact(transact) {
			return true
		}
	}
	return false
}

// Get queue length
func (obj *Queue) GetLength() int {
	return obj.tb.GetLen()
}

func (obj *Queue) HandleTransacts(wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()
		// Queue not empty
		if obj.tb.GetLen() > 0 {
			tr := obj.tb.GetFirstItem()
			for tr != nil {
				if obj.IsObjectAfterMeEmpty(tr.transact) {
					obj.sum_timequeue += float64(tr.transact.GetQueueTime())
					obj.tb.Pop()
					tr = obj.tb.GetFirstItem()
				} else {
					break
				}
			}
		}
		transacts := obj.tb.GetItems()
		for _, tr := range transacts {
			obj.HandleTransact(tr.transact)
		}
		obj.sum_content += float64(obj.tb.GetLen())
	}()
}

func (obj *Queue) AppendTransact(transact ITransaction) bool {
	obj.GetLogger().GetTrace().Println("Append transact ", transact.GetId(), " to Queue")
	transact.SetHolderName(obj.name)
	if !obj.IsObjectAfterMeEmpty(transact) {
		transact.ResetQueueTime()
		obj.tb.Push(transact)
		transact.InqQueueTime()
		if obj.max_content < obj.tb.GetLen() {
			obj.max_content = obj.tb.GetLen()
		}
	} else {
		obj.sum_zeroEntries++
	}
	obj.sum_Entries++
	return true
}

func (obj *Queue) PrintReport() {
	obj.BaseObj.PrintReport()
	fmt.Println("Max content", obj.max_content)
	fmt.Println("Total entries", obj.sum_Entries)
	fmt.Println("Zero entries", obj.sum_zeroEntries)
	fmt.Printf("Persent zero entries %.2f%%\n", 100*obj.sum_zeroEntries/obj.sum_Entries)
	fmt.Println("Current contents", obj.tb.GetLen())
	fmt.Printf("Average content %.2f\n", obj.sum_content/float64(obj.GetPipeline().GetSimTime()))
	fmt.Printf("Average time/trans %.2f\n", obj.sum_timequeue/obj.sum_Entries)
	if obj.sum_Entries-obj.sum_zeroEntries > 0 {
		fmt.Printf("Average time/trans without zero entries %.2f\n", obj.sum_timequeue/(obj.sum_Entries-obj.sum_zeroEntries))
	}
	fmt.Println()
}
