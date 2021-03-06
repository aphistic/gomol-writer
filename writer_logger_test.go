package gomolwriter

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aphistic/gomol"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type GomolSuite struct{}

var _ = Suite(&GomolSuite{})

func (s *GomolSuite) TestWriterSetTemplate(c *C) {
	var b bytes.Buffer
	wl, err := NewWriterLogger(&b, nil)
	c.Check(wl.tpl, NotNil)

	err = wl.SetTemplate(nil)
	c.Check(err, NotNil)

	tpl, err := gomol.NewTemplate("")
	c.Assert(err, IsNil)
	err = wl.SetTemplate(tpl)
	c.Check(err, IsNil)
}

func (s *GomolSuite) TestWriterInitLoggerNoWriter(c *C) {
	wl, err := NewWriterLogger(nil, nil)
	c.Assert(err, NotNil)
	c.Check(err.Error(), Equals, "An io.Writer must be provided")
	c.Check(wl, IsNil)
}

func (s *GomolSuite) TestWriterInitLoggerNoConfig(c *C) {
	var b bytes.Buffer
	wl, err := NewWriterLogger(&b, nil)
	c.Check(err, IsNil)
	c.Check(wl, NotNil)
	c.Assert(wl.config, NotNil)
	c.Check(wl.config.BufferSize, Equals, 1000)
}

func (s *GomolSuite) TestWriterInitLogger(c *C) {
	var b bytes.Buffer
	wl, err := NewWriterLogger(&b, nil)
	c.Assert(err, IsNil)
	c.Check(wl.IsInitialized(), Equals, false)
	wl.InitLogger()
	c.Check(wl.IsInitialized(), Equals, true)
}

func (s *GomolSuite) TestWriterShutdownLogger(c *C) {
	var b bytes.Buffer
	wl, err := NewWriterLogger(&b, nil)
	c.Assert(err, IsNil)
	c.Check(wl.IsInitialized(), Equals, false)
	wl.InitLogger()
	c.Check(wl.IsInitialized(), Equals, true)
	wl.ShutdownLogger()
	c.Check(wl.IsInitialized(), Equals, false)
}

func (s *GomolSuite) TestWriterWithConfig(c *C) {
	var b bytes.Buffer
	cfg := NewWriterLoggerConfig()
	cfg.BufferSize = 1
	wl, err := NewWriterLogger(&b, cfg)
	c.Assert(err, IsNil)
	c.Check(wl.config.BufferSize, Equals, 1)
}

func (s *GomolSuite) TestWriterMultipleMessages(c *C) {
	var b bytes.Buffer
	cfg := NewWriterLoggerConfig()
	wl, err := NewWriterLogger(&b, cfg)
	c.Assert(err, IsNil)
	wl.Logm(time.Now(), gomol.LevelDebug, nil, "dbg 1234")
	wl.Logm(time.Now(), gomol.LevelWarning, nil, "warn 4321")

	wl.flushMessages()

	c.Check(b.String(), Equals, "[DEBUG] dbg 1234\n[WARN] warn 4321\n")
}

func (s *GomolSuite) TestWriterFlushOnBufferSize(c *C) {
	var b bytes.Buffer
	cfg := NewWriterLoggerConfig()
	cfg.BufferSize = 2
	wl, err := NewWriterLogger(&b, cfg)
	c.Assert(err, IsNil)

	c.Check(wl.buffer, HasLen, 0)

	wl.Logm(time.Now(), gomol.LevelDebug, nil, "Message 1")
	c.Check(wl.buffer, HasLen, 1)

	wl.Logm(time.Now(), gomol.LevelDebug, nil, "Message 2")
	c.Check(wl.buffer, HasLen, 0)

	c.Check(strings.Count(b.String(), "\n"), Equals, 2)
	c.Check(b.String(), Equals, "[DEBUG] Message 1\n[DEBUG] Message 2\n")
}

func (s *GomolSuite) TestWriterToFile(c *C) {
	f, err := ioutil.TempFile("", "gomol_test_")
	if err != nil {
		c.Fatal("Unable to create temp file to test writer logger")
	}
	defer f.Close()
	defer os.Remove(f.Name())

	cfg := NewWriterLoggerConfig()
	wl, _ := NewWriterLogger(f, cfg)
	wl.InitLogger()
	wl.Logm(time.Now(), gomol.LevelDebug, nil, "Message 1")
	wl.Logm(time.Now(), gomol.LevelFatal, nil, "Message 2")
	wl.ShutdownLogger()

	fData, err := ioutil.ReadFile(f.Name())
	if err != nil {
		c.Fatal("Could not read from writer logger test file")
	}
	fStr := string(fData)
	c.Check(fStr, Equals, "[DEBUG] Message 1\n[FATAL] Message 2\n")
}

func (s *GomolSuite) TestWriterLogmNoAttrs(c *C) {
	var b bytes.Buffer
	cfg := NewWriterLoggerConfig()
	wl, err := NewWriterLogger(&b, cfg)
	c.Assert(err, IsNil)
	wl.Logm(time.Now(), gomol.LevelDebug, nil, "test")

	wl.flushMessages()

	c.Check(b.String(), Equals, "[DEBUG] test\n")
}

func (s *GomolSuite) TestWriterLogmAttrs(c *C) {
	var b bytes.Buffer
	cfg := NewWriterLoggerConfig()
	wl, err := NewWriterLogger(&b, cfg)
	c.Assert(err, IsNil)
	wl.Logm(
		time.Now(),
		gomol.LevelDebug,
		map[string]interface{}{
			"attr1": 4321,
		},
		"test 1234")

	wl.flushMessages()

	c.Check(b.String(), Equals, "[DEBUG] test 1234\n")
}

func (s *GomolSuite) TestWriterBaseAttrs(c *C) {
	var buf bytes.Buffer
	b := gomol.NewBase()
	b.SetAttr("attr1", 7890)
	b.SetAttr("attr2", "val2")

	cfg := NewWriterLoggerConfig()
	wl, err := NewWriterLogger(&buf, cfg)
	c.Check(err, IsNil)
	b.AddLogger(wl)
	wl.Logm(
		time.Now(),
		gomol.LevelDebug,
		map[string]interface{}{
			"attr1": 4321,
			"attr3": "val3",
		},
		"test 1234")

	wl.flushMessages()

	c.Check(buf.String(), Equals, "[DEBUG] test 1234\n")
}
