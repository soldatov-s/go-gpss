// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package gpss

import (
	"math/rand"
	"time"
)

// Generate random between min and max
func GetRandom(min, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max-min+1) + min
}

// Get random bool
func GetRandomBool() bool {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Float32() < 0.5
}
