// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package gpss

import (
	"fmt"
	"sync"
)

// IQueue implements Queue interface
type IQueue interface {
	IsObjectAfterMeEmpty(transact *Transaction) bool // Check that after queue exist empty object
	GetLength() int                                  // Get queue length
}

// Queue of transaction
type Queue struct {
	BaseObj
	sum_timequeue   float64 // Sum all transact queue time
	sum_zeroEntries float64 // Sum zero entrise
	sum_Entries     float64 // Sum all entries
	max_content     int     // Max content in queue
	sum_content     float64 // Sum content in queue
}

// NewQueue creates new Queue.
// name - name of object
func NewQueue(name string) *Queue {
	obj := &Queue{}
	obj.BaseObj.Init(name)
	return obj
}

// HandleTransact handle transact
func (obj *Queue) HandleTransact(transact *Transaction) {
	transact.InqQueueTime()
	transact.PrintInfo()
}

// IsObjectAfterMeEmpty check that after queue exist empty object
func (obj *Queue) IsObjectAfterMeEmpty(transact *Transaction) bool {
	for _, o := range obj.GetDst() {
		if o.AppendTransact(transact) {
			return true
		}
	}
	return false
}

// GetLength get queue length
func (obj *Queue) GetLength() int {
	return obj.tb.Len()
}

// HandleTransacts handle transacts in goroutine
func (obj *Queue) HandleTransacts(wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()
		// Queue not empty
		if obj.tb.Len() > 0 {
			tr := obj.tb.First()
			for tr != nil {
				if obj.IsObjectAfterMeEmpty(tr.transact) {
					obj.sum_timequeue += float64(tr.transact.GetQueueTime())
					obj.tb.Pop()
					tr = obj.tb.First()
				} else {
					break
				}
			}
		}
		transacts := obj.tb.Items()
		for _, tr := range transacts {
			obj.HandleTransact(tr.transact)
		}
		obj.sum_content += float64(obj.tb.Len())
	}()
}

// AppendTransact append transact to object
func (obj *Queue) AppendTransact(transact *Transaction) bool {
	obj.BaseObj.AppendTransact(transact)
	transact.SetHolder(obj.name)
	if !obj.IsObjectAfterMeEmpty(transact) {
		transact.ResetQueueTime()
		obj.tb.Push(transact)
		transact.InqQueueTime()
		if obj.max_content < obj.tb.Len() {
			obj.max_content = obj.tb.Len()
		}
	} else {
		obj.sum_zeroEntries++
	}
	obj.sum_Entries++
	return true
}

// Report - print report about object
func (obj *Queue) Report() {
	obj.BaseObj.Report()
	fmt.Printf("Max content \t%d\tTotal entries \t%2.f\tZero entries \t%2.f\tPersent zero entries \t%.2f%%\n",
		obj.max_content, obj.sum_Entries, obj.sum_zeroEntries, 100*obj.sum_zeroEntries/obj.sum_Entries)
	fmt.Printf("Current contents \t%d\tAverage content \t%.2f\tAverage time/trans \t%.2f\n", obj.tb.Len(),
		obj.sum_content/float64(obj.Pipe.SimTime), obj.sum_timequeue/obj.sum_Entries)
	if obj.sum_Entries-obj.sum_zeroEntries > 0 {
		fmt.Printf("Average time/trans without zero entries \t%.2f\n", obj.sum_timequeue/(obj.sum_Entries-obj.sum_zeroEntries))
	}
	fmt.Println()
}
