// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package objects

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
	sumTimequeue   float64 // Sum all transact queue time
	sumZeroEntries float64 // Sum zero entrise
	sumEntries     float64 // Sum all entries
	maxContent     int     // Max content in queue
	sumContent     float64 // Sum content in queue
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
					obj.sumTimequeue += float64(tr.transact.GetQueueTime())
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
		obj.sumContent += float64(obj.tb.Len())
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
		if obj.maxContent < obj.tb.Len() {
			obj.maxContent = obj.tb.Len()
		}
	} else {
		obj.sumZeroEntries++
	}
	obj.sumEntries++
	return true
}

// Report - print report about object
func (obj *Queue) Report() {
	obj.BaseObj.Report()
	fmt.Printf("Max content \t%d\tTotal entries \t%2.f\tZero entries \t%2.f\tPersent zero entries \t%.2f%%\n",
		obj.maxContent, obj.sumEntries, obj.sumZeroEntries, 100*obj.sumZeroEntries/obj.sumEntries)
	fmt.Printf("Current contents \t%d\tAverage content \t%.2f\tAverage time/trans \t%.2f\n", obj.tb.Len(),
		obj.sumContent/float64(obj.Pipe.SimTime), obj.sumTimequeue/obj.sumEntries)
	if obj.sumEntries-obj.sumZeroEntries > 0 {
		fmt.Printf("Average time/trans without zero entries \t%.2f\n", obj.sumTimequeue/(obj.sumEntries-obj.sumZeroEntries))
	}
	fmt.Println()
}
