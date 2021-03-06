// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package objects

import (
	"testing"
)

func TestTransaction_GetId(t *testing.T) {
	pipe := NewPipeline("pipe")
	id := 1
	transact := NewTransaction(pipe)
	if transact.GetID() != id {
		t.Error("Transact id, expected", id, "got", transact.GetID())
	}
}
