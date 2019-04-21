// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package gpss

import (
	"io"
	"log"
)

type ILogger interface {
	GetTrace() *log.Logger
	GetInfo() *log.Logger
	GetWarning() *log.Logger
	GetError() *log.Logger
}

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
	logger := &Logger{}

	logger.Trace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	logger.Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	logger.Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	logger.Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	return logger
}

func (logger *Logger) GetTrace() *log.Logger {
	return logger.Trace
}

func (logger *Logger) GetInfo() *log.Logger {
	return logger.Info
}

func (logger *Logger) GetWarning() *log.Logger {
	return logger.Warning
}

func (logger *Logger) GetError() *log.Logger {
	return logger.Error
}
