// hole
package gpss

import (
	"fmt"
	"sync"
)

type IHole interface {
}

type Hole struct {
	BaseObj
	sum_life     float64
	sum_advance  float64
	cnt_transact float64 // How much killed
}

func NewHole(name string) *Hole {
	obj := &Hole{}
	obj.BaseObj.Init(name)
	return obj
}

func (obj *Hole) HandleTransact(transact ITransaction) {
	if !transact.IsKilled() {
		transact.Kill()
		transact.PrintInfo()
		obj.sum_life += float64(transact.GetLife())
		obj.sum_advance += float64(transact.GetAdvanceTime())
		obj.cnt_transact++
	}
}

func (obj *Hole) HandleTransacts(wg *sync.WaitGroup) {
	if obj.tb.GetLen() == 0 {
		wg.Done()
		return
	}
	go func() {
		defer wg.Done()
		transacts := obj.tb.GetItems()
		defer obj.tb.UnlockTable()
		for _, tr := range transacts {
			obj.HandleTransact(tr.transact)
		}
	}()
}

func (obj *Hole) AppendTransact(transact ITransaction) bool {
	PrintlnVerbose(obj.GetPipeline().IsVerbose(), "Append transact ", transact.GetId(), " to Hole")
	transact.SetHolderName(obj.name)
	obj.tb.Push(transact)
	return true
}

func (obj *Hole) PrintReport() {
	obj.BaseObj.PrintReport()
	fmt.Println("Killed", obj.cnt_transact)
	fmt.Printf("Average advance %.2f\n", obj.sum_advance/obj.cnt_transact)
	fmt.Printf("Average life %.2f\n", obj.sum_life/obj.cnt_transact)
	fmt.Println()
}
