// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package gpss

import (
	"fmt"
	"sync"
)

type InFacility struct {
	BaseObj
	HoldedTransactID int
	cnt_transact     float64
	sum_advance      float64
	timeOfInput      int
}

type OutFacility struct {
	BaseObj
	inFacility *InFacility
	isOut      bool
}

func NewBifacility(name string) (*InFacility, *OutFacility) {
	inObj := &InFacility{}
	inObj.BaseObj.Init(name)
	outObj := &OutFacility{}
	outObj.name = name + "_OUT"
	outObj.tb = inObj.tb
	outObj.inFacility = inObj
	return inObj, outObj
}

func (obj *InFacility) HandleTransact(transact ITransaction) {
	transact.PrintInfo()
	for _, v := range obj.GetDst() {
		if v.AppendTransact(transact) {
			break
		}
	}
}

func (obj *InFacility) HandleTransacts(wg *sync.WaitGroup) {
	wg.Done()
	return
}

func (obj *InFacility) AppendTransact(transact ITransaction) bool {
	if obj.tb.GetLen() != 0 {
		// Facility is busy
		return false
	}
	obj.GetLogger().GetTrace().Println("Append transact ", transact.GetId(), " to Facility")
	transact.SetHolderName(obj.name)
	obj.HoldedTransactID = transact.GetId()
	obj.tb.Push(transact)
	obj.cnt_transact++
	obj.timeOfInput = obj.GetPipeline().GetModelTime()
	obj.HandleTransact(transact)
	return true
}

func (obj *InFacility) PrintReport() {
	obj.BaseObj.PrintReport()
	avr := obj.sum_advance / obj.cnt_transact
	fmt.Printf("Average advance %.2f\n", avr)
	fmt.Printf("Average utilization %.2f\n", 100*avr*obj.cnt_transact/float64(obj.GetPipeline().GetSimTime()))
	fmt.Printf("Number entries %.2f\n", obj.cnt_transact)
	if obj.HoldedTransactID != 0 {
		fmt.Println("Transact", obj.HoldedTransactID, "in facility")
	} else {
		fmt.Println("Facility is empty")
	}
	fmt.Println()
}

func (obj *OutFacility) HandleTransact(transact ITransaction) {
	transact.PrintInfo()
	for _, v := range obj.GetDst() {
		if v.AppendTransact(transact) {
			advance := obj.GetPipeline().GetModelTime() - obj.inFacility.timeOfInput
			obj.inFacility.sum_advance += float64(advance)
			obj.tb.Remove(transact)
		}
	}
}

func (obj *OutFacility) HandleTransacts(wg *sync.WaitGroup) {
	wg.Done()
	return
}

func (obj *OutFacility) AppendTransact(transact ITransaction) bool {
	if obj.inFacility.HoldedTransactID != transact.GetId() {
		return false
	}
	obj.GetLogger().GetTrace().Println("Append transact ", transact.GetId(), " to Facility")
	obj.HandleTransact(transact)
	if obj.tb.GetLen() == 0 {
		return true
	}
	return false
}

func (obj *OutFacility) PrintReport() {
	return
}
