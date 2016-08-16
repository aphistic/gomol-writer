gomol-writer
============

[![GoDoc](https://godoc.org/github.com/aphistic/gomol-writer?status.svg)](https://godoc.org/github.com/aphistic/gomol-writer)
[![Build Status](https://img.shields.io/travis/aphistic/gomol-writer.svg)](https://travis-ci.org/aphistic/gomol-writer)
[![Code Coverage](https://img.shields.io/codecov/c/github/aphistic/gomol-writer.svg)](http://codecov.io/github/aphistic/gomol-writer?branch=master)

gomol-writer is a logger for [gomol](https://github.com/aphistic/gomol) to support logging to an `io.Writer`.

Installation
============

The recommended way to install is via http://gopkg.in

    go get gopkg.in/aphistic/gomol-writer.v0
    ...
    import "gopkg.in/aphistic/gomol-writer.v0"

gomol-writer can also be installed the standard way as well

    go get github.com/aphistic/gomol-writer
    ...
    import "github.com/aphistic/gomol-writer"

Examples
========

For brevity a lot of error checking has been omitted, be sure you do your checks!

This is a super basic example of adding an io.Writer logger to gomol and then logging a few messages:

```go
package main

import (
	"os"
	"github.com/aphistic/gomol"
	gw "github.com/aphistic/gomol-writer"
)

func main() {
	// Add an io.Writer logger
	writerCfg := gw.NewWriterLoggerConfig()
	writerLogger, _ := gw.NewWriterLogger(os.Stdout, writerCfg)
	gomol.AddLogger(writerLogger)

	// Set some global attrs that will be added to all
	// messages automatically
	gomol.SetAttr("facility", "gomol.example")
	gomol.SetAttr("another_attr", 1234)

	// Initialize the loggers
	gomol.InitLoggers()
	defer gomol.ShutdownLoggers()

	// Log some debug messages with message-level attrs
	// that will be sent only with that message
	for idx := 1; idx <= 10; idx++ {
		gomol.Dbgm(
			gomol.NewAttrs().
				SetAttr("msg_attr1", 4321),
			"Test message %v", idx)
	}
}

```
