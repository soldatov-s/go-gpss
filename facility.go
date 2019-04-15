// facility
package gpss

import (
	"fmt"
	"sync"
)

type IFacility interface {
	GenerateAdvance() int
}

type Facility struct {
	BaseObj
	Interval         int
	Modificator      int
	HoldedTransactID int
	sum_advance      float64
	cnt_transact     float64
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
				obj.tb.Remove(transact)
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
	PrintlnVerbose(obj.GetPipeline().IsVerbose(), "Append transact ", transact.GetId(), " to Facility")
	transact.SetHolderName(obj.name)
	advance := obj.GenerateAdvance()
	obj.sum_advance += float64(advance)
	transact.SetTiсks(advance)
	obj.HoldedTransactID = transact.GetId()
	obj.tb.Push(transact)
	obj.cnt_transact++
	return true
}

func (obj *Facility) PrintReport() {
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
