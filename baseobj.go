// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package gpss

import (
	"fmt"
	"sync"
)

type IBaseObj interface {
	SetID(int)                          // Set object ID
	GetID() int                         // Get object ID
	GetName() string                    // Get object name
	SetDst([]IBaseObj)                  // Set dst for object
	GetDst() []IBaseObj                 // Get dst for object
	SetPipeline(pipe *Pipeline)         // Set pipeline for object
	GetPipeline() *Pipeline             // Get pipeline for object
	AppendTransact(*Transaction) bool   // Append transact to object
	HandleTransacts(wg *sync.WaitGroup) // Handle all transacts of object
	PrintReport()                       // Print report
}

type BaseObj struct {
	name    string
	objTime int
	dst     []IBaseObj
	pipe    *Pipeline
	tb      *TransactTable
	id      int
}

func (obj *BaseObj) Init(name string) {
	obj.name = name
	obj.tb = NewTransactTable()
}

func (obj *BaseObj) GetName() string {
	return obj.name
}

func (obj *BaseObj) SetDst(dst []IBaseObj) {
	obj.dst = dst
}

func (obj *BaseObj) GetDst() []IBaseObj {
	return obj.dst
}

func (obj *BaseObj) SetPipeline(pipe *Pipeline) {
	obj.pipe = pipe
}

func (obj *BaseObj) GetPipeline() *Pipeline {
	return obj.pipe
}

func (obj *BaseObj) SetID(id int) {
	obj.id = id
}

func (obj *BaseObj) GetID() int {
	return obj.id
}

func (obj *BaseObj) GetTransactTable() *TransactTable {
	return obj.tb
}

func (obj *BaseObj) AppendTransact(t *Transaction) bool {
	obj.tb.Push(t)
	return true
}

func (obj *BaseObj) HandleTransacts(wg *sync.WaitGroup) {
	wg.Done()
}

func (obj *BaseObj) PrintReport() {
	fmt.Println("Object name \"", obj.name, "\"")
}
