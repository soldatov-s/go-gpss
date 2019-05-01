// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package gpss

import (
	"fmt"
	"sync"
)

// IGenerator implements Generator interface
type IGenerator interface {
	GenerateBorn(obj *Generator, modelTime int) int
	GenerateTransact()
}

type HandleBornFunc func(obj *Generator) int

// A Generator sequentially generates transactions
type Generator struct {
	BaseObj
	Interval    int            // Mean inter generation time
	Modificator int            // Inter generation time half-range
	Start       int            // Start delay time
	Count       int            // Creation limit. Max count of transactions.
	id          int            // ID of new transaction
	nextborn    int            // The time when will create new transaction
	HandleBorn  HandleBornFunc // Function for generate transaction
}

// Default function for generate transaction
func GenerateBorn(obj *Generator) int {
	var born int
	born += obj.Interval
	if obj.Modificator > 0 {
		born += GetRandom(-obj.Modificator, obj.Modificator)
	}
	if obj.GetPipeline() != nil {
		born += obj.GetPipeline().GetModelTime()
	}
	return born
}

// Creates new Generator.
// name - name of object; interval - mean inter generation time;
// modificator - inter generation time half-range; start - start delay time;
// count - creation limit, max count of transactions; hndl - function for generate
// transaction
func NewGenerator(name string, interval, modificator, start, count int, hndl HandleBornFunc) *Generator {
	obj := &Generator{}
	obj.name = name
	obj.Interval = interval
	obj.Modificator = modificator
	obj.Start = start
	obj.Count = count
	obj.id = 1
	if hndl != nil {
		obj.HandleBorn = hndl
	} else {
		obj.HandleBorn = GenerateBorn
	}
	obj.nextborn = obj.HandleBorn(obj)
	return obj
}

// Generates transaction and it send into the simulation
func (obj *Generator) GenerateTransact() {
	var isTransactSended bool
	obj.GetLogger().GetTrace().Println("Generate transact ", obj.id)
	t := NewTransaction(obj.id, obj.GetPipeline())
	t.SetHolderName(obj.name)
	for _, v := range obj.GetDst() {
		isTransactSended = isTransactSended || v.AppendTransact(t)
	}
	if isTransactSended {
		obj.id++
	}
}

func (obj *Generator) HandleTransacts(wg *sync.WaitGroup) {
	if (obj.Count != 0 && obj.id > obj.Count) ||
		(obj.nextborn != obj.GetPipeline().GetModelTime()) {
		wg.Done()
		return
	}
	go func() {
		defer func() {
			obj.nextborn = obj.HandleBorn(obj)
			wg.Done()
		}()
		// Generate trasact one by one
		if obj.Count == 0 {
			obj.GenerateTransact()
			return
		}
		// Generate all transact at once
		for {
			obj.GenerateTransact()
			if obj.id > obj.Count {
				obj.GetLogger().GetTrace().Println("Stop generate")
				return
			}
		}
	}()
}

func (obj *Generator) PrintReport() {
	obj.BaseObj.PrintReport()
	fmt.Println("Generated", obj.id-1)
	fmt.Println()
}
