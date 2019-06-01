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
	sum_transact float64 // Counter of all fully aggregated transactions
}

// Creates new Aggregate
// name - name of object
func NewAggregate(name string) *Aggregate {
	obj := &Aggregate{}
	obj.BaseObj.Init(name)
	return obj
}

func (obj *Aggregate) SendToDst(transact *Transaction) bool {
	for _, v := range obj.GetDst() {
		if v.AppendTransact(transact) {
			obj.tb.Remove(transact)
			obj.sum_transact++
			return true
		}
	}
	return false
}

func (obj *Aggregate) HandleTransact(transact *Transaction) bool {
	transact.PrintInfo()
	_, parts, parent_id := transact.GetParts()
	if parent_id == 0 {
		return obj.SendToDst(transact)
	}
	holded_tr := obj.tb.Item(parent_id)
	if holded_tr == nil {
		tr := transact.Copy()
		tr.SetID(parent_id)
		tr.SetParts(0, parts-1, 0)
		if parts-1 == 0 {
			return obj.SendToDst(tr)
		}
		obj.tb.Push(tr)
	} else {
		// Update Advance
		if holded_tr.transact.GetAdvanceTime() < transact.GetAdvanceTime() {
			holded_tr.transact.SetTiсks(transact.GetAdvanceTime())
			holded_tr.transact.SetTiсks(0)
		}
		_, holded_parts, _ := holded_tr.transact.GetParts()
		if holded_parts-1 == 0 {
			// We aggregate all parts
			holded_tr.transact.SetParts(0, 0, 0)
			return obj.SendToDst(holded_tr.transact)
		} else {
			holded_tr.transact.SetParts(0, holded_parts-1, 0)
		}
	}
	return true
}

func (obj *Aggregate) HandleTransacts(wg *sync.WaitGroup) {
	wg.Done()
	return
}

func (obj *Aggregate) AppendTransact(transact *Transaction) bool {
	Logger.Trace.Println("Append transact ", transact.GetID(), " to Aggregate")
	transact.SetHolder(obj.name)
	return obj.HandleTransact(transact)
}

func (obj *Aggregate) PrintReport() {
	obj.BaseObj.PrintReport()
	fmt.Printf("Number of aggregated transact %.2f\n", obj.sum_transact)
	if obj.tb.Len() > 0 {
		fmt.Println("Await end aggregate:")
		for _, item := range obj.tb.Items() {
			_, parts, _ := item.transact.GetParts()
			fmt.Printf("transact %d wait %d parts\n", item.transact.GetID(), parts)
		}
	}
	fmt.Println()
}
