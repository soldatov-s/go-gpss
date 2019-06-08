// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package gpss

import (
	"fmt"
	"sync"
)

// Aggregate multiple sub-transactions in Transaction
type Aggregate struct {
	BaseObj
	sumTransact float64 // Counter of all fully aggregated transactions
}

// NewAggregate creates new Aggregate
// name - name of object
func NewAggregate(name string) *Aggregate {
	obj := &Aggregate{}
	obj.BaseObj.Init(name)
	return obj
}

// SendToDst - send transact to sedtination
func (obj *Aggregate) SendToDst(transact *Transaction) bool {
	for _, v := range obj.GetDst() {
		if v.AppendTransact(transact) {
			obj.tb.Remove(transact)
			obj.sumTransact++
			return true
		}
	}
	return false
}

// HandleTransact handle transact
func (obj *Aggregate) HandleTransact(transact *Transaction) bool {
	transact.PrintInfo()
	_, parts, parentID := transact.GetParts()
	if parentID == 0 {
		return obj.SendToDst(transact)
	}
	holdedTr := obj.tb.Item(parentID)
	if holdedTr == nil {
		tr := transact.Copy()
		tr.SetID(parentID)
		tr.SetParts(0, parts-1, 0)
		if parts-1 == 0 {
			return obj.SendToDst(tr)
		}
		obj.tb.Push(tr)
	} else {
		// Update Advance
		if holdedTr.transact.GetAdvanceTime() < transact.GetAdvanceTime() {
			holdedTr.transact.SetTiсks(transact.GetAdvanceTime())
			holdedTr.transact.SetTiсks(0)
		}
		_, holdedParts, _ := holdedTr.transact.GetParts()
		if holdedParts-1 == 0 {
			// We aggregate all parts
			holdedTr.transact.SetParts(0, 0, 0)
			return obj.SendToDst(holdedTr.transact)
		}
		holdedTr.transact.SetParts(0, holdedParts-1, 0)
	}
	return true
}

// HandleTransacts handle transacts
func (obj *Aggregate) HandleTransacts(wg *sync.WaitGroup) {
	wg.Done()
	return
}

// AppendTransact append transact to object
func (obj *Aggregate) AppendTransact(transact *Transaction) bool {
	obj.BaseObj.AppendTransact(transact)
	transact.SetHolder(obj.name)
	return obj.HandleTransact(transact)
}

// Report - print report about object
func (obj *Aggregate) Report() {
	obj.BaseObj.Report()
	fmt.Printf("Number of aggregated transact %.2f\n", obj.sumTransact)
	if obj.tb.Len() > 0 {
		fmt.Println("Await end aggregate:")
		for _, item := range obj.tb.Items() {
			_, parts, _ := item.transact.GetParts()
			fmt.Printf("transact %d wait %d parts\n", item.transact.GetID(), parts)
		}
	}
	fmt.Println()
}
