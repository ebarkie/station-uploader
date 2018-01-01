// Copyright (c) 2016-2018 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"log"
	"os"
)

// Loggers
var (
	Trace = log.New(ioutil.Discard, "[TRCE]", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	Debug = log.New(ioutil.Discard, "[DBUG]", log.LstdFlags|log.Lshortfile)
	Info  = log.New(os.Stdout, "[INFO]", log.LstdFlags)
	Warn  = log.New(os.Stderr, "[WARN]", log.LstdFlags|log.Lshortfile)
	Error = log.New(os.Stderr, "[ERRO]", log.LstdFlags|log.Lshortfile)
)
