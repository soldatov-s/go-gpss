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

func (obj *Aggregate) SendToDst(transact ITransaction) bool {
	for _, v := range obj.GetDst() {
		if v.AppendTransact(transact) {
			obj.tb.Remove(transact)
			obj.sum_transact++
			return true
		}
	}
	return false
}

func (obj *Aggregate) HandleTransact(transact ITransaction) {
	transact.PrintInfo()
	part, parts := transact.GetParts()
	if parts == 0 {
		if obj.SendToDst(transact) {
			return
		}
	}
	holded_tr := obj.tb.GetItem(transact.GetId())
	if holded_tr == nil {
		transact.SetParts(0, parts-1)
		if parts-1 == 0 {
			if obj.SendToDst(transact) {
				return
			}
		}
		obj.tb.Push(transact)
		return
	} else {
		holded_part, holded_parts := holded_tr.transact.GetParts()
		if holded_part != part {
			holded_tr.transact.SetParts(0, holded_parts-1)
			if holded_tr.transact.GetAdvanceTime() < transact.GetAdvanceTime() {
				holded_tr.transact.SetTiсks(transact.GetAdvanceTime())
				holded_tr.transact.SetTiсks(0)
			}
		}
		if holded_parts-1 == 0 {
			if obj.SendToDst(transact) {
				return
			}
		}

	}
}

func (obj *Aggregate) HandleTransacts(wg *sync.WaitGroup) {
	wg.Done()
	return
}

func (obj *Aggregate) AppendTransact(transact ITransaction) bool {
	obj.GetLogger().GetTrace().Println("Append transact ", transact.GetId(), " to Advance")
	transact.SetHolderName(obj.name)
	obj.HandleTransact(transact)
	return true
}

func (obj *Aggregate) PrintReport() {
	obj.BaseObj.PrintReport()
	fmt.Printf("Number of aggregated transact %.2f\n", obj.sum_transact)
	if obj.tb.GetLen() > 0 {
		fmt.Println("Await end aggregate:")
		for _, item := range obj.tb.GetItems() {
			_, parts := item.transact.GetParts()
			fmt.Printf("transact %d wait %d parts\n", item.transact.GetId(), parts)
		}
	}
	fmt.Println()
}
