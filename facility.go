// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package gpss

import (
	"fmt"
	"sync"
)

type IFacility interface {
	IsEmpty() bool
}

type Facility struct {
	BaseObj
	Interval    int
	Modificator int
	// Holded transast ID
	HoldedTransactID int
	// For backuping Facility/Bifacility name if we includes Facility in Bifacility
	bakupFacilityName string
	// For counting the advance of transact
	sum_advance float64
	// For counting the transacts that go through Bifacility
	cnt_transact float64
}

func NewFacility(name string, interval, modificator int) *Facility {
	obj := &Facility{}
	obj.BaseObj.Init(name)
	obj.Interval = interval
	obj.Modificator = modificator
	obj.HoldedTransactID = -1
	return obj
}

func (obj *Facility) GenerateAdvance() int {
	advance := obj.Interval
	if obj.Modificator > 0 {
		advance += GetRandom(-obj.Modificator, obj.Modificator)
	}
	return advance
}

func (obj *Facility) HandleTransact(transact ITransaction) {
	transact.DecTiсks()
	transact.PrintInfo()
	if transact.IsTheEnd() {
		for _, v := range obj.GetDst() {
			if v.AppendTransact(transact) {
				transact.SetParameters([]Parameter{{Name: "Facility", Value: nil}})
				if obj.bakupFacilityName != "" {
					transact.SetParameters([]Parameter{{Name: "Facility",
						Value: obj.bakupFacilityName}})
				} else {
					transact.SetParameters([]Parameter{{Name: "Facility", Value: nil}})
				}
				obj.tb.Remove(transact)
				obj.HoldedTransactID = -1
				break
			}
		}
	}
}
func (obj *Facility) HandleTransacts(wg *sync.WaitGroup) {
	if obj.tb.GetLen() == 0 {
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

func (obj *Facility) AppendTransact(transact ITransaction) bool {
	if obj.tb.GetLen() != 0 {
		// Facility is busy
		return false
	}
	obj.GetLogger().GetTrace().Println("Append transact ", transact.GetId(), " to Facility")
	transact.SetHolderName(obj.name)
	advance := obj.GenerateAdvance()
	obj.sum_advance += float64(advance)
	transact.SetTiсks(advance)
	transact.SetParameters([]Parameter{{Name: "Facility", Value: obj.name}})
	obj.HoldedTransactID = transact.GetId()
	obj.bakupFacilityName = transact.GetParameterByName("Facility").(string)
	obj.tb.Push(transact)
	obj.cnt_transact++
	return true
}

func (obj *Facility) PrintReport() {
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

func (obj *Facility) IsEmpty() bool {
	if obj.tb.GetLen() != 0 {
		// Facility is busy
		return false
	}
	return true
}
