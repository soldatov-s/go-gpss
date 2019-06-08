// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package gpss

import (
	"fmt"
	"sync"
)

// IAdvance implements Advance interface
type IAdvance interface {
	GenerateAdvance() int
}

// Advance block delays the progress of a Transaction for a specified amount
// of simulated time
type Advance struct {
	BaseObj
	Interval     int     // The mean time increment
	Modificator  int     // The time half-range
	sum_advance  float64 // Totalize advance for all transacts
	sum_transact float64 // Counter of transacts
}

// NewAdvance creates new Advance.
// name - name of object; interval - the mean time increment;
// modificator - the time half-range
func NewAdvance(name string, interval, modificator int) *Advance {
	obj := &Advance{}
	obj.BaseObj.Init(name)
	obj.Interval = interval
	obj.Modificator = modificator
	return obj
}

// GenerateAdvance generate advance
func (obj *Advance) GenerateAdvance() int {
	advance := obj.Interval
	if obj.Modificator > 0 {
		advance += GetRandom(-obj.Modificator, obj.Modificator)
	}
	return advance
}

// HandleTransact handle transact
func (obj *Advance) HandleTransact(transact *Transaction) {
	transact.DecTiсks()
	transact.PrintInfo()
	if transact.IsTheEnd() {
		for _, v := range obj.GetDst() {
			if v.AppendTransact(transact) {
				obj.tb.Remove(transact)
				break
			}
		}
	}
}

// HandleTransacts handle transacts in goroutine
func (obj *Advance) HandleTransacts(wg *sync.WaitGroup) {
	if obj.Interval == 0 ||
		obj.tb.Len() == 0 {
		wg.Done()
		return
	}
	go func() {
		defer wg.Done()
		transacts := obj.tb.Items()
		for _, tr := range transacts {
			obj.HandleTransact(tr.transact)
		}
	}()
}

// AppendTransact append transact to object
func (obj *Advance) AppendTransact(transact *Transaction) bool {
	obj.BaseObj.AppendTransact(transact)
	transact.SetHolder(obj.name)
	advance := obj.GenerateAdvance()
	obj.sum_advance += float64(advance)
	transact.SetTiсks(advance)
	obj.tb.Push(transact)
	obj.sum_transact++
	return true
}

// Report - print report about object
func (obj *Advance) Report() {
	obj.BaseObj.Report()
	fmt.Printf("Average advance %.2f\n", obj.sum_advance/obj.sum_transact)
	fmt.Println()
}
