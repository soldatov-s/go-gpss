// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package gpss

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

var Log *Logger = NewLogger(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

type Logger struct {
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
}

func NewLogger(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) *Logger {
	return &Logger{
		Trace: log.New(traceHandle, "TRACE: ",
			log.Ldate|log.Ltime|log.Lshortfile),
		Info: log.New(infoHandle, "INFO: ",
			log.Ldate|log.Ltime|log.Lshortfile),
		Warning: log.New(warningHandle, "WARNING: ",
			log.Ldate|log.Ltime|log.Lshortfile),
		Error: log.New(errorHandle, "ERROR: ",
			log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func EnableVerbose() {
	Log = NewLogger(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
}
