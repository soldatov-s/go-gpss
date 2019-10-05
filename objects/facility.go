// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package objects

import (
	"fmt"
	"sync"

	"github.com/soldatov-s/go-gpss/internal"
)

// IFacility implements Facility interface
type IFacility interface {
	IsEmpty() bool
}

// Facility entity with advance in it
type Facility struct {
	BaseObj
	// The mean time increment
	Interval int
	// The time half-range
	Modificator int
	// Holded transast ID
	HoldedTransactID int
	// For backuping Facility/Bifacility name if we includes Facility in Bifacility
	bakupFacilityName string
	// For counting the advance of transact
	sumAdvance float64
	// For counting the transacts that go through Bifacility
	cntTransact float64
}

// NewFacility creates new Facility.
// name - name of object; interval - the mean time increment;
// modificator - the time half-range
func NewFacility(name string, interval, modificator int) *Facility {
	obj := &Facility{}
	obj.BaseObj.Init(name)
	obj.Interval = interval
	obj.Modificator = modificator
	obj.HoldedTransactID = -1
	return obj
}

// GenerateAdvance generate advance for facility
func (obj *Facility) GenerateAdvance() int {
	advance := obj.Interval
	if obj.Modificator > 0 {
		advance += utils.GetRandom(-obj.Modificator, obj.Modificator)
	}
	return advance
}

// HandleTransact handle transact
func (obj *Facility) HandleTransact(transact *Transaction) {
	transact.DecTiсks()
	transact.PrintInfo()
	if transact.IsTheEnd() {
		if obj.bakupFacilityName != "" {
			transact.SetParameter("Facility",
				obj.bakupFacilityName)
		} else {
			transact.SetParameter("Facility", nil)
		}
		for _, v := range obj.GetDst() {
			if v.AppendTransact(transact) {
				obj.tb.Remove(transact)
				obj.HoldedTransactID = -1
				return
			}
		}
		transact.SetParameter("Facility", obj.name)
	}
}

// HandleTransacts handle transacts in goroutine
func (obj *Facility) HandleTransacts(wg *sync.WaitGroup) {
	if obj.tb.Len() == 0 {
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
func (obj *Facility) AppendTransact(transact *Transaction) bool {
	if obj.tb.Len() != 0 {
		// Facility is busy
		return false
	}
	obj.BaseObj.AppendTransact(transact)
	transact.SetHolder(obj.name)
	advance := obj.GenerateAdvance()
	obj.sumAdvance += float64(advance)
	transact.SetTiсks(advance)
	if transact.GetParameter("Facility") != nil {
		obj.bakupFacilityName = transact.GetParameter("Facility").(string)
	}
	transact.SetParameter("Facility", obj.name)
	obj.HoldedTransactID = transact.GetID()
	obj.tb.Push(transact)
	obj.cntTransact++
	return true
}

// Report - print report about object
func (obj *Facility) Report() {
	obj.BaseObj.Report()
	avr := obj.sumAdvance / obj.cntTransact
	fmt.Printf("Average advance %.2f \tAverage utilization %.2f%%\tNumber entries %.2f \t", avr,
		100*avr*obj.cntTransact/float64(obj.Pipe.SimTime), obj.cntTransact)
	if obj.HoldedTransactID > 0 {
		fmt.Print("Transact ", obj.HoldedTransactID, " in facility")
		part, _, parentID := obj.tb.Item(obj.HoldedTransactID).transact.GetParts()
		if parentID > 0 {
			fmt.Print(", parent transact ", parentID, " part ", part)
		}
	} else {
		fmt.Print("Facility is empty")
	}
	fmt.Printf("\n\n")
}

// IsEmpty check that facility is empty
func (obj *Facility) IsEmpty() bool {
	if obj.tb.Len() != 0 {
		// Facility is busy
		return false
	}
	return true
}
