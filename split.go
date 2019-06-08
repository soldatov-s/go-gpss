// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package gpss

import (
	"fmt"
)

// HandleSplittingFunc is a splitting function signature
type HandleSplittingFunc func(obj *Split, transact *Transaction)

// A Split creates assembly set of sub-transactions of a Transaction
type Split struct {
	BaseObj
	Cntsplit        int                 // Number of related Transactions to be created
	Modificator     int                 // The count half-range
	sumSplit        float64             // Counter of sub-transactions
	sumTransact     float64             // Counter of transactions
	HandleSplitting HandleSplittingFunc // Function for splitting transaction
}

// Splitting - default splitting function
func Splitting(obj *Split, transact *Transaction) {
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

	obj.sumSplit += float64(cntsplit)
	if cntsplit == len(obj.GetDst()) {
		// Default case, cntsplit equal to length of GetDst()
		for i, v := range obj.GetDst() {
			tr := transact.Copy()
			parentID := tr.GetID()
			tr.SetID(obj.Pipe.NewID())
			tr.SetParts(i+1, cntsplit, parentID)
			v.AppendTransact(tr) // Take in mind that after Split must be only Queues
		}
	} else {
		// Another case, cntsplit can be smaller than length of GetDst()
		// Randomized selections of dst for send transact
		dsts := make([]bool, cntsplit)
		partID := 1
		for {
			for _, v := range obj.GetDst() {
				if GetRandomBool() && !dsts[partID-1] {
					tr := transact.Copy()
					parentID := tr.GetID()
					tr.SetID(obj.Pipe.NewID())
					tr.SetParts(partID, cntsplit, parentID)
					v.AppendTransact(tr)
					dsts[partID-1] = true
					partID++
					if partID > cntsplit {
						return
					}
				}
			}
		}
	}
}

// NewSplit creates new Split.
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

// HandleTransact handle transact
func (obj *Split) HandleTransact(transact *Transaction) {
	transact.PrintInfo()
	obj.HandleSplitting(obj, transact)
}

// AppendTransact append transact to object
func (obj *Split) AppendTransact(transact *Transaction) bool {
	obj.BaseObj.AppendTransact(transact)
	transact.SetHolder(obj.name)
	obj.sumTransact++
	obj.HandleTransact(transact)
	return true
}

// Report - print report about object
func (obj *Split) Report() {
	obj.BaseObj.Report()
	fmt.Printf("Average split %.2f\n", obj.sumSplit/obj.sumTransact)
	fmt.Println()
}
