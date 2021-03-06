// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package objects

import (
	"sync"
)

// TableItem struct for TransactTable item
type TableItem struct {
	transact   *Transaction
	prevoiseID int
	nextID     int
}

// TransactTable is structure for description a table of transactions
type TransactTable struct {
	firstID int                // ID of first transact in table
	lastID  int                // ID of last transact in table
	mp      map[int]*TableItem // Map with items
	mu      *sync.Mutex
}

// NewTransactTable create new TransactTable
func NewTransactTable() *TransactTable {
	obj := &TransactTable{
		lastID:  -1,
		firstID: -1,
		mp:      make(map[int]*TableItem),
		mu:      &sync.Mutex{},
	}
	return obj
}

// Remove transact from table
func (obj *TransactTable) Remove(transact *Transaction) {
	defer obj.mu.Unlock()
	obj.mu.Lock()
	item := obj.mp[transact.GetID()]
	if item == nil {
		return
	}
	if item.prevoiseID != -1 {
		prevoiseItem := obj.mp[item.prevoiseID]
		if prevoiseItem != nil {
			prevoiseItem.nextID = item.nextID
		}
	}
	delete(obj.mp, transact.GetID())
}

// Items get all items of table
func (obj *TransactTable) Items() map[int]*TableItem {
	defer obj.mu.Unlock()
	obj.mu.Lock()
	items := make(map[int]*TableItem)
	for k, v := range obj.mp {
		items[k] = v
	}
	return items
}

// Push transact to end table
func (obj *TransactTable) Push(transact *Transaction) {
	defer obj.mu.Unlock()
	obj.mu.Lock()
	if obj.firstID == -1 {
		obj.firstID = transact.GetID()
	} else {
		lastTransact := obj.mp[obj.lastID]
		if lastTransact != nil {
			lastTransact.nextID = transact.GetID()
		}
	}
	obj.mp[transact.GetID()] = &TableItem{transact: transact, nextID: -1, prevoiseID: obj.lastID}
	obj.lastID = transact.GetID()
}

// Pop - return first transact from table and remove it from table
func (obj *TransactTable) Pop() *Transaction {
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

// Len - return length table
func (obj *TransactTable) Len() int {
	return len(obj.mp)
}

// Item - get item in table by ID
func (obj *TransactTable) Item(id int) *TableItem {
	defer obj.mu.Unlock()
	obj.mu.Lock()
	return obj.mp[id]
}

// First - get first item in table
func (obj *TransactTable) First() *TableItem {
	defer obj.mu.Unlock()
	obj.mu.Lock()
	return obj.mp[obj.firstID]
}
