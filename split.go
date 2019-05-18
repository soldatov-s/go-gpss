// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package gpss

import (
	"fmt"
	"sync"
)

type HandleSplittingFunc func(obj *Split, transact ITransaction)

// A Split creates assembly set of sub-transactions of a Transaction
type Split struct {
	BaseObj
	Cntsplit        int                 // Number of related Transactions to be created
	Modificator     int                 // The count half-range
	sum_split       float64             // Counter of sub-transactions
	sum_transact    float64             // Counter of transactions
	HandleSplitting HandleSplittingFunc // Function for splitting transaction
}

// Default splitting function
func Splitting(obj *Split, transact ITransaction) {
	cntsplit := obj.Cntsplit
	if obj.Modificator > 0 {
		cntsplit += GetRandom(-obj.Modificator, obj.Modificator)
	}

	if cntsplit <= 0 {
		cntsplit = 1
	}
	if cntsplit > len(obj.GetDst()) {
		cntsplit = len(obj.GetDst())
	}

	obj.sum_split += float64(cntsplit)
	if cntsplit == len(obj.GetDst()) {
		// Default case, cntsplit equal to length of GetDst()
		for i, v := range obj.GetDst() {
			tr := transact.Copy()
			parent_id := tr.GetId()
			tr.SetID(obj.GetPipeline().GetIDNewTransaction())
			tr.SetParts(i+1, cntsplit, parent_id)
			v.AppendTransact(tr) // Take in mind that after Split must be only Queues
		}
	} else {
		// Another case, cntsplit can be smaller than length of GetDst()
		// Randomized selections of dst for send transact
		dsts := make([]bool, cntsplit)
		part_id := 1
		for {
			for _, v := range obj.GetDst() {
				if GetRandomBool() && !dsts[part_id-1] {
					tr := transact.Copy()
					parent_id := tr.GetId()
					tr.SetID(obj.GetPipeline().GetIDNewTransaction())
					tr.SetParts(part_id, cntsplit, parent_id)
					v.AppendTransact(tr)
					dsts[part_id-1] = true
					part_id++
					if part_id > cntsplit {
						return
					}
				}
			}
		}
	}
}

// Creates new Split.
// name - name of object; cntsplit - number of related Transactions to be created;
// modificator - the count half-range; hndl - function for splitting transaction
func NewSplit(name string, cntsplit, modificator int, hndl HandleSplittingFunc) *Split {
	obj := &Split{}
	obj.name = name
	obj.Cntsplit = cntsplit
	obj.Modificator = modificator
	if hndl != nil {
		obj.HandleSplitting = hndl
	} else {
		obj.HandleSplitting = Splitting
	}
	return obj
}

func (obj *Split) HandleTransact(transact ITransaction) {
	transact.PrintInfo()
	obj.HandleSplitting(obj, transact)
}

func (obj *Split) HandleTransacts(wg *sync.WaitGroup) {
	wg.Done()
	return
}

func (obj *Split) AppendTransact(transact ITransaction) bool {
	Logger.Trace.Println("Append transact ", transact.GetId(), " to Split")
	transact.SetHolderName(obj.name)
	obj.sum_transact++
	obj.HandleTransact(transact)
	return true
}

func (obj *Split) PrintReport() {
	obj.BaseObj.PrintReport()
	fmt.Printf("Average split %.2f\n", obj.sum_split/obj.sum_transact)
	fmt.Println()
}
