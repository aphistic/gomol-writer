package gomolwriter

import (
	"os"

	gw "."
	"github.com/aphistic/gomol"
)

// Code for the README example to make sure it still builds!
func Example() {
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
