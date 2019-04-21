// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package gpss

import (
	"fmt"
	"sync"
)

type IAdvance interface {
	GenerateAdvance() int
}

type Advance struct {
	BaseObj
	Interval     int
	Modificator  int
	sum_advance  float64
	sum_transact float64
}

func NewAdvance(name string, interval, modificator int) *Advance {
	obj := &Advance{}
	obj.BaseObj.Init(name)
	obj.Interval = interval
	obj.Modificator = modificator
	return obj
}

func (obj *Advance) GenerateAdvance() int {
	advance := obj.Interval
	if obj.Modificator > 0 {
		advance += GetRandom(-obj.Modificator, obj.Modificator)
	}
	return advance
}

func (obj *Advance) HandleTransact(transact ITransaction) {
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

func (obj *Advance) HandleTransacts(wg *sync.WaitGroup) {
	if obj.Interval == 0 ||
		obj.tb.GetLen() == 0 {
		wg.Done()
		return
	}
	go func() {
		defer wg.Done()
		transacts := obj.tb.GetItems()
		for _, tr := range transacts {
			obj.HandleTransact(tr.transact)
		}
	}()
}

func (obj *Advance) AppendTransact(transact ITransaction) bool {
	obj.GetLogger().GetTrace().Println("Append transact ", transact.GetId(), " to Advance")
	transact.SetHolderName(obj.name)
	advance := obj.GenerateAdvance()
	obj.sum_advance += float64(advance)
	transact.SetTiсks(advance)
	obj.tb.Push(transact)
	obj.sum_transact++
	return true
}

func (obj *Advance) PrintReport() {
	obj.BaseObj.PrintReport()
	fmt.Printf("Average advance %.2f\n", obj.sum_advance/obj.sum_transact)
	fmt.Println()
}
