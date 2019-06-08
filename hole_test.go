// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package gpss

import (
	"testing"
)

func TestHole_HandleTransact(t *testing.T) {
	pipe := NewPipeline("pipe")
	hole := NewHole("hole")
	modeltime := 5
	advance := 3
	pipe.Append(hole)
	transact := NewTransaction(pipe)
	transact.SetTiсks(advance)
	pipe.ModelTime = modeltime
	hole.HandleTransact(transact)
	if !transact.IsKilled() {
		t.Error("Transact killing, expected", true, "got", transact.IsKilled())
	}
	if hole.cntTransact != 1 {
		t.Error("Transact cnt_transact, expected", 1, "got", hole.cntTransact)
	}
	if hole.sumLife != float64(modeltime) {
		t.Error("Transact sum_life, expected", modeltime, "got", hole.sumLife)
	}
	if hole.sumAdvance != float64(advance) {
		t.Error("Transact sum_life, expected", advance, "got", hole.sumAdvance)
	}
}
