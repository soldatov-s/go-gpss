// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package objects

import (
	// stdlib
	"fmt"
	"sync"

	"github.com/soldatov-s/go-gpss/internal"
)

// IBaseObj implements BaseObj interface
type IBaseObj interface {
	SetID(int)                          // Set object ID
	GetID() int                         // Get object ID
	GetName() string                    // Get object name
	SetDst([]IBaseObj)                  // Set dst for object
	GetDst() []IBaseObj                 // Get dst for object
	SetPipeline(pipe *Pipeline)         // Set pipeline for object
	AppendTransact(*Transaction) bool   // Append transact to object
	HandleTransacts(wg *sync.WaitGroup) // Handle all transacts of object
	Report()                            // Print report
	AddObject(obj IBaseObj) IBaseObj    // Add object to pipeline
}

// BaseObj is the base object of simulation system
type BaseObj struct {
	name    string
	objTime int
	dst     []IBaseObj
	Pipe    *Pipeline
	tb      *TransactTable
	id      int
}

// Add object to pipeline
func (o *BaseObj) AddObject(obj IBaseObj) IBaseObj {
	fmt.Printf("%+v\n", o)
	fmt.Printf("%+v\n", obj)
	o.Pipe.Append(o, obj)
	return obj
}

// Init - initializate BaseObj
func (obj *BaseObj) Init(name string) {
	obj.name = name
	obj.tb = NewTransactTable()
}

// GetName - get name of BaseObj
func (obj *BaseObj) GetName() string {
	return obj.name
}

// SetDst - set destination of BaseObj
func (obj *BaseObj) SetDst(dst []IBaseObj) {
	obj.dst = dst
}

// GetDst - get destination of BaseObj
func (obj *BaseObj) GetDst() []IBaseObj {
	return obj.dst
}

// SetPipeline - set pipeline of BaseObj
func (obj *BaseObj) SetPipeline(pipe *Pipeline) {
	obj.Pipe = pipe
}

// SetID set ID of BaseObj
func (obj *BaseObj) SetID(id int) {
	obj.id = id
}

// GetID get id of BaseObj
func (obj *BaseObj) GetID() int {
	return obj.id
}

// AppendTransact append transact to object
func (obj *BaseObj) AppendTransact(t *Transaction) bool {
	utils.Log.Trace.Println("Append transact ", t.GetID(), " to ", obj.name)
	return true
}

// HandleTransacts handle transacts
func (obj *BaseObj) HandleTransacts(wg *sync.WaitGroup) {
	wg.Done()
}

// Report - print report about object
func (obj *BaseObj) Report() {
	fmt.Println("Object name \"", obj.name, "\"")
}
