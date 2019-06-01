// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package gpss

// A Bifacility as facility, but without advance in it and present in two parts,
// first for takes ownership of a Facility, second for release ownership of a Facility

import (
	"fmt"
)

// The first part of a Bifacility, it takes ownership of a Facility
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

// The second part of a Bifacility, for release ownership of a Facility
type OutFacility struct {
	BaseObj
	// Pointer to inFacility structure
	inFacility *InFacility
}

// Creates new Bifacility (InFacility + OutFacility).
// name - name of object
func NewBifacility(name string) (*InFacility, *OutFacility) {
	inObj := &InFacility{}
	inObj.BaseObj.Init(name)
	outObj := &OutFacility{}
	outObj.name = name + "_OUT"
	outObj.tb = inObj.tb
	outObj.inFacility = inObj
	return inObj, outObj
}

func (obj *InFacility) HandleTransact(transact *Transaction) {
	transact.PrintInfo()
	for _, v := range obj.GetDst() {
		if v.AppendTransact(transact) {
			return
		}
	}
}

func (obj *InFacility) AppendTransact(transact *Transaction) bool {
	if obj.tb.Len() != 0 {
		// Facility is busy
		return false
	}
	Logger.Trace.Println("Append transact ", transact.GetID(), " to Facility")
	transact.SetHolder(obj.name)
	if transact.GetParameter("Facility") != nil {
		obj.bakupFacilityName = transact.GetParameter("Facility").(string)
	}
	transact.SetParameter("Facility", obj.name)
	obj.HoldedTransactID = transact.GetID()
	obj.tb.Push(transact)
	obj.cnt_transact++
	obj.timeOfInput = obj.GetPipeline().ModelTime
	obj.HandleTransact(transact)
	return true
}

func (obj *InFacility) PrintReport() {
	obj.BaseObj.PrintReport()
	avr := obj.sum_advance / obj.cnt_transact
	fmt.Printf("Average advance %.2f \tAverage utilization %.2f%%\tNumber entries %.2f \t", avr,
		100*avr*obj.cnt_transact/float64(obj.GetPipeline().SimTime), obj.cnt_transact)
	if obj.HoldedTransactID > 0 {
		fmt.Print("Transact ", obj.HoldedTransactID, " in facility")
	} else {
		fmt.Print("Facility is empty")
	}
	fmt.Printf("\n\n")
}

func (obj *InFacility) IsEmpty() bool {
	if obj.tb.Len() != 0 {
		// Facility is busy
		return false
	}
	return true
}

func (obj *OutFacility) HandleTransact(transact *Transaction) {
	transact.PrintInfo()
	if obj.inFacility.bakupFacilityName != "" {
		transact.SetParameter("Facility",
			obj.inFacility.bakupFacilityName)
	} else {
		transact.SetParameter("Facility", nil)
	}

	for _, v := range obj.GetDst() {
		if v.AppendTransact(transact) {
			advance := obj.GetPipeline().ModelTime - obj.inFacility.timeOfInput
			obj.inFacility.sum_advance += float64(advance)
			obj.tb.Remove(transact)
			obj.inFacility.HoldedTransactID = -1
			return
		}
	}
	transact.SetParameters([]Parameter{{Name: "Facility", Value: obj.name}})
}

func (obj *OutFacility) AppendTransact(transact *Transaction) bool {
	if obj.inFacility.HoldedTransactID != transact.GetID() {
		return false
	}
	Logger.Trace.Println("Append transact ", transact.GetID(), " to Facility")
	obj.HandleTransact(transact)
	if obj.tb.Len() == 0 {
		return true
	}
	return false
}

func (obj *OutFacility) PrintReport() {
	return
}
