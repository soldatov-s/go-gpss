// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package reports

import (
	"fmt"

	"github.com/soldatov-s/go-gpss/objects"
)

type Config struct {
}

type Report struct {
	Name   string
	Config Config
}

type ReportOfObject struct {
	Report
}

func (rp *ReportOfObject) Build(obj objects.IBaseObj) {
	fmt.Println("Object name \"", obj.GetName(), "\"")
}
