// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package gpss

import (
	"fmt"
)

type InFacility struct {
	BaseObj
	// Holded transast ID
	HoldedTransactID int
	// For backuping Facility/Bifacility name if we includes Bifacility in Bifacility
	bakupFacilityName string
	// For counting the transacts that go through Bifacility
	cnt_transact float64
	// For counting the advance of transact
	sum_advance float64
	// For saving time of input transact in Bifacility
	timeOfInput int
}

type OutFacility struct {
	BaseObj
	// Pointer to inFacility structure
	inFacility *InFacility
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

func (obj *InFacility) AppendTransact(transact ITransaction) bool {
	if obj.tb.GetLen() != 0 {
		// Facility is busy
		return false
	}
	obj.GetLogger().GetTrace().Println("Append transact ", transact.GetId(), " to Facility")
	transact.SetHolderName(obj.name)
	transact.SetParameters([]Parameter{{Name: "Facility", Value: obj.name}})
	obj.HoldedTransactID = transact.GetId()
	obj.bakupFacilityName = transact.GetParameterByName("Facility").(string)
	obj.tb.Push(transact)
	obj.cnt_transact++
	obj.timeOfInput = obj.GetPipeline().GetModelTime()
	obj.HandleTransact(transact)
	return true
}

func (obj *InFacility) PrintReport() {
	obj.BaseObj.PrintReport()
	avr := obj.sum_advance / obj.cnt_transact
	fmt.Printf("Average advance %.2f \tAverage utilization %.2f%%\tNumber entries %.2f \t", avr,
		100*avr*obj.cnt_transact/float64(obj.GetPipeline().GetSimTime()), obj.cnt_transact)
	if obj.HoldedTransactID > 0 {
		fmt.Print("Transact ", obj.HoldedTransactID, " in facility")
	} else {
		fmt.Print("Facility is empty")
	}
	fmt.Printf("\n\n")
}

func (obj *InFacility) IsEmpty() bool {
	if obj.tb.GetLen() != 0 {
		// Facility is busy
		return false
	}
	return true
}

func (obj *OutFacility) HandleTransact(transact ITransaction) {
	transact.PrintInfo()
	for _, v := range obj.GetDst() {
		if v.AppendTransact(transact) {
			advance := obj.GetPipeline().GetModelTime() - obj.inFacility.timeOfInput
			obj.inFacility.sum_advance += float64(advance)
			if obj.inFacility.bakupFacilityName != "" {
				transact.SetParameters([]Parameter{{Name: "Facility",
					Value: obj.inFacility.bakupFacilityName}})
			} else {
				transact.SetParameters([]Parameter{{Name: "Facility", Value: nil}})
			}
			obj.tb.Remove(transact)
		}
	}
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
