// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package objects

import (
	"fmt"
	"reflect"
	"sort"
	"sync"

	"github.com/soldatov-s/go-gpss/internal"
)

// Pipeline is structure for pipeline
type Pipeline struct {
	Name      string              // Pipeline name
	objects   map[string]IBaseObj // Maps of objects
	lstObject []IBaseObj          // Last object in map of objects
	ModelTime int                 // Current Model Time
	Done      chan struct{}       // Chan for done
	SimTime   int                 // Simulation time
	id        int                 // ID of new transaction
	doneHndl  []func(p *Pipeline) // Handler of done event
}

type IPipeline interface {
	// Add object to pipeline
	AddObject(obj IBaseObj) IBaseObj
}

// NewPipeline create new Pipeline
func NewPipeline(name string, doneHndl ...func(p *Pipeline)) *Pipeline {
	return &Pipeline{
		objects:  make(map[string]IBaseObj),
		Name:     name,
		Done:     make(chan struct{}),
		doneHndl: doneHndl,
	}
}

// Add object to pipeline
func (p *Pipeline) AddObject(obj ...IBaseObj) *Pipeline {
	if p.lstObject != nil {
		for _, item := range p.lstObject {
			item.SetDst(obj...)
		}
	}
	for _, item := range obj {
		item.SetPipeline(p)
		item.SetID(len(p.objects))
		p.objects[item.GetName()] = item
	}
	p.lstObject = obj
	return p
}

// Loop pipeline to selected object
func (p *Pipeline) Loop(objName string) *Pipeline {
	for _, item := range p.lstObject {
		item.SetDst(p.GetObjByName(objName))
		fmt.Println(item.GetName())
	}
	return p
}

// Loop pipeline to selected object
func (p *Pipeline) LoopByObj(obj IBaseObj) *Pipeline {
	for _, item := range p.lstObject {
		item.SetDst(obj)
		fmt.Println(obj.GetName())
	}
	return p
}

// Append object to pipeline. Src is multiple sources of transact for appended
// object.
func (p *Pipeline) Append(obj IBaseObj, dst ...IBaseObj) {
	p.AppendISlice(obj, dst)
}

// AppendMultiple - append multiple objects to pipeline.  Src is multiple sources of transact
// for appended object.
func (p *Pipeline) AppendMultiple(obj []IBaseObj, dst ...IBaseObj) {
	for _, o := range obj {
		p.AppendISlice(o, dst)
	}
}

// AppendISlice - append slice IBaseObj
func (p *Pipeline) AppendISlice(obj IBaseObj, dst []IBaseObj) {
	obj.SetDst(dst...)
	obj.SetPipeline(p)
	obj.SetID(len(p.objects))
	p.objects[obj.GetName()] = obj
}

// Delete object from pipeline
func (p *Pipeline) Delete(obj IBaseObj) {
	delete(p.objects, obj.GetName())
}

// PrintObjects - print list of objects ib pipeline
func (p *Pipeline) PrintObjects() {
	keys := make([]string, 0, len(p.objects))
	for k := range p.objects {
		keys = append(keys, k)
	}
	fmt.Println("Pipeline ", p.Name)
	for _, k := range keys {
		fmt.Println("Key:", k, "Value:", reflect.TypeOf(p.objects[k]))
	}
}

// Start simulation
func (p *Pipeline) Start(value int) {
	var wg sync.WaitGroup

	p.SimTime = value
	go func() {
		for {
			select {
			case <-p.Done:
				for _, f := range p.doneHndl {
					f(p)
				}
				return
			default:
				utils.Log.Trace.Println("ModelTime ", p.ModelTime)
				wg.Add(len(p.objects))
				for _, o := range p.objects {
					o.HandleTransacts(&wg)
				}
				wg.Wait()
				if p.ModelTime++; p.ModelTime == value {
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

// Report - print report about work of pipeline
func (p *Pipeline) Report() {
	fmt.Println("Pipeline name \"", p.Name, "\"")
	fmt.Println("Simulation time", p.ModelTime)
	sortedObjects := make([]IBaseObj, 0, len(p.objects))
	for _, v := range p.objects {
		sortedObjects = append(sortedObjects, v)
	}

	id := func(p1, p2 IBaseObj) bool {
		return p1.GetID() < p2.GetID()
	}

	By(id).Sort(sortedObjects)
	for _, v := range sortedObjects {
		v.Report()
	}
}

// GetObjByName get object from pipeline by name
func (p *Pipeline) GetObjByName(name string) IBaseObj {
	return p.objects[name]
}

// By is a signature for sort
type By func(p1, p2 IBaseObj) bool

// Sort - sort object in pipeline
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

// NewID - get ID for new transaction
func (p *Pipeline) NewID() int {
	p.id++
	return p.id
}

// EnableTrace enable trace log
// TODO: refactoring required, because logging enabled for all pipelines, not for selected
func (p *Pipeline) EnableTrace() {
	utils.EnableVerbose()
}
