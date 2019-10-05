// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package objects

import (
	"fmt"
)

// Count counts all Transactions which pass through the block, it present in two parts,
// first for increment Count value, second for decrement Count value
type Count struct {
	BaseObj
	value  *int // Value of counter
	incDec int  // Value of increment/decrement
}

// NewCount creates two objects, for incremet and decrement. After enter transact in
// inc_obj, value incremented by inc_value. After enter transact in dec_obj,
// value decremented by dec_value.
func NewCount(name string, incValue, decValue int) (*Count, *Count) {
	value := 0
	inc := &Count{}
	dec := &Count{}
	inc.name = name + "_INC"
	inc.value = &value
	inc.incDec = incValue
	dec.name = name + "_DEC"
	dec.value = inc.value
	dec.incDec = decValue
	return inc, dec
}

// AppendTransact append transact to object
func (obj *Count) AppendTransact(transact *Transaction) bool {
	for _, v := range obj.GetDst() {
		if v.AppendTransact(transact) {
			*obj.value += obj.incDec
			obj.BaseObj.AppendTransact(transact)
			return true
		}
	}
	return false
}

// Report - print report about object
func (obj *Count) Report() {
	fmt.Printf("Count value %d\n", obj.value)
}
