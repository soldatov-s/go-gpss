// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package gpss

import (
	"sync"
)

type ITransactTable interface {
	Push(item ITransaction)       // Push item to table
	Pop() ITransaction            // Pop item from table
	Remove(item ITransaction)     // Remove item from table
	GetItems() map[int]*TableItem // Get all items from table
	GetLen() int                  // Return length table
	GetItem(int) *TableItem       // Get item in table by ID
	GetFirstItem() *TableItem     // Get first item in table
	LockTable()
	UnlockTable()
}

type TableItem struct {
	transact   ITransaction
	prevoiseID int
	nextID     int
}

type TransactTable struct {
	firstID int                // ID of first transact in table
	lastID  int                // ID of last transact in table
	mp      map[int]*TableItem //Map with items
	mu      *sync.Mutex
}

// Create new TransactTable
func NewTransactTable() *TransactTable {
	obj := &TransactTable{}
	obj.lastID = -1
	obj.firstID = -1
	obj.mp = make(map[int]*TableItem)
	obj.mu = &sync.Mutex{}
	return obj
}

// Remove transact from table
func (obj *TransactTable) Remove(transact ITransaction) {
	defer obj.mu.Unlock()
	obj.mu.Lock()
	item := obj.mp[transact.GetId()]
	if item == nil {
		return
	}
	if item.prevoiseID != -1 {
		prevoiseItem := obj.mp[item.prevoiseID]
		if prevoiseItem != nil {
			prevoiseItem.nextID = item.nextID
		}
	}
	delete(obj.mp, transact.GetId())

}

// Get all items of table
func (obj *TransactTable) GetItems() map[int]*TableItem {
	defer obj.mu.Unlock()
	obj.mu.Lock()
	items := make(map[int]*TableItem)
	for k, v := range obj.mp {
		items[k] = v
	}
	return items //obj.mp
}

// Push transact to end table
func (obj *TransactTable) Push(transact ITransaction) {
	defer obj.mu.Unlock()
	obj.mu.Lock()
	if obj.firstID == -1 {
		obj.firstID = transact.GetId()
	} else {
		last_transact := obj.mp[obj.lastID]
		if last_transact != nil {
			last_transact.nextID = transact.GetId()
		}
	}
	obj.mp[transact.GetId()] = &TableItem{transact: transact, nextID: -1, prevoiseID: obj.lastID}
	obj.lastID = transact.GetId()
}

// Return first transact from table and remove it from table
func (obj *TransactTable) Pop() ITransaction {
	defer obj.mu.Unlock()
	obj.mu.Lock()
	item := obj.mp[obj.firstID]
	if item != nil {
		r := obj.firstID
		obj.firstID = item.nextID
		delete(obj.mp, r)
	}
	return item.transact
}

func (obj *TransactTable) GetLen() int {
	return len(obj.mp)
}

func (obj *TransactTable) GetItem(id int) *TableItem {
	defer obj.mu.Unlock()
	obj.mu.Lock()
	return obj.mp[id]
}

func (obj *TransactTable) GetFirstItem() *TableItem {
	defer obj.mu.Unlock()
	obj.mu.Lock()
	return obj.mp[obj.firstID]
}

func (obj *TransactTable) LockTable() {
	obj.mu.Lock()
}

func (obj *TransactTable) UnlockTable() {
	//obj.mu.Unlock()
}
