// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package gpss

import (
	"fmt"
	"reflect"
	"sort"
	"sync"
)

type IPipeline interface {
	Append(obj IBaseObj, src ...IBaseObj)           // Append  object to pipeline
	AppendMultiple(obj []IBaseObj, dst ...IBaseObj) // Append  multiple objects to pipeline
	AppendISlice(obj IBaseObj, dst []IBaseObj)      // Append slice IBaseObj
	Delete(obj IBaseObj)                            // Delete object from pipeline
	Start(value int)                                // Start simulation
	Stop()                                          // Stop simulation
	GetSimTime() int                                // Get Simulation time
	GetModelTime() int                              // Get current model time
	GetObjByName(name string) IBaseObj              // Get object from pipeline by name
	GetIDNewTransaction() int                       // Get ID for new transaction
	PrintReport()                                   // Print report
}

type Pipeline struct {
	name      string              // Pipeline name
	objects   map[string]IBaseObj // Maps of objects
	modelTime int                 // Current Model Time
	Done      chan struct{}       // Chan for done
	simTime   int                 // Simulation time
	id        int                 // ID of new transaction
}

// Create new Pipeline
func NewPipeline(name string) *Pipeline {
	p := &Pipeline{}
	p.objects = make(map[string]IBaseObj)
	p.name = name
	p.Done = make(chan struct{})
	p.modelTime = 0
	p.id = 0
	return p
}

// Append object to pipeline. Src is multiple sources of transact for appended
// object.
func (p *Pipeline) Append(obj IBaseObj, dst ...IBaseObj) {
	obj.SetDst(dst)
	obj.SetPipeline(p)
	obj.SetID(len(p.objects))
	p.objects[obj.GetName()] = obj
}

// Append multiple objects to pipeline.  Src is multiple sources of transact
// for appended object.
func (p *Pipeline) AppendMultiple(obj []IBaseObj, dst ...IBaseObj) {
	for _, o := range obj {
		o.SetDst(dst)
		o.SetPipeline(p)
		o.SetID(len(p.objects))
		p.objects[o.GetName()] = o
	}
}

func (p *Pipeline) AppendISlice(obj IBaseObj, dst []IBaseObj) {
	obj.SetDst(dst)
	obj.SetPipeline(p)
	obj.SetID(len(p.objects))
	p.objects[obj.GetName()] = obj
}

// Delete object from pipeline
func (p *Pipeline) Delete(obj IBaseObj) {
	o := obj.(IBaseObj)
	delete(p.objects, o.GetName())
}

// Print list of objects ib pipeline
func (p *Pipeline) PrintObjects() {
	keys := make([]string, 0, len(p.objects))
	for k := range p.objects {
		keys = append(keys, k)
	}
	fmt.Println("Pipeline ", p.name)
	for _, k := range keys {
		fmt.Println("Key:", k, "Value:", reflect.TypeOf(p.objects[k]))
	}
}

// Start simulation
func (p *Pipeline) Start(value int) {
	var wg sync.WaitGroup

	p.simTime = value
	go func() {
		for {
			select {
			case <-p.Done:
				return
			default:
				Logger.Trace.Println("ModelTime ", p.modelTime)
				wg.Add(len(p.objects))
				for _, o := range p.objects {
					o.HandleTransacts(&wg)
				}
				wg.Wait()
				if p.modelTime++; p.modelTime == value {
					p.Stop()
				}
			}
		}
	}()
}

// Stop simulation
func (p *Pipeline) Stop() {
	close(p.Done)
}

// Print report about work of pipeline
func (p *Pipeline) PrintReport() {
	fmt.Println("Pipeline name \"", p.name, "\"")
	fmt.Println("Simulation time", p.modelTime)
	sortedObjects := make([]IBaseObj, 0, len(p.objects))
	for _, v := range p.objects {
		sortedObjects = append(sortedObjects, v)
	}

	id := func(p1, p2 IBaseObj) bool {
		return p1.GetID() < p2.GetID()
	}

	By(id).Sort(sortedObjects)
	for _, v := range sortedObjects {
		v.PrintReport()
	}

}

// Get value of simulation time
func (p *Pipeline) GetSimTime() int {
	return p.simTime
}

// Get current model time
func (p *Pipeline) GetModelTime() int {
	return p.modelTime
}

// Get object from pipeline by name
func (p *Pipeline) GetObjByName(name string) IBaseObj {
	return p.objects[name]
}

type By func(p1, p2 IBaseObj) bool

func (by By) Sort(objects []IBaseObj) {
	objs := &objectSorter{
		objects: objects,
		by:      by,
	}
	sort.Sort(objs)
}

type objectSorter struct {
	objects []IBaseObj
	by      By
}

// Len is part of sort.Interface.
func (s *objectSorter) Len() int {
	return len(s.objects)
}

// Swap is part of sort.Interface.
func (s *objectSorter) Swap(i, j int) {
	s.objects[i], s.objects[j] = s.objects[j], s.objects[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *objectSorter) Less(i, j int) bool {
	return s.by(s.objects[i], s.objects[j])
}

func (p *Pipeline) GetIDNewTransaction() int {
	p.id++
	return p.id
}
