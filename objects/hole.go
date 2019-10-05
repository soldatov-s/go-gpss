// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package objects

import (
	"fmt"
	"sync"
)

// Hole in which fall in transactions
type Hole struct {
	BaseObj
	sumLife     float64 // For count average transact life
	sumAdvance  float64 // For count average advance
	cntTransact float64 // How much killed
}

// NewHole creates new Hole
// name - name of object
func NewHole(name string) *Hole {
	obj := &Hole{}
	obj.BaseObj.Init(name)
	return obj
}

// HandleTransact handle transact
func (obj *Hole) HandleTransact(transact *Transaction) {
	if !transact.IsKilled() {
		transact.Kill()
		transact.PrintInfo()
		obj.sumLife += float64(transact.GetLife())
		obj.sumAdvance += float64(transact.GetAdvanceTime())
		obj.cntTransact++
	}
}

// HandleTransacts handle transacts in goroutine
func (obj *Hole) HandleTransacts(wg *sync.WaitGroup) {
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
func (obj *Hole) AppendTransact(transact *Transaction) bool {
	obj.BaseObj.AppendTransact(transact)
	transact.SetHolder(obj.name)
	obj.tb.Push(transact)
	return true
}

// Report - print report about object
func (obj *Hole) Report() {
	obj.BaseObj.Report()
	fmt.Println("Killed", obj.cntTransact)
	fmt.Printf("Average advance %.2f\n", obj.sumAdvance/obj.cntTransact)
	fmt.Printf("Average life %.2f\n", obj.sumLife/obj.cntTransact)
	fmt.Println()
}
