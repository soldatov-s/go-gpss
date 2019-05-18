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
	pipe.modelTime = modeltime
	hole.HandleTransact(transact)
	if !transact.IsKilled() {
		t.Error("Transact killing, expected", true, "got", transact.IsKilled())
	}
	if hole.cnt_transact != 1 {
		t.Error("Transact cnt_transact, expected", 1, "got", hole.cnt_transact)
	}
	if hole.sum_life != float64(modeltime) {
		t.Error("Transact sum_life, expected", modeltime, "got", hole.sum_life)
	}
	if hole.sum_advance != float64(advance) {
		t.Error("Transact sum_life, expected", advance, "got", hole.sum_advance)
	}
}
